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

package wal

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var l *Log

func TestMain(m *testing.M) {
	var err error
	l, err = Open("test")
	if err != nil {
		log.Fatal(err)
	}

	m.Run()

	if err = l.Close(); err != nil {
		log.Fatal(err)
	}
}

func TestWAL(t *testing.T) {
	var err error
	err = l.Write([]byte("test"))
	assert.NoError(t, err)

	var data []byte
	data, err = l.Read(1)
	assert.NoError(t, err)
	assert.NotNil(t, data)
}

func BenchmarkWAL(b *testing.B) {
	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = l.Write([]byte("test"))
		assert.NoError(b, err)
	}
}
