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

type IndexMeta struct {
	Name        string `json:"name"`
	StorageType string `json:"storage_type"`
	Version     string `json:"version"`
}

func NewIndexMeta() *IndexMeta {
	return &IndexMeta{}
}

func (t *IndexMeta) GetName() string {
	return t.Name
}

func (t *IndexMeta) GetStorageType() string {
	return t.StorageType
}

func (t *IndexMeta) GetVersion() string {
	return t.Version
}

func (t *IndexMeta) SetVersion(ver string) {
	t.Version = ver
}

func (t *IndexMeta) Copy() *IndexMeta {
	return &IndexMeta{
		Name:        t.Name,
		StorageType: t.StorageType,
		Version:     t.Version,
	}
}
