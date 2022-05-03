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
	"os"
	"strconv"
	"strings"
)

// GetEnv returns the value of the environment variable named by the key and returns the default value if the environment variable is not set.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvToLower(key, fallback string) string {
	return strings.ToLower(GetEnv(key, fallback))
}

func GetEnvToUpper(key, fallback string) string {
	return strings.ToUpper(GetEnv(key, fallback))
}

func GetEnvToBool(key, fallback string) bool {
	enabled := false
	if v := GetEnv(key, fallback); v != "" {
		enabled, _ = strconv.ParseBool(v)
	}
	return enabled
}
