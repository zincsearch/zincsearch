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
	"fmt"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis/analyzer"
	querystr "github.com/blugelabs/query_string"
	v1 "github.com/zinclabs/zinc/pkg/meta/v1"
)

func QueryStringQuery(iQuery *v1.ZincQuery) (bluge.SearchRequest, error) {
	options := querystr.DefaultOptions()
	options.WithDefaultAnalyzer(analyzer.NewStandardAnalyzer())
	userQuery, err := querystr.ParseQueryString(iQuery.Query.Term, options)
	if err != nil {
		return nil, fmt.Errorf("error parsing query string '%s': %s", iQuery.Query.Term, err.Error())
	}

	dateQuery := bluge.NewDateRangeQuery(iQuery.Query.StartTime, iQuery.Query.EndTime).SetField("@timestamp")
	finalQuery := bluge.NewBooleanQuery().AddMust(dateQuery).AddMust(userQuery)

	// sortFields := []string{"-@timestamp"} // adding a - (minus) before the field name will sort the field in descending order

	searchRequest := buildRequest(iQuery, finalQuery)

	return searchRequest, nil

}
