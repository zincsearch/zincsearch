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

// Package testutil holds helpers shared by the analysis test packages.
package testutil

import (
	"github.com/vcaesar/riot/analysis"
	"github.com/vcaesar/riot/analysis/tokenizer"
)

// Tokens splits text on whitespace into a TokenStream.
func Tokens(text string) analysis.TokenStream {
	return tokenizer.NewWhitespaceTokenizer().Tokenize([]byte(text))
}

// Terms returns the term text of each token in ts.
func Terms(ts analysis.TokenStream) []string {
	out := make([]string, 0, len(ts))
	for _, t := range ts {
		out = append(out, string(t.Term))
	}
	return out
}
