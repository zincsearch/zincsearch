package zutils

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zincsearch/zincsearch/test/utils"
)

func TestGetRenderer(t *testing.T) {
	type args struct {
		qparams map[string]string
	}
	type want struct {
		pretty bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "pretty without value should use IndentedJSON",
			args: args{
				qparams: map[string]string{"pretty": ""},
			},
			want: want{pretty: true},
		},
		{
			name: "param pretty with true value should use IndentedJSON",
			args: args{
				qparams: map[string]string{"pretty": "true"},
			},
			want: want{pretty: true},
		},
		{
			name: "param pretty with true value should be case insensitive",
			args: args{
				qparams: map[string]string{"pretty": "TrUe"},
			},
			want: want{pretty: true},
		},
		{
			name: "param pretty with any other value then empty or true should use JSON",
			args: args{
				qparams: map[string]string{"pretty": "No"},
			},
			want: want{pretty: false},
		},
		{
			name: "param pretty missing should use JSON",
			args: args{
				qparams: map[string]string{},
			},
			want: want{pretty: false},
		},
	}
	type Response struct {
		Name string `json:"name"`
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := utils.NewGinContext()
			utils.SetGinRequestURL(c, "/", tt.args.qparams)
			GinRenderJSON(c, http.StatusOK, Response{Name: "zinc"})
			if tt.want.pretty {
				assert.Contains(t, w.Body.String(), "\n")
			} else {
				assert.NotContains(t, w.Body.String(), "\n")
			}
		})
	}
}
