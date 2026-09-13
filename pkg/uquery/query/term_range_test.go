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

package query

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vcaesar/riot"

	"github.com/zincsearch/zincsearch/pkg/meta"
)

// testMappings returns a mapping with one field per type used across the query tests.
func testMappings() *meta.Mappings {
	m := meta.NewMappings()
	m.SetProperty("name", meta.NewProperty("text"))
	title := meta.NewProperty("text")
	title.Analyzer = "standard"
	m.SetProperty("title", title)
	m.SetProperty("tag", meta.NewProperty("keyword"))
	m.SetProperty("age", meta.NewProperty("numeric"))
	m.SetProperty("active", meta.NewProperty("bool"))
	m.SetProperty(meta.TimeFieldName, meta.NewProperty("date"))
	custom := meta.NewProperty("date")
	custom.Format = "2006-01-02"
	custom.TimeZone = "+08:00"
	m.SetProperty("day", custom)
	return m
}

func TestTermQuery(t *testing.T) {
	mappings := testMappings()

	t.Run("text string", func(t *testing.T) {
		q, err := TermQuery(map[string]interface{}{"tag": "go"}, mappings)
		assert.NoError(t, err)
		tq := q.(*riot.TermQuery)
		assert.Equal(t, "tag", tq.Field())
		assert.Equal(t, "go", tq.Term())
		assert.Equal(t, 1.0, tq.Boost())
	})

	t.Run("text object with boost", func(t *testing.T) {
		q, err := TermQuery(map[string]interface{}{"tag": map[string]interface{}{"VALUE": "go", "boost": 2.5, "case_insensitive": true, "x": 1}}, mappings)
		assert.NoError(t, err)
		tq := q.(*riot.TermQuery)
		assert.Equal(t, "go", tq.Term())
		assert.Equal(t, 2.5, tq.Boost())
	})

	t.Run("unmapped field is text", func(t *testing.T) {
		q, err := TermQuery(map[string]interface{}{"other": 3.0}, mappings)
		assert.NoError(t, err)
		assert.Equal(t, "3", q.(*riot.TermQuery).Term())
	})

	t.Run("numeric", func(t *testing.T) {
		q, err := TermQuery(map[string]interface{}{"age": 18.0}, mappings)
		assert.NoError(t, err)
		nq := q.(*riot.NumericRangeQuery)
		min, incMin := nq.Min()
		max, incMax := nq.Max()
		assert.Equal(t, 18.0, min)
		assert.Equal(t, 18.0, max)
		assert.True(t, incMin)
		assert.True(t, incMax)
		assert.Equal(t, "age", nq.Field())
	})

	t.Run("numeric from string with boost", func(t *testing.T) {
		q, err := TermQuery(map[string]interface{}{"age": map[string]interface{}{"value": "21", "boost": 3.0}}, mappings)
		assert.NoError(t, err)
		nq := q.(*riot.NumericRangeQuery)
		min, _ := nq.Min()
		assert.Equal(t, 21.0, min)
		assert.Equal(t, 3.0, nq.Boost())
	})

	t.Run("numeric invalid", func(t *testing.T) {
		_, err := TermQuery(map[string]interface{}{"age": "abc"}, mappings)
		assert.ErrorContains(t, err, "numeric")
	})

	t.Run("bool", func(t *testing.T) {
		for _, v := range []interface{}{true, "true"} {
			q, err := TermQuery(map[string]interface{}{"active": v}, mappings)
			assert.NoError(t, err)
			tq := q.(*riot.TermQuery)
			assert.Equal(t, "true", tq.Term())
			assert.Equal(t, "active", tq.Field())
		}
		q, err := TermQuery(map[string]interface{}{"active": map[string]interface{}{"value": false, "boost": 0.5}}, mappings)
		assert.NoError(t, err)
		assert.Equal(t, "false", q.(*riot.TermQuery).Term())
		assert.Equal(t, 0.5, q.(*riot.TermQuery).Boost())
	})

	t.Run("bool invalid", func(t *testing.T) {
		_, err := TermQuery(map[string]interface{}{"active": "maybe"}, mappings)
		assert.ErrorContains(t, err, "boolean")
	})

	t.Run("errors", func(t *testing.T) {
		_, err := TermQuery(map[string]interface{}{"a": "x", "b": "y"}, mappings)
		assert.ErrorContains(t, err, "multiple fields")
		_, err = TermQuery(map[string]interface{}{"tag": []interface{}{"x"}}, mappings)
		assert.ErrorContains(t, err, "doesn't support values of type")
	})
}

func TestTermsQuery(t *testing.T) {
	mappings := testMappings()

	t.Run("strings", func(t *testing.T) {
		q, err := TermsQuery(map[string]interface{}{"tag": []interface{}{"a", "b"}}, mappings)
		assert.NoError(t, err)
		bq := q.(*riot.BooleanQuery)
		assert.Len(t, bq.Shoulds(), 2)
		assert.Equal(t, "a", bq.Shoulds()[0].(*riot.TermQuery).Term())
		assert.Equal(t, "b", bq.Shoulds()[1].(*riot.TermQuery).Term())
		assert.Equal(t, 1.0, bq.Boost())
	})

	t.Run("typed slices and boost", func(t *testing.T) {
		q, err := TermsQuery(map[string]interface{}{"tag": []string{"a"}, "BOOST": 2.0}, mappings)
		assert.NoError(t, err)
		assert.Len(t, q.(*riot.BooleanQuery).Shoulds(), 1)
		assert.Equal(t, 2.0, q.(*riot.BooleanQuery).Boost())

		q, err = TermsQuery(map[string]interface{}{"age": []float64{1, 2, 3}}, mappings)
		assert.NoError(t, err)
		assert.Len(t, q.(*riot.BooleanQuery).Shoulds(), 3)
		assert.IsType(t, &riot.NumericRangeQuery{}, q.(*riot.BooleanQuery).Shoulds()[0])

		q, err = TermsQuery(map[string]interface{}{"age": []int{1, 2}}, mappings)
		assert.NoError(t, err)
		assert.Len(t, q.(*riot.BooleanQuery).Shoulds(), 2)

		q, err = TermsQuery(map[string]interface{}{"active": []bool{true, false}}, mappings)
		assert.NoError(t, err)
		assert.Equal(t, "true", q.(*riot.BooleanQuery).Shoulds()[0].(*riot.TermQuery).Term())
		assert.Equal(t, "false", q.(*riot.BooleanQuery).Shoulds()[1].(*riot.TermQuery).Term())
	})

	t.Run("mixed interface values", func(t *testing.T) {
		q, err := TermsQuery(map[string]interface{}{"x": []interface{}{"a", 1.0, 2, true}}, mappings)
		assert.NoError(t, err)
		assert.Len(t, q.(*riot.BooleanQuery).Shoulds(), 4)
	})

	t.Run("errors", func(t *testing.T) {
		_, err := TermsQuery(map[string]interface{}{"a": []interface{}{}, "b": []interface{}{}, "c": []interface{}{}}, mappings)
		assert.ErrorContains(t, err, "multiple fields")
		_, err = TermsQuery(map[string]interface{}{"tag": "a"}, mappings)
		assert.ErrorContains(t, err, "doesn't support values of type")
		_, err = TermsQuery(map[string]interface{}{"tag": []interface{}{map[string]interface{}{}}}, mappings)
		assert.ErrorContains(t, err, "doesn't support values of type")
	})
}

func TestIdsQuery(t *testing.T) {
	mappings := testMappings()
	for name, q := range map[string]map[string]interface{}{
		"string slice":    {"values": []string{"1", "2"}},
		"interface slice": {"values": []interface{}{"1", "2"}},
		"object":          {"values": map[string]interface{}{"VALUE": []interface{}{"1", "2"}, "other": 1}},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := IdsQuery(q, mappings)
			assert.NoError(t, err)
			bq := got.(*riot.BooleanQuery)
			assert.Len(t, bq.Shoulds(), 2)
			assert.Equal(t, "_id", bq.Shoulds()[0].(*riot.TermQuery).Field())
			assert.Equal(t, "1", bq.Shoulds()[0].(*riot.TermQuery).Term())
		})
	}

	_, err := IdsQuery(map[string]interface{}{"a": []string{}, "b": []string{}}, mappings)
	assert.ErrorContains(t, err, "multiple fields")
	_, err = IdsQuery(map[string]interface{}{"values": "1"}, mappings)
	assert.ErrorContains(t, err, "doesn't support values of type")
	_, err = IdsQuery(map[string]interface{}{"values": map[string]interface{}{"value": "1"}}, mappings)
	assert.ErrorContains(t, err, "doesn't support values of type")
}

func TestRangeQuery(t *testing.T) {
	mappings := testMappings()

	t.Run("errors", func(t *testing.T) {
		_, err := RangeQuery(map[string]interface{}{"a": map[string]interface{}{}, "b": map[string]interface{}{}}, mappings)
		assert.ErrorContains(t, err, "multiple fields")
		_, err = RangeQuery(map[string]interface{}{"age": "x"}, mappings)
		assert.ErrorContains(t, err, "doesn't support values of type")
		_, err = RangeQuery(map[string]interface{}{"tag": map[string]interface{}{"gte": 1}}, mappings)
		assert.ErrorContains(t, err, "only support values of [numeric, time]")
	})

	t.Run("empty query returns nil", func(t *testing.T) {
		q, err := RangeQuery(map[string]interface{}{}, mappings)
		assert.NoError(t, err)
		assert.Nil(t, q)
	})

	t.Run("numeric", func(t *testing.T) {
		tests := []struct {
			name             string
			query            map[string]interface{}
			min, max         float64
			incMin, incMax   bool
			boost            float64
		}{
			{"gt/lt", map[string]interface{}{"gt": 1.0, "lt": 10.0}, 1, 10, false, false, 1},
			{"gte/lte with boost", map[string]interface{}{"GTE": 1.0, "LTE": 10.0, "boost": 2.0}, 1, 10, true, true, 2},
			{"strings", map[string]interface{}{"gte": "5", "lt": "7"}, 5, 7, true, false, 1},
			{"no upper bound", map[string]interface{}{"gte": 5.0}, 5, float64(math.MaxInt64), true, false, 1},
			{"no bounds", map[string]interface{}{"unknown": 1}, 0, float64(math.MaxInt64), false, false, 1},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				q, err := RangeQuery(map[string]interface{}{"age": tt.query}, mappings)
				assert.NoError(t, err)
				nq := q.(*riot.NumericRangeQuery)
				min, incMin := nq.Min()
				max, incMax := nq.Max()
				assert.Equal(t, tt.min, min)
				assert.Equal(t, tt.max, max)
				assert.Equal(t, tt.incMin, incMin)
				assert.Equal(t, tt.incMax, incMax)
				assert.Equal(t, tt.boost, nq.Boost())
				assert.Equal(t, "age", nq.Field())
			})
		}
	})

	t.Run("time", func(t *testing.T) {
		t1 := time.Date(2022, 3, 4, 5, 6, 7, 0, time.UTC)
		t2 := time.Date(2022, 3, 5, 5, 6, 7, 0, time.UTC)
		tests := []struct {
			name           string
			field          string
			query          map[string]interface{}
			start, end     time.Time
			incMin, incMax bool
			boost          float64
		}{
			{"rfc3339 gte/lt", meta.TimeFieldName, map[string]interface{}{"gte": t1.Format(time.RFC3339), "lt": t2.Format(time.RFC3339)}, t1, t2, true, false, 1},
			{"rfc3339 gt/lte boost", meta.TimeFieldName, map[string]interface{}{"gt": t1.Format(time.RFC3339), "lte": t2.Format(time.RFC3339), "boost": 3.0}, t1, t2, false, true, 3},
			{"epoch_millis", meta.TimeFieldName, map[string]interface{}{"format": "epoch_millis", "gte": float64(t1.UnixMilli()), "lte": float64(t2.UnixMilli())}, t1, t2, true, true, 1},
			{"epoch_millis gt/lt", meta.TimeFieldName, map[string]interface{}{"format": "epoch_millis", "gt": t1.UnixMilli(), "lt": t2.UnixMilli()}, t1, t2, false, false, 1},
			{"mapping format and time_zone", "day", map[string]interface{}{"gte": "2022-03-04", "lt": "2022-03-05"},
				time.Date(2022, 3, 3, 16, 0, 0, 0, time.UTC), time.Date(2022, 3, 4, 16, 0, 0, 0, time.UTC), true, false, 1},
			{"query overrides format and time_zone", "day", map[string]interface{}{"format": "2006-01-02 15:04", "time_zone": "UTC", "gte": "2022-03-04 05:06", "lt": "2022-03-05 05:06"},
				time.Date(2022, 3, 4, 5, 6, 0, 0, time.UTC), time.Date(2022, 3, 5, 5, 6, 0, 0, time.UTC), true, false, 1},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				q, err := RangeQuery(map[string]interface{}{tt.field: tt.query}, mappings)
				assert.NoError(t, err)
				dq := q.(*riot.DateRangeQuery)
				start, incMin := dq.Start()
				end, incMax := dq.End()
				assert.True(t, tt.start.Equal(start), "start %v != %v", tt.start, start)
				assert.True(t, tt.end.Equal(end), "end %v != %v", tt.end, end)
				assert.Equal(t, tt.incMin, incMin)
				assert.Equal(t, tt.incMax, incMax)
				assert.Equal(t, tt.boost, dq.Boost())
				assert.Equal(t, tt.field, dq.Field())
			})
		}

		t.Run("missing upper bound defaults to now", func(t *testing.T) {
			before := time.Now()
			q, err := RangeQuery(map[string]interface{}{meta.TimeFieldName: map[string]interface{}{"gte": t1.Format(time.RFC3339)}}, mappings)
			assert.NoError(t, err)
			end, _ := q.(*riot.DateRangeQuery).End()
			assert.False(t, end.Before(before))
		})

		t.Run("nil mappings falls back to rfc3339", func(t *testing.T) {
			q, err := RangeQueryTime("ts", map[string]interface{}{"gte": t1.Format(time.RFC3339)}, nil)
			assert.NoError(t, err)
			start, _ := q.(*riot.DateRangeQuery).Start()
			assert.True(t, t1.Equal(start))
		})

		for _, key := range []string{"gt", "gte", "lt", "lte"} {
			t.Run("bad "+key, func(t *testing.T) {
				_, err := RangeQuery(map[string]interface{}{meta.TimeFieldName: map[string]interface{}{key: "nope"}}, mappings)
				assert.ErrorContains(t, err, "range."+key+" format err")
			})
		}
		t.Run("bad time_zone", func(t *testing.T) {
			_, err := RangeQuery(map[string]interface{}{meta.TimeFieldName: map[string]interface{}{"gte": t1.Format(time.RFC3339), "time_zone": "Nowhere/City"}}, mappings)
			assert.ErrorContains(t, err, "time_zone parse err")
		})
	})
}
