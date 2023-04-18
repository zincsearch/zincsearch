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

package index

import (
	"fmt"

	"github.com/blugelabs/bluge/analysis"

	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/meta"
	zincanalysis "github.com/zincsearch/zincsearch/pkg/uquery/analysis"
	"github.com/zincsearch/zincsearch/pkg/uquery/mappings"
	"github.com/zincsearch/zincsearch/pkg/zutils/json"
)

func Request(data map[string]interface{}) (*meta.Index, error) {
	if len(data) == 0 {
		return nil, nil
	}

	index := new(meta.Index)

	// parse settings
	var err error
	var analyzers map[string]*analysis.Analyzer
	if v, ok := data["settings"]; ok {
		v, ok := v.(map[string]interface{})
		if !ok {
			return nil, errors.New(errors.ErrorTypeParsingException, "[index] settings should be an object")
		}
		vjson, _ := json.Marshal(v)
		settings := new(meta.IndexSettings)
		if err := json.Unmarshal(vjson, settings); err != nil {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[index] settings parse error: %s", err.Error()))
		}
		if analyzers, err = zincanalysis.RequestAnalyzer(settings.Analysis); err != nil {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[index] settings.analysis parse error: %s", err.Error()))
		}
		if settings != nil && (settings.NumberOfShards > 0 || settings.NumberOfReplicas > 0 || settings.Analysis != nil) {
			index.Settings = settings
		}
	}

	// parse mappings
	if v, ok := data["mappings"]; ok {
		v, ok := v.(map[string]interface{})
		if !ok {
			return nil, errors.New(errors.ErrorTypeParsingException, "[index] mappings should be an object")
		}
		mappings, err := mappings.Request(analyzers, v)
		if err != nil {
			return nil, err
		}
		index.Mappings = mappings
	}

	// parse alias
	// if _, ok := data["alias"]; ok {
	//	 noop
	// }

	return index, nil
}
