package index

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/zutils"
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
// @Tags    Index
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} meta.HTTPResponseError
// @Router /es/_aliases [post]
func AddOrRemoveESAlias(c *gin.Context) {
	var alias Alias
	err := zutils.GinBindJSON(c, &alias)
	if err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	addMap := map[string][]string{}
	removeMap := map[string][]string{}

	indexList := core.ZINC_INDEX_LIST.List()

outerLoop:
	for _, action := range alias.Actions {
		if action.Add != nil {
			if action.Add.Index != "" {
				matchAndAddToMap(indexList, action.Add.Index, addMap, action.Add)
				continue outerLoop
			}

			// index is empty, try the indices field
			for _, indexName := range action.Add.Indices {
				matchAndAddToMap(indexList, indexName, addMap, action.Add)
			}

			continue outerLoop // this was an add action, don't bother checking action.Remove
		}

		if action.Remove != nil {
			if action.Remove.Index != "" {
				matchAndAddToMap(indexList, action.Remove.Index, removeMap, action.Remove)
				continue outerLoop
			}

			// index is empty, try the indices field
			for _, indexName := range action.Remove.Indices {
				matchAndAddToMap(indexList, indexName, removeMap, action.Remove)
			}
		}
	}

	for alias, indexes := range addMap {
		core.ZINC_INDEX_ALIAS_LIST.AddIndexesToAlias(alias, indexes)
	}

	for alias, indexes := range removeMap {
		core.ZINC_INDEX_ALIAS_LIST.RemoveIndexesFromAlias(alias, indexes)
	}

	c.JSON(http.StatusOK, gin.H{"acknowledged": true})
}

type M map[string]interface{}

// @Id GetESAliases
// @Summary Get index alias for compatible ES
// @Tags    Index
// @Produce json
// @Param   index path  string  false  "Index"
// @Param   alias path  string  false  "Alias"
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

	c.JSON(http.StatusOK, m)
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
	for _, index := range indexList {
		n := index.GetName()
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
