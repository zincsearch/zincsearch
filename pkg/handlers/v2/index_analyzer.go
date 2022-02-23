package v2

import (
	"net/http"

	"github.com/blugelabs/bluge/analysis"
	"github.com/gin-gonic/gin"

	"github.com/prabhatsharma/zinc/pkg/core"
	"github.com/prabhatsharma/zinc/pkg/uquery/v2/analyzer"
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
	index, exists := core.GetIndex(indexName)
	if exists {
		ana, err = analyzer.Query(index.CachedAnalysis, query.Analyzer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "analyzer " + query.Analyzer + " does not exists"})
			return
		}
	} else {
		ana, err = analyzer.Query(nil, query.Analyzer)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "analyzer " + query.Analyzer + " does not exists"})
			return
		}
	}

	charFilters := make([]analysis.CharFilter, 0)
	if query.CharFilter != nil {
		switch v := query.CharFilter.(type) {
		case string:
			filter, err := analyzer.RequestCharFilterSingle(v, nil)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			charFilters = append(charFilters, filter)
		case []interface{}:
			filters, err := analyzer.RequestCharFilterSlice(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			charFilters = append(charFilters, filters...)
		case map[string]interface{}:
			filters, err := analyzer.RequestCharFilter(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			for _, filter := range filters {
				charFilters = append(charFilters, filter)
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
			filter, err := analyzer.RequestTokenFilterSingle(v, nil)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			tokenFilters = append(tokenFilters, filter)
		case []interface{}:
			filters, err := analyzer.RequestTokenFilterSlice(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			tokenFilters = append(tokenFilters, filters...)
		case map[string]interface{}:
			filters, err := analyzer.RequestTokenFilter(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			for _, filter := range filters {
				tokenFilters = append(tokenFilters, filter)
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
			zer, err := analyzer.RequestTokenizerSingle(v, nil)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			tokenizers = append(tokenizers, zer)
		case []interface{}:
			zers, err := analyzer.RequestTokenizerSlice(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			tokenizers = append(tokenizers, zers...)
		case map[string]interface{}:
			zers, err := analyzer.RequestTokenizer(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			for _, zer := range zers {
				tokenizers = append(tokenizers, zer)
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
	Text        string      `json:"text"`
	Tokenizer   interface{} `json:"tokenizer"`
	CharFilter  interface{} `json:"char_filter"`
	TokenFilter interface{} `json:"token_filter"`
}
