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

package meta

type Analyzer struct {
	CharFilter  []string `json:"char_filter,omitempty"`
	Tokenizer   string   `json:"tokenizer,omitempty"`
	TokenFilter []string `json:"token_filter,omitempty"`
	Filter      []string `json:"filter,omitempty"` // compatibility with es, alias for TokenFilter

	// options for compatible
	Type      string   `json:"type,omitempty"`
	Pattern   string   `json:"pattern,omitempty"`   // for type=pattern
	Lowercase bool     `json:"lowercase,omitempty"` // for type=pattern
	Stopwords []string `json:"stopwords,omitempty"` // for type=pattern,standard,stop
}

type Tokenizer struct {
	Type string `json:"type"`
}
type TokenFilter struct {
	Type string `json:"type"`
}
