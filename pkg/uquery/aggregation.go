package uquery

import (
	"fmt"

	"github.com/blugelabs/bluge/search"
	"github.com/blugelabs/bluge/search/aggregations"

	zincaggregation "github.com/zinclabs/zinc/pkg/bluge/aggregation"
	v1 "github.com/zinclabs/zinc/pkg/meta/v1"
	meta "github.com/zinclabs/zinc/pkg/meta/v2"
)

func AddAggregations(req zincaggregation.SearchAggregation, aggs map[string]v1.AggregationParams, mapping *meta.Mappings) error {
	if len(aggs) == 0 {
		return nil // not need aggregation
	}
	if mapping == nil {
		return nil // mapping is empty, return
	}

	// handle aggregation
	for name, agg := range aggs {
		if agg.Size == 0 || agg.Size >= 100 {
			agg.Size = 100 // default returns top 100 aggregation results
		}
		switch agg.AggType {
		case "term", "terms":
			var subreq *zincaggregation.TermsAggregation
			switch mapping.Properties[agg.Field].Type {
			case "text", "keyword":
				subreq = zincaggregation.NewTermsAggregation(search.Field(agg.Field), zincaggregation.TextValueSource, agg.Size)
			case "numeric":
				subreq = zincaggregation.NewTermsAggregation(search.Field(agg.Field), zincaggregation.NumericValueSource, agg.Size)
			default:
				return fmt.Errorf("terms aggregation not supported type [%s:[%s]]", agg.Field, mapping.Properties[agg.Field].Type)
			}
			if len(agg.Aggregations) > 0 {
				if err := AddAggregations(subreq, agg.Aggregations, mapping); err != nil {
					return err
				}
			}
			req.AddAggregation(name, subreq)
		case "range":
			if len(agg.Ranges) == 0 {
				return fmt.Errorf("range aggregation needs ranges")
			}
			var subreq *aggregations.RangeAggregation
			switch mapping.Properties[agg.Field].Type {
			case "numeric":
				subreq = aggregations.Ranges(search.Field(agg.Field))
				for _, v := range agg.Ranges {
					subreq.AddRange(aggregations.Range(v.From, v.To))
				}
				req.AddAggregation(name, subreq)
			default:
				return fmt.Errorf("range aggregation only support type numeric")
			}
		case "date_range":
			if len(agg.DateRanges) == 0 {
				return fmt.Errorf("date_range aggregation needs date_ranges")
			}
			var subreq *aggregations.DateRangeAggregation
			switch mapping.Properties[agg.Field].Type {
			case "time":
				subreq = aggregations.DateRanges(search.Field(agg.Field))
				// time format: 2022-01-21T09:22:50.604Z
				for _, v := range agg.DateRanges {
					subreq.AddRange(aggregations.NewDateRange(v.From, v.To))
				}
				req.AddAggregation(name, subreq)
			default:
				return fmt.Errorf("date_range aggregation only support type datetime")
			}
		case "max":
			req.AddAggregation(name, aggregations.Max(search.Field(agg.Field)))
		case "min":
			req.AddAggregation(name, aggregations.Min(search.Field(agg.Field)))
		case "avg":
			req.AddAggregation(name, aggregations.Avg(search.Field(agg.Field)))
		case "weighted_avg":
			req.AddAggregation(name, aggregations.WeightedAvg(search.Field(agg.Field), search.Field(agg.WeightField)))
		case "sum":
			req.AddAggregation(name, aggregations.Sum(search.Field(agg.Field)))
		case "count":
			req.AddAggregation(name, aggregations.CountMatches())
		default:
			return fmt.Errorf("aggregation not supported %s", agg.AggType)
		}
	}

	return nil
}

func ParseAggregations(bucket *search.Bucket) (map[string]v1.AggregationResponse, error) {
	resp := make(map[string]v1.AggregationResponse)
	aggs := bucket.Aggregations()
	for name, v := range aggs {
		switch v := v.(type) {
		case search.MetricCalculator:
			resp[name] = v1.AggregationResponse{Value: v.Value()}
		case search.DurationCalculator:
			resp[name] = v1.AggregationResponse{Value: v.Duration().Milliseconds()}
		case search.BucketCalculator:
			buckets := v.Buckets()
			aggResp := v1.AggregationResponse{Buckets: make([]v1.AggregationBucket, 0)}
			for _, bucket := range buckets {
				aggBucket := v1.AggregationBucket{Key: bucket.Name(), DocCount: bucket.Count()}
				if subAggs := bucket.Aggregations(); len(subAggs) > 1 {
					subResp, err := ParseAggregations(bucket)
					if err != nil {
						return nil, err
					}
					delete(subResp, "count")
					aggBucket.Aggregations = subResp
				}
				aggResp.Buckets = append(aggResp.Buckets, aggBucket)
			}
			resp[name] = aggResp
		default:
			return nil, fmt.Errorf("aggregation not supported agg type: [%s:%T]", name, v)
		}
	}

	return resp, nil
}
