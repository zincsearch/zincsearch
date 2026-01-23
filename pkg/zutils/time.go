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

package zutils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	dayDuration   = time.Hour * 24
	monthDuration = dayDuration * 30
	yearDuration  = monthDuration * 12
)

func ParseDuration(s string) (time.Duration, error) {
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	if strings.HasSuffix(s, "d") {
		h := strings.TrimSuffix(s, "d")
		hour, err := strconv.Atoi(h)
		if err != nil {
			return 0, fmt.Errorf("invalid day format: %s", s)
		}
		return dayDuration * time.Duration(hour), nil
	}

	dv, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: %s", s)
	}
	return time.Duration(dv), nil
}

func FormatDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	var sb strings.Builder
	sb.Grow(16)

	// Process years
	if d >= yearDuration {
		years := int64(d / yearDuration)
		sb.WriteString(strconv.FormatInt(years, 10))
		sb.WriteByte('y')
		d -= yearDuration * time.Duration(years)
	}

	// Process months
	if d >= monthDuration {
		months := int64(d / monthDuration)
		sb.WriteString(strconv.FormatInt(months, 10))
		sb.WriteByte('M')
		d -= monthDuration * time.Duration(months)
	}

	// Process days
	if d >= dayDuration {
		days := int64(d / dayDuration)
		sb.WriteString(strconv.FormatInt(days, 10))
		sb.WriteByte('d')
		d -= dayDuration * time.Duration(days)
	}

	// Process hours
	if d >= time.Hour {
		hours := int64(d / time.Hour)
		sb.WriteString(strconv.FormatInt(hours, 10))
		sb.WriteByte('h')
		d -= time.Hour * time.Duration(hours)
	}

	// Process minutes
	if d >= time.Minute {
		minutes := int64(d / time.Minute)
		sb.WriteString(strconv.FormatInt(minutes, 10))
		sb.WriteByte('m')
		d -= time.Minute * time.Duration(minutes)
	}

	// Process seconds
	if d > 0 {
		seconds := int64(d / time.Second)
		sb.WriteString(strconv.FormatInt(seconds, 10))
		sb.WriteByte('s')
	}

	return sb.String()
}

func Unix(n int64) time.Time {
	switch {
	case n == 0:
		return time.Unix(0, 0)
	case n > 1e18: // Nanoseconds (19+ digits)
		return time.Unix(0, n)
	case n > 1e15: // Microseconds (16-18 digits)
		return time.UnixMicro(n)
	case n > 1e12: // Milliseconds (13-15 digits)
		return time.UnixMilli(n)
	default: // Seconds (1-12 digits)
		return time.Unix(n, 0)
	}
}

func ParseTime(value interface{}, format, timeZone string) (time.Time, error) {
	switch v := value.(type) {
	case float64:
		return parseNumericTime(int64(v))
	case int64:
		return parseNumericTime(v)
	case int:
		return parseNumericTime(int64(v))
	case int32:
		return parseNumericTime(int64(v))
	case string:
		return parseStringTime(v, format, timeZone)
	default:
		return time.Time{}, fmt.Errorf("value type of time must be string or numeric, got %T", value)
	}
}

func parseNumericTime(vInt int64) (time.Time, error) {
	t := Unix(vInt)
	if t.IsZero() && vInt != 0 {
		return time.Time{}, fmt.Errorf("value [%d] is not a valid timestamp", vInt)
	}
	return t, nil
}

func parseStringTime(vStr, format, timeZone string) (time.Time, error) {
	if vStr == "" {
		return time.Time{}, fmt.Errorf("time value is empty")
	}

	if format == "epoch_millis" {
		v, err := strconv.ParseInt(vStr, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("time format is [epoch_millis] but value [%s] cannot be converted to int", vStr)
		}
		return parseNumericTime(v)
	}

	loc, err := parseTimeZone(timeZone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time zone: %s", timeZone)
	}

	timFormat := time.RFC3339
	if format != "" {
		timFormat = format
	}

	t, err := time.ParseInLocation(timFormat, vStr, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("time format [%s] value [%s] parse error: %w", timFormat, vStr, err)
	}
	return t, nil
}

func parseTimeZone(timeZone string) (*time.Location, error) {
	if timeZone == "" {
		return time.UTC, nil
	}

	switch strings.ToUpper(timeZone) {
	case "UTC", "":
		return time.UTC, nil
	case "LOCAL", "SYSTEM":
		return time.Local, nil
	}

	return time.LoadLocation(timeZone)
}
