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

package uquery

import (
	"fmt"
	"time"

	"github.com/zinclabs/zincsearch/pkg/meta"
	v1 "github.com/zinclabs/zincsearch/pkg/meta/v1"
)

// ParseQueryDSLFromV1 parse query DSL from search v1 and return new zinc query DSL
func ParseQueryDSLFromV1(q *v1.ZincQuery) (*meta.ZincQuery, error) {
	newquery := new(meta.ZincQuery)
	newquery.From = q.From
	newquery.Size = q.MaxResults
	newquery.Explain = q.Explain
	newquery.Highlight = q.Highlight
	newquery.Source = q.Source

	if q.SortFields != nil {
		sort := make([]interface{}, 0, len(q.SortFields))
		for _, field := range q.SortFields {
			sort = append(sort, field)
		}
		newquery.Sort = sort
	}

	query := make(map[string]interface{})
	boolQuery := make(map[string]interface{})
	mustQuery := make([]interface{}, 0)
	// time range
	if !q.Query.StartTime.IsZero() || !q.Query.EndTime.IsZero() {
		timeRangeQuery := make(map[string]interface{})
		timeRangeQuery["format"] = "epoch_millis"
		if !q.Query.StartTime.IsZero() {
			timeRangeQuery["gte"] = q.Query.StartTime.UnixMilli()
		}
		if !q.Query.EndTime.IsZero() {
			timeRangeQuery["lt"] = q.Query.EndTime.UnixMilli()
		}
		mustQuery = append(mustQuery, map[string]interface{}{
			"range": map[string]interface{}{
				"@timestamp": timeRangeQuery,
			},
		})
	}
	// search type
	switch q.SearchType {
	case "alldocuments":
		mustQuery = append(mustQuery, map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	case "wildcard":
		mustQuery = append(mustQuery, map[string]interface{}{
			"wildcard": map[string]interface{}{
				q.Query.Field: q.Query.Term,
			},
		})
	case "fuzzy":
		mustQuery = append(mustQuery, map[string]interface{}{
			"fuzzy": map[string]interface{}{
				q.Query.Field: q.Query.Term,
			},
		})
	case "term":
		mustQuery = append(mustQuery, map[string]interface{}{
			"term": map[string]interface{}{
				q.Query.Field: q.Query.Term,
			},
		})
	case "daterange":
		if q.Query.Field != "" && q.Query.Field != "@timestamp" {
			mustQuery = append(mustQuery, map[string]interface{}{
				"range": map[string]interface{}{
					q.Query.Field: map[string]interface{}{
						"gte":    q.Query.StartTime.UnixMilli(),
						"lt":     q.Query.EndTime.UnixMilli(),
						"format": "epoch_millis",
					},
				},
			})
		}
	case "matchall":
		mustQuery = append(mustQuery, map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	case "match":
		mustQuery = append(mustQuery, map[string]interface{}{
			"match": map[string]interface{}{
				q.Query.Field: q.Query.Term,
			},
		})
	case "matchphrase":
		mustQuery = append(mustQuery, map[string]interface{}{
			"match_phrase": map[string]interface{}{
				q.Query.Field: q.Query.Term,
			},
		})
	case "prefix":
		mustQuery = append(mustQuery, map[string]interface{}{
			"prefix": map[string]interface{}{
				q.Query.Field: q.Query.Term,
			},
		})
	case "querystring":
		mustQuery = append(mustQuery, map[string]interface{}{
			"query_string": map[string]interface{}{
				"query": q.Query.Term,
			},
		})
	default:
		// default use alldocuments search
		mustQuery = append(mustQuery, map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	}
	boolQuery["must"] = mustQuery
	query["bool"] = boolQuery
	newquery.Query = query

	// aggs
	if len(q.Aggregations) > 0 {
		newquery.Aggregations = make(map[string]meta.Aggregations)
	}
	for name, agg := range q.Aggregations {
		newagg, err := covertAggregationFromV1(agg)
		if err != nil {
			return nil, err
		}
		newquery.Aggregations[name] = newagg
	}
	return newquery, nil
}

func covertAggregationFromV1(agg v1.AggregationParams) (meta.Aggregations, error) {
	newagg := meta.Aggregations{}
	if len(agg.Aggregations) > 0 {
		for name, subagg := range agg.Aggregations {
			newsubagg, err := covertAggregationFromV1(subagg)
			if err != nil {
				return newagg, err
			}
			newagg.Aggregations[name] = newsubagg
		}
	}
	switch agg.AggType {
	case "term", "terms":
		newagg.Terms = &meta.AggregationsTerms{
			Field: agg.Field,
			Size:  agg.Size,
		}
	case "range":
		newagg.Range = &meta.AggregationRange{
			Field: agg.Field,
		}
		for _, v := range agg.Ranges {
			newagg.Range.Ranges = append(newagg.Range.Ranges, meta.Range{
				To:   v.To,
				From: v.From,
			})
		}
	case "date_range":
		newagg.DateRange = &meta.AggregationDateRange{
			Field: agg.Field,
		}
		for _, v := range agg.DateRanges {
			newagg.DateRange.Ranges = append(newagg.DateRange.Ranges, meta.DateRange{
				To:   v.To.Format(time.RFC3339),
				From: v.From.Format(time.RFC3339),
			})
		}
	case "max":
		newagg.Max = &meta.AggregationMetric{Field: agg.Field}
	case "min":
		newagg.Min = &meta.AggregationMetric{Field: agg.Field}
	case "avg":
		newagg.Avg = &meta.AggregationMetric{Field: agg.Field}
	case "weighted_avg":
		newagg.WeightedAvg = &meta.AggregationMetric{Field: agg.Field, WeightField: agg.WeightField}
	case "sum":
		newagg.Sum = &meta.AggregationMetric{Field: agg.Field}
	case "count":
		newagg.Count = &meta.AggregationMetric{Field: agg.Field}
	default:
		return newagg, fmt.Errorf("aggregation not supported %s", agg.AggType)
	}

	return newagg, nil
}
