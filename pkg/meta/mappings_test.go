package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProperty_DeepClone(t *testing.T) {
	prop := NewProperty("text")
	prop.AddField("sub_field", NewProperty("text"))
	prop.Sortable = true
	prop.Aggregatable = true
	prop.Index = true
	prop.Store = true
	prop.Aggregatable = true
	prop.Highlightable = true
	prop.Analyzer = "some_nice_analyzer"
	prop.SearchAnalyzer = "search_analyzer"
	prop.Format = "2009-11-10T23:00:00Z"

	clone := prop.DeepClone()
	assert.Equal(t, prop, clone)
}
