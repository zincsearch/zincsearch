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

// @Summary Create update user
// @Tags    User
// @Produce json
// @Param   user body meta.User true "User data"
// @Success 200 {object} meta.HTTPResponse
// @Failure 400 {object} meta.HTTPResponse
// @Failure 500 {object} meta.HTTPResponse
// @Router /api/user [post]
func CreateUpdate(c *gin.Context) {
	var user meta.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: err.Error()})
		return
	}

	if user.ID == "" {
		c.JSON(http.StatusBadRequest, meta.HTTPResponse{Error: "user.id should be not empty"})
		return
	}

	newUser, err := auth.CreateUser(user.ID, user.Name, user.Password, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, meta.HTTPResponse{Message: "ok", Data: gin.H{"id": newUser.ID}})
}
