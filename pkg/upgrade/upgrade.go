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

package upgrade

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/zinclabs/zinc/pkg/meta"
)

func Do(oldVersion string, index *meta.Index) error {
	var err error
	log.Info().Msgf("Begin upgrade[%s] from version[%s]", index.Meta.Name, oldVersion)
	columns := strings.Split(oldVersion, "-")
	oldVersion = columns[0]
	switch oldVersion {
	case "v0.2.4":
		if err = UpgradeFromV024T025(index); err != nil {
			return err
		}
		return Do("v0.2.5", index)
	case "v0.2.5":
		if err = UpgradeFromV025T026(index); err != nil {
			return err
		}
		return Do("v0.2.6", index)
	case "v0.2.6":
		if err = UpgradeFromV026T027(index); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported upgrade from version: %s", oldVersion)
	}
}
