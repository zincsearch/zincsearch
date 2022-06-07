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
	"fmt"
	"strings"
	"time"

	"github.com/blugelabs/bluge"

	"github.com/zinclabs/zinc/pkg/errors"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/zutils"
)

func RangeQuery(query map[string]interface{}, mappings *meta.Mappings) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[range] query doesn't support multiple fields")
	}

	for field, v := range query {
		vv, ok := v.(map[string]interface{})
		if !ok {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[range] query doesn't support values of type: %T", v))
		}
		prop, _ := mappings.GetProperty(field)
		switch prop.Type {
		case "numeric":
			return RangeQueryNumeric(field, vv, mappings)
		case "date", "time":
			return RangeQueryTime(field, vv, mappings)
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[range] %s only support values of [numeric, time]", field))
		}
	}

	return nil, nil
}

func RangeQueryNumeric(field string, query map[string]interface{}, mappings *meta.Mappings) (bluge.Query, error) {
	value := new(meta.RangeQuery)
	value.Boost = -1.0
	for k, v := range query {
		k := strings.ToLower(k)
		switch k {
		case "gt":
			value.GT = v.(float64)
		case "gte":
			value.GTE = v.(float64)
		case "lt":
			value.LT = v.(float64)
		case "lte":
			value.LTE = v.(float64)
		case "format":
			value.Format = v.(string)
		case "time_zone":
			value.TimeZone = v.(string)
		case "boost":
			value.Boost = v.(float64)
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[range] unknown field [%s]", k))
		}
	}

	min := 0.0
	max := 0.0
	minInclusive := false
	maxInclusive := false
	if value.GT != nil && value.GT.(float64) > 0 {
		min = value.GT.(float64)

	}
	if value.GTE != nil && value.GTE.(float64) > 0 {
		min = value.GTE.(float64)
		minInclusive = true
	}
	if value.LT != nil && value.LT.(float64) > 0 {
		max = value.LT.(float64)
	}
	if value.LTE != nil && value.LTE.(float64) > 0 {
		max = value.LTE.(float64)
		maxInclusive = true
	}
	subq := bluge.NewNumericRangeInclusiveQuery(min, max, minInclusive, maxInclusive).SetField(field)
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}

func RangeQueryTime(field string, query map[string]interface{}, mappings *meta.Mappings) (bluge.Query, error) {
	value := new(meta.RangeQuery)
	value.Boost = -1.0
	for k, v := range query {
		k := strings.ToLower(k)
		switch k {
		case "gt":
			value.GT = v
		case "gte":
			value.GTE = v
		case "lt":
			value.LT = v
		case "lte":
			value.LTE = v
		case "format":
			value.Format = v.(string)
		case "time_zone":
			value.TimeZone = v.(string)
		case "boost":
			value.Boost = v.(float64)
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[range] unknown field [%s]", k))
		}
	}

	var err error
	format := time.RFC3339
	if mappings != nil {
		if prop, ok := mappings.GetProperty(field); ok {
			if prop.Format != "" {
				format = prop.Format
			}
		}
	}
	if value.Format != "" {
		format = value.Format
	}
	timeZone := time.UTC
	if value.TimeZone != "" {
		timeZone, err = zutils.ParseTimeZone(value.TimeZone)
		if err != nil {
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[range] %s time_zone parse err %s", field, err.Error()))
		}
	}

	min := time.Time{}
	max := time.Time{}
	minInclusive := false
	maxInclusive := false
	if value.GT != nil {
		if format == "epoch_millis" {
			min = time.UnixMilli(int64(value.GT.(float64)))
		} else {
			min, err = time.ParseInLocation(format, value.GT.(string), timeZone)
		}
		if err != nil {
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[range] %s range.gt format err %s", field, err.Error()))
		}
	}
	if value.GTE != nil {
		minInclusive = true
		if format == "epoch_millis" {
			min = time.UnixMilli(int64(value.GTE.(float64)))
		} else {
			min, err = time.ParseInLocation(format, value.GTE.(string), timeZone)
		}
		if err != nil {
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[range] %s range.gte format err %s", field, err.Error()))
		}
	}
	if value.LT != nil {
		if format == "epoch_millis" {
			max = time.UnixMilli(int64(value.LT.(float64)))
		} else {
			max, err = time.ParseInLocation(format, value.LT.(string), timeZone)
		}
		if err != nil {
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[range] %s range.lt format err %s", field, err.Error()))
		}
	}
	if value.LTE != nil {
		maxInclusive = true
		if format == "epoch_millis" {
			max = time.UnixMilli(int64(value.LTE.(float64)))
		} else {
			max, err = time.ParseInLocation(format, value.LTE.(string), timeZone)
		}
		if err != nil {
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[range] %s range.lte format err %s", field, err.Error()))
		}
	}
	subq := bluge.NewDateRangeInclusiveQuery(min.UTC(), max.UTC(), minInclusive, maxInclusive).SetField(field)
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}
