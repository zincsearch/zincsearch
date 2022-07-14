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
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var l *Log

func TestMain(m *testing.M) {
	var err error
	l, err = Open("data/redoTest", nil)
	if err != nil {
		log.Fatal(err)
	}

	m.Run()

	if err = l.Close(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func TestOpenClose(t *testing.T) {
	l, err := Open("data/redoTest2", nil)
	assert.NoError(t, err)
	err = l.Close()
	assert.NoError(t, err)
}

func TestLog(t *testing.T) {
	type args struct {
		index uint64
		data  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				index: 1,
				data:  []byte("test1"),
			},
		},
		{
			name: "2",
			args: args{
				index: 2,
				data:  []byte("test2"),
			},
		},
		{
			name: "2",
			args: args{
				index: 2,
				data:  []byte("test2-2"),
			},
		},
		{
			name: "3",
			args: args{
				index: 3,
				data:  []byte("test3"),
			},
		},
		{
			name: "3",
			args: args{
				index: 3,
				data:  []byte("test3-2"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := l.Write(tt.args.index, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Log.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := l.Read(tt.args.index)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.data, got)
			}
		})
	}
}

func BenchmarkLogWrite(b *testing.B) {
	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = l.Write(1, []byte("test"))
		assert.NoError(b, err)
	}
}

func BenchmarkLogRead(b *testing.B) {
	var err error
	var data []byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, err = l.Read(1)
		assert.NoError(b, err)
		assert.NotNil(b, data)
	}
}
