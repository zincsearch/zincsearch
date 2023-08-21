package cache

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

type structTestFile struct {
}

func TestLru(t *testing.T) {
	temp := t.TempDir()
	defer os.Remove(temp)
	m := CacheManager{
		caches:   make(map[string]*CachedFile),
		maxSize:  1000,
		rootPath: temp,
	}

	// this will be removed, it's the first access file
	cache1 := &CachedFile{
		refCount: 0,
		path:     path.Join(temp, "t1.seg"),
	}
	// this will be remained, cause refCount
	cache2 := &CachedFile{
		refCount: 1,
		path:     path.Join(temp, "t2.seg"),
	}

	// this will be removed
	cache3 := &CachedFile{
		refCount: 0,
		path:     path.Join(temp, "t3.seg"),
	}
	// this will be remained, it's under 70% max
	cache4 := &CachedFile{
		refCount: 0,
		path:     path.Join(temp, "t4.seg"),
	}

	m.caches[cache1.path] = cache1
	m.caches[cache2.path] = cache2
	m.caches[cache3.path] = cache3
	m.caches[cache4.path] = cache4

	for _, c := range m.caches {
		fi, _ := os.OpenFile(c.path, os.O_CREATE|os.O_RDWR, 0777)
		b := bytes.Buffer{}
		for i := 0; i < 250; i++ {
			b.Write([]byte("1"))
		}
		_, _ = fi.Write(b.Bytes())
		_ = fi.Close()
	}

	err := m.doCleanup()
	assert.Nil(t, err)

	assert.False(t, m.Exists(cache1.path))
	assert.True(t, m.Exists(cache2.path))

	assert.False(t, m.Exists(cache3.path))
	assert.True(t, m.Exists(cache4.path))
}
