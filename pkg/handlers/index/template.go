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

	"github.com/zinclabs/zinc/pkg/core"
	"github.com/zinclabs/zinc/pkg/meta"
	"github.com/zinclabs/zinc/pkg/uquery/template"
)

// @Summary List index teplates
// @Tags    Index
// @Produce json
// @Success 200 {object} []meta.Template
// @Failure 400 {object} meta.HTTPResponse
// @Router /es/_index_template [get]
func ListTemplate(c *gin.Context) {
	pattern := c.Query("pattern")
	templates, err := core.ListTemplates(pattern)
	if err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, templates)
}

// @Summary Get index template
// @Tags    Index
// @Produce json
// @Param   target path  string  true  "Template"
// @Success 200 {object} meta.IndexTemplate
// @Failure 400 {object} meta.HTTPResponse
// @Router /es/_index_template/:target [get]
func GetTemplate(c *gin.Context) {
	name := c.Param("target")
	if name == "" {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: "template.name should be not empty"})
		return
	}
	template, exists, err := core.LoadTemplate(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: "template " + name + " does not exists"})
		return
	}
	c.JSON(http.StatusOK, template)
}

// @Summary Create update index template
// @Tags    Index
// @Produce json
// @Param   template body meta.IndexTemplate true "Template data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} meta.HTTPResponse
// @Router /es/_index_template [post]
func UpdateTemplate(c *gin.Context) {
	name := c.Param("target")
	if name == "" {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: "template.name should be not empty"})
		return
	}

	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: err.Error()})
		return
	}

	template, err := template.Request(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: err.Error()})
		return
	}

	err = core.NewTemplate(name, template)
	if err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok", "template": name})
}

// @Summary Delete template
// @Tags    Index
// @Param   target path  string  true  "Template"
// @Success 200 {object} meta.HTTPResponse
// @Failure 400 {object} meta.HTTPResponse
// @Router /es/_index_template/:target [delete]
func DeleteTemplate(c *gin.Context) {
	name := c.Param("target")
	err := core.DeleteTemplate(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, meta.HTTPResponse{Message: "ok"})
}
