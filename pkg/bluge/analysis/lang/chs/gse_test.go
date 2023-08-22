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
	"os"
	"strings"
	"testing"

	"github.com/blugelabs/bluge/analysis"
	"github.com/stretchr/testify/assert"

	"github.com/zincsearch/zincsearch/pkg/config"
)

func TestLoadDict(t *testing.T) {
	type args struct {
		enable     bool
		enableStop bool
		embed      string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "enable=false,embed=small",
			args: args{
				enable:     false,
				enableStop: false,
				embed:      "SMALL",
			},
		},
		{
			name: "enable=true,embed=small",
			args: args{
				enable:     true,
				enableStop: true,
				embed:      "SMALL",
			},
		},
		{
			name: "enable=true,embed=big",
			args: args{
				enable:     true,
				enableStop: true,
				embed:      "BIG",
			},
		},
	}

	t.Run("prepare dict", func(t *testing.T) {
		_ = os.Mkdir("data", 0755)
		config.Global.Plugin.GSE.DictPath = "./data"
		err := writeFile("./data/user.txt", "你若安好便是晴天 100 n\n")
		assert.NoError(t, err)
		err = writeFile("./data/stop.txt", "你好\n")
		assert.NoError(t, err)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadDict(tt.args.enable, tt.args.enableStop, tt.args.embed)
		})
	}

	t.Run("clean dict", func(t *testing.T) {
		os.RemoveAll("data")
	})
}

func TestNewGseStandardAnalyzer(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "default",
			text: "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄制作的科幻片.",
			want: "[复仇者 联盟 3 无限 战争 全片 使用 imax 摄影机 拍摄 制作 科幻片]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGseStandardAnalyzer().Analyze([]byte(tt.text))
			assert.Equal(t, tt.want, collectToken(got))
		})
	}
}

func TestNewGseSearchAnalyzer(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "default",
			text: "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄制作的科幻片.",
			want: "[复仇 仇者 复仇者 联盟 3 无限 战争 全片 使用 imax 摄影 摄影机 拍摄 制作 科幻 科幻片]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGseSearchAnalyzer().Analyze([]byte(tt.text))
			assert.Equal(t, tt.want, collectToken(got))
		})
	}
}

func TestNewGseStandardTokenizer(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "default",
			text: "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄制作的科幻片.",
			want: "[《 复仇者 联盟 3 ： 无限 战争 》 是 全片 使用 imax 摄影机 拍摄 制作 的 科幻片 .]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGseStandardTokenizer().Tokenize([]byte(tt.text))
			assert.Equal(t, tt.want, collectToken(got))
		})
	}
}

func TestNewGseSearchTokenizer(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "default",
			text: "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄制作的科幻片.",
			want: "[《 复仇 仇者 复仇者 联盟 3 ： 无限 战争 》 是 全片 使用 imax 摄影 摄影机 拍摄 制作 的 科幻 科幻片 .]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGseSearchTokenizer().Tokenize([]byte(tt.text))
			assert.Equal(t, tt.want, collectToken(got))
		})
	}
}

func TestNewGseStopTokenFilter(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "default",
			text: "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄制作的科幻片.",
			want: "[复仇 仇者 复仇者 联盟 3 无限 战争 全片 使用 imax 摄影 摄影机 拍摄 制作 科幻 科幻片]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGseSearchTokenizer().Tokenize([]byte(tt.text))
			got = NewGseStopTokenFilter().Filter(got)
			assert.Equal(t, tt.want, collectToken(got))
		})
	}
}

func writeFile(path string, content string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(content))
	return err
}

func collectToken(tokens analysis.TokenStream) string {
	str := make([]string, 0, len(tokens))
	for _, token := range tokens {
		str = append(str, string(token.Term))
	}
	return "[" + strings.Join(str, " ") + "]"
}
