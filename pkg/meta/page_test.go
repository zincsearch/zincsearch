package meta

import (
	"reflect"
	"testing"

	"github.com/zinclabs/zinc/test/utils"
)

func TestNewPage(t *testing.T) {
	type args struct {
		pageNum  string
		pageSize string
	}
	tests := []struct {
		name string
		args args
		want *Page
	}{
		{
			name: "normal",
			args: args{
				pageNum:  "1",
				pageSize: "10",
			},
			want: &Page{
				PageNum:  1,
				PageSize: 10,
			},
		},
		{
			name: "zero",
			args: args{
				pageNum:  "0",
				pageSize: "0",
			},
			want: &Page{
				PageNum:  0,
				PageSize: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := utils.NewGinContext()
			params := map[string]string{
				"page_num":  tt.args.pageNum,
				"page_size": tt.args.pageSize,
			}
			utils.SetGinRequestURL(c, "/", params)
			got := NewPage(c)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPage_GetStartEndIndex(t *testing.T) {
	type fields struct {
		PageNum  int64
		PageSize int64
		Total    int64
	}
	tests := []struct {
		name           string
		fields         fields
		wantStartIndex int64
		wantEndIndex   int64
	}{
		{
			name: "normal",
			fields: fields{
				PageNum:  0,
				PageSize: 10,
				Total:    10,
			},
			wantStartIndex: 0,
			wantEndIndex:   10,
		},
		{
			name: "zero page size",
			fields: fields{
				PageNum:  0,
				PageSize: 0,
				Total:    10,
			},
			wantStartIndex: 0,
			wantEndIndex:   10,
		},
		{
			name: "zero page num",
			fields: fields{
				PageNum:  0,
				PageSize: 0,
				Total:    10,
			},
			wantStartIndex: 0,
			wantEndIndex:   10,
		},
		{
			name: "over total",
			fields: fields{
				PageNum:  2,
				PageSize: 10,
				Total:    10,
			},
			wantStartIndex: 10,
			wantEndIndex:   10,
		},
		{
			name: "over total",
			fields: fields{
				PageNum:  3,
				PageSize: 5,
				Total:    18,
			},
			wantStartIndex: 15,
			wantEndIndex:   18,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Page{
				PageNum:  tt.fields.PageNum,
				PageSize: tt.fields.PageSize,
				Total:    tt.fields.Total,
			}
			gotStartIndex, gotEndIndex := e.GetStartEndIndex()
			if gotStartIndex != tt.wantStartIndex {
				t.Errorf("Page.GetStartEndIndex() gotStartIndex = %v, want %v", gotStartIndex, tt.wantStartIndex)
			}
			if gotEndIndex != tt.wantEndIndex {
				t.Errorf("Page.GetStartEndIndex() gotEndIndex = %v, want %v", gotEndIndex, tt.wantEndIndex)
			}
		})
	}
}
