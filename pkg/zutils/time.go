package zutils

import (
	"strconv"
	"strings"
	"time"
)

func ParseDuration(s string) (time.Duration, error) {
	d, err := time.ParseDuration(s)
	if err == nil {
		return d, nil
	}
	if !strings.HasSuffix(s, "d") {
		return 0, err
	}

	h := strings.TrimSuffix(s, "d")
	hour, _ := strconv.Atoi(h)
	d = time.Hour * time.Duration(hour) * 24
	return d, nil
}
