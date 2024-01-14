/* Copyright 2022 Zinc Labs Inc. and Contributors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package redo

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path"
	"sync"
)

/*
 * Log file format
 * |---------------------------------------------------------------|
 * |     uint64      | ... |       uint64       | 64 FIXED LENGTH  |
 * |-- key length -- | key | -- value length -- | value            |
 * |---------------------------------------------------------------|
 */

var (
	ErrNotFound      = errors.New("not found")
	ErrValueTooLarge = errors.New("value too large")
)

const ValueFixedLength = 64

type Log struct {
	index map[uint64]int
	data  []byte
	f     *os.File
	opt   *Options
	lock  sync.RWMutex
}

type Options struct {
	NoSync bool
	NoCopy bool
}

func DefaultOptions() *Options {
	return &Options{}
}

func Open(name string, opt *Options) (*Log, error) {
	err := os.MkdirAll(path.Dir(name), 0o755)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	n := 0
	index := make(map[uint64]int)
	for n < len(data) {
		key := binary.LittleEndian.Uint64(data[n : n+8])
		n += 8 // + key length
		index[key] = n
		n += 8                // + value length
		n += ValueFixedLength // + value
	}

	if len(data) == 0 {
		data = make([]byte, 0, 256)
	}

	if opt == nil {
		opt = DefaultOptions()
	}

	return &Log{f: f, opt: opt, data: data, index: index}, nil
}

func (l *Log) Write(index uint64, data []byte) error {
	dataLen := len(data)
	if dataLen > ValueFixedLength {
		return ErrValueTooLarge
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	offset, ok := l.index[index]
	if ok {
		// overwrite
		binary.LittleEndian.PutUint64(l.data[offset:offset+8], uint64(dataLen))
		copy(l.data[offset+8:offset+8+ValueFixedLength], data)
	} else {
		// need to append
		offset = len(l.data)
		l.data = append(l.data, make([]byte, 8+8+ValueFixedLength)...)
		binary.LittleEndian.PutUint64(l.data[offset:offset+8], index)
		offset += 8
		binary.LittleEndian.PutUint64(l.data[offset:offset+8], uint64(dataLen))
		copy(l.data[offset+8:offset+8+ValueFixedLength], data)
		l.index[index] = offset
	}
	// sync to disk
	_, _ = l.f.Seek(0, io.SeekStart)
	_, err := l.f.Write(l.data)
	if err != nil {
		return err
	}

	if l.opt.NoSync {
		return nil
	}
	return l.f.Sync()
}

func (l *Log) Read(index uint64) ([]byte, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	offset, ok := l.index[index]
	if !ok {
		return nil, ErrNotFound
	}
	valueLen := binary.LittleEndian.Uint64(l.data[offset : offset+8])
	if l.opt.NoCopy {
		return l.data[offset+8 : offset+8+int(valueLen)], nil
	}
	data := make([]byte, valueLen)
	copy(data, l.data[offset+8:offset+8+int(valueLen)])
	return data, nil
}

func (l *Log) Close() error {
	l.lock.Lock()
	l.index = nil
	l.data = nil
	l.lock.Unlock()
	if err := l.f.Sync(); err != nil {
		return err
	}
	return l.f.Close()
}
