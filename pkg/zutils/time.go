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
	"strconv"
	"strings"
	"time"
)

func ParseDuration(s string) (time.Duration, error) {
	d, err := time.ParseDuration(s)
	if err == nil {
		return d, nil
	}
	if !strings.HasSuffix(s, "d") {
		return 0, err
	}

	h := strings.TrimSuffix(s, "d")
	hour, _ := strconv.Atoi(h)
	d = time.Hour * time.Duration(hour) * 24
	return d, nil
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
