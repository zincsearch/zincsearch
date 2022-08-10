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

package meta

import "sync/atomic"

type IndexStat struct {
	DocNum      uint64 `json:"doc_num"`
	DocTimeMin  int64  `json:"doc_time_min"`
	DocTimeMax  int64  `json:"doc_time_max"`
	StorageSize uint64 `json:"storage_size"`
	WALSize     uint64 `json:"wal_size"`
}

func NewIndexStat() *IndexStat {
	return &IndexStat{}
}

func (t *IndexStat) GetDocNum() uint64 {
	return atomic.LoadUint64(&t.DocNum)
}

func (t *IndexStat) SetDocNum(num uint64) {
	atomic.StoreUint64(&t.DocNum, num)
}

func (t *IndexStat) GetDocTimeMin() int64 {
	return atomic.LoadInt64(&t.DocTimeMin)
}

func (t *IndexStat) SetDocTimeMin(val int64) {
	atomic.StoreInt64(&t.DocTimeMin, val)
}

func (t *IndexStat) GetDocTimeMax() int64 {
	return atomic.LoadInt64(&t.DocTimeMax)
}

func (t *IndexStat) SetDocTimeMax(val int64) {
	atomic.StoreInt64(&t.DocTimeMax, val)
}

func (t *IndexStat) GetStorageSize() uint64 {
	return atomic.LoadUint64(&t.StorageSize)
}

func (t *IndexStat) SetStorageSize(num uint64) {
	atomic.StoreUint64(&t.StorageSize, num)
}

func (t *IndexStat) GetWALSize() uint64 {
	return atomic.LoadUint64(&t.WALSize)
}

func (t *IndexStat) SetWALSize(num uint64) {
	atomic.StoreUint64(&t.WALSize, num)
}

func (t *IndexStat) Copy() *IndexStat {
	return &IndexStat{
		DocNum:      atomic.LoadUint64(&t.DocNum),
		DocTimeMin:  atomic.LoadInt64(&t.DocTimeMin),
		DocTimeMax:  atomic.LoadInt64(&t.DocTimeMax),
		StorageSize: atomic.LoadUint64(&t.StorageSize),
		WALSize:     atomic.LoadUint64(&t.WALSize),
	}
}
