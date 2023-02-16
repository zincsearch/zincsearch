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

package index

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zincsearch/pkg/core"
	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/pkg/uquery/template"
	"github.com/zinclabs/zincsearch/pkg/zutils"
)

// @Id ListTemplates
// @Summary List index teplates
// @security BasicAuth
// @Tags    Index
// @Produce json
// @Success 200 {object} []meta.Template
// @Failure 400 {object} meta.HTTPResponseError
// @Router /es/_index_template [get]
func ListTemplate(c *gin.Context) {
	pattern := c.Query("pattern")
	templates, err := core.ListTemplates(pattern)
	if err != nil {
		zutils.GinRenderJSON(c, http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}
	zutils.GinRenderJSON(c, http.StatusOK, templates)
}

// @Id GetTemplate
// @Summary Get index template
// @security BasicAuth
// @Tags    Index
// @Produce json
// @Param   name path  string  true  "Template"
// @Success 200 {object} meta.IndexTemplate
// @Failure 400 {object} meta.HTTPResponseError
// @Router /es/_index_template/{name} [get]
func GetTemplate(c *gin.Context) {
	name := c.Param("target")
	if name == "" {
		zutils.GinRenderJSON(c, http.StatusBadRequest, meta.HTTPResponseError{Error: "template.name should be not empty"})
		return
	}
	template, exists, err := core.LoadTemplate(name)
	if err != nil {
		zutils.GinRenderJSON(c, http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}
	if !exists {
		zutils.GinRenderJSON(c, http.StatusNotFound, meta.HTTPResponseError{Error: "template " + name + " does not exists"})
		return
	}
	zutils.GinRenderJSON(c, http.StatusOK, template)
}

// @Id CreateTemplate
// @Summary Create update index template
// @security BasicAuth
// @Tags    Index
// @Accept  json
// @Produce json
// @Param   template body meta.IndexTemplate true "Template data"
// @Success 200 {object} meta.HTTPResponseTemplate
// @Failure 400 {object} meta.HTTPResponseError
// @Router /es/_index_template [post]
func CreateTemplate(c *gin.Context) {
	data := make(map[string]interface{})
	if err := zutils.GinBindJSON(c, &data); err != nil {
		zutils.GinRenderJSON(c, http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	name := c.Param("target")
	if name == "" {
		if v, ok := data["name"]; ok {
			name, _ = v.(string)
		}
	}
	if name == "" {
		zutils.GinRenderJSON(c, http.StatusBadRequest, meta.HTTPResponseError{Error: "template.name should be not empty"})
		return
	}

	template, err := template.Request(data)
	if err != nil {
		zutils.GinRenderJSON(c, http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	err = core.NewTemplate(name, template)
	if err != nil {
		zutils.GinRenderJSON(c, http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	zutils.GinRenderJSON(c, http.StatusOK, meta.HTTPResponseTemplate{Message: "ok", Template: name})
}

// @Id UpdateTemplate
// @Summary Create update index template
// @security BasicAuth
// @Tags    Index
// @Accept  json
// @Produce json
// @Param   name     path string  true  "Template"
// @Param   template body meta.IndexTemplate true "Template data"
// @Success 200 {object} meta.HTTPResponseTemplate
// @Failure 400 {object} meta.HTTPResponseError
// @Router /es/_index_template/{name} [put]
func UpdateTemplateForSDK() {}

// @Id DeleteTemplate
// @Summary Delete template
// @security BasicAuth
// @Tags    Index
// @Produce json
// @Param   name  path  string  true  "Template"
// @Success 200 {object} meta.HTTPResponse
// @Failure 400 {object} meta.HTTPResponseError
// @Router /es/_index_template/{name} [delete]
func DeleteTemplate(c *gin.Context) {
	name := c.Param("target")
	err := core.DeleteTemplate(name)
	if err != nil {
		zutils.GinRenderJSON(c, http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}
	zutils.GinRenderJSON(c, http.StatusOK, meta.HTTPResponse{Message: "ok"})
}
