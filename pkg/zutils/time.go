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

func FormatDuration(d time.Duration) string {
	if d.Hours() >= 24*30*12 {
		return strconv.FormatInt(int64(d.Hours())/24/30/12, 10) + "y"
	}
	if d.Hours() >= 24*30 {
		return strconv.FormatInt(int64(d.Hours())/24/30, 10) + "M"
	}
	if d.Hours() >= 24 {
		return strconv.FormatInt(int64(d.Hours())/24, 10) + "d"
	}
	if d.Hours() >= 1 {
		return strconv.FormatInt(int64(d.Hours()), 10) + "h"
	}
	if d.Minutes() >= 1 {
		return strconv.FormatInt(int64(d.Minutes()), 10) + "m"
	}
	return strconv.FormatInt(int64(d.Seconds()), 10) + "s"
}
