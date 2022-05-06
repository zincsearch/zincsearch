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
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	v1 "github.com/zinclabs/zinc/pkg/meta/v1"
)

func TestIndex_UpdateDocument(t *testing.T) {
	type fields struct {
		Name string
	}
	type args struct {
		docID    string
		doc      map[string]interface{}
		mintedID bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "UpdateDocument with generated ID",
			args: args{
				docID: "test1",
				doc: map[string]interface{}{
					"name": "Hello",
				},
				mintedID: true,
			},
		},
		{
			name: "UpdateDocument with provided ID",
			args: args{
				docID: "test2",
				doc: map[string]interface{}{
					"test": "Hello",
				},
				mintedID: false,
			},
		},
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			rand.Seed(time.Now().UnixNano())
			id := rand.Intn(1000)
			indexName := "TestUpdateDocument.index_" + strconv.Itoa(id)

			index, _ := NewIndex(indexName, "disk", UseNewIndexMeta, nil)

			err := index.UpdateDocument(tt.args.docID, tt.args.doc, tt.args.mintedID)

			assert.Nil(t, err)

			if err == nil {
				query := &v1.ZincQuery{
					SearchType: "match",
					Query: v1.QueryParams{
						Term: "Hello",
					},
				}
				res, err := index.Search(query)
				assert.Nil(t, err)
				assert.Equal(t, 1, res.Hits.Total.Value)
			}
		})
	}

	os.RemoveAll("data") // cleanup data folder
}
