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
	"github.com/zincsearch/zincsearch/pkg/errors"
	"github.com/zincsearch/zincsearch/pkg/meta"
	"github.com/zincsearch/zincsearch/pkg/zutils"
)

type UpdateAccountRequest struct {
	ID          string `json:"_id"`
	Password    string `json:"password"`
	Name        string `json:"name,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

// @Id UpdateAccount
// @Summary Update own name and/or password
// @Tags    User
// @Accept  json
// @Produce json
// @Param   account body UpdateAccountRequest true "Current password plus new name and/or new password"
// @Success 200 {object} LoginUser
// @Failure 400 {object} meta.HTTPResponseError
// @Failure 401 {object} meta.HTTPResponseError
// @Failure 500 {object} meta.HTTPResponseError
// @Router /api/account [put]
func UpdateAccount(c *gin.Context) {
	var req UpdateAccountRequest
	if err := zutils.GinBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}
	if req.ID == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: "_id and password should be not empty"})
		return
	}
	if req.Name == "" && req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: "name or new_password should be not empty"})
		return
	}

	user, err := auth.UpdateAccount(req.ID, req.Password, req.Name, req.NewPassword)
	if errors.Is(err, auth.ErrInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, meta.HTTPResponseError{Error: err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, meta.HTTPResponseError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, LoginUser{ID: user.ID, Name: user.Name, Role: user.Role})
}
