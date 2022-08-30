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

func ParseDuration(s string) (time.Duration, error) {
	d, err := time.ParseDuration(s)
	if err == nil {
		return d, nil
	}

	if strings.HasSuffix(s, "d") {
		h := strings.TrimSuffix(s, "d")
		hour, _ := strconv.Atoi(h)
		d = time.Hour * 24 * time.Duration(hour)
		return d, nil
	}

	dv, err := strconv.ParseInt(s, 10, 64)
	return time.Duration(dv), err
}

func FormatDuration(d time.Duration) string {
	var s string
	if d.Hours() >= 24*30*12 {
		t := int64(d.Hours()) / 24 / 30 / 12
		s += strconv.FormatInt(t, 10) + "y"
		d -= time.Hour * 24 * 30 * 12 * time.Duration(t)
	}
	if d.Hours() >= 24*30 {
		t := int64(d.Hours()) / 24 / 30
		s += strconv.FormatInt(t, 10) + "M"
		d -= time.Hour * 24 * 30 * time.Duration(t)
	}
	if d.Hours() >= 24 {
		t := int64(d.Hours()) / 24
		s += strconv.FormatInt(t, 10) + "d"
		d -= time.Hour * 24 * time.Duration(t)
	}
	if d.Hours() >= 1 {
		t := int64(d.Hours())
		s += strconv.FormatInt(t, 10) + "h"
		d -= time.Hour * time.Duration(t)
	}
	if d.Minutes() >= 1 {
		t := int64(d.Minutes())
		s += strconv.FormatInt(t, 10) + "m"
		d -= time.Minute * time.Duration(t)
	}
	if d > 0 {
		s += strconv.FormatInt(int64(d.Seconds()), 10) + "s"
	}
	return s
}

func Unix(n int64) time.Time {
	if n > 1e18 {
		return time.Unix(0, n)
	}
	if n > 1e15 {
		return time.UnixMicro(n)
	}
	if n > 1e12 {
		return time.UnixMilli(n)
	}
	return time.Unix(n, 0)
}

func ParseTime(value interface{}, format, timeZone string) (time.Time, error) {
	var vInt int64
	var vStr string
	switch v := value.(type) {
	case float64:
		vInt = int64(v)
	case int64:
		vInt = v
	case string:
		vStr = v
	default:
		return time.Time{}, fmt.Errorf("value type of time must be string / float64 / int64")
	}

	if vInt != 0 {
		t := Unix(vInt)
		if t.IsZero() {
			return time.Time{}, fmt.Errorf("time format is [epoch_millis] but the value [%d] is not a valid timestamp", vInt)
		}
		return t, nil
	}

	if vStr == "" {
		return time.Time{}, fmt.Errorf("time value is empty")
	}

	var err error
	timFormat := time.RFC3339
	timZone := time.UTC
	if format != "" {
		timFormat = format
	}
	if timeZone != "" {
		timZone, err = ParseTimeZone(timeZone)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid time zone: %s", timeZone)
		}
	}

	if timFormat == "epoch_millis" {
		v, err := ToInt(vStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("time format is [epoch_millis] but the value [%s] can't convert to int", vStr)
		}
		if t := Unix(int64(v)); t.IsZero() {
			return time.Time{}, fmt.Errorf("time format is [epoch_millis] but the value [%s] is not a valid timestamp", vStr)
		} else {
			return t, nil
		}
	}

	t, err := time.ParseInLocation(timFormat, vStr, timZone)
	if err != nil {
		return time.Time{}, fmt.Errorf("time format is [%s] but the value [%s] parse err: %s", timFormat, vStr, err.Error())
	}
	return t, nil
}
