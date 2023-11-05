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

package aggregation

import (
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/blugelabs/bluge/search"
	"github.com/blugelabs/bluge/search/aggregations"
)

type DateHistogramAggregation struct {
	src              search.FieldSource
	size             int
	calendarInterval string
	fixedInterval    int64 // unit: time.Nanosecond
	minDocCount      int
	format           string
	timeZone         *time.Location

	extendedBounds *HistogramBound
	hardBounds     *HistogramBound

	aggregations map[string]search.Aggregation

	lessFunc func(a, b *search.Bucket) bool
	desc     bool
	sortFunc func(p sort.Interface)
}

// NewDateHistogramAggregation returns a termsAggregation
// field use to set the field use to terms aggregation
func NewDateHistogramAggregation(
	field search.FieldSource,
	calendarInterval string,
	fixedInterval int64,
	format string,
	timeZone *time.Location,
	extendedBounds,
	hardBounds *HistogramBound,
	minDocCount,
	size int,
) *DateHistogramAggregation {
	rv := &DateHistogramAggregation{
		src:              field,
		size:             size,
		calendarInterval: calendarInterval,
		fixedInterval:    fixedInterval,
		minDocCount:      minDocCount,
		format:           format,
		timeZone:         timeZone,
		extendedBounds:   extendedBounds,
		hardBounds:       hardBounds,
		desc:             false,
		lessFunc: func(a, b *search.Bucket) bool {
			return a.Name() < b.Name()
		},
		aggregations: make(map[string]search.Aggregation),
		sortFunc:     sort.Sort,
	}
	rv.aggregations["count"] = aggregations.CountMatches()
	return rv
}

func (t *DateHistogramAggregation) Fields() []string {
	rv := t.src.Fields()
	for _, agg := range t.aggregations {
		rv = append(rv, agg.Fields()...)
	}
	return rv
}

func (t *DateHistogramAggregation) Calculator() search.Calculator {
	return &DateHistogramCalculator{
		src:              t.src,
		size:             t.size,
		calendarInterval: t.calendarInterval,
		fixedInterval:    t.fixedInterval,
		minDocCount:      t.minDocCount,
		format:           t.format,
		timeZone:         t.timeZone,
		minValue:         math.MaxInt64,
		maxValue:         math.MinInt64,
		extendedBounds:   t.extendedBounds,
		hardBounds:       t.hardBounds,
		aggregations:     t.aggregations,
		desc:             t.desc,
		lessFunc:         t.lessFunc,
		sortFunc:         t.sortFunc,
		bucketsMap:       make(map[string]*search.Bucket),
	}
}

func (t *DateHistogramAggregation) AddAggregation(name string, aggregation search.Aggregation) {
	t.aggregations[name] = aggregation
}

type DateHistogramCalculator struct {
	src              interface{}
	size             int
	calendarInterval string
	fixedInterval    int64
	minDocCount      int
	format           string
	timeZone         *time.Location

	minValue       int64
	maxValue       int64
	extendedBounds *HistogramBound
	hardBounds     *HistogramBound

	aggregations map[string]search.Aggregation

	bucketsList []*search.Bucket
	bucketsMap  map[string]*search.Bucket
	total       int
	other       int

	desc     bool
	lessFunc func(a, b *search.Bucket) bool
	sortFunc func(p sort.Interface)
}

func (a *DateHistogramCalculator) Consume(d *search.DocumentMatch) {
	a.total++
	src := a.src.(search.DateValuesSource)
	for _, term := range src.Dates(d) {
		if term.UnixNano() < a.minValue {
			a.minValue = term.UnixNano()
		}
		if term.UnixNano() > a.maxValue {
			a.maxValue = term.UnixNano()
		}
		termStr := a.bucketKey(term.UnixNano())
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

func (a *DateHistogramCalculator) Merge(other search.Calculator) {
	if other, ok := other.(*DateHistogramCalculator); ok {
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

func (a *DateHistogramCalculator) Finish() {
	// re calculate min max
	if a.extendedBounds != nil {
		min := int64(a.extendedBounds.Min * 1e6)
		max := int64(a.extendedBounds.Max * 1e6)
		if a.minValue > min {
			a.minValue = min
		}
		if a.maxValue < max {
			a.maxValue = max
		}
	}
	if a.hardBounds != nil {
		a.minValue = int64(a.hardBounds.Min * 1e6)
		a.maxValue = int64(a.hardBounds.Max * 1e6)
	}
	// Replenish bucket
	if a.minDocCount == 0 {
		if a.calendarInterval != "" {
			for value := a.minValue; value < a.maxValue; {
				termStr := a.bucketKey(value)
				if _, ok := a.bucketsMap[termStr]; !ok {
					a.bucketsList = append(a.bucketsList, search.NewBucket(termStr, a.aggregations))
				}
				t := time.Unix(0, value).In(a.timeZone)
				switch a.calendarInterval {
				case "week", "1w":
					t = time.Date(t.Year(), t.Month(), t.Day()+7, 0, 0, 0, 0, t.Location())
				case "month", "1M":
					t = time.Date(t.Year(), t.Month()+1, t.Day(), 0, 0, 0, 0, t.Location())
				case "quarter", "1q":
					t = time.Date(t.Year(), t.Month()+3, t.Day(), 0, 0, 0, 0, t.Location())
				case "year", "1y":
					t = time.Date(t.Year()+1, t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
				default:
					// noop
				}
				value = t.UnixNano()
			}
		} else {
			for value := a.minValue; value < a.maxValue; value += a.fixedInterval {
				termStr := a.bucketKey(value)
				if _, ok := a.bucketsMap[termStr]; !ok {
					a.bucketsList = append(a.bucketsList, search.NewBucket(termStr, a.aggregations))
				}
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

func (a *DateHistogramCalculator) Buckets() []*search.Bucket {
	return a.bucketsList
}

func (a *DateHistogramCalculator) Other() int {
	return a.other
}

func (a *DateHistogramCalculator) Len() int {
	return len(a.bucketsList)
}

func (a *DateHistogramCalculator) Less(i, j int) bool {
	return a.lessFunc(a.bucketsList[i], a.bucketsList[j])
}

func (a *DateHistogramCalculator) Swap(i, j int) {
	a.bucketsList[i], a.bucketsList[j] = a.bucketsList[j], a.bucketsList[i]
}

func (a *DateHistogramCalculator) bucketKey(value int64) string {
	var nsec int64
	if a.calendarInterval != "" {
		t := time.Unix(0, value).In(a.timeZone)
		switch a.calendarInterval {
		case "week", "1w":
			t = time.Date(t.Year(), t.Month(), t.Day()-int(t.Weekday()), 0, 0, 0, 0, t.Location())
		case "month", "1M":
			t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		case "quarter", "1q":
			switch t.Month() {
			case 1, 2, 3:
				t = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
			case 4, 5, 6:
				t = time.Date(t.Year(), 4, 1, 0, 0, 0, 0, t.Location())
			case 7, 8, 9:
				t = time.Date(t.Year(), 7, 1, 0, 0, 0, 0, t.Location())
			case 10, 11, 12:
				t = time.Date(t.Year(), 10, 1, 0, 0, 0, 0, t.Location())
			}
		case "year", "1y":
			t = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
		default:
			// noop
		}
		nsec = t.UnixNano()
	} else {
		nsec = (value / a.fixedInterval) * a.fixedInterval
	}

	if a.format == "epoch_millis" {
		return strconv.FormatInt(time.Unix(0, nsec).In(a.timeZone).UnixMilli(), 10)
	}

	return time.Unix(0, nsec).In(a.timeZone).Format(a.format)
}
