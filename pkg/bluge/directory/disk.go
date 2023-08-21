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

package directory

import (
	"github.com/zincsearch/zincsearch/pkg/config"
	"github.com/zincsearch/zincsearch/pkg/objstore"
	"path"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/index"
)

// GetDiskConfig returns a bluge config that will store index data in local disk
// rootPath: the root path of data
// indexName: the name of the index to use.
func GetDiskConfig(rootPath string, indexName string, timeRange ...int64) bluge.Config {
	cfg := index.DefaultConfig(path.Join(rootPath, indexName))
	cfg = cfg.WithPersisterNapTimeMSec(50)

	if config.Global.StorageType == config.S3Storage {
		cfg.DirectoryFunc = func() index.Directory {
			dir, err := objstore.New(rootPath, indexName, config.Global.S3ConfigPath)
			if err != nil {
				panic(err)
			}
			return dir
		}
	}
	if len(timeRange) == 2 {
		if timeRange[0] <= timeRange[1] {
			cfg = cfg.WithTimeRange(timeRange[0], timeRange[1])
		}
	}
	return bluge.DefaultConfigWithIndexConfig(cfg)
}
