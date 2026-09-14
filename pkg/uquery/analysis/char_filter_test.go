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

package analysis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestCharFilter(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		got, err := RequestCharFilter(nil)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
	t.Run("ok", func(t *testing.T) {
		got, err := RequestCharFilter(map[string]interface{}{
			"my_html": map[string]interface{}{"type": "html_strip"},
			"my_re":   map[string]interface{}{"type": "regexp", "pattern": "a", "replacement": "b"},
		})
		assert.NoError(t, err)
		assert.Len(t, got, 2)
		assert.Equal(t, "bbc", string(got["my_re"].Filter([]byte("abc"))))
	})
	t.Run("missing type", func(t *testing.T) {
		got, err := RequestCharFilter(map[string]interface{}{"x": map[string]interface{}{}})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("invalid options", func(t *testing.T) {
		got, err := RequestCharFilter(map[string]interface{}{"x": map[string]interface{}{"type": "regexp"}})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestRequestCharFilterSlice(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		got, err := RequestCharFilterSlice(nil)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})
	t.Run("string and object", func(t *testing.T) {
		got, err := RequestCharFilterSlice([]interface{}{
			"html_strip",
			map[string]interface{}{"type": "mapping", "mappings": []interface{}{"a => b"}},
		})
		assert.NoError(t, err)
		assert.Len(t, got, 2)
		assert.Equal(t, "bbc", string(got[1].Filter([]byte("abc"))))
	})
	t.Run("object missing type", func(t *testing.T) {
		got, err := RequestCharFilterSlice([]interface{}{map[string]interface{}{}})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("unknown string", func(t *testing.T) {
		got, err := RequestCharFilterSlice([]interface{}{"nope"})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
	t.Run("unsupported element type", func(t *testing.T) {
		got, err := RequestCharFilterSlice([]interface{}{1})
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestRequestCharFilterSingle(t *testing.T) {
	names := []string{
		"ascii_folding", "asciifolding", "html", "html_strip", "zero_width_non_joiner",
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			got, err := RequestCharFilterSingle(name, nil)
			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}

	t.Run("case insensitive", func(t *testing.T) {
		got, err := RequestCharFilterSingle("HTML_STRIP", nil)
		assert.NoError(t, err)
		assert.Equal(t, " hi ", string(got.Filter([]byte("<b>hi</b>"))))
	})
	t.Run("regexp aliases", func(t *testing.T) {
		for _, name := range []string{"regexp", "pattern", "pattern_replace"} {
			got, err := RequestCharFilterSingle(name, map[string]interface{}{"pattern": "x"})
			assert.NoError(t, err, name)
			assert.NotNil(t, got, name)
		}
	})
	t.Run("mapping", func(t *testing.T) {
		got, err := RequestCharFilterSingle("mapping", map[string]interface{}{"mappings": []interface{}{"a => b"}})
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
	t.Run("stconvert", func(t *testing.T) {
		got, err := RequestCharFilterSingle("stconvert", map[string]interface{}{"convert_type": "t2s"})
		assert.NoError(t, err)
		assert.Equal(t, "汉字", string(got.Filter([]byte("漢字"))))
	})
	t.Run("unknown", func(t *testing.T) {
		got, err := RequestCharFilterSingle("nope", nil)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
