package zutils

import "fmt"

func GetStringFromMap(m interface{}, key string) (string, error) {
	v, err := GetAnyFromMap(m, key)
	if err != nil {
		return "", fmt.Errorf("GetStringFromMap: key [%s] not found", key)
	}
	vs, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("GetStringFromMap: value [%s] should be a string", key)
	}

	return vs, nil
}

func GetBoolFromMap(m interface{}, key string) (bool, error) {
	v, err := GetAnyFromMap(m, key)
	if err != nil {
		return false, fmt.Errorf("GetBoolFromMap: key [%s] not found", key)
	}
	vs, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("GetBoolFromMap: value [%s] shuld be a bool", key)
	}

	return vs, nil
}
func GetFloatFromMap(m interface{}, key string) (float64, error) {
	v, err := GetAnyFromMap(m, key)
	if err != nil {
		return 0, fmt.Errorf("GetFloatFromMap: key [%s] not found", key)
	}
	vs, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("GetFloatFromMap: value [%s] should be a float64", key)
	}

	return vs, nil
}

func GetStringSliceFromMap(m interface{}, key string) ([]string, error) {
	v, err := GetAnyFromMap(m, key)
	if err != nil {
		return nil, fmt.Errorf("GetStringSliceFromMap: key [%s] not found", key)
	}
	var vs []interface{}
	switch v := v.(type) {
	case []string:
		return v, nil
	case []interface{}:
		vs = v
	default:
		return nil, fmt.Errorf("GetStringSliceFromMap: value [%s] should be an array of string", key)
	}

	ss := make([]string, 0, len(vs))
	for _, v := range vs {
		sv, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("GetStringSliceFromMap: value [%s] should be an array of string", key)
		}
		ss = append(ss, sv)
	}

	return ss, nil
}

func GetAnyFromMap(m interface{}, key string) (interface{}, error) {
	if m == nil {
		return nil, fmt.Errorf("GetAnyFromMap: map is nil")
	}
	mm, ok := m.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("GetAnyFromMap: map should be a map / object")
	}
	v, ok := mm[key]
	if !ok {
		return nil, fmt.Errorf("GetAnyFromMap: key [%s] not found", key)
	}

	return v, nil
}
