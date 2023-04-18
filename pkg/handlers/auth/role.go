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

package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zincsearch/zincsearch/pkg/auth"
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

// @Id CreateRole
// @Summary Create role
// @security BasicAuth
// @Tags    Role
// @Accept  json
// @Produce json
// @Param   role body meta.Role true "Role data"
// @Success 200 {object} meta.HTTPResponseID
// @Failure 400 {object} meta.HTTPResponseError
// @Failure 500 {object} meta.HTTPResponseError
// @Router /api/role [post]
func CreateUpdateRole(c *gin.Context) {
	var role meta.Role
	if err := zutils.GinBindJSON(c, &role); err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	if role.ID == "" {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: "role.id should be not empty"})
		return
	}

	newRole, err := auth.CreateRole(role.ID, role.Name, role.Permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponseError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, meta.HTTPResponseID{Message: "ok", ID: newRole.ID})
}

// @Id ListRoles
// @Summary List role
// @security BasicAuth
// @Tags    Role
// @Produce json
// @Success 200 {object} []meta.Role
// @Failure 500 {object} meta.HTTPResponseError
// @Router /api/role [get]
func ListRole(c *gin.Context) {
	roles, err := auth.GetRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// @Id DeleteRole
// @Summary Delete role
// @security BasicAuth
// @Tags    Role
// @Produce json
// @Param   id  path  string  true  "Role id"
// @Success 200 {object} meta.HTTPResponseID
// @Success 500 {object} meta.HTTPResponseError
// @Router /api/role/{id} [delete]
func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	err := auth.DeleteRole(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponseError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, meta.HTTPResponseID{Message: "deleted", ID: id})
}
