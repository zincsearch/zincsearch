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

type Index struct {
	Meta     *IndexMeta     `json:"-"`
	Stats    *IndexStat     `json:"-"`
	Shards   *IndexShards   `json:"-"`
	Settings *IndexSettings `json:"-"`
	Mappings *Mappings      `json:"-"`
}

type IndexSimple struct {
	Name        string                 `json:"name"`
	StorageType string                 `json:"storage_type"`
	ShardNum    int64                  `json:"shard_num"`
	Settings    *IndexSettings         `json:"settings,omitempty"`
	Mappings    map[string]interface{} `json:"mappings,omitempty"`
}

func NewIndex(name, storageType, version string) *Index {
	return &Index{
		Meta: &IndexMeta{
			Name:        name,
			StorageType: storageType,
			Version:     version,
		},
		Stats:    NewIndexStat(),
		Shards:   NewIndexShards(),
		Settings: NewIndexSettings(),
		Mappings: NewMappings(),
	}
}

func (t *Index) GetName() string {
	return t.Meta.GetName()
}

func (t *Index) GetStorageType() string {
	return t.Meta.GetStorageType()
}

func (t *Index) GetVersion() string {
	return t.Meta.GetVersion()
}

func (t *Index) GetShardNum() int64 {
	return t.Shards.GetShardNum()
}
