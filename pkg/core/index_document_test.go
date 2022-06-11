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

	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zinc/pkg/meta"
)

func TestIndex_CreateUpdateDocument(t *testing.T) {
	type fields struct {
		Name string
	}
	type args struct {
		docID  string
		doc    map[string]interface{}
		update bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Document with generated ID",
			args: args{
				docID: "test1",
				doc: map[string]interface{}{
					"name": "Hello",
				},
				update: false,
			},
		},
		{
			name: "Document with provided ID",
			args: args{
				docID: "test1",
				doc: map[string]interface{}{
					"test": "Hello",
				},
				update: true,
			},
		},
		{
			name: "Document with type conflict",
			args: args{
				docID: "test1",
				doc: map[string]interface{}{
					"test": true,
				},
				update: true,
			},
			wantErr: true,
		},
	}

	indexName := "TestDocument.index_1"
	var index *Index
	var err error
	t.Run("prepare", func(t *testing.T) {
		index, err = NewIndex(indexName, "disk", nil)
		assert.NoError(t, err)
		assert.NotNil(t, index)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := index.CreateDocument(tt.args.docID, tt.args.doc, tt.args.update)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			query := &meta.ZincQuery{
				Query: &meta.Query{
					Match: map[string]*meta.MatchQuery{
						"_all": {
							Query: "Hello",
						},
					},
				},
			}
			res, err := index.Search(query)
			assert.NoError(t, err)
			assert.Equal(t, 1, res.Hits.Total.Value)
		})
	}

	t.Run("cleanup", func(t *testing.T) {
		err = DeleteIndex(indexName)
		assert.NoError(t, err)
	})
}

func TestIndex_UpdateDocument(t *testing.T) {
	type args struct {
		docID string
		doc   map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "update",
			args: args{
				docID: "1",
				doc: map[string]interface{}{
					"name": "HelloUpdate",
					"time": float64(1579098983),
				},
			},
			wantErr: false,
		},
	}

	var index *Index
	var err error
	t.Run("prepare", func(t *testing.T) {
		index, err = NewIndex("TestIndex_UpdateDocument.index_1", "disk", nil)
		assert.NoError(t, err)
		assert.NotNil(t, index)

		err = StoreIndex(index)
		assert.NoError(t, err)
		prop := meta.NewProperty("date")
		index.Mappings.SetProperty("time", prop)

		err = index.CreateDocument("1", map[string]interface{}{
			"name": "Hello",
			"time": float64(1579098983),
		}, false)
		assert.NoError(t, err)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := index.UpdateDocument(tt.args.docID, tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("Index.UpdateDocument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
