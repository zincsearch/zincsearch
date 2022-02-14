package zutils

import "strings"

// CamelCase camel case a string
func CamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	ss := strings.Split(s, " ")
	for k, v := range ss {
		ss[k] = strings.Title(v)
	}
	return strings.Join(ss, "")
}
