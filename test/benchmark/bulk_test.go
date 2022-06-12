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

package benchmark

import (
	"os"
	"testing"

	"github.com/zinclabs/zinc/pkg/handlers/document"
)

func BenchmarkBulk(b *testing.B) {
	f, err := os.Open("../../tmp/olympics.ndjson")
	if err != nil {
		b.Error(err)
	}

	target := "olympics"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err = document.BulkWorker(target, f)
		if err != nil {
			b.Error(err)
		}
	}
}
