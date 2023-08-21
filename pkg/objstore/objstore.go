package objstore

import (
	"bytes"
	"fmt"
	blugeindex "github.com/blugelabs/bluge/index"
	segment "github.com/blugelabs/bluge_segment_api"
	"github.com/rs/zerolog/log"
	"github.com/zincsearch/zincsearch/pkg/objstore/cache"
	"golang.org/x/sys/unix"
	"gopkg.in/ini.v1"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
)

const pidFilename = "bluge.pid"

type ObjStore struct {
	backend storageBackend
	prefix  string
	path    string
	pid     *os.File
	lock    sync.Mutex
	cache   *cache.CacheManager
}

// ObjectInfo contains metadata for an object
type ObjectInfo struct {
	Key          string
	Size         int64
	LastModified int64
}

// storageBackend is the interface implemented by storage backends.
// An object store may have one or multiple storage backends.
type storageBackend interface {
	// Read an object from backend and write the contents into w.
	read(key string, w io.Writer) error
	// Write the contents from r to the object.
	write(key string, r io.Reader) error
	// check if the bucket exists
	bucketExists() (bool, error)
	// create new bucket
	createBucket() error
	// listObjects returns object info under the prefix.
	// The prefix allows following formats:
	// 1. "", all objects will be listed.
	// 2. "{UUID}/", objects that have the given UUID as the first part of key
	//    will be listed. Otherwise, the returned result is depend on backend
	//    type. For example, the FS backend may returns a not exist error, the
	//    OSS and S3 backend may return objects which key matches the prefix.
	// 3. "{UUID}/*", objects that have the given prefix will be listed. The
	//    rest part after UUID should follow consistent rules while using.
	// If accurateCtime is true, the LastModified field of ObjectInfo will be
	// set to:
	//   * the ctime value set by writeWithCtime()
	//   * the last modified time of the object, otherwise
	listObjects(prefix string, accurateCtime bool) ([]ObjectInfo, error)
	// remove single object
	remove(key string) error
}

// New returns a new object store.
func New(dataPath string, indexName string, confPath string) (*ObjStore, error) {
	config, err := ini.Load(confPath)
	if err != nil {
		err := fmt.Errorf("failed to load dtable.conf: %v", err)
		return nil, err
	}

	var backendType = "s3"
	var section *ini.Section
	section, err = config.GetSection("storage backend")
	if err != nil {
		return nil, err
	}

	if key, err := section.GetKey("type"); err == nil {
		backendType = key.String()
	}

	obj := new(ObjStore)
	if backendType == "oss" {
		obj.backend, err = createOSSBackend(section)
		if err != nil {
			err := fmt.Errorf("failed to create OSS backend: %v", err)
			return nil, err
		}
	} else if backendType == "s3" {
		obj.backend, err = createS3Backend(section)
		if err != nil {
			err := fmt.Errorf("failed to create S3 backend: %v", err)
			return nil, err
		}
	} else {
		err := fmt.Errorf("unsupported backend type: %s", backendType)
		return nil, err
	}
	obj.cache = cache.Manager
	obj.prefix = indexName
	obj.path = path.Join(dataPath, indexName)
	return obj, nil
}

func (o *ObjStore) Setup(readOnly bool) error {
	dirExists, err := o.dirExists()
	if err != nil {
		return fmt.Errorf("error checking if directory exists '%s': %w", o.path, err)
	}
	if !dirExists {
		if readOnly {
			return fmt.Errorf("readOnly, directory does not exist")
		}
		err = os.MkdirAll(o.path, 0777)
		if err != nil {
			return fmt.Errorf("error creating directory '%s': %w", o.path, err)
		}
	}

	bucketExists, err := o.backend.bucketExists()
	if err != nil {
		return fmt.Errorf("error checking if bucket exists: %w", err)
	}
	if !bucketExists {
		if readOnly {
			return fmt.Errorf("readOnly, bucket does not exist")
		}
		err := o.backend.createBucket()
		if err != nil {
			return fmt.Errorf("error creating bucket '%s': %w", err)
		}
	}

	return nil
}

func (o *ObjStore) dirExists() (bool, error) {
	_, err := os.Stat(o.path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (o *ObjStore) List(kind string) ([]uint64, error) {
	list, err := o.backend.listObjects(o.prefix, false)
	if err != nil {
		return nil, err
	}
	var itemList []uint64

	for _, item := range list {
		if filepath.Ext(item.Key) != kind {
			continue
		}
		stringID := filepath.Base(item.Key)
		stringID = stringID[:len(stringID)-len(kind)]
		parsedID, err := strconv.ParseUint(stringID, 16, 64)
		if err != nil {
			log.Error().Err(err).Msg("List: failed to parse object id: ")
			continue
		}
		itemList = append(itemList, parsedID)
	}
	return itemList, nil
}

func (o *ObjStore) Load(kind string, id uint64) (*segment.Data, io.Closer, error) {
	key := fileName(kind, id)

	err := o.PrepareFile(key)
	if err != nil {
		return nil, nil, err
	}
	return o.cache.Load(path.Join(o.path, key))
}

// PrepareFile
// Prevent duplicate downloads from obj_store
func (o *ObjStore) PrepareFile(fileName string) error {
	fipath := path.Join(o.path, fileName)

	if o.cache.Exists(fipath) {
		return nil
	}

	o.lock.Lock()
	defer o.lock.Unlock()

	if o.cache.Exists(fipath) {
		return nil
	}

	// clear exists file, sync data from obj store.
	_ = o.cache.Remove(fipath)
	cf := o.cache.CreateCache(fipath)

	fi, err := cf.GetWriter()
	defer cf.Close()

	if err != nil {
		o.cache.Remove(fipath)
		return err
	}
	defer fi.Close()

	err = o.backend.read(path.Join(o.prefix, fileName), fi)
	if err != nil {
		o.cache.Remove(fipath)
		return err
	}

	cf.UpdateSize()
	o.cache.AddCache(cf)
	return nil

}

func (o *ObjStore) Persist(kind string, id uint64, w blugeindex.WriterTo, closeCh chan struct{}) error {
	key := fileName(kind, id)
	fipath := path.Join(o.path, key)
	backendKey := path.Join(o.prefix, key)
	errCh1 := make(chan error, 1)
	errCh2 := make(chan error, 1)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	var f *os.File

	cleanup := func() {
		_ = (*f).Close()
		o.cache.Remove(fipath)
		_ = o.backend.remove(backendKey)
	}

	go func(ch chan error) {
		defer wg.Done()
		buffer := bytes.Buffer{}
		_, err := w.WriteTo(&buffer, closeCh)
		if err != nil {
			ch <- err
			return
		}
		err = o.backend.write(backendKey, &buffer)
		if err != nil {
			ch <- err
			return
		}
		close(ch)
	}(errCh1)

	go func(ch chan error) {
		defer wg.Done()
		exists := o.cache.Exists(fipath)
		var cf *cache.CachedFile
		var err error

		if !exists {
			cf = o.cache.CreateCache(fipath)
			f, err = cf.GetWriter()
		} else {
			cf = o.cache.GetCache(fipath)
			f, err = cf.GetWriter()
		}
		if err != nil {
			ch <- err
			return
		}

		defer func() {
			_ = (*f).Close()
			cf.Close()
		}()

		_, err = w.WriteTo(f, closeCh)
		if err != nil {
			ch <- err
			return
		}

		err = (*f).Sync()
		if err != nil {
			ch <- err
			return
		}

		if !exists {
			o.cache.AddCache(cf)
		}
		cf.UpdateSize()
		close(ch)
	}(errCh2)

	wg.Wait()

	err1, ok := <-errCh1
	if ok {
		log.Warn().Err(err1).Msg("persist to obj store failed")
	}

	err2, ok := <-errCh2
	if ok {
		log.Warn().Err(err2).Msg("persist to fs cache failed")
	}

	if err1 != nil || err2 != nil {
		cleanup()
	}

	return nil
}

func (o *ObjStore) Remove(kind string, id uint64) error {
	key := fileName(kind, id)

	err := o.backend.remove(path.Join(o.prefix, key))

	err2 := o.cache.Remove(path.Join(o.path, key))
	if err == nil {
		err = err2
	}
	return err
}

func (o *ObjStore) Stats() (numItems uint64, numBytes uint64) {
	objs, err := o.backend.listObjects(o.prefix, false)
	if err != nil {
		log.Warn().Err(err).Msg("could not get obj_store stats")
		return 0, 0
	}
	objectCount := uint64(0)
	sizeOfObjects := uint64(0)

	for _, obj := range objs {
		size := uint64(obj.Size)
		objectCount++
		sizeOfObjects += size
	}
	return objectCount, sizeOfObjects
}

func (o *ObjStore) Sync() error {
	dir, err := os.Open(o.path)
	if err != nil {
		return fmt.Errorf("error opening directory for sync: %w", err)
	}
	err = dir.Sync()
	if err != nil {
		_ = dir.Close()
		return fmt.Errorf("error syncing directory: %w", err)
	}
	err = dir.Close()
	if err != nil {
		return fmt.Errorf("error closing directing after sync: %w", err)
	}
	return nil
}

func (o *ObjStore) Lock() error {
	pidPath := filepath.Join(o.path, pidFilename)
	var err error
	o.pid, err = os.OpenFile(pidPath, os.O_CREATE|os.O_RDWR, 0777)
	err = unix.Flock(int(o.pid.Fd()), unix.LOCK_EX|unix.LOCK_NB)
	if err != nil {
		_ = o.pid.Close()
		return err
	}
	if err != nil {
		return fmt.Errorf("unable to obtain exclusive access: %w", err)
	}
	err = o.pid.Truncate(0)
	if err != nil {
		return fmt.Errorf("error truncating pid file: %w", err)
	}
	_, err = o.pid.Write([]byte(fmt.Sprintf("%d\n", os.Getpid())))
	if err != nil {
		return fmt.Errorf("error writing pid: %w", err)
	}
	err = o.pid.Sync()
	if err != nil {
		return fmt.Errorf("error syncing pid file: %w", err)
	}
	return nil
}

func (o *ObjStore) Unlock() error {
	pidPath := filepath.Join(o.path, pidFilename)
	var err error
	err = o.pid.Close()
	if err != nil {
		return fmt.Errorf("error closing pid file: %w", err)
	}
	err = os.RemoveAll(pidPath)
	if err != nil {
		return fmt.Errorf("error removing pid file: %w", err)
	}
	return err
}

func fileName(kind string, id uint64) string {
	return fmt.Sprintf("%012x", id) + kind
}
