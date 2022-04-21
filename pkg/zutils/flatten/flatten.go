package flatten

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Unflatten takes a map where dot-delimited keys are replaced by nested maps
func Unflatten(flat map[string]interface{}) (map[string]interface{}, error) {
	unflat := map[string]interface{}{}

	for key, value := range flat {
		keyParts := strings.Split(key, ".")

		// Walk the keys until we get to a leaf node.
		m := unflat
		for i, k := range keyParts[:len(keyParts)-1] {
			v, exists := m[k]
			if !exists {
				newMap := map[string]interface{}{}
				m[k] = newMap
				m = newMap
				continue
			}

			innerMap, ok := v.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("key=%s is not an object", strings.Join(keyParts[0:i+1], "."))
			}
			m = innerMap
		}

		leafKey := keyParts[len(keyParts)-1]
		if _, exists := m[leafKey]; exists {
			return nil, fmt.Errorf("key=%s already exists", key)
		}
		m[keyParts[len(keyParts)-1]] = value
	}

	return unflat, nil
}

// ErrNotValidInput Nested input must be a map or slice
var ErrNotValidInput = errors.New("not a valid input: map or slice")

// Flatten generates a flat map from a nested one.  The original may include values of type map, slice and scalar,
// but not struct.  Keys in the flat map will be a compound of descending map keys and slice iterations.
// The presentation of keys is set by style.  A prefix is joined to each key.
func Flatten(nested map[string]interface{}, prefix string) (map[string]interface{}, error) {
	flatmap := make(map[string]interface{})

	err := flatten(true, flatmap, nested, prefix)
	if err != nil {
		return nil, err
	}

	return flatmap, nil
}

func flatten(top bool, flatMap map[string]interface{}, nested interface{}, prefix string) error {
	assign := func(newKey string, v interface{}) error {
		switch v.(type) {
		case map[string]interface{}, []interface{}:
			if err := flatten(false, flatMap, v, newKey); err != nil {
				return err
			}
		default:
			flatMap[newKey] = v
		}

		return nil
	}

	switch v := nested.(type) {
	case map[string]interface{}:
		for k, v := range v {
			newKey := enkey(top, prefix, k)
			assign(newKey, v)
		}
	case []interface{}:
		needFlat := true
		for _, v := range v {
			if err := checkTypeIsMapOrSlice(v); err != nil {
				needFlat = false
			}
		}
		if needFlat {
			for i, v := range v {
				newKey := enkey(top, prefix, strconv.Itoa(i))
				assign(newKey, v)
			}
		} else {
			flatMap[prefix] = v
		}
	default:
		return ErrNotValidInput
	}

	return nil
}

func checkTypeIsMapOrSlice(v interface{}) error {
	switch v.(type) {
	case map[string]interface{}, []interface{}:
		return nil
	default:
		return ErrNotValidInput
	}
}

func enkey(top bool, prefix, subkey string) string {
	key := prefix

	if top {
		key += subkey
	} else {
		key += "." + subkey
	}

	return key
}
