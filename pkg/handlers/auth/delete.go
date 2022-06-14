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

	"github.com/zinclabs/zinc/pkg/auth"
	"github.com/zinclabs/zinc/pkg/meta"
)

// @Summary Delete user
// @Tags    User
// @Param   id  path  string  true  "User id"
// @Success 200 {object} meta.HTTPResponse
// @Success 500 {object} meta.HTTPResponse
// @Router /api/user/:id [delete]
func Delete(c *gin.Context) {
	id := c.Param("id")
	err := auth.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, meta.HTTPResponse{Message: "deleted", ID: id})
}
