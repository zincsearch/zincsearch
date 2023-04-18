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

package utils

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/zincsearch/zincsearch/pkg/zutils/json"
)

func NewGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.ReleaseMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	c.Request = req
	return c, w
}

func SetGinRequestData(c *gin.Context, data interface{}) {
	c.Request.Header.Set("Content-Type", "application/json;charset=utf-8")
	switch v := data.(type) {
	case string:
		c.Request.Body = io.NopCloser(bytes.NewBufferString(v))
	case map[string]interface{}:
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return
		}
		buf := bytes.NewBuffer(jsonBytes)
		c.Request.Body = io.NopCloser(buf)
	default:
	}
}

func SetGinRequestURL(c *gin.Context, path string, params map[string]string) {
	q := c.Request.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	c.Request.URL.Path = path
	c.Request.URL.RawQuery = q.Encode()
}

func SetGinRequestParams(c *gin.Context, params map[string]string) {
	p := gin.Params{}
	for k, v := range params {
		p = append(p, gin.Param{Key: k, Value: v})
	}
	c.Params = p
}
