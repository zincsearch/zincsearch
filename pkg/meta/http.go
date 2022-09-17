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

type HTTPResponse struct {
	Message string `json:"message"`
}

type HTTPResponseID struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

type HTTPResponseDocument struct {
	Message string `json:"message"`
	Index   string `json:"index"`
	ID      string `json:"id,omitempty"`
}

type HTTPResponseIndex struct {
	Message     string `json:"message"`
	Index       string `json:"index"`
	StorageType string `json:"storage_type,omitempty"`
}

type HTTPResponseTemplate struct {
	Message  string `json:"message"`
	Template string `json:"template"`
}

type HTTPResponseRecordCount struct {
	Message     string `json:"message"`
	RecordCount int64  `json:"record_count"`
}

type HTTPResponseError struct {
	Error string `json:"error"`
}

type HTTPResponseESID struct {
	Message     string `json:"message"`
	ID          string `json:"id"`
	ESID        string `json:"_id"`
	Index       string `json:"_index"`
	Version     int    `json:"_version"`
	SeqNo       int    `json:"_seq_no"`
	PrimaryTerm int    `json:"_primary_term"`
	Result      string `json:"result"` // created, updated, deleted
}

type HttpRetriesResponse struct {
	Bulk   int `json:"bulk"`
	Search int `json:"search"`
}

type HTTPResponseDeleteByQuery struct {
	Took                 int64               `json:"took"`
	TimedOut             bool                `json:"time_out"`
	Total                int                 `json:"total"`
	Deleted              int                 `json:"deleted"`
	Batches              int                 `json:"batches"`
	VersionConflicts     int                 `json:"version_conflicts"`
	Noops                int                 `json:"noops"`
	Failures             []string            `json:"failures"`
	Retries              HttpRetriesResponse `json:"retries"`
	ThrottledMillis      int                 `json:"throttled_millis"`
	RequestsPerSecond    int                 `json:"requests_per_second"`
	ThrottledUntilMillis int                 `json:"throttled_until_millis"`
}
