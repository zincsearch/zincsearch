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

package index

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zincsearch/pkg/core"
	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/pkg/zutils"
)

type Alias struct {
	Actions []Action `json:"actions"`
}

type Action struct {
	Add    *base `json:"add"`
	Remove *base `json:"remove"`
}

type base struct {
	Index   string   `json:"index"`
	Alias   string   `json:"alias"`
	Indices []string `json:"indices"`
	Aliases []string `json:"aliases"`
}

// @Id AddOrRemoveESAlias
// @Summary Add or remove index alias for compatible ES
// @security BasicAuth
// @Tags    Index
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} meta.HTTPResponseError
// @Router /es/_aliases [post]
func AddOrRemoveESAlias(c *gin.Context) {
	var alias Alias
	err := zutils.GinBindJSON(c, &alias)
	if err != nil {
		zutils.GinRenderJSON(c, http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	addMap := map[string][]string{}
	removeMap := map[string][]string{}

	indexList := core.ZINC_INDEX_LIST.List()

	for _, action := range alias.Actions {
		if action.Add != nil {
			if action.Add.Index != "" {
				matchAndAddToMap(indexList, action.Add.Index, addMap, action.Add)
				continue
			}

			// index is empty, try the indices field
			for _, indexName := range action.Add.Indices {
				matchAndAddToMap(indexList, indexName, addMap, action.Add)
			}

			continue // this was an add action, don't bother checking action.Remove
		}

		if action.Remove != nil {
			if action.Remove.Index != "" {
				matchAndAddToMap(indexList, action.Remove.Index, removeMap, action.Remove)
				continue
			}

			// index is empty, try the indices field
			for _, indexName := range action.Remove.Indices {
				matchAndAddToMap(indexList, indexName, removeMap, action.Remove)
			}
		}
	}

	for alias, indexes := range addMap {
		_ = core.ZINC_INDEX_ALIAS_LIST.AddIndexesToAlias(alias, indexes)
	}

	for alias, indexes := range removeMap {
		_ = core.ZINC_INDEX_ALIAS_LIST.RemoveIndexesFromAlias(alias, indexes)
	}

	zutils.GinRenderJSON(c, http.StatusOK, gin.H{"acknowledged": true})
}

// @Id GetESAliases
// @Summary Get index alias for compatible ES
// @security BasicAuth
// @Tags    Index
// @Produce json
// @Param   target path  string  false  "Target Index"
// @Param   target_alias path  string  false  "Target Alias"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} meta.HTTPResponseError
// @Router /es/{target}/_alias/{target_alias} [get]
func GetESAliases(c *gin.Context) {
	targetIndex := c.Param("target")

	var targetIndexes []string
	if targetIndex != "" {
		targetIndexes = strings.Split(targetIndex, ",")
	}

	targetAlias := c.Param("target_alias")

	var targetAliases []string
	if targetAlias != "" {
		targetAliases = strings.Split(targetAlias, ",")
	}

	m := core.ZINC_INDEX_ALIAS_LIST.GetAliasMap(targetIndexes, targetAliases)

	zutils.GinRenderJSON(c, http.StatusOK, m)
}

func indexNameMatches(name, indexName string) bool {
	if name == indexName {
		return true
	}

	if strings.Contains(name, "*") {
		p, err := getRegex(name)
		if err != nil {
			log.Err(err).Msg("failed to compile regex")
			return false
		}

		return p.MatchString(indexName)
	}

	return false
}

func getRegex(s string) (*regexp.Regexp, error) {
	parts := strings.Split(s, "*")
	pattern := ""
	for i, part := range parts {
		pattern += part
		if i < len(parts)-1 {
			pattern += "[a-zA-Z0-9_.-]+"
		}
	}

	p, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func matchAndAddToMap(indexList []*core.Index, indexName string, m map[string][]string, b *base) {
	var n string // reuse same string variable

	if !strings.Contains(indexName, "*") {
		x, ok := core.ZINC_INDEX_LIST.Get(indexName)
		if !ok {
			return
		}

		n = x.GetName()

		if b.Alias != "" { // alias takes precedence over aliases
			m[b.Alias] = append(m[b.Alias], n)
		} else {
			for _, a := range b.Aliases {
				m[a] = append(m[a], n)
			}
		}
		return
	}

	// indexName contains a wildcard(*) r, range over the entire indexlist looking for matches
	for _, index := range indexList {
		n = index.GetName()
		if indexNameMatches(indexName, n) {
			if b.Alias != "" { // alias takes precedence over aliases
				m[b.Alias] = append(m[b.Alias], n)
			} else {
				for _, a := range b.Aliases {
					m[a] = append(m[a], n)
				}
			}
		}
	}
}
