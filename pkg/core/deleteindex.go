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

package core

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zincsearch/pkg/config"
	"github.com/zinclabs/zincsearch/pkg/metadata"
)

func DeleteIndex(name string) error {
	// 1. Check if index exists
	index, exists := GetIndex(name)
	if !exists {
		return errors.New("index " + name + " does not exists")
	}

	// 2. Close and Delete from cache
	ZINC_INDEX_LIST.Delete(name)

	// 3. Physically delete the index
	dataPath := config.Global.DataPath
	err := os.RemoveAll(dataPath + "/" + index.GetName())
	if err != nil {
		log.Error().Err(err).Msg("failed to delete index")
	}

	// 4. Delete form metadata
	return metadata.Index.Delete(name)
}
