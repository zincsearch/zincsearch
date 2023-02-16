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
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/blugelabs/bluge/search"
	"github.com/blugelabs/bluge/search/aggregations"

	zincaggregation "github.com/zinclabs/zincsearch/pkg/bluge/aggregation"
	"github.com/zinclabs/zincsearch/pkg/config"
	"github.com/zinclabs/zincsearch/pkg/errors"
	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/pkg/zutils"
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
			req.AddAggregation(name, aggregations.Min(search.Field(agg.Min.Field)))
		case agg.Sum != nil:
			req.AddAggregation(name, aggregations.Sum(search.Field(agg.Sum.Field)))
		case agg.Count != nil:
			req.AddAggregation(name, aggregations.CountMatches())
		case agg.Cardinality != nil:
			req.AddAggregation(name, aggregations.Cardinality(search.Field(agg.Cardinality.Field)))
		case agg.Terms != nil:
			if agg.Terms.Size == 0 {
				agg.Terms.Size = config.Global.AggregationTermsSize
			}
			var subreq *zincaggregation.TermsAggregation
			prop, _ := mappings.GetProperty(agg.Terms.Field)
			switch prop.Type {
			case "text", "keyword":
				subreq = zincaggregation.NewTermsAggregation(search.Field(agg.Terms.Field), zincaggregation.TextValueSource, agg.Terms.Size)
			case "numeric":
				subreq = zincaggregation.NewTermsAggregation(search.Field(agg.Terms.Field), zincaggregation.NumericValueSource, agg.Terms.Size)
			case "bool", "boolean":
				subreq = zincaggregation.NewTermsAggregation(search.Field(agg.Terms.Field), zincaggregation.BooleanValueSource, agg.Terms.Size)
			default:
				return errors.New(
					errors.ErrorTypeParsingException,
					fmt.Sprintf("[terms] aggregation doesn't support values of type: [%s:[%s]]", agg.Terms.Field, prop.Type),
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
			prop, _ := mappings.GetProperty(agg.Range.Field)
			switch prop.Type {
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
			prop, ok := mappings.GetProperty(agg.DateRange.Field)
			if ok {
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
			switch prop.Type {
			case "date", "time":
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
				agg.Histogram.Size = config.Global.AggregationTermsSize
			}
			if agg.Histogram.Interval <= 0 {
				return errors.New(errors.ErrorTypeParsingException, "[histogram] aggregation interval must be a positive decimal")
			}
			if agg.Histogram.Offset >= agg.Histogram.Interval {
				return errors.New(errors.ErrorTypeParsingException, "[histogram] aggregation offset must be in [0, interval)")
			}
			var subreq *zincaggregation.HistogramAggregation
			prop, _ := mappings.GetProperty(agg.Histogram.Field)
			switch prop.Type {
			case "numeric":
				subreq = zincaggregation.NewHistogramAggregation(
					search.Field(agg.Histogram.Field),
					agg.Histogram.Interval,
					agg.Histogram.Offset,
					agg.Histogram.ExtendedBounds,
					agg.Histogram.HardBounds,
					agg.Histogram.MinDocCount,
					agg.Histogram.Size,
				)
			default:
				return errors.New(
					errors.ErrorTypeParsingException,
					fmt.Sprintf("[histogram] aggregation doesn't support values of type: [%s:[%s]]", agg.Histogram.Field, prop.Type),
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
				agg.DateHistogram.Size = config.Global.AggregationTermsSize
			}
			if agg.DateHistogram.Interval != "" {
				agg.DateHistogram.FixedInterval = agg.DateHistogram.Interval
			}
			if agg.DateHistogram.CalendarInterval == "" && agg.DateHistogram.FixedInterval == "" {
				return errors.New(errors.ErrorTypeParsingException, "[date_histogram] aggregation calendar_interval or fixed_interval must be set one")
			}

			// format interval
			var interval int64
			var calendarInterval = agg.DateHistogram.CalendarInterval
			if agg.DateHistogram.CalendarInterval != "" {
				switch agg.DateHistogram.CalendarInterval {
				case "second", "1s":
					interval = int64(time.Second)
					calendarInterval = ""
				case "minute", "1m":
					interval = int64(time.Minute)
					calendarInterval = ""
				case "hour", "1h":
					interval = int64(time.Hour)
					calendarInterval = ""
				case "day", "1d":
					interval = int64(time.Hour * 24)
					calendarInterval = ""
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
			prop, _ := mappings.GetProperty(agg.DateHistogram.Field)
			switch prop.Type {
			case "date", "time":
				subreq = zincaggregation.NewDateHistogramAggregation(
					search.Field(agg.DateHistogram.Field),
					calendarInterval,
					interval,
					agg.DateHistogram.Format,
					timeZone,
					agg.DateHistogram.ExtendedBounds,
					agg.DateHistogram.HardBounds,
					agg.DateHistogram.MinDocCount,
					agg.DateHistogram.Size,
				)
			default:
				return errors.New(
					errors.ErrorTypeParsingException,
					fmt.Sprintf(
						"[date_histogram] aggregation doesn't support values of type: [%s:[%s]]",
						agg.DateHistogram.Field,
						prop.Type,
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
			prop, _ := mappings.GetProperty(agg.AutoDateHistogram.Field)
			switch prop.Type {
			case "date", "time":
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
						prop.Type,
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
			f := v.Value()
			if math.IsNaN(f) {
				f = 0
			}
			resp[name] = meta.AggregationResponse{Value: f}
		case search.DurationCalculator:
			resp[name] = meta.AggregationResponse{Value: v.Duration().Milliseconds()}
		case search.BucketCalculator:
			buckets := v.Buckets()
			aggResp := meta.AggregationResponse{Buckets: make([]map[string]interface{}, 0)}
			aggRespBuckets := make([]map[string]interface{}, 0)
			for _, bucket := range buckets {
				aggBucket := map[string]interface{}{"key": bucket.Name(), "doc_count": bucket.Count()}
				if zutils.IsNumeric(bucket.Name()) {
					key, _ := strconv.ParseInt(bucket.Name(), 10, 64)
					aggBucket["key"] = key
					aggBucket["key_as_string"] = bucket.Name()
				}
				if subAggs := bucket.Aggregations(); len(subAggs) > 1 {
					subResp, err := Response(bucket)
					if err != nil {
						return nil, err
					}
					delete(subResp, "count")
					for k, v := range subResp {
						aggBucket[k] = v
					}
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
