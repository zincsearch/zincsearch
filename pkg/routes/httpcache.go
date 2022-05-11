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

package routes

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func HTTPCacheForUI(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
			if strings.Contains(c.Request.RequestURI, "/ui/assets/") {
				c.Writer.Header().Set("cache-control", "public, max-age=2592000")
				c.Writer.Header().Set("expires", time.Now().Add(time.Hour*24*30).Format(time.RFC1123))
				if strings.Contains(c.Request.RequestURI, ".js") {
					c.Writer.Header().Set("content-type", "application/javascript")
				}
				if strings.Contains(c.Request.RequestURI, ".css") {
					c.Writer.Header().Set("content-type", "text/css; charset=utf-8")
				}
			}
		}

		c.Next()
	})
}
