package v2

import (
	"net/http"

	"github.com/blugelabs/bluge/analysis"
	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/core"
	zincanalysis "github.com/prabhatsharma/zinc/pkg/uquery/v2/analysis"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func Analyze(c *gin.Context) {
	var query analyzeRequest
	if err := c.BindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var err error
	var ana *analysis.Analyzer
	indexName := c.Param("target")
	if indexName != "" {
		index, exists := core.GetIndex(indexName)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "index " + indexName + " does not exists"})
			return
		}
		if query.Filed != "" && query.Analyzer == "" {
			if index.CachedMappings != nil && index.CachedMappings.Properties != nil {
				if prop, ok := index.CachedMappings.Properties[query.Filed]; ok {
					query.Analyzer = prop.Analyzer
				}
			}
		}
		ana, err = zincanalysis.Query(index.CachedAnalysis, query.Analyzer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "analyzer " + query.Analyzer + " does not exists"})
			return
		}
	} else {
		// none index specified
		ana, err = zincanalysis.Query(nil, query.Analyzer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "analyzer " + query.Analyzer + " does not exists"})
			return
		}
	}

	charFilters := make([]analysis.CharFilter, 0)
	if query.CharFilter != nil {
		switch v := query.CharFilter.(type) {
		case string:
			filter, err := zincanalysis.RequestCharFilterSingle(v, nil)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			charFilters = append(charFilters, filter)
		case []interface{}:
			filters, err := zincanalysis.RequestCharFilterSlice(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			charFilters = append(charFilters, filters...)
		case map[string]interface{}:
			typ, err := zutils.GetStringFromMap(v, "type")
			if typ != "" && err == nil {
				filter, err := zincanalysis.RequestCharFilterSingle(typ, v)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				charFilters = append(charFilters, filter)
			} else {
				filters, err := zincanalysis.RequestCharFilter(v)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				for _, filter := range filters {
					charFilters = append(charFilters, filter)
				}
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "char_filter unsuported type"})
			return
		}
	}

	tokenFilters := make([]analysis.TokenFilter, 0)
	if query.TokenFilter != nil {
		switch v := query.TokenFilter.(type) {
		case string:
			filter, err := zincanalysis.RequestTokenFilterSingle(v, nil)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			tokenFilters = append(tokenFilters, filter)
		case []interface{}:
			filters, err := zincanalysis.RequestTokenFilterSlice(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			tokenFilters = append(tokenFilters, filters...)
		case map[string]interface{}:
			typ, err := zutils.GetStringFromMap(v, "type")
			if typ != "" && err == nil {
				filter, err := zincanalysis.RequestTokenFilterSingle(typ, v)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				tokenFilters = append(tokenFilters, filter)
			} else {
				filters, err := zincanalysis.RequestTokenFilter(v)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				for _, filter := range filters {
					tokenFilters = append(tokenFilters, filter)
				}
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "token_filter unsuported type"})
			return
		}
	}

	tokenizers := make([]analysis.Tokenizer, 0)
	if query.Tokenizer != nil {
		switch v := query.Tokenizer.(type) {
		case string:
			zer, err := zincanalysis.RequestTokenizerSingle(v, nil)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			tokenizers = append(tokenizers, zer)
		case []interface{}:
			zers, err := zincanalysis.RequestTokenizerSlice(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			tokenizers = append(tokenizers, zers...)
		case map[string]interface{}:
			typ, err := zutils.GetStringFromMap(v, "type")
			if typ != "" && err == nil {
				zer, err := zincanalysis.RequestTokenizerSingle(typ, v)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				tokenizers = append(tokenizers, zer)
			} else {
				zers, err := zincanalysis.RequestTokenizer(v)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				for _, zer := range zers {
					tokenizers = append(tokenizers, zer)
				}
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "tokenizer unsuported type"})
			return
		}
	}

	if len(charFilters) > 0 {
		ana.CharFilters = append(ana.CharFilters, charFilters...)
	}

	if len(tokenFilters) > 0 {
		ana.TokenFilters = append(ana.TokenFilters, tokenFilters...)
	}

	if len(tokenizers) > 0 {
		ana.Tokenizer = tokenizers[0]
	}

	tokens := ana.Analyze([]byte(query.Text))
	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}

type analyzeRequest struct {
	Analyzer    string      `json:"analyzer"`
	Filed       string      `json:"field"`
	Text        string      `json:"text"`
	Tokenizer   interface{} `json:"tokenizer"`
	CharFilter  interface{} `json:"char_filter"`
	TokenFilter interface{} `json:"token_filter"`
}
