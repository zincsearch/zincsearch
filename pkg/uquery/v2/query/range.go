package query

import (
	"fmt"
	"strings"
	"time"

	"github.com/blugelabs/bluge"

	"github.com/zinclabs/zinc/pkg/errors"
	meta "github.com/zinclabs/zinc/pkg/meta/v2"
	"github.com/zinclabs/zinc/pkg/zutils"
)

func RangeQuery(query map[string]interface{}, mappings *meta.Mappings) (bluge.Query, error) {
	if len(query) > 1 {
		return nil, errors.New(errors.ErrorTypeParsingException, "[range] query doesn't support multiple fields")
	}

	for field, v := range query {
		vv, ok := v.(map[string]interface{})
		if !ok {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[range] query doesn't support values of type: %T", v))
		}
		switch mappings.Properties[field].Type {
		case "numeric":
			return RangeQueryNumeric(field, vv, mappings)
		case "time":
			return RangeQueryTime(field, vv, mappings)
		default:
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[range] %s only support values of [numeric, time]", field))
		}
	}

	return nil, nil
}

func RangeQueryNumeric(field string, query map[string]interface{}, mappings *meta.Mappings) (bluge.Query, error) {
	value := new(meta.RangeQuery)
	value.Boost = -1.0
	for k, v := range query {
		k := strings.ToLower(k)
		switch k {
		case "gt":
			value.GT = v.(float64)
		case "gte":
			value.GTE = v.(float64)
		case "lt":
			value.LT = v.(float64)
		case "lte":
			value.LTE = v.(float64)
		case "format":
			value.Format = v.(string)
		case "time_zone":
			value.TimeZone = v.(string)
		case "boost":
			value.Boost = v.(float64)
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[range] unknown field [%s]", k))
		}
	}

	min := 0.0
	max := 0.0
	minInclusive := false
	maxInclusive := false
	if value.GT != nil && value.GT.(float64) > 0 {
		min = value.GT.(float64)

	}
	if value.GTE != nil && value.GTE.(float64) > 0 {
		min = value.GTE.(float64)
		minInclusive = true
	}
	if value.LT != nil && value.LT.(float64) > 0 {
		max = value.LT.(float64)
	}
	if value.LTE != nil && value.LTE.(float64) > 0 {
		max = value.LTE.(float64)
		maxInclusive = true
	}
	subq := bluge.NewNumericRangeInclusiveQuery(min, max, minInclusive, maxInclusive).SetField(field)
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}

func RangeQueryTime(field string, query map[string]interface{}, mappings *meta.Mappings) (bluge.Query, error) {
	value := new(meta.RangeQuery)
	value.Boost = -1.0
	for k, v := range query {
		k := strings.ToLower(k)
		switch k {
		case "gt":
			value.GT = v.(string)
		case "gte":
			value.GTE = v.(string)
		case "lt":
			value.LT = v.(string)
		case "lte":
			value.LTE = v.(string)
		case "format":
			value.Format = v.(string)
		case "time_zone":
			value.TimeZone = v.(string)
		case "boost":
			value.Boost = v.(float64)
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[range] unknown field [%s]", k))
		}
	}

	var err error
	format := time.RFC3339
	if prop, ok := mappings.Properties[field]; ok {
		if prop.Format != "" {
			format = prop.Format
		}
	}
	if value.Format != "" {
		format = value.Format
	}
	timeZone := time.UTC
	if value.TimeZone != "" {
		timeZone, err = zutils.ParseTimeZone(value.TimeZone)
		if err != nil {
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[range] %s time_zone parse err %s", field, err.Error()))
		}
	}

	min := time.Time{}
	max := time.Time{}
	minInclusive := false
	maxInclusive := false
	if value.GT != nil && value.GT.(string) != "" {
		min, err = time.ParseInLocation(format, value.GT.(string), timeZone)
		if err != nil {
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[range] %s range.gt format err %s", field, err.Error()))
		}
	}
	if value.GTE != nil && value.GTE.(string) != "" {
		minInclusive = true
		min, err = time.ParseInLocation(format, value.GTE.(string), timeZone)
		if err != nil {
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[range] %s range.gte format err %s", field, err.Error()))
		}
	}
	if value.LT != nil && value.LT.(string) != "" {
		max, err = time.ParseInLocation(format, value.LT.(string), timeZone)
		if err != nil {
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[range] %s range.lt format err %s", field, err.Error()))
		}
	}
	if value.LTE != nil && value.LTE.(string) != "" {
		maxInclusive = true
		max, err = time.ParseInLocation(format, value.LTE.(string), timeZone)
		if err != nil {
			return nil, errors.New(errors.ErrorTypeXContentParseException, fmt.Sprintf("[range] %s range.lte format err %s", field, err.Error()))
		}
	}
	subq := bluge.NewDateRangeInclusiveQuery(min.UTC(), max.UTC(), minInclusive, maxInclusive).SetField(field)
	if value.Boost >= 0 {
		subq.SetBoost(value.Boost)
	}

	return subq, nil
}
