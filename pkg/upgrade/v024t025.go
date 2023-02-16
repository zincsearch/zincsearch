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

package upgrade

import (
	"os"
	"path"

	"github.com/zinclabs/zincsearch/pkg/config"
	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/pkg/zutils"
)

// UpgradeFromV024T025 upgrades from version <= 0.2.4
// upgrade steps:
// -- mv    index index_old
// -- mkdir index
// -- mv    index_old index/000000
func UpgradeFromV024T025(index *meta.Index) error {
	indexName := index.Name
	rootPath := config.Global.DataPath
	if ok, _ := zutils.IsExist(path.Join(rootPath, indexName)); !ok {
		return nil // if index does not exist, skip
	}
	if ok, _ := zutils.IsExist(path.Join(rootPath, indexName, "000000")); ok {
		return nil // if index already upgraded, skip
	}
	if err := os.Rename(path.Join(rootPath, indexName), path.Join(rootPath, indexName+"_old")); err != nil {
		return err
	}
	if err := os.Mkdir(path.Join(rootPath, indexName), 0755); err != nil {
		return err
	}
	if err := os.Rename(path.Join(rootPath, indexName+"_old"), path.Join(rootPath, indexName, "000000")); err != nil {
		return err
	}
	return nil
}
