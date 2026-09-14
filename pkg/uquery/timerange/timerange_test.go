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

package timerange

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/pkg/meta"
)

var (
	t1  = time.Date(2022, 3, 4, 5, 6, 7, 0, time.UTC)
	t2  = time.Date(2022, 3, 5, 5, 6, 7, 0, time.UTC)
	ts1 = t1.Format(time.RFC3339)
	ts2 = t2.Format(time.RFC3339)
)

func tsRange(kv map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{"range": map[string]interface{}{meta.TimeFieldName: kv}}
}

func TestRangeQueryTime(t *testing.T) {
	tests := []struct {
		name         string
		query        map[string]interface{}
		wantMin, max int64
	}{
		{"gte/lt rfc3339", map[string]interface{}{"gte": ts1, "lt": ts2}, t1.UnixNano(), t2.UnixNano()},
		{"gt/lte rfc3339", map[string]interface{}{"gt": ts1, "lte": ts2}, t1.UnixNano(), t2.UnixNano()},
		{"only gte", map[string]interface{}{"gte": ts1}, t1.UnixNano(), time.Time{}.UTC().UnixNano()},
		{"only lt", map[string]interface{}{"lt": ts2}, time.Time{}.UTC().UnixNano(), t2.UnixNano()},
		{"epoch_millis float", map[string]interface{}{"format": "epoch_millis", "gte": float64(t1.UnixMilli()), "lte": float64(t2.UnixMilli())},
			t1.UnixNano(), t2.UnixNano()},
		{"epoch_millis string", map[string]interface{}{"format": "epoch_millis", "gt": "1646370367000", "lt": "1646456767000"},
			t1.UnixNano(), t2.UnixNano()},
		{"custom format with time_zone", map[string]interface{}{"format": "2006-01-02 15:04:05", "time_zone": "+08:00", "gte": "2022-03-04 13:06:07"},
			t1.UnixNano(), time.Time{}.UTC().UnixNano()},
		{"boost accepted", map[string]interface{}{"gte": ts1, "boost": 2.0}, t1.UnixNano(), time.Time{}.UTC().UnixNano()},
		{"unknown key", map[string]interface{}{"gte": ts1, "foo": 1}, 0, 0},
		{"bad time_zone", map[string]interface{}{"gte": ts1, "time_zone": "Nowhere/City"}, 0, 0},
		{"bad gt", map[string]interface{}{"gt": "nope"}, 0, 0},
		{"bad gte", map[string]interface{}{"gte": "nope"}, 0, 0},
		{"bad lt", map[string]interface{}{"lt": "nope"}, 0, 0},
		{"bad lte", map[string]interface{}{"lte": "nope"}, 0, 0},
		{"bad epoch gt", map[string]interface{}{"format": "epoch_millis", "gt": "x"}, 0, 0},
		{"bad epoch gte", map[string]interface{}{"format": "epoch_millis", "gte": "x"}, 0, 0},
		{"bad epoch lt", map[string]interface{}{"format": "epoch_millis", "lt": "x"}, 0, 0},
		{"bad epoch lte", map[string]interface{}{"format": "epoch_millis", "lte": "x"}, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := RangeQueryTime(meta.TimeFieldName, tt.query)
			assert.Equal(t, tt.wantMin, gotMin)
			assert.Equal(t, tt.max, gotMax)
		})
	}
}

func TestRangeQuery(t *testing.T) {
	min, max := RangeQuery(map[string]interface{}{meta.TimeFieldName: map[string]interface{}{"gte": ts1, "lt": ts2}})
	assert.Equal(t, t1.UnixNano(), min)
	assert.Equal(t, t2.UnixNano(), max)

	min, max = RangeQuery(map[string]interface{}{"other": map[string]interface{}{"gte": ts1}})
	assert.Zero(t, min)
	assert.Zero(t, max)

	min, max = RangeQuery(map[string]interface{}{meta.TimeFieldName: "not-a-map"})
	assert.Zero(t, min)
	assert.Zero(t, max)
}

func TestQuery(t *testing.T) {
	rng := tsRange(map[string]interface{}{"gte": ts1, "lt": ts2})
	noRange := map[string]interface{}{"match_all": map[string]interface{}{}}
	tests := []struct {
		name  string
		query interface{}
		found bool
	}{
		{"nil", nil, false},
		{"not a map", "x", false},
		{"value not a map", map[string]interface{}{"range": 1}, false},
		{"unknown query type", noRange, false},
		{"range", rng, true},
		{"RANGE upper case key", map[string]interface{}{"RANGE": rng["range"]}, true},
		{"bool must object", map[string]interface{}{"bool": map[string]interface{}{"must": rng}}, true},
		{"bool must array", map[string]interface{}{"bool": map[string]interface{}{"must": []interface{}{noRange, rng}}}, true},
		{"bool should object", map[string]interface{}{"bool": map[string]interface{}{"should": rng}}, true},
		{"bool should array", map[string]interface{}{"bool": map[string]interface{}{"should": []interface{}{rng}}}, true},
		{"bool filter object", map[string]interface{}{"bool": map[string]interface{}{"filter": rng}}, true},
		{"bool filter array", map[string]interface{}{"bool": map[string]interface{}{"filter": []interface{}{rng}}}, true},
		{"bool must_not ignored", map[string]interface{}{"bool": map[string]interface{}{"must_not": rng}}, false},
		{"bool nested", map[string]interface{}{"bool": map[string]interface{}{"must": map[string]interface{}{"bool": map[string]interface{}{"filter": []interface{}{rng}}}}}, true},
		{"bool array without range", map[string]interface{}{"bool": map[string]interface{}{"must": []interface{}{noRange}}}, false},
		{"meta.Query struct", &meta.Query{Range: map[string]*meta.RangeQuery{meta.TimeFieldName: {GTE: ts1, LT: ts2}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max := Query(tt.query)
			if !tt.found {
				assert.Zero(t, min)
				assert.Zero(t, max)
				return
			}
			assert.Equal(t, t1.UnixNano(), min)
			assert.Equal(t, t2.UnixNano(), max)
		})
	}
}
