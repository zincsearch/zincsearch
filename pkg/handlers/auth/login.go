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

	"github.com/zinclabs/zincsearch/pkg/auth"
	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/pkg/zutils"
)

// @Id Login
// @Summary Login
// @Tags    User
// @Accept  json
// @Produce json
// @Param   login body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} meta.HTTPResponseError
// @Router /api/login [post]
func Login(c *gin.Context) {
	// Read login input
	var loginInput LoginRequest
	if err := zutils.GinBindJSON(c, &loginInput); err != nil {
		c.JSON(http.StatusBadRequest, meta.HTTPResponseError{Error: err.Error()})
		return
	}

	loggedInUser, validationResult := auth.VerifyCredentials(loginInput.ID, loginInput.Password)
	var resUser LoginUser
	if validationResult {
		resUser = LoginUser{
			ID:   loggedInUser.ID,
			Name: loggedInUser.Name,
			Role: loggedInUser.Role,
		}
	}
	c.JSON(http.StatusOK, LoginResponse{
		Validated: validationResult,
		User:      resUser,
	})
}

type LoginUser struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type LoginRequest struct {
	ID       string `json:"_id"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Validated bool      `json:"validated"`
	User      LoginUser `json:"user"`
}
