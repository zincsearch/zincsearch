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
	"strings"
	"time"

	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/pkg/zutils"
)

func RangeQuery(query map[string]interface{}) (int64, int64) {
	for field, v := range query {
		if field == meta.TimeFieldName {
			vv, ok := v.(map[string]interface{})
			if !ok {
				return 0, 0
			}
			return RangeQueryTime(field, vv)
		}
	}
	return 0, 0
}

func RangeQueryTime(field string, query map[string]interface{}) (int64, int64) {
	value := new(meta.RangeQuery)
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
			return 0, 0
		}
	}

	var err error
	format := time.RFC3339
	if value.Format != "" {
		format = value.Format
	}
	timeZone := time.UTC
	if value.TimeZone != "" {
		timeZone, err = zutils.ParseTimeZone(value.TimeZone)
		if err != nil {
			return 0, 0
		}
	}

	min := time.Time{}
	max := time.Time{}
	if value.GT != nil {
		if format == "epoch_millis" {
			num, err := zutils.ToInt(value.GT)
			if err != nil {
				return 0, 0
			}

			min = time.UnixMilli(int64(num))
		} else {
			min, err = time.ParseInLocation(format, value.GT.(string), timeZone)
		}
		if err != nil {
			return 0, 0
		}
	}
	if value.GTE != nil {
		if format == "epoch_millis" {
			num, err := zutils.ToInt(value.GTE)
			if err != nil {
				return 0, 0
			}

			min = time.UnixMilli(int64(num))
		} else {
			min, err = time.ParseInLocation(format, value.GTE.(string), timeZone)
		}
		if err != nil {
			return 0, 0
		}
	}
	if value.LT != nil {
		if format == "epoch_millis" {
			num, err := zutils.ToInt(value.LT)
			if err != nil {
				return 0, 0
			}

			max = time.UnixMilli(int64(num))
		} else {
			max, err = time.ParseInLocation(format, value.LT.(string), timeZone)
		}
		if err != nil {
			return 0, 0
		}
	}
	if value.LTE != nil {
		if format == "epoch_millis" {
			num, err := zutils.ToInt(value.LTE)
			if err != nil {
				return 0, 0
			}

			max = time.UnixMilli(int64(num))
		} else {
			max, err = time.ParseInLocation(format, value.LTE.(string), timeZone)
		}
		if err != nil {
			return 0, 0
		}
	}

	return min.UTC().UnixNano(), max.UTC().UnixNano()
}
