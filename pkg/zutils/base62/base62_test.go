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

package base62

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase62(t *testing.T) {
	n := int64(1517153236107137025)
	s := "1O4zPvQMvmh"
	t.Run("base62:Encode", func(t *testing.T) {
		assert.Equal(t, s, Encode(n))
	})

	t.Run("base62:Decode", func(t *testing.T) {
		assert.Equal(t, n, Decode(s))
	})
}

func BenchmarkBase62Encode(b *testing.B) {
	n := int64(1517153236107137025)
	for i := 0; i < b.N; i++ {
		Encode(n)
	}
}

func BenchmarkBase62Decode(b *testing.B) {
	s := "1O4zPvQMvmh"
	for i := 0; i < b.N; i++ {
		Decode(s)
	}
}

func BenchmarkBase62EncodeParallel(b *testing.B) {
	n := int64(1517153236107137025)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Encode(n)
		}
	})
}

func BenchmarkBase62DecodeParallel(b *testing.B) {
	s := "1O4zPvQMvmh"
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Decode(s)
		}
	})
}
