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

package core

// UpdateDocument inserts or updates a document in the zinc index
func (index *Index) UpdateDocument(docID string, doc map[string]interface{}, mintedID bool) error {
	bdoc, err := index.BuildBlugeDocumentFromJSON(docID, doc)
	if err != nil {
		return err
	}

	// Finally update the document on disk
	writer := index.Writer
	if !mintedID {
		err = writer.Update(bdoc.ID(), bdoc)
	} else {
		err = writer.Insert(bdoc)
		index.GainDocsCount(1)
	}
	return err
}
