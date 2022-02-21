package zutils

import "fmt"

func GetStringFromMap(m interface{}, key string) (string, error) {
	if m == nil {
		return "", fmt.Errorf("GetStringFromMap: map is nil")
	}
	mm, ok := m.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("GetStringFromMap: map is not a map (object)")
	}
	v, ok := mm[key]
	if !ok {
		return "", fmt.Errorf("GetStringFromMap: key [%s] not found", key)
	}
	vs, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("GetStringFromMap: value [%s] is not a string", key)
	}

	return vs, nil
}

func GetStringSliceFromMap(m interface{}, key string) ([]string, error) {
	if m == nil {
		return nil, fmt.Errorf("GetStringSliceFromMap: map is nil")
	}
	mm, ok := m.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("GetStringSliceFromMap: map is not a map (object)")
	}
	v, ok := mm[key]
	if !ok {
		return nil, fmt.Errorf("GetStringSliceFromMap: key [%s] not found", key)
	}
	vs, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("GetStringSliceFromMap: value [%s] is not an array of string", key)
	}

	ss := make([]string, 0, len(vs))
	for _, v := range vs {
		sv, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("GetStringSliceFromMap: value [%s] is not an array of string", key)
		}
		ss = append(ss, sv)
	}

	return ss, nil
}
