// Copyright 2022 Zinc Labs Inc. and Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package routes

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// SetPrometheus sets up prometheus metrics for gin
func SetPrometheus(r *gin.Engine) {
	enable := false
	if v := os.Getenv("ZINC_PROMETHEUS_ENABLE"); v != "" {
		enable, _ = strconv.ParseBool(v)
	}
	if !enable {
		return
	}

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)
}
