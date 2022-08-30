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
	"github.com/spf13/cast"
	"strconv"
)

func ToString(v interface{}) (string, error) {
	return cast.ToStringE(v)
}

func ToFloat64(v interface{}) (float64, error) {
	return cast.ToFloat64E(v)
}

func ToUint64(v interface{}) (uint64, error) {
	return cast.ToUint64E(v)
}

func ToInt(v interface{}) (int, error) {
	return cast.ToIntE(v)
}

func ToBool(v interface{}) (bool, error) {
	switch v := v.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	case float64:
		return v != 0, nil
	case uint64:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case int:
		return v != 0, nil
	default:
		return false, fmt.Errorf("ToInt: unknown supported type %T", v)
	}
}
