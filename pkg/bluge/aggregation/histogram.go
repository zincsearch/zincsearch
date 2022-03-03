package aggregation

import (
	"math"
	"sort"
	"strconv"

	"github.com/blugelabs/bluge/search"
	"github.com/blugelabs/bluge/search/aggregations"
)

type HistogramAggregation struct {
	src         search.FieldSource
	size        int
	interval    float64
	offset      float64
	minDocCount int

	aggregations map[string]search.Aggregation

	lessFunc func(a, b *search.Bucket) bool
	desc     bool
	sortFunc func(p sort.Interface)
}

// NewHistogramAggregation returns a termsAggregation
// field use to set the field use to terms aggregation
func NewHistogramAggregation(field search.FieldSource, interval, offset float64, minDocCount, size int) *HistogramAggregation {
	rv := &HistogramAggregation{
		src:         field,
		size:        size,
		interval:    interval,
		offset:      offset,
		minDocCount: minDocCount,
		desc:        false,
		lessFunc: func(a, b *search.Bucket) bool {
			return a.Name() < b.Name()
		},
		aggregations: make(map[string]search.Aggregation),
		sortFunc:     sort.Sort,
	}
	rv.aggregations["count"] = aggregations.CountMatches()
	return rv
}

func (t *HistogramAggregation) Fields() []string {
	rv := t.src.Fields()
	for _, agg := range t.aggregations {
		rv = append(rv, agg.Fields()...)
	}
	return rv
}

func (t *HistogramAggregation) Calculator() search.Calculator {
	return &HistogramCalculator{
		src:          t.src,
		size:         t.size,
		interval:     t.interval,
		offset:       t.offset,
		minDocCount:  t.minDocCount,
		minValue:     math.MaxFloat64,
		maxValue:     math.SmallestNonzeroFloat64,
		aggregations: t.aggregations,
		desc:         t.desc,
		lessFunc:     t.lessFunc,
		sortFunc:     t.sortFunc,
		bucketsMap:   make(map[string]*search.Bucket),
	}
}

func (t *HistogramAggregation) AddAggregation(name string, aggregation search.Aggregation) {
	t.aggregations[name] = aggregation
}

type HistogramCalculator struct {
	src         interface{}
	size        int
	interval    float64
	offset      float64
	minDocCount int

	minValue float64
	maxValue float64

	aggregations map[string]search.Aggregation

	bucketsList []*search.Bucket
	bucketsMap  map[string]*search.Bucket
	total       int
	other       int

	desc     bool
	lessFunc func(a, b *search.Bucket) bool
	sortFunc func(p sort.Interface)
}

func (a *HistogramCalculator) Consume(d *search.DocumentMatch) {
	a.total++
	src := a.src.(search.NumericValuesSource)
	for _, term := range src.Numbers(d) {
		if term < a.minValue {
			a.minValue = term
		}
		if term > a.maxValue {
			a.maxValue = term
		}
		termStr := a.bucketKey(term)
		bucket, ok := a.bucketsMap[termStr]
		if ok {
			bucket.Consume(d)
		} else {
			newBucket := search.NewBucket(termStr, a.aggregations)
			newBucket.Consume(d)
			a.bucketsMap[termStr] = newBucket
			a.bucketsList = append(a.bucketsList, newBucket)
		}
	}
}

func (a *HistogramCalculator) Merge(other search.Calculator) {
	if other, ok := other.(*HistogramCalculator); ok {
		// first sum to the totals and others
		a.total += other.total
		// now, walk all of the other buckets
		// if we have a local match, merge otherwise append
		for i := range other.bucketsList {
			var foundLocal bool
			for j := range a.bucketsList {
				if other.bucketsList[i].Name() == a.bucketsList[j].Name() {
					a.bucketsList[j].Merge(other.bucketsList[i])
					foundLocal = true
				}
			}
			if !foundLocal {
				a.bucketsList = append(a.bucketsList, other.bucketsList[i])
			}
		}
		// now re-invoke finish, this should trim to correct size again
		// and recalculate other
		a.Finish()
	}
}

func (a *HistogramCalculator) Finish() {
	// check bucket
	if a.minDocCount == 0 {
		for value := a.minValue; value < a.maxValue; value += a.interval {
			termStr := a.bucketKey(value)
			if _, ok := a.bucketsMap[termStr]; !ok {
				a.bucketsList = append(a.bucketsList, search.NewBucket(termStr, a.aggregations))
			}
		}
	} else {
		for i := 0; i < len(a.bucketsList); i++ {
			if a.bucketsList[i].Count() >= uint64(a.minDocCount) {
				continue
			}
			if i == 0 {
				a.bucketsList = a.bucketsList[1:]
			} else {
				a.bucketsList = append(a.bucketsList[:i], a.bucketsList[i+1:]...)
			}
			i--
		}
	}

	// sort the buckets
	if a.desc {
		a.sortFunc(sort.Reverse(a))
	} else {
		a.sortFunc(a)
	}

	trimTopN := a.size
	if trimTopN > len(a.bucketsList) {
		trimTopN = len(a.bucketsList)
	}
	a.bucketsList = a.bucketsList[:trimTopN]

	var notOther int
	for _, bucket := range a.bucketsList {
		notOther += int(bucket.Aggregations()["count"].(search.MetricCalculator).Value())
	}
	a.other = a.total - notOther
}

func (a *HistogramCalculator) Buckets() []*search.Bucket {
	return a.bucketsList
}

func (a *HistogramCalculator) Other() int {
	return a.other
}

func (a *HistogramCalculator) Len() int {
	return len(a.bucketsList)
}

func (a *HistogramCalculator) Less(i, j int) bool {
	return a.lessFunc(a.bucketsList[i], a.bucketsList[j])
}

func (a *HistogramCalculator) Swap(i, j int) {
	a.bucketsList[i], a.bucketsList[j] = a.bucketsList[j], a.bucketsList[i]
}

func (a *HistogramCalculator) bucketKey(value float64) string {
	f := math.Floor((value-a.offset)/a.interval)*a.interval + a.offset
	return strconv.FormatFloat(f, 'f', -1, 64)
}
