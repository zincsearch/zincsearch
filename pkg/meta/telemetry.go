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
	"io"
	"log"
	"time"

	"github.com/segmentio/analytics-go/v3"
)

var SEGMENT_CLIENT analytics.Client

func init() {
	cf := analytics.Config{
		Interval:  15 * time.Second,
		BatchSize: 10,
		Endpoint:  "https://e1.zinclabs.dev",
		Verbose:   false,
		Logger:    analytics.StdLogger(log.New(io.Discard, "marker ", log.LstdFlags)), // discard any logs
	}

	SEGMENT_CLIENT, _ = analytics.NewWithConfig("", cf)
}
