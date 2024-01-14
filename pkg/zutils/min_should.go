package zutils

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type combination struct {
	condition int
	count     int
}

var regex = regexp.MustCompile(`(\d+)<([-+]?\d+%?)`)

// CalculateMin
// calculate the MinimumShouldMatch value with given expr and sub query count.
func CalculateMin(subCount int, v interface{}) (res int, err error) {
	if subCount == 0 {
		return 1, nil
	}
	defer func() {
		if err != nil {
			return
		}
		if res <= 1 {
			res = 1
			return
		}
		if res >= subCount {
			res = subCount
			return
		}
	}()
	switch x := v.(type) {
	case int64, int, float64:
		m := 0
		switch val := v.(type) {
		case int:
			m = val
		case int64:
			m = int(val)
		case float64:
			m = int(math.Floor(val))
		}
		if m < 0 {
			m = subCount + m
		}
		return m, nil
	case []string:
		conditions := make([]combination, len(x))
		for i, str := range x {
			match := regex.FindStringSubmatch(str)
			if match != nil {
				leftPart := match[1]
				rightPart := match[2]
				condition, err := strconv.ParseInt(leftPart, 10, 64)
				if err != nil {
					return 0, fmt.Errorf("cannot parse the condition value: %w", err)
				}
				count, err := getPartValue(subCount, rightPart)
				if err != nil {
					return 0, fmt.Errorf("cannot parse the clauses count: %w", err)
				}
				conditions[i] = combination{
					condition: int(condition),
					count:     count,
				}
			} else {
				return 0, fmt.Errorf("invalid MinimumShould value: %s", x)
			}
		}

		sort.Slice(conditions, func(i, j int) bool {
			return conditions[i].condition < conditions[j].condition
		})

		for i, condition := range conditions {
			// only match first
			if subCount <= condition.condition {
				return subCount, nil
			}
			// we are the last one, matched
			if i == len(conditions)-1 {
				return condition.count, nil
			}
			// less than next, we matched
			if subCount <= conditions[i+1].condition {
				return condition.count, nil
			}
		}
		return 0, fmt.Errorf("invalid MinimumShould value: %v", x)
	case string:
		combinations := strings.Split(x, " ")
		if len(combinations) > 1 {
			return CalculateMin(subCount, combinations)
		}
		// simple expr
		if res, err := getPartValue(subCount, x); err == nil {
			return res, nil
		}

		// complex, we use regex
		match := regex.FindStringSubmatch(x)
		if match != nil {
			leftPart := match[1]
			rightPart := match[2]
			condition, err := strconv.ParseInt(leftPart, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("cannot parse the condition value: %w", err)
			}
			count, err := getPartValue(subCount, rightPart)
			if err != nil {
				return 0, fmt.Errorf("cannot parse the clauses count: %w", err)
			}
			if subCount <= int(condition) {
				return subCount, nil
			}
			return count, nil
		} else {
			return 0, fmt.Errorf("invalid MinimumShould value: %s", x)
		}
	default:
		return 0, fmt.Errorf("invalid MinimumShouldMatch value")
	}
}

func getPartValue(termCount int, part string) (int, error) {
	if strings.Contains(part, "%") {
		proportion, err := strconv.ParseInt(part[0:len(part)-1], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot parse a percent value: %w", err)
		}
		if proportion < 0 {
			count := float64(termCount) * float64(-proportion) / 100
			return termCount - int(count), nil
		}
		count := float64(termCount) * float64(proportion) / 100
		return int(count), nil
	}
	count, err := strconv.ParseInt(part, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse a int value: %w", err)
	}
	if count < 0 {
		return termCount + int(count), nil
	}
	return int(count), nil
}
