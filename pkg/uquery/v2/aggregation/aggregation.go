package aggregation

import (
	"fmt"
	"time"

	"github.com/blugelabs/bluge/search"
	"github.com/blugelabs/bluge/search/aggregations"

	zincaggregation "github.com/zinclabs/zinc/pkg/bluge/aggregation"
	"github.com/zinclabs/zinc/pkg/errors"
	meta "github.com/zinclabs/zinc/pkg/meta/v2"
	"github.com/zinclabs/zinc/pkg/startup"
	"github.com/zinclabs/zinc/pkg/zutils"
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
					fmt.Sprintf("[terms] aggregation doesn't support values of type: [%s:[%s]]", agg.Terms.Field, mappings.Properties[agg.Terms.Field].Type),
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
					return errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[date_range] time_zone parse err %s", err.Error()))
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
							return errors.New(errors.ErrorTypeIllegalArgumentException, fmt.Sprintf("[date_range] range value from parse err %s", err.Error()))
						}
					}
					if v.To != "" {
						to, err = time.ParseInLocation(format, v.To, timeZone)
						if err != nil {
							return errors.New(errors.ErrorTypeIllegalArgumentException, fmt.Sprintf("[date_range] range value to parse err %s", err.Error()))
						}
					}
					subreq.AddRange(aggregations.NewDateRange(from, to))
				}
				req.AddAggregation(name, subreq)
			default:
				return errors.New(errors.ErrorTypeParsingException, "[date_range] aggregation only support type datetime")
			}
		case agg.Histogram != nil:
			if agg.Histogram.Size == 0 {
				agg.Histogram.Size = startup.LoadAggregationTermsSize()
			}
			if agg.Histogram.Interval <= 0 {
				return errors.New(errors.ErrorTypeParsingException, "[histogram] aggregation interval must be a positive decimal")
			}
			if agg.Histogram.Offset >= agg.Histogram.Interval {
				return errors.New(errors.ErrorTypeParsingException, "[histogram] aggregation offset must be in [0, interval)")
			}
			var subreq *zincaggregation.HistogramAggregation
			switch mappings.Properties[agg.Histogram.Field].Type {
			case "numeric":
				subreq = zincaggregation.NewHistogramAggregation(
					search.Field(agg.Histogram.Field),
					agg.Histogram.Interval,
					agg.Histogram.Offset,
					agg.Histogram.MinDocCount,
					agg.Histogram.Size,
				)
			default:
				return errors.New(
					errors.ErrorTypeParsingException,
					fmt.Sprintf("[histogram] aggregation doesn't support values of type: [%s:[%s]]", agg.Histogram.Field, mappings.Properties[agg.Histogram.Field].Type),
				)
			}
			if len(agg.Aggregations) > 0 {
				if err := Request(subreq, agg.Aggregations, mappings); err != nil {
					return err
				}
			}
			req.AddAggregation(name, subreq)
		case agg.DateHistogram != nil:
			if agg.DateHistogram.Size == 0 {
				agg.DateHistogram.Size = startup.LoadAggregationTermsSize()
			}
			if agg.DateHistogram.CalendarInterval == "" && agg.DateHistogram.FixedInterval == "" {
				return errors.New(errors.ErrorTypeParsingException, "[date_histogram] aggregation calendar_interval or fixed_interval must be set one")
			}

			// format interval
			var interval int64
			if agg.DateHistogram.CalendarInterval != "" {
				switch agg.DateHistogram.CalendarInterval {
				case "second", "1s":
					interval = int64(time.Second)
					agg.DateHistogram.CalendarInterval = ""
				case "minute", "1m":
					interval = int64(time.Minute)
					agg.DateHistogram.CalendarInterval = ""
				case "hour", "1h":
					interval = int64(time.Hour)
					agg.DateHistogram.CalendarInterval = ""
				case "day", "1d":
					interval = int64(time.Hour * 24)
					agg.DateHistogram.CalendarInterval = ""
				case "week", "1w", "month", "1M", "quarter", "1q", "year", "1y":
					// calendar
				default:
					return errors.New(
						errors.ErrorTypeParsingException,
						"[date_histogram] aggregation calendar_interval must be Date Calendar, such as: second, minute, hour, day, week, month, quarter, year",
					)
				}
			} else if agg.DateHistogram.FixedInterval != "" {
				if duration, err := zutils.ParseDuration(agg.DateHistogram.FixedInterval); err != nil {
					return errors.New(errors.ErrorTypeParsingException, "[date_histogram] aggregation fixed_interval must be time duration, such as: 1s, 1m, 1h, 1d")
				} else {
					interval = int64(duration)
				}
			}

			timeZone := time.UTC
			if agg.DateHistogram.TimeZone != "" {
				timeZone, err = zutils.ParseTimeZone(agg.DateHistogram.TimeZone)
				if err != nil {
					return errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[date_histogram] time_zone parse err %s", err.Error()))
				}
			}
			if agg.DateHistogram.Format == "" {
				agg.DateHistogram.Format = time.RFC3339
			}
			var subreq *zincaggregation.DateHistogramAggregation
			switch mappings.Properties[agg.DateHistogram.Field].Type {
			case "time":
				subreq = zincaggregation.NewDateHistogramAggregation(
					search.Field(agg.DateHistogram.Field),
					agg.DateHistogram.CalendarInterval,
					interval,
					agg.DateHistogram.Format,
					timeZone,
					agg.DateHistogram.MinDocCount,
					agg.DateHistogram.Size,
				)
			default:
				return errors.New(
					errors.ErrorTypeParsingException,
					fmt.Sprintf(
						"[date_histogram] aggregation doesn't support values of type: [%s:[%s]]",
						agg.DateHistogram.Field,
						mappings.Properties[agg.DateHistogram.Field].Type,
					),
				)
			}
			if len(agg.Aggregations) > 0 {
				if err := Request(subreq, agg.Aggregations, mappings); err != nil {
					return err
				}
			}
			req.AddAggregation(name, subreq)
		case agg.AutoDateHistogram != nil:
			if agg.AutoDateHistogram.Buckets <= 0 {
				agg.AutoDateHistogram.Buckets = 10
			}
			if agg.AutoDateHistogram.MinimumInterval == "" {
				agg.AutoDateHistogram.MinimumInterval = "second"
			}
			if agg.AutoDateHistogram.MinimumInterval != "" {
				switch agg.AutoDateHistogram.MinimumInterval {
				case "second", "minute", "hour", "day", "month", "year":
					// calendar
				default:
					return errors.New(
						errors.ErrorTypeParsingException,
						"[auto_date_histogram] aggregation minimum_interval must be Date Calendar, such as: second, minute, hour, day, month, year",
					)
				}
			}

			timeZone := time.UTC
			if agg.AutoDateHistogram.TimeZone != "" {
				timeZone, err = zutils.ParseTimeZone(agg.AutoDateHistogram.TimeZone)
				if err != nil {
					return errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[auto_date_histogram] time_zone parse err %s", err.Error()))
				}
			}
			if agg.AutoDateHistogram.Format == "" {
				agg.AutoDateHistogram.Format = time.RFC3339
			}
			var subreq *zincaggregation.AutoDateHistogramAggregation
			switch mappings.Properties[agg.AutoDateHistogram.Field].Type {
			case "time":
				subreq = zincaggregation.NewAutoDateHistogramAggregation(
					search.Field(agg.AutoDateHistogram.Field),
					agg.AutoDateHistogram.Buckets,
					agg.AutoDateHistogram.MinimumInterval,
					agg.AutoDateHistogram.Format,
					timeZone,
				)
			default:
				return errors.New(
					errors.ErrorTypeParsingException,
					fmt.Sprintf(
						"[auto_date_histogram] aggregation doesn't support values of type: [%s:[%s]]",
						agg.AutoDateHistogram.Field,
						mappings.Properties[agg.AutoDateHistogram.Field].Type,
					),
				)
			}
			if len(agg.Aggregations) > 0 {
				if err := Request(subreq, agg.Aggregations, mappings); err != nil {
					return err
				}
			}
			req.AddAggregation(name, subreq)
		case agg.IPRange != nil:
			return errors.New(errors.ErrorTypeNotImplemented, "[ip_range] aggregation doesn't support")
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

			// hack: auto_date_histogram aggregation
			if v, ok := aggs[name].(*zincaggregation.AutoDateHistogramCalculator); ok {
				aggResp.Interval = v.Interval()
			}

			resp[name] = aggResp
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[%s:%T] aggregation doesn't support", name, v))
		}
	}

	return resp, nil
}
