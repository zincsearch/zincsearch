package aggregation

import (
	"fmt"
	"time"

	"github.com/blugelabs/bluge/search"
	"github.com/blugelabs/bluge/search/aggregations"

	zincaggregation "github.com/prabhatsharma/zinc/pkg/bluge/aggregation"
	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
	"github.com/prabhatsharma/zinc/pkg/startup"
	"github.com/prabhatsharma/zinc/pkg/zutils"
)

func Request(req zincaggregation.SearchAggregation, aggs map[string]meta.Aggregations, mappings *meta.Mappings) error {
	if len(aggs) == 0 {
		return nil // not need aggregation
	}
	if mappings == nil {
		return nil // mapping is empty
	}

	var err error
	// handle aggregation
	for name, agg := range aggs {
		switch {
		case agg.Avg != nil:
			req.AddAggregation(name, aggregations.Avg(search.Field(agg.Avg.Field)))
		case agg.WeightedAvg != nil:
			req.AddAggregation(name, aggregations.WeightedAvg(search.Field(agg.WeightedAvg.Field), search.Field(agg.WeightedAvg.WeightField)))
		case agg.Max != nil:
			req.AddAggregation(name, aggregations.Max(search.Field(agg.Max.Field)))
		case agg.Min != nil:
			req.AddAggregation(name, aggregations.Min(search.Field(agg.Max.Field)))
		case agg.Sum != nil:
			req.AddAggregation(name, aggregations.Sum(search.Field(agg.Sum.Field)))
		case agg.Count != nil:
			req.AddAggregation(name, aggregations.CountMatches())
		case agg.Terms != nil:
			if agg.Terms.Size == 0 {
				agg.Terms.Size = startup.LoadAggregationTermsSize()
			}
			var subreq *zincaggregation.TermsAggregation
			switch mappings.Properties[agg.Terms.Field].Type {
			case "text", "keyword":
				subreq = zincaggregation.NewTermsAggregation(search.Field(agg.Terms.Field), zincaggregation.TextValueSource, agg.Terms.Size)
			case "numeric":
				subreq = zincaggregation.NewTermsAggregation(search.Field(agg.Terms.Field), zincaggregation.NumericValueSource, agg.Terms.Size)
			default:
				return errors.New(
					errors.ErrorTypeParsingException,
					fmt.Sprintf("[terms] aggregation doesn't support values of type: [%s:[%v]]", agg.Terms.Field, mappings.Properties[agg.Terms.Field].Type),
				)
			}
			if len(agg.Aggregations) > 0 {
				if err := Request(subreq, agg.Aggregations, mappings); err != nil {
					return err
				}
			}
			req.AddAggregation(name, subreq)
		case agg.Range != nil:
			if len(agg.Range.Ranges) == 0 {
				return errors.New(errors.ErrorTypeParsingException, "[range] aggregation needs ranges")
			}
			var subreq *aggregations.RangeAggregation
			switch mappings.Properties[agg.Range.Field].Type {
			case "numeric":
				subreq = aggregations.Ranges(search.Field(agg.Range.Field))
				for _, v := range agg.Range.Ranges {
					subreq.AddRange(aggregations.Range(v.From, v.To))
				}
				req.AddAggregation(name, subreq)
			default:
				return errors.New(errors.ErrorTypeParsingException, "[range] aggregation only support type numeric")
			}
		case agg.DateRange != nil:
			if len(agg.DateRange.Ranges) == 0 {
				return errors.New(errors.ErrorTypeParsingException, "[date_range] aggregation needs ranges")
			}
			var subreq *aggregations.DateRangeAggregation
			format := time.RFC3339
			if prop, ok := mappings.Properties[agg.DateRange.Field]; ok {
				if prop.Format != "" {
					format = prop.Format
				}
			}
			if agg.DateRange.Format != "" {
				format = agg.DateRange.Format
			}
			timeZone := time.UTC
			if agg.DateRange.TimeZone != "" {
				timeZone, err = zutils.ParseTimeZone(agg.DateRange.TimeZone)
				if err != nil {
					return errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[date_range] time_zone parse err %v", err))
				}
			}
			switch mappings.Properties[agg.DateRange.Field].Type {
			case "time":
				subreq = aggregations.DateRanges(search.Field(agg.DateRange.Field))
				for _, v := range agg.DateRange.Ranges {
					from := time.Time{}
					to := time.Time{}
					if v.From != "" {
						from, err = time.ParseInLocation(format, v.From, timeZone)
						if err != nil {
							return errors.New(errors.ErrorTypeIllegalArgumentException, "[date_range] range value from parse error "+err.Error())
						}
					}
					if v.To != "" {
						to, err = time.ParseInLocation(format, v.To, timeZone)
						if err != nil {
							return errors.New(errors.ErrorTypeIllegalArgumentException, "[date_range] range value to parse error "+err.Error())
						}
					}
					subreq.AddRange(aggregations.NewDateRange(from, to))
				}
				req.AddAggregation(name, subreq)
			default:
				return errors.New(errors.ErrorTypeParsingException, "[date_range] aggregation only support type datetime")
			}
		case agg.IPRange != nil:
			return errors.New(errors.ErrorTypeNotImplemented, "[ip_range] aggregation doesn't support")
		case agg.Histogram != nil:
			return errors.New(errors.ErrorTypeNotImplemented, "[histogram] aggregation doesn't support")
		case agg.DateHistogram != nil:
			return errors.New(errors.ErrorTypeNotImplemented, "[date_histogram] aggregation doesn't support")
		default:
			// nothing
		}
	}

	return nil
}

func Response(bucket *search.Bucket) (map[string]meta.AggregationResponse, error) {
	resp := make(map[string]meta.AggregationResponse)
	aggs := bucket.Aggregations()
	for name, v := range aggs {
		switch v := v.(type) {
		case search.MetricCalculator:
			resp[name] = meta.AggregationResponse{Value: v.Value()}
		case search.DurationCalculator:
			resp[name] = meta.AggregationResponse{Value: v.Duration().Milliseconds()}
		case search.BucketCalculator:
			buckets := v.Buckets()
			aggResp := meta.AggregationResponse{Buckets: make([]meta.AggregationBucket, 0)}
			aggRespBuckets := make([]meta.AggregationBucket, 0)
			for _, bucket := range buckets {
				aggBucket := meta.AggregationBucket{Key: bucket.Name(), DocCount: bucket.Count()}
				if subAggs := bucket.Aggregations(); len(subAggs) > 1 {
					subResp, err := Response(bucket)
					if err != nil {
						return nil, err
					}
					delete(subResp, "count")
					aggBucket.Aggregations = subResp
				}
				aggRespBuckets = append(aggRespBuckets, aggBucket)
			}
			aggResp.Buckets = aggRespBuckets
			resp[name] = aggResp
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[%s:%T] aggregation doesn't support", name, v))
		}
	}

	return resp, nil
}
