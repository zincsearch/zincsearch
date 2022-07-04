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
	"sort"
	"strconv"

	"github.com/blugelabs/bluge/search"
	"github.com/blugelabs/bluge/search/aggregations"
)

type TermsAggregation struct {
	src     search.FieldSource
	srcType int
	size    int

	aggregations map[string]search.Aggregation

	lessFunc func(a, b *search.Bucket) bool
	desc     bool
	sortFunc func(p sort.Interface)
}

// NewTermsAggregation returns a termsAggregation
// field use to set the field use to terms aggregation
// valueType use to set the value type, can be diy.TextValueSource / diy.TextValuesSource / diy.NumericValueSource / diy.NumericValuesSource
func NewTermsAggregation(field search.FieldSource, valueType int, size int) *TermsAggregation {
	rv := &TermsAggregation{
		src:     field,
		srcType: valueType,
		size:    size,
		desc:    true,
		lessFunc: func(a, b *search.Bucket) bool {
			return a.Aggregations()["count"].(search.MetricCalculator).Value() < b.Aggregations()["count"].(search.MetricCalculator).Value()
		},
		aggregations: make(map[string]search.Aggregation),
		sortFunc:     sort.Sort,
	}
	rv.aggregations["count"] = aggregations.CountMatches()
	return rv
}

func (t *TermsAggregation) Fields() []string {
	rv := t.src.Fields()
	for _, agg := range t.aggregations {
		rv = append(rv, agg.Fields()...)
	}
	return rv
}

func (t *TermsAggregation) Calculator() search.Calculator {
	return &TermsCalculator{
		src:          t.src,
		srcType:      t.srcType,
		size:         t.size,
		aggregations: t.aggregations,
		desc:         t.desc,
		lessFunc:     t.lessFunc,
		sortFunc:     t.sortFunc,
		bucketsMap:   make(map[string]*search.Bucket),
	}
}

func (t *TermsAggregation) AddAggregation(name string, aggregation search.Aggregation) {
	t.aggregations[name] = aggregation
}

type TermsCalculator struct {
	src     interface{}
	srcType int
	size    int

	aggregations map[string]search.Aggregation

	bucketsList []*search.Bucket
	bucketsMap  map[string]*search.Bucket
	total       int
	other       int

	desc     bool
	lessFunc func(a, b *search.Bucket) bool
	sortFunc func(p sort.Interface)
}

func (a *TermsCalculator) Consume(d *search.DocumentMatch) {
	switch a.srcType {
	case TextValueSource:
		a.consumeTextValueSource(d)
	case TextValuesSource:
		a.consumeTextValuesSource(d)
	case NumericValueSource:
		a.consumeNumericValueSource(d)
	case NumericValuesSource:
		a.consumeNumericValuesSource(d)
	case BooleanValueSource:
		a.consumeBooleanValueSource(d)
	case BooleanValuesSource:
		a.consumeBooleanValuesSource(d)
	default:
		// not supoort
	}
}

func (a *TermsCalculator) consumeTextValueSource(d *search.DocumentMatch) {
	a.total++
	src := a.src.(search.TextValueSource)
	term := src.Value(d)
	termStr := string(term)
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

func (a *TermsCalculator) consumeTextValuesSource(d *search.DocumentMatch) {
	a.total++
	src := a.src.(search.TextValuesSource)
	for _, term := range src.Values(d) {
		termStr := string(term)
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

func (a *TermsCalculator) consumeNumericValueSource(d *search.DocumentMatch) {
	a.total++
	src := a.src.(search.NumericValueSource)
	term := src.Number(d)
	termStr := strconv.FormatFloat(term, 'f', -1, 64)
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

func (a *TermsCalculator) consumeNumericValuesSource(d *search.DocumentMatch) {
	a.total++
	src := a.src.(search.NumericValuesSource)
	for _, term := range src.Numbers(d) {
		termStr := strconv.FormatFloat(term, 'f', -1, 64)
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

func (a *TermsCalculator) consumeBooleanValueSource(d *search.DocumentMatch) {
	a.total++
	src := a.src.(search.NumericValueSource)
	term := src.Number(d)
	termStr := "false"
	if term != 0 {
		termStr = "true"
	}

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

func (a *TermsCalculator) consumeBooleanValuesSource(d *search.DocumentMatch) {
	a.total++
	src := a.src.(search.NumericValuesSource)
	for _, term := range src.Numbers(d) {
		termStr := "false"
		if term != 0 {
			termStr = "true"
		}

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

func (a *TermsCalculator) Merge(other search.Calculator) {
	if other, ok := other.(*TermsCalculator); ok {
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

func (a *TermsCalculator) Finish() {
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

func (a *TermsCalculator) Buckets() []*search.Bucket {
	return a.bucketsList
}

func (a *TermsCalculator) Other() int {
	return a.other
}

func (a *TermsCalculator) Len() int {
	return len(a.bucketsList)
}

func (a *TermsCalculator) Less(i, j int) bool {
	return a.lessFunc(a.bucketsList[i], a.bucketsList[j])
}

func (a *TermsCalculator) Swap(i, j int) {
	a.bucketsList[i], a.bucketsList[j] = a.bucketsList[j], a.bucketsList[i]
}
