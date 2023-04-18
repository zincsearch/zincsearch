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

package elastic

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zincsearch/zincsearch/pkg/zutils"
)

func GetDataStream(c *gin.Context) {
	target := c.Param("target")
	zutils.GinRenderJSON(c, http.StatusOK, gin.H{
		"data_streams": []gin.H{
			{
				"name": target,
				"timestamp_field": gin.H{
					"name": "@timestamp",
				},
			},
		},
	})
}

func PutDataStream(c *gin.Context) {
	target := c.Param("target")
	zutils.GinRenderJSON(c, http.StatusOK, gin.H{
		"name":    target,
		"message": "ok",
	})
}
