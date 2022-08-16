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

package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAliasList_AddIndexesToAlias(t *testing.T) {
	type args struct {
		alias   string
		indexes []string
	}
	tests := []struct {
		name        string
		nFn         func(al *AliasList)
		args        args
		wantErr     bool
		wantIndexes []string
	}{
		{
			name: "should_add_indexes_to_alias",
			nFn: func(al *AliasList) {
				al.Aliases["alias_1"] = append(al.Aliases["alias_1"], "index_0")
			},
			args: args{
				alias:   "alias_1",
				indexes: []string{"index_1", "index_2"},
			},
			wantErr:     false,
			wantIndexes: []string{"index_0", "index_1", "index_2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			al := NewAliasList()

			if tt.nFn != nil {
				tt.nFn(al)
			}

			err := al.AddIndexesToAlias(tt.args.alias, tt.args.indexes)
			if tt.wantErr {
				require.NotNil(t, err)
				return
			}

			require.Equal(t, tt.wantIndexes, al.Aliases[tt.args.alias])
		})
	}
}

func TestAliasList_RemoveIndexesFromAlias(t *testing.T) {
	type args struct {
		alias         string
		removeIndexes []string
	}
	tests := []struct {
		name        string
		nFn         func(al *AliasList)
		args        args
		wantErr     bool
		wantIndexes []string
	}{
		{
			name: "should_remove_indexes_from_alias",
			nFn: func(al *AliasList) {
				al.Aliases["alias_1"] = append(al.Aliases["alias_1"], "index_0", "index_1", "index_2", "index_3")
			},
			args: args{
				alias:         "alias_1",
				removeIndexes: []string{"index_1", "index_3"},
			},
			wantErr:     false,
			wantIndexes: []string{"index_0", "index_2"},
		},
		{
			name: "should_not_find_alias",
			nFn:  nil,
			args: args{
				alias:         "alias_1",
				removeIndexes: []string{"index_1"},
			},
			wantErr:     false,
			wantIndexes: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			al := NewAliasList()

			if tt.nFn != nil {
				tt.nFn(al)
			}

			err := al.RemoveIndexesFromAlias(tt.args.alias, tt.args.removeIndexes)
			if tt.wantErr {
				require.NotNil(t, err)
				return
			}

			require.Equal(t, tt.wantIndexes, al.Aliases[tt.args.alias])
		})
	}
}

func TestAliasList_GetIndexesForAlias(t *testing.T) {
	type args struct {
		aliasName string
	}

	tests := []struct {
		name        string
		nFn         func(al *AliasList)
		args        args
		wantOk      bool
		wantIndexes []string
	}{
		{
			name: "should_get_indexes_for_alias",
			nFn: func(al *AliasList) {
				al.Aliases["alias_1"] = append(al.Aliases["alias_1"], "index_0", "index_1", "index_2")
			},
			args: args{
				aliasName: "alias_1",
			},
			wantOk:      true,
			wantIndexes: []string{"index_0", "index_1", "index_2"},
		},
		{
			name: "should_get_no_indexes_for_alias",
			nFn:  nil,
			args: args{
				aliasName: "alias_1",
			},
			wantOk:      false,
			wantIndexes: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			al := NewAliasList()

			if tt.nFn != nil {
				tt.nFn(al)
			}

			indexes, ok := al.GetIndexesForAlias(tt.args.aliasName)

			require.Equal(t, tt.wantOk, ok)
			require.Equal(t, tt.wantIndexes, indexes)
		})
	}
}

func TestAliasList_GetAliasesForIndex(t *testing.T) {
	type args struct {
		indexName string
	}

	tests := []struct {
		name        string
		nFn         func(al *AliasList)
		args        args
		wantAliases []string
	}{
		{
			name: "should_get_aliases_for_index",
			nFn: func(al *AliasList) {
				al.Aliases["alias_1"] = append(al.Aliases["alias_1"], "index_0", "index_1", "index_2")
				al.Aliases["alias_2"] = append(al.Aliases["alias_2"], "index_0", "index_1", "index_2")
			},
			args: args{
				indexName: "index_0",
			},
			wantAliases: []string{"alias_1", "alias_2"},
		},
		{
			name: "should_get_no_indexes_for_alias",
			nFn: func(al *AliasList) {
				al.Aliases["alias_1"] = append(al.Aliases["alias_1"], "index_0", "index_1", "index_2")
				al.Aliases["alias_2"] = append(al.Aliases["alias_1"], "index_0", "index_1", "index_2")
			},
			args: args{
				indexName: "index_6",
			},
			wantAliases: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			al := NewAliasList()

			if tt.nFn != nil {
				tt.nFn(al)
			}

			indexes := al.GetAliasesForIndex(tt.args.indexName)
			require.ElementsMatch(t, tt.wantAliases, indexes)
		})
	}
}

func TestAliasList_GetAliasMap(t *testing.T) {
	type args struct {
		targetIndexes []string
		targetAliases []string
	}
	tests := []struct {
		name string
		nFn  func(al *AliasList)
		args args
		want M
	}{
		{
			name: "should_get_alias_map",
			nFn: func(al *AliasList) {
				al.Aliases["alias_1"] = append(al.Aliases["alias_1"], "index_0")
				al.Aliases["alias_2"] = append(al.Aliases["alias_2"], "index_1")
			},
			args: args{
				targetIndexes: nil,
				targetAliases: nil,
			},
			want: M{
				"index_0": M{
					"aliases": M{
						"alias_1": struct{}{},
					},
				},
				"index_1": M{
					"aliases": M{
						"alias_2": struct{}{},
					},
				},
			},
		},
		{
			name: "should_get_alias_map_with_targets",
			nFn: func(al *AliasList) {
				al.Aliases["alias_1"] = append(al.Aliases["alias_1"], "index_0")
				al.Aliases["alias_2"] = append(al.Aliases["alias_2"], "index_1")
			},
			args: args{
				targetIndexes: []string{"index_1"},
				targetAliases: []string{"alias_2"},
			},
			want: M{
				"index_1": M{
					"aliases": M{
						"alias_2": struct{}{},
					},
				},
			},
		},
		{
			name: "should_get_alias_map_with_targets",
			nFn: func(al *AliasList) {
				al.Aliases["alias_1"] = append(al.Aliases["alias_1"], "index_0")
				al.Aliases["alias_2"] = append(al.Aliases["alias_2"], "index_1")
			},
			args: args{
				targetIndexes: []string{"index_0"},
				targetAliases: []string{"alias_1"},
			},
			want: M{
				"index_0": M{
					"aliases": M{
						"alias_1": struct{}{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			al := NewAliasList()

			if tt.nFn != nil {
				tt.nFn(al)
			}

			m := al.GetAliasMap(tt.args.targetIndexes, tt.args.targetAliases)
			require.Equal(t, tt.want, m)
		})
	}
}
