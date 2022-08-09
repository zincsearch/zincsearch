package core

import (
	"sync"

	"github.com/zinclabs/zinc/pkg/metadata"
	"github.com/zinclabs/zinc/pkg/zutils"
)

var ZINC_INDEX_ALIAS_LIST AliasList

type AliasList struct {
	// Name    string
	// Indices []string
	lock    sync.RWMutex
	Aliases map[string][]string
}

func NewAliasList() *AliasList {
	return &AliasList{Aliases: map[string][]string{}}
}

func (al *AliasList) AddIndexesToAlias(alias string, indexes []string) error {
	al.lock.Lock()
	if al.Aliases == nil {
		al.Aliases = map[string][]string{}
	}

	al.Aliases[alias] = append(al.Aliases[alias], indexes...)

	err := metadata.Alias.Set(al.Aliases)
	if err != nil {
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

outer:
	for _, removeIndex := range removeIndexes {
		for i, s := range indexes {
			if s == removeIndex {
				indexes = append(indexes[:i], indexes[i+1:]...)
				continue outer
			}
		}
	}

	al.Aliases[alias] = indexes

	err := metadata.Alias.Set(al.Aliases)
	if err != nil {
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
