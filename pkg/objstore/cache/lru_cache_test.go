package cache

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
	"time"
)

func TestLru(t *testing.T) {
	temp := t.TempDir()
	defer os.Remove(temp)
	m := CacheManager{
		caches:   make(map[string]*CachedFile),
		maxSize:  1000,
		rootPath: temp,
	}
	// this will be remained because it's already under 70%
	cache1 := &CachedFile{
		size:     400,
		aTime:    time.Now(),
		refCount: 0,
		path:     path.Join(temp, "t1"),
	}
	// this will be removed, cause atime
	cache2 := &CachedFile{
		size:     200,
		aTime:    time.Now().AddDate(0, 0, -1),
		refCount: 0,
		path:     path.Join(temp, "t2"),
	}

	// this will be remained
	cache3 := &CachedFile{
		size:     200,
		aTime:    time.Now().AddDate(0, 0, -1),
		refCount: 1,
		path:     path.Join(temp, "t3"),
	}
	// this will be removed
	cache4 := &CachedFile{
		size:     300,
		aTime:    time.Now().AddDate(0, 0, -2),
		refCount: 0,
		path:     path.Join(temp, "t4"),
	}

	m.caches[cache1.path] = cache1
	m.caches[cache2.path] = cache2
	m.caches[cache3.path] = cache3
	m.caches[cache4.path] = cache4

	for _, c := range m.caches {
		fi, _ := os.OpenFile(c.path, os.O_CREATE|os.O_RDWR, 0777)
		_, _ = fi.Write([]byte("test file"))
		_ = fi.Close()
	}

	err := m.doCleanup()
	assert.Nil(t, err)

	assert.True(t, m.Exists(cache1.path))
	assert.True(t, m.Exists(cache3.path))

	assert.False(t, m.Exists(cache2.path))
	assert.False(t, m.Exists(cache4.path))
}
