package zutils

import "strconv"

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
