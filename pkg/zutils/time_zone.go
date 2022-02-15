package zutils

import (
	"strings"
	"time"
)

func ParseTimeZone(name string) (*time.Location, error) {
	offset := 0
	ln := len(name)
	if ln > 0 && (name[0] == '+' || name[0] == '-') {
		if ln >= 3 {
			offset = 60 * 60 * StringToInt(name[1:3])
		}
		if ln == 5 {
			offset += 60 * StringToInt(name[3:5])
		}
		if ln == 6 && name[3] == ':' {
			offset += 60 * StringToInt(name[4:6])
		}
		if name[0] == '-' {
			offset = -offset
		}
		return time.FixedZone(name, offset), nil
	}

	upperName := strings.ToUpper(name)
	if upperName == "" || upperName == "UTC" {
		return time.UTC, nil
	}
	if upperName == "LOCAL" {
		return time.Local, nil
	}

	return time.LoadLocation(name)
}
