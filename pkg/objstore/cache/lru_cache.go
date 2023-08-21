package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/blevesearch/mmap-go"
	"github.com/blugelabs/bluge/index"
	segment "github.com/blugelabs/bluge_segment_api"
	"github.com/rs/zerolog/log"
	"github.com/zincsearch/zincsearch/pkg/config"
	"golang.org/x/sys/unix"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"
	"syscall"
	"time"
)

var (
	ErrFileNotBeCached = errors.New("the file is not be cached")
	Manager            *CacheManager
)

type CacheManager struct {
	pid      *os.File
	caches   map[string]*CachedFile
	rootPath string
	maxSize  int64
	lock     sync.RWMutex
	closer   context.Context
	close    context.CancelFunc
	wg       sync.WaitGroup
}

type CachedFile struct {
	path     string
	refCount uint
	manager  *CacheManager
}

type closerFunc func() error

func (c closerFunc) Close() error {
	return c()
}

func init() {
	ctx, closer := context.WithCancel(context.Background())
	Manager = &CacheManager{
		caches:   make(map[string]*CachedFile),
		maxSize:  config.Global.MaxCacheSize,
		rootPath: config.Global.DataPath,
		closer:   ctx,
		close:    closer,
	}
	err := Manager.InitCache()
	if err != nil {
		panic(fmt.Sprintf("init cache error: %v", err))
	}
	go Manager.cleanup()
}

func Close() {
	if Manager == nil {
		return
	}
	if config.Global.StorageType != config.S3Storage {
		return
	}

	Manager.close()
	Manager.wg.Wait()
}

func (c *CacheManager) InitCache() error {
	err := filepath.Walk(c.rootPath, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return err
		}

		ext := path.Ext(info.Name())
		if ext != index.ItemKindSegment && ext != index.ItemKindSnapshot {
			return err
		}

		c.TraceCache(p)

		return err
	})
	return err
}

func (c *CacheManager) cleanup() {
	c.wg.Add(1)
	defer c.wg.Done()

	tick := time.NewTicker(time.Minute * 5)
	for {
		select {
		case <-tick.C:
			err := c.doCleanup()
			if err != nil {
				log.Error().Err(err).Msgf("cleanup error: %s", err)
			}
		case <-c.closer.Done():
			return
		}
	}
}

type tempFile struct {
	aTime time.Time
	path  string
	size  int64
}

func (c *CacheManager) doCleanup() error {
	var tempFiles []tempFile
	var curSize int64 = 0
	var targetSize = float64(c.maxSize) * 0.7

	err := filepath.Walk(c.rootPath, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return err
		}

		ext := path.Ext(info.Name())
		if ext != index.ItemKindSegment && ext != index.ItemKindSnapshot {
			return err
		}
		fi, err1 := os.Stat(p)
		if err1 != nil {
			return err1
		}

		curSize += fi.Size()

		statT := fi.Sys().(*syscall.Stat_t)

		tempFiles = append(tempFiles, tempFile{
			aTime: timespecToTime(statT.Atimespec),
			size:  fi.Size(),
			path:  p,
		})

		return nil
	})

	if err != nil {
		return err
	}

	if len(tempFiles) == 0 {
		log.Debug().Msgf("clean up finished: there are not unused caches")
		return nil
	}

	if float64(curSize) < targetSize {
		log.Debug().Msgf("clean up finished: cache size is below max")
		c.lock.RUnlock()
		return nil
	}

	sort.Slice(tempFiles, func(i, j int) bool {
		return tempFiles[i].aTime.Before(tempFiles[j].aTime)
	})

	c.lock.Lock()
	defer c.lock.Unlock()

	for _, f := range tempFiles {
		if cf, ok := c.caches[f.path]; ok {
			if cf.refCount > 0 {
				continue
			}
		}
		delete(c.caches, f.path)
		err := os.Remove(f.path)
		if err != nil {
			return err
		}
		curSize -= f.size

		if float64(curSize) <= targetSize {
			break
		}
	}
	return nil
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(ts.Sec, ts.Nsec)
}

func (c *CacheManager) Load(filepath string) (*segment.Data, io.Closer, error) {
	c.lock.Lock()
	if f, ok := c.caches[filepath]; ok {
		f.refCount++
		c.lock.Unlock()
		return f.loadReadOnlyData()
	}
	return nil, nil, ErrFileNotBeCached
}

func (c *CacheManager) Remove(filepath string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.caches, filepath)
	return os.Remove(filepath)
}

func (c *CacheManager) Exists(fileName string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.caches[fileName]
	return ok
}

func (c *CacheManager) GetCache(fileName string) *CachedFile {
	cf, _ := c.caches[fileName]
	return cf
}

func (c *CacheManager) TraceCache(filepath string) *CachedFile {
	cf := &CachedFile{
		path:     filepath,
		refCount: 0,
		manager:  c,
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	c.caches[filepath] = cf

	return cf
}

func (c *CachedFile) loadReadOnlyData() (*segment.Data, io.Closer, error) {

	f, err := os.OpenFile(c.path, os.O_RDONLY, 0)
	if err != nil {
		return nil, nil, err
	}
	err = unix.Flock(int(f.Fd()), unix.LOCK_SH|unix.LOCK_NB)
	if err != nil {
		_ = f.Close()
		return nil, nil, err
	}

	mm, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		// mmap failed, try to close the file
		_ = f.Close()
		return nil, nil, err
	}

	return segment.NewDataBytes(mm), c.getCloser(mm, f), nil
}

func (c *CachedFile) getCloser(mm mmap.MMap, f *os.File) closerFunc {
	cf := func() error {
		err := mm.Unmap()
		// try to close file even if unmap failed
		err2 := f.Close()
		if err == nil {
			// try to return first error
			err = err2
		}
		c.Close()
		return err
	}
	return cf
}

func (c *CachedFile) GetWriter() (*os.File, error) {
	c.manager.lock.Lock()
	c.refCount++
	c.manager.lock.Unlock()

	fi, err := os.OpenFile(c.path, os.O_CREATE|os.O_RDWR, 0777)
	err = unix.Flock(int(fi.Fd()), unix.LOCK_EX|unix.LOCK_NB)
	if err != nil {
		_ = fi.Close()
		return nil, err
	}
	return fi, err
}

func (c *CachedFile) Close() {
	c.manager.lock.Lock()
	c.refCount--
	c.manager.lock.Unlock()
}
