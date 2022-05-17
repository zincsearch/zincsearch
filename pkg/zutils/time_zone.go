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
	"strings"
	"time"
)

func ParseTimeZone(name string) (*time.Location, error) {
	offset := 0
	ln := len(name)
	if ln > 0 && (name[0] == '+' || name[0] == '-') {
		if ln >= 3 {
			offset = 60 * 60 * StringToInt(name[1:3]) // +08:00  +0800
		}
		if ln == 5 {
			offset += 60 * StringToInt(name[3:5]) // +0830
		}
		if ln == 6 && name[3] == ':' {
			offset += 60 * StringToInt(name[4:6]) // +08:30
		}
		if name[0] == '-' { // -01:00
			offset = -offset
		}
		return time.FixedZone(name, offset), nil
	}

	upperName := strings.ToUpper(name)
	if upperName == "" || upperName == "UTC" {
		return time.UTC, nil
	}
	if upperName == "LOCAL" {
		return time.Local, nil
	}

	return time.LoadLocation(name)
}
