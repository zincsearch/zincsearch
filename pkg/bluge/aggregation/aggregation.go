package aggregation

import "github.com/blugelabs/bluge/search"

const (
	TextValueSource = iota
	TextValuesSource
	NumericValueSource
	NumericValuesSource
)

type SearchAggregation interface {
	AddAggregation(name string, aggregation search.Aggregation)
}
