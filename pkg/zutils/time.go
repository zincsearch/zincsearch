// Copyright 2022 Zinc Labs Inc. and Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	if d.Hours() >= 24*30*12 {
		return strconv.FormatInt(int64(d.Hours())/24/30/12, 10) + "y"
	}
	if d.Hours() >= 24*30 {
		return strconv.FormatInt(int64(d.Hours())/24/30, 10) + "M"
	}
	if d.Hours() >= 24 {
		return strconv.FormatInt(int64(d.Hours())/24, 10) + "d"
	}
	if d.Hours() >= 1 {
		return strconv.FormatInt(int64(d.Hours()), 10) + "h"
	}
	if d.Minutes() >= 1 {
		return strconv.FormatInt(int64(d.Minutes()), 10) + "m"
	}
	return strconv.FormatInt(int64(d.Seconds()), 10) + "s"
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
