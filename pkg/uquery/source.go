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

package uquery

import (
	"github.com/goccy/go-json"

	v1 "github.com/zinclabs/zinc/pkg/meta/v1"
)

func HandleSource(source *v1.Source, data []byte) map[string]interface{} {
	ret := make(map[string]interface{})
	// return empty
	if !source.Enable {
		return ret
	}

	err := json.Unmarshal(data, &ret)
	if err != nil {
		return nil
	}

	// return all fields
	if len(source.Fields) == 0 {
		return ret
	}

	// delete field not in source.Fields
	for field := range ret {
		if _, ok := source.Fields[field]; ok {
			continue
		}
		delete(ret, field)
	}

	return ret
}
