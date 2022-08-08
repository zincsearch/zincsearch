package index

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/test/utils"
)

func TestAddOrRemoveESAlias(t *testing.T) {
	indexName := "TestAddOrRemoveESAlias.index_1"
	type args struct {
		data   string
		params map[string]string
		result string
	}
	tests := []struct {
		name        string
		nFn         func(*core.Index)
		args        args
		wantErr     bool
		wantCode    int
		wantAliases []string
	}{
		{
			name: "should_add_es_alias",
			args: args{
				data:   `{"actions": [{"add": {"index": "TestAddOrRemoveESAlias.index_1","alias": "test_alias_1"}}]}`,
				result: `{"acknowledged":true}`,
			},
			wantCode:    http.StatusOK,
			wantErr:     false,
			wantAliases: []string{"test_alias_1"},
		},
		{
			name: "should_add_es_alias_with_wildcard_index_name",
			args: args{
				data:   `{"actions": [{"add": {"index": "TestAdd*.index_1","alias": "test_alias_2"}}]}`,
				result: `{"acknowledged":true}`,
			},
			wantCode:    http.StatusOK,
			wantErr:     false,
			wantAliases: []string{"test_alias_2"},
		},
		{
			name: "should_add_es_alias_with_aliases_field",
			args: args{
				data:   `{"actions": [{"add": {"index": "TestAddOrRemoveESAlias.index_1","aliases": ["test_alias_3","test_alias_4"]}}]}`,
				result: `{"acknowledged":true}`,
			},
			nFn: func(index *core.Index) {
				index.AddAliases([]string{"existing_alias"})
			},
			wantCode:    http.StatusOK,
			wantErr:     false,
			wantAliases: []string{"existing_alias", "test_alias_3", "test_alias_4"},
		},
		{
			name: "should_add_es_alias_with_indices_field",
			args: args{
				data:   `{"actions": [{"add": {"indices": ["TestAddOrRemoveESAlias.index_1"],"aliases": ["test_alias_3","test_alias_4"]}}]}`,
				result: `{"acknowledged":true}`,
			},
			nFn: func(index *core.Index) {
				index.AddAliases([]string{"existing_alias"})
			},
			wantCode:    http.StatusOK,
			wantErr:     false,
			wantAliases: []string{"existing_alias", "test_alias_3", "test_alias_4"},
		},
		{
			name: "should_remove_es_alias",
			args: args{
				data:   `{"actions": [{"remove": {"index": "TestAddOrRemoveESAlias.index_1","alias": "existing_alias_2"}}]}`,
				result: `{"acknowledged":true}`,
			},
			nFn: func(index *core.Index) {
				index.AddAliases([]string{"existing_alias_1", "existing_alias_2"})
			},
			wantCode:    http.StatusOK,
			wantErr:     false,
			wantAliases: []string{"existing_alias_1"},
		},
		{
			name: "should_remove_es_alias_with_aliases_field",
			args: args{
				data:   `{"actions": [{"remove": {"index": "TestAddOrRemoveESAlias.index_1","aliases": ["existing_alias_1","existing_alias_2"]}}]}`,
				result: `{"acknowledged":true}`,
			},
			nFn: func(index *core.Index) {
				index.AddAliases([]string{"existing_alias_1", "existing_alias_2", "existing_alias_3"})
			},
			wantCode:    http.StatusOK,
			wantErr:     false,
			wantAliases: []string{"existing_alias_3"},
		},
		{
			name: "should_remove_es_alias_with_indices_field",
			args: args{
				data:   `{"actions": [{"remove": {"indices": ["TestAddOrRemoveESAlias.index_1"],"aliases": ["existing_alias_1","existing_alias_2"]}}]}`,
				result: `{"acknowledged":true}`,
			},
			nFn: func(index *core.Index) {
				index.AddAliases([]string{"existing_alias_1", "existing_alias_2", "existing_alias_3"})
			},
			wantCode:    http.StatusOK,
			wantErr:     false,
			wantAliases: []string{"existing_alias_3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, closeFn := newIndex(t, indexName)
			defer closeFn()

			if tt.nFn != nil {
				tt.nFn(index)
			}

			c, w := utils.NewGinContext()
			utils.SetGinRequestData(c, tt.args.data)
			utils.SetGinRequestParams(c, tt.args.params)
			AddOrRemoveESAlias(c)

			require.Equal(t, tt.wantCode, w.Code)
			require.Equal(t, tt.args.result, w.Body.String())
			require.Equal(t, tt.wantAliases, index.GetAliases())
		})
	}
}

func TestGetESAliases(t *testing.T) {
	indexName := "TestAddOrRemoveESAlias.index_1"
	type args struct {
		data   string
		result string
	}
	tests := []struct {
		name     string
		nFn      func(*core.Index)
		params   map[string]string
		args     args
		wantErr  bool
		wantCode int
	}{
		{
			name: "should_get_es_alias",
			args: args{
				data:   `{"actions": [{"add": {"index": "TestAddOrRemoveESAlias.index_1","alias": "test_alias_1"}}]}`,
				result: `{"TestAddOrRemoveESAlias.index_1":{"aliases":{"existing_alias_1":{},"existing_alias_2":{},"existing_alias_3":{}}}}`,
			},
			nFn: func(index *core.Index) {
				index.AddAliases([]string{"existing_alias_1", "existing_alias_2", "existing_alias_3"})
			},
			wantCode: http.StatusOK,
			wantErr:  false,
		},
		{
			name: "should_get_es_alias_of_target_index",
			args: args{
				data:   `{"actions": [{"add": {"index": "TestAddOrRemoveESAlias.index_1","alias": "test_alias_1"}}]}`,
				result: `{"TestAddOrRemoveESAlias.index_1":{"aliases":{"existing_alias_1":{},"existing_alias_2":{},"existing_alias_3":{}}}}`,
			},
			params: map[string]string{
				"target": "TestAddOrRemoveESAlias.index_1",
			},
			nFn: func(index *core.Index) {
				index.AddAliases([]string{"existing_alias_1", "existing_alias_2", "existing_alias_3"})
			},
			wantCode: http.StatusOK,
			wantErr:  false,
		},
		{
			name: "should_get_es_alias_of_target_alias",
			args: args{
				data:   `{"actions": [{"add": {"index": "TestAddOrRemoveESAlias.index_1","alias": "test_alias_1"}}]}`,
				result: `{"TestAddOrRemoveESAlias.index_1":{"aliases":{"existing_alias_1":{}}}}`,
			},
			params: map[string]string{
				"target_alias": "existing_alias_1",
			},
			nFn: func(index *core.Index) {
				index.AddAliases([]string{"existing_alias_1", "existing_alias_2", "existing_alias_3"})
			},
			wantCode: http.StatusOK,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, closeFn := newIndex(t, indexName)
			defer closeFn()

			if tt.nFn != nil {
				tt.nFn(index)
			}

			c, w := utils.NewGinContext()
			utils.SetGinRequestData(c, tt.args.data)
			utils.SetGinRequestParams(c, tt.params)
			GetESAliases(c)

			require.Equal(t, tt.wantCode, w.Code)
			require.Equal(t, tt.args.result, w.Body.String())
		})
	}
}

func newIndex(t *testing.T, indexName string) (*core.Index, func()) {
	index, err := core.NewIndex(indexName, "disk", 2)
	require.NoError(t, err)

	err = core.StoreIndex(index)
	require.NoError(t, err)

	return index, func() {
		require.NoError(t, core.DeleteIndex(indexName))
	}
}
