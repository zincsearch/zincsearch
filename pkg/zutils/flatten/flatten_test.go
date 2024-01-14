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

package flatten

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlattern(t *testing.T) {
	data := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": map[string]interface{}{
				"oxx": "cbd",
				"xxo": "dba",
			},
			"arr": []interface{}{"a", "b", "c"},
			"arm": []interface{}{
				map[string]interface{}{
					"a1": "b1",
					"a2": "b2",
				},
				map[string]interface{}{
					"b1": "a1",
					"b2": "a2",
				},
			},
		},
	}
	fdata, err := Flatten(data, "")
	assert.NoError(t, err)
	assert.Equal(t, 7, len(fdata))
	assert.Equal(t, "cbd", fdata["foo.bar.oxx"].(string))
	assert.Equal(t, "a1", fdata["foo.arm.1.b1"].(string))
	assert.Equal(t, 3, len(fdata["foo.arr"].([]interface{})))
}

func TestUnflatten(t *testing.T) {
	data := map[string]interface{}{
		"foo.bar.coo": "abc",
		"foo.bar.oxx": "cbd",
		"foo.bcc.xox": "bdc",
	}
	undata, err := Unflatten(data)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(undata))
	assert.Equal(t, 2, len(undata["foo"].(map[string]interface{})))
	assert.Equal(t, 2, len(undata["foo"].(map[string]interface{})["bar"].(map[string]interface{})))
	assert.Equal(t, "abc", undata["foo"].(map[string]interface{})["bar"].(map[string]interface{})["coo"])
}
