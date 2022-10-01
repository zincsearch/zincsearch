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
	"github.com/zinclabs/zinc/pkg/zutils"
)

// @Id CreateUser
// @Summary Create user
// @Tags    User
// @Accept  json
// @Produce json
// @Param   user body meta.User true "User data"
// @Success 200 {object} meta.HTTPResponseID
// @Failure 400 {object} meta.HTTPResponseError
// @Failure 500 {object} meta.HTTPResponseError
// @Router /api/user [post]
func CreateUpdateUser(c *gin.Context) {
	var user meta.User
	if err := zutils.GinBindJSON(c, &user); err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	if user.ID == "" {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: "user.id should be not empty"})
		return
	}

	newUser, err := auth.CreateUser(user.ID, user.Name, user.Password, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponseError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, meta.HTTPResponseID{Message: "ok", ID: newUser.ID})
}

// @Id UpdateUser
// @Summary Update user
// @Tags    User
// @Accept  json
// @Produce json
// @Param   user body meta.User true "User data"
// @Success 200 {object} meta.HTTPResponseID
// @Failure 400 {object} meta.HTTPResponseError
// @Failure 500 {object} meta.HTTPResponseError
// @Router /api/user [put]
func UpdateForSDK() {}

// @Id ListUsers
// @Summary List user
// @Tags    User
// @Produce json
// @Success 200 {object} []meta.User
// @Failure 500 {object} meta.HTTPResponseError
// @Router /api/user [get]
func ListUser(c *gin.Context) {
	users, err := auth.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	for _, u := range users {
		// remove password and salt from response
		u.Salt = ""
		u.Password = ""

	}
	c.JSON(http.StatusOK, users)
}

// @Id DeleteUser
// @Summary Delete user
// @Tags    User
// @Produce json
// @Param   id  path  string  true  "User id"
// @Success 200 {object} meta.HTTPResponseID
// @Success 500 {object} meta.HTTPResponseError
// @Router /api/user/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := auth.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponseError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, meta.HTTPResponseID{Message: "deleted", ID: id})
}
