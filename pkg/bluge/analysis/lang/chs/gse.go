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

package chs

import (
	"strings"

	"github.com/blugelabs/bluge/analysis"
	"github.com/go-ego/gse"
	"github.com/rs/zerolog/log"

	"github.com/zincsearch/zincsearch/pkg/bluge/analysis/lang/chs/analyzer"
	"github.com/zincsearch/zincsearch/pkg/bluge/analysis/lang/chs/token"
	"github.com/zincsearch/zincsearch/pkg/bluge/analysis/lang/chs/tokenizer"
	"github.com/zincsearch/zincsearch/pkg/config"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func NewGseStandardAnalyzer() *analysis.Analyzer {
	return analyzer.NewStandardAnalyzer(seg)
}

func NewGseSearchAnalyzer() *analysis.Analyzer {
	return analyzer.NewSearchAnalyzer(seg)
}

func NewGseStandardTokenizer() analysis.Tokenizer {
	return tokenizer.NewStandardTokenizer(seg)
}

func NewGseSearchTokenizer() analysis.Tokenizer {
	return tokenizer.NewSearchTokenizer(seg)
}

func NewGseStopTokenFilter() analysis.TokenFilter {
	return token.NewStopTokenFilter(seg, nil)
}

var seg *gse.Segmenter

func init() {
	seg = new(gse.Segmenter)
	enable := config.Global.Plugin.GSE.Enable         // true / false
	enableStop := config.Global.Plugin.GSE.EnableStop // true / false
	embed := config.Global.Plugin.GSE.DictEmbed       // small / big
	embed = strings.ToUpper(embed)
	loadDict(enable, enableStop, embed)
}

func loadDict(enable, enableStop bool, embed string) {
	if enable {
		// load default dict
		if embed == "BIG" {
			_ = seg.LoadDictEmbed("zh_s")
			if enableStop {
				_ = seg.LoadStopEmbed()
			}
		} else {
			_ = seg.LoadDictStr(_dictCHS)
			if enableStop {
				_ = seg.LoadStopStr(_dictStop)
			}
		}
	} else {
		// load empty dict
		_ = seg.LoadDictStr(`zinc`)
		if enableStop {
			_ = seg.LoadStopStr(_dictStop)
		}
	}

	seg.Load = true
	seg.SkipLog = true
	if !enable {
		return
	}

	// load user dict
	dataPath := config.Global.Plugin.GSE.DictPath
	userDict := dataPath + "/user.txt"
	log.Info().Msgf("Loading  Gse user dict... %s", userDict)
	if ok, _ := zutils.IsExist(userDict); ok {
		_ = seg.LoadDict(userDict)
	}
	stopDict := dataPath + "/stop.txt"
	log.Info().Msgf("Loading  Gse user stop... %s", stopDict)
	if ok, _ := zutils.IsExist(stopDict); ok {
		_ = seg.LoadStop(stopDict)
	}
}
