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
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/zincsearch/zincsearch/pkg/metadata"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

var ZINC_INDEX_ALIAS_LIST AliasList

type AliasList struct {
	lock    sync.RWMutex
	Aliases map[string][]string
}

func NewAliasList() *AliasList {
	return &AliasList{Aliases: map[string][]string{}}
}

func (al *AliasList) AddIndexesToAlias(alias string, indexes []string) error {
	al.lock.Lock()
	al.Aliases[alias] = append(al.Aliases[alias], indexes...)

	err := metadata.Alias.Set(al.Aliases)
	if err != nil {
		log.Err(err).Msg("failed to save alias in metadata after add operation")
		al.lock.Unlock()
		return err
	}

	al.lock.Unlock()
	return nil
}

func (al *AliasList) RemoveIndexesFromAlias(alias string, removeIndexes []string) error {
	al.lock.Lock()

	indexes, ok := al.Aliases[alias]
	if !ok {
		al.lock.Unlock()
		return nil
	}

	removeIndexesMap := make(map[string]bool)
	for _, index := range removeIndexes {
		removeIndexesMap[index] = true
	}

	lastIndex := len(indexes)
	for i := 0; i < lastIndex; i++ {
		if _, ok := removeIndexesMap[indexes[i]]; ok {
			indexes[lastIndex-1], indexes[i] = indexes[i], indexes[lastIndex-1]
			i--
			lastIndex--
		}
	}

	al.Aliases[alias] = indexes[:lastIndex]

	err := metadata.Alias.Set(al.Aliases)
	if err != nil {
		log.Err(err).Msg("failed to save alias in metadata after remove operation")
		al.lock.Unlock()
		return err
	}

	al.lock.Unlock()
	return nil
}

func (al *AliasList) GetIndexesForAlias(aliasName string) ([]string, bool) {
	al.lock.RLock()
	idx, ok := al.Aliases[aliasName]
	if !ok {
		al.lock.RUnlock()
		return nil, false
	}

	v := make([]string, len(idx))
	copy(v, idx)

	al.lock.RUnlock()
	return v, ok
}

func (al *AliasList) GetAliasesForIndex(indexName string) []string {
	al.lock.RLock()
	var aliases []string
	for alias, indexes := range al.Aliases {
		if zutils.SliceExists(indexes, indexName) {
			aliases = append(aliases, alias)
		}
	}

	al.lock.RUnlock()
	return aliases
}

type M map[string]interface{}

// GetAliasMap returns an ES compatible map of indexes to their aliases
// In the form:
//
//	{"gitea_issues":{"aliases":{}},"gitea_codes.v1":{"aliases":{"gitea_codes":{}}}}
func (al *AliasList) GetAliasMap(targetIndexes, targetAliases []string) M {
	al.lock.RLock()
	top := M{}

outerLoop:
	for alias, indexes := range al.Aliases {
		if len(targetAliases) > 0 && !zutils.SliceExists(targetAliases, alias) { // check if this is one of the aliased we're looking for
			continue outerLoop
		}

	innerLoop:
		for _, index := range indexes {
			if len(targetIndexes) > 0 && !zutils.SliceExists(targetIndexes, index) { // check if this is one of the indexes we're looking for
				continue innerLoop
			}

			indexMap, _ := top[index].(M)
			if indexMap == nil {
				indexMap = M{}
				top[index] = indexMap
			}

			aliases, _ := indexMap["aliases"].(M)
			if aliases == nil {
				aliases = M{}
				indexMap["aliases"] = aliases
			}

			aliases[alias] = struct{}{}
		}
	}

	al.lock.RUnlock()
	return top
}
