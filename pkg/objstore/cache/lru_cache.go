package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/blevesearch/mmap-go"
	segment "github.com/blugelabs/bluge_segment_api"
	"github.com/rs/zerolog/log"
	"github.com/zincsearch/zincsearch/pkg/config"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"sort"
	"sync"
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
	aTime    time.Time
	manager  *CacheManager
	size     int64
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

func (c *CacheManager) doCleanup() error {
	var cacheSlice []*CachedFile
	var curSize int64 = 0
	var targetSize = float64(c.maxSize) * 0.7

	c.lock.RLock()
	for _, cache := range c.caches {
		curSize += cache.size
		if cache.refCount != 0 {
			continue
		}
		cacheSlice = append(cacheSlice, cache)
	}

	if len(cacheSlice) == 0 {
		log.Debug().Msgf("clean up finished: there are not unused caches")
		c.lock.RUnlock()
		return nil
	}

	if float64(curSize) < targetSize {
		log.Debug().Msgf("clean up finished: cache size is below max")
		c.lock.RUnlock()
		return nil
	}

	sort.Slice(cacheSlice, func(i, j int) bool {
		return cacheSlice[i].aTime.Before(cacheSlice[j].aTime)
	})
	c.lock.RUnlock()

	c.lock.Lock()
	defer c.lock.Unlock()

	for _, cache := range cacheSlice {
		delete(c.caches, cache.path)
		err := os.Remove(cache.path)
		if err != nil {
			return err
		}
		curSize -= cache.size

		if float64(curSize) <= targetSize {
			break
		}
	}
	return nil
}

func (c *CacheManager) TryLoad(filepath string) (*segment.Data, io.Closer, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if f, ok := c.caches[filepath]; ok {
		return f.loadReadOnlyData()
	}
	return nil, nil, ErrFileNotBeCached
}

func (c *CacheManager) Load(filepath string) (*segment.Data, io.Closer, error) {
	c.lock.Lock()
	if f, ok := c.caches[filepath]; ok {
		f.aTime = time.Now()
		f.refCount++
		c.lock.Unlock()
		return f.loadReadOnlyData()
	}
	panic(fmt.Sprintf("try to load not cached file %s", filepath))
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

func (c *CacheManager) AddCache(cf *CachedFile) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.caches[cf.path] = cf
}

func (c *CacheManager) GetCache(fileName string) *CachedFile {
	cf, _ := c.caches[fileName]
	return cf
}

func (c *CacheManager) CreateCache(filepath string) *CachedFile {
	return &CachedFile{
		path:     filepath,
		refCount: 0,
		manager:  c,
		aTime:    time.Now(),
	}
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
	c.aTime = time.Now()
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

func (c *CachedFile) UpdateSize() {
	fi, err := os.Stat(c.path)
	if err != nil {
		panic(fmt.Errorf("cache file error: %w", err))
	}
	c.manager.lock.Lock()
	c.size = fi.Size()
	c.manager.lock.Unlock()
}
