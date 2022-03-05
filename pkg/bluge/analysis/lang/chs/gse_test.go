package chs

import (
	"strings"
	"testing"

	"github.com/blugelabs/bluge/analysis"
)

func TestGseAnalyzer(t *testing.T) {
	seg.LoadDictEmbed("zh_s")
	seg.LoadStopEmbed()

	text := "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄制作的科幻片."
	standard := NewGseStandardAnalyzer()
	search := NewGseSearchAnalyzer()

	tokens1 := standard.Analyze([]byte(text))
	result1 := "[复仇者 联盟 3 无限 战争 全片 使用 imax 摄影机 拍摄 制作 科幻片]"
	if result1 != collectToken(tokens1) {
		t.Error(collectToken(tokens1), "should equal", result1)
	}

	tokens2 := search.Analyze([]byte(text))
	result2 := "[复仇 仇者 复仇者 联盟 3 无限 战争 全片 使用 imax 摄影 摄影机 拍摄 制作 科幻 科幻片]"
	if result2 != collectToken(tokens2) {
		t.Error(collectToken(tokens2), "should equal", result2)
	}
}

func collectToken(tokens analysis.TokenStream) string {
	str := make([]string, 0, len(tokens))
	for _, token := range tokens {
		str = append(str, string(token.Term))
	}
	return "[" + strings.Join(str, " ") + "]"
}
