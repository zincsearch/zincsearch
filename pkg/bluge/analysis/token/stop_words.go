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

package token

import (
	"github.com/blugelabs/bluge/analysis"
	"github.com/blugelabs/bluge/analysis/lang/ar"
	"github.com/blugelabs/bluge/analysis/lang/bg"
	"github.com/blugelabs/bluge/analysis/lang/ca"
	"github.com/blugelabs/bluge/analysis/lang/ckb"
	"github.com/blugelabs/bluge/analysis/lang/cs"
	"github.com/blugelabs/bluge/analysis/lang/da"
	"github.com/blugelabs/bluge/analysis/lang/de"
	"github.com/blugelabs/bluge/analysis/lang/el"
	"github.com/blugelabs/bluge/analysis/lang/en"
	"github.com/blugelabs/bluge/analysis/lang/es"
	"github.com/blugelabs/bluge/analysis/lang/eu"
	"github.com/blugelabs/bluge/analysis/lang/fa"
	"github.com/blugelabs/bluge/analysis/lang/fi"
	"github.com/blugelabs/bluge/analysis/lang/fr"
	"github.com/blugelabs/bluge/analysis/lang/gl"
	"github.com/blugelabs/bluge/analysis/lang/hi"
	"github.com/blugelabs/bluge/analysis/lang/hu"
	"github.com/blugelabs/bluge/analysis/lang/hy"
	"github.com/blugelabs/bluge/analysis/lang/id"
	"github.com/blugelabs/bluge/analysis/lang/it"
	"github.com/blugelabs/bluge/analysis/lang/nl"
	"github.com/blugelabs/bluge/analysis/lang/no"
	"github.com/blugelabs/bluge/analysis/lang/pt"
	"github.com/blugelabs/bluge/analysis/lang/ro"
	"github.com/blugelabs/bluge/analysis/lang/ru"
	"github.com/blugelabs/bluge/analysis/lang/sv"
	"github.com/blugelabs/bluge/analysis/lang/tr"

	"github.com/zincsearch/zincsearch/pkg/bluge/analysis/lang/bn"
	"github.com/zincsearch/zincsearch/pkg/bluge/analysis/lang/br"
	"github.com/zincsearch/zincsearch/pkg/bluge/analysis/lang/et"
	"github.com/zincsearch/zincsearch/pkg/bluge/analysis/lang/lv"
	"github.com/zincsearch/zincsearch/pkg/bluge/analysis/lang/th"
)

func StopWords(stopwords []string) analysis.TokenMap {
	if len(stopwords) == 0 {
		stopwords = []string{"_english_"}
	}

	rv := analysis.NewTokenMap()
	for _, word := range stopwords {
		if ok := loadLanguageStopWords(&rv, word); !ok {
			rv.AddToken(word)
		}
	}

	return rv
}

func loadLanguageStopWords(rv *analysis.TokenMap, language string) bool {
	var dict analysis.TokenMap
	switch language {
	case "_ar_", "_arabic_":
		dict = ar.StopWords()
	case "_bg_", "_bulgarian_":
		dict = bg.StopWords()
	case "_bn_", "_bengali_":
		dict = bn.StopWords()
	case "_br_", "_brazilian_": // _brazilian_ (Brazilian Portuguese)
		dict = br.StopWords()
	case "_ca_", "catalan_":
		dict = ca.StopWords()
	case "_cjk_": // _cjk_ (Chinese, Japanese, and Korean)
		// none
	case "_ckb_", "_sorani_":
		dict = ckb.StopWords()
	case "_cs_", "_czech_":
		dict = cs.StopWords()
	case "_da_", "_danish_":
		dict = da.StopWords()
	case "_de_", "_german_":
		dict = de.StopWords()
	case "_el_", "_greek_":
		dict = el.StopWords()
	case "_en_", "_english_":
		dict = en.StopWords()
	case "_es_", "_spanish_":
		dict = es.StopWords()
	case "_et_", "_estonian_":
		dict = et.StopWords()
	case "_eu_", "_basque_":
		dict = eu.StopWords()
	case "_fa_", "_persian_":
		dict = fa.StopWords()
	case "_fi_", "_finnish_":
		dict = fi.StopWords()
	case "_fr_", "_french_":
		dict = fr.StopWords()
	case "_ga_", "_irish_":
		dict = fa.StopWords()
	case "_gl_", "_galician_":
		dict = gl.StopWords()
	case "_hi_", "_hindi_":
		dict = hi.StopWords()
	case "_hu_", "_hungarian_":
		dict = hu.StopWords()
	case "_hy_", "_armenian_":
		dict = hy.StopWords()
	case "_id_", "_indonesian_":
		dict = id.StopWords()
	case "_it_", "_italian_":
		dict = it.StopWords()
	case "_lv", "_latvian_":
		dict = lv.StopWords()
	case "_nl_", "_dutch_":
		dict = nl.StopWords()
	case "_no_", "_norwegian_":
		dict = no.StopWords()
	case "_pt_", "_portuguese_":
		dict = pt.StopWords()
	case "_ro_", "_romanian_":
		dict = ro.StopWords()
	case "_ru_", "_russian_":
		dict = ru.StopWords()
	case "_sv_", "_swedish_":
		dict = sv.StopWords()
	case "_tr_", "_turkish_":
		dict = tr.StopWords()
	case "_th_", "_thai_":
		dict = th.StopWords()
	default:
		return false
	}

	for token := range dict {
		rv.AddToken(token)
	}

	return true
}
