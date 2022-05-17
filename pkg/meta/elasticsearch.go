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

package meta

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/zinclabs/zinc/pkg/zutils"
)

func NewESInfo(c *gin.Context) *ESInfo {
	version := strings.TrimLeft(Version, "v")
	userAgent := c.Request.UserAgent()
	if strings.Contains(userAgent, "Elastic") {
		reg := regexp.MustCompile(`([0-9]+\.[0-9]+\.[0-9]+)`)
		matches := reg.FindAllString(userAgent, 1)
		if len(matches) > 0 {
			version = matches[0]
		}
	}
	if v := strings.ToUpper(zutils.GetEnv("ZINC_PLUGIN_ES_VERSION", "")); v != "" {
		version = v
	}
	return &ESInfo{
		Name:        "zinc",
		ClusterName: "N/A",
		ClusterUUID: "N/A",
		Version: ESInfoVersion{
			Number:                    version,
			BuildFlavor:               "default",
			BuildHash:                 CommitHash,
			BuildDate:                 BuildDate,
			BuildSnapshot:             false,
			LuceneVersion:             "N/A",
			MinimumWireVersion:        "N/A",
			MinimumIndexCompatibility: "N/A",
		},
		Tagline: "You Know, for Search",
	}
}

func NewESLicense(_ *gin.Context) *ESLicense {
	return &ESLicense{
		Status: "active",
	}
}

func NewESXPack(_ *gin.Context) *ESXPack {
	return &ESXPack{
		Build:    make(map[string]bool),
		Features: make(map[string]bool),
		License: ESLicense{
			Status: "active",
		},
	}
}

type ESInfo struct {
	Name        string        `json:"name"`
	ClusterName string        `json:"cluster_name"`
	ClusterUUID string        `json:"cluster_uuid"`
	Version     ESInfoVersion `json:"version"`
	Tagline     string        `json:"tagline"`
}

type ESInfoVersion struct {
	Number                    string `json:"number"`
	BuildFlavor               string `json:"build_flavor"`
	BuildHash                 string `json:"build_hash"`
	BuildDate                 string `json:"build_date"`
	BuildSnapshot             bool   `json:"build_snapshot"`
	LuceneVersion             string `json:"lucene_version"`
	MinimumWireVersion        string `json:"minimum_wire_version"`
	MinimumIndexCompatibility string `json:"minimum_index_compatibility"`
}

type ESLicense struct {
	Status string `json:"status"`
}

type ESXPack struct {
	Build    map[string]bool `json:"build"`
	Features map[string]bool `json:"features"`
	License  ESLicense       `json:"license"`
}
