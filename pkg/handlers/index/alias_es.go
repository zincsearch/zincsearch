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
	Index string `json:"index"`
	Alias string `json:"alias"`
}

func AddOrRemoveESAlias(c *gin.Context) {
	var alias Alias
	err := zutils.GinBindJSON(c, &alias)
	if err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	addMap := map[string][]string{}
	removeMap := map[string][]string{}

	var indexList []*core.Index

	target := c.Param("target")
	if target != "" {
		index, ok := core.ZINC_INDEX_LIST.Get(target)
		if !ok {
			c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: "index not found"})
			return
		}
		indexList = []*core.Index{index}
	} else {
		indexList = core.ZINC_INDEX_LIST.List()
	}

	c.Params[0].Value = ""
	for _, action := range alias.Actions {
		if action.Add != nil {
			for _, index := range indexList {
				indexName := index.GetName()
				if indexNameMatches(action.Add.Index, indexName) {
					addMap[indexName] = append(addMap[indexName], action.Add.Alias) // append the alias to add list for this index
				}
			}

			continue
		}

		if action.Remove != nil {
			for _, index := range indexList {
				indexName := index.GetName()
				if indexNameMatches(action.Add.Index, indexName) {
					removeMap[indexName] = append(addMap[indexName], action.Add.Alias) // append the alias to remove list for this index
				}
			}
		}
	}

	var aliases []string
	for _, index := range indexList {
		aliases = addMap[index.GetName()]
		if aliases != nil {
			index.AddAliases(aliases)
		}

		aliases = removeMap[index.GetName()]
		if aliases != nil {
			index.RemoveAliases(aliases)
		}
	}
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
