package analysis

import (
	"fmt"
	"plugin"

	"github.com/blugelabs/bluge/analysis"
)

type Analysis struct {
	Analyzer    map[string]*analysis.Analyzer
	Tokenizer   map[string]analysis.Tokenizer
	CharFilter  map[string]analysis.CharFilter
	TokenFilter map[string]analysis.TokenFilter
}

const prefix = "plugin_"

var store = new(Analysis)

func init() {
	store.Analyzer = make(map[string]*analysis.Analyzer)
	store.Tokenizer = make(map[string]analysis.Tokenizer)
	store.CharFilter = make(map[string]analysis.CharFilter)
	store.TokenFilter = make(map[string]analysis.TokenFilter)
}

func GetAnalyzer(name string) (*analysis.Analyzer, bool) {
	v, ok := store.Analyzer[prefix+name]
	return v, ok
}

func GetTokenizer(name string) (analysis.Tokenizer, bool) {
	v, ok := store.Tokenizer[prefix+name]
	return v, ok
}

func GetCharFilter(name string) (analysis.CharFilter, bool) {
	v, ok := store.CharFilter[prefix+name]
	return v, ok
}

func GetTokenFilter(name string) (analysis.TokenFilter, bool) {
	v, ok := store.TokenFilter[prefix+name]
	return v, ok
}

func Load(p *plugin.Plugin) error {
	loader, err := p.Lookup("Load")
	if err != nil {
		return err
	}
	fn := loader.(func() *Analysis)
	analysis := fn()

	for k, v := range analysis.Analyzer {
		if _, ok := store.Analyzer[prefix+k]; ok {
			return fmt.Errorf("duplicate analyzer name: %s", k)
		}
		store.Analyzer[prefix+k] = v
	}

	for k, v := range analysis.Tokenizer {
		if _, ok := store.Tokenizer[prefix+k]; ok {
			return fmt.Errorf("duplicate tokenizer name: %s", k)
		}
		store.Tokenizer[prefix+k] = v
	}

	for k, v := range analysis.CharFilter {
		if _, ok := store.CharFilter[prefix+k]; ok {
			return fmt.Errorf("duplicate char_filter name: %s", k)
		}
		store.CharFilter[prefix+k] = v
	}

	for k, v := range analysis.TokenFilter {
		if _, ok := store.TokenFilter[prefix+k]; ok {
			return fmt.Errorf("duplicate token_filter name: %s", k)
		}
		store.TokenFilter[prefix+k] = v
	}

	return nil
}
