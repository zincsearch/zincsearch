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
)

func TestDeleteIndex(t *testing.T) {
	var indexName = "TestDeleteIndex.index_1"
	var indexNameS3 = "TestDeleteIndex.index_s3"
	var indexNameMinIO = "TestDeleteIndex.index_minio"
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "exist",
			args: args{
				name: indexName,
			},
			wantErr: false,
		},
		{
			name: "not exist",
			args: args{
				name: "my-index-not-exist",
			},
			wantErr: true,
		},
		{
			name: "s3",
			args: args{
				name: indexNameS3,
			},
			wantErr: false,
		},
		{
			name: "minio",
			args: args{
				name: indexNameMinIO,
			},
			wantErr: false,
		},
	}

	t.Run("prepare", func(t *testing.T) {
		index, _, err := GetOrCreateIndex(indexName, "disk", 1)
		assert.NoError(t, err)
		assert.NotNil(t, index)

		indexS3, _, err := GetOrCreateIndex(indexNameS3, "disk", 1)
		assert.NoError(t, err)
		assert.NotNil(t, indexS3)
		indexS3.ref.Meta.StorageType = "s3"

		indexMinio, _, err := GetOrCreateIndex(indexNameMinIO, "disk", 1)
		assert.NoError(t, err)
		assert.NotNil(t, indexMinio)
		indexS3.ref.Meta.StorageType = "minio"
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteIndex(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DeleteIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
