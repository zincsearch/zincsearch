package meta

import (
	"errors"
	"strconv"
	"strings"

	"github.com/zincsearch/zincsearch/pkg/zutils"
)

// point is a geo point
// format1: {"type": "Point", "coordinates": [-71.34, 41.12]}
// format2: POINT (-71.34 41.12)
// format3: {"lat": 41.12, "lon": -71.34}
// format4: [-71.34, 41.12]
// foramt5: "41.12,-71.34"
type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func ParseGeoPoint(v interface{}) (GeoPoint, error) {
	switch v := v.(type) {
	case string:
		v = strings.ToLower(v)
		v = strings.TrimSpace(v)
		if strings.HasPrefix(v, "point") {
			// POINT (-71.34 41.12)
			v = strings.TrimPrefix(v, "point")
			v = strings.TrimSpace(v)
			v = strings.TrimPrefix(v, "(")
			v = strings.TrimSuffix(v, ")")
			v = strings.TrimSpace(v)
			point := strings.Split(v, " ")
			if len(point) != 2 {
				return GeoPoint{}, errors.New("invalid point")
			}
			lon, _ := strconv.ParseFloat(point[0], 64)
			lat, _ := strconv.ParseFloat(point[1], 64)
			return GeoPoint{
				Lat: lat,
				Lon: lon,
			}, nil
		} else {
			// "41.12,-71.34"
			point := strings.Split(v, ",")
			if len(point) != 2 {
				return GeoPoint{}, errors.New("invalid point")
			}
			lat, _ := strconv.ParseFloat(point[0], 64)
			lon, _ := strconv.ParseFloat(point[1], 64)
			return GeoPoint{
				Lat: lat,
				Lon: lon,
			}, nil
		}
	case []float64:
		// [-71.34, 41.12]
		if len(v) != 2 {
			return GeoPoint{}, errors.New("invalid point")
		}
		return GeoPoint{
			Lat: v[1],
			Lon: v[0],
		}, nil
	case map[string]interface{}:
		if v["type"] != nil && v["coordinates"] != nil {
			// {"type": "Point", "coordinates": [-71.34, 41.12]}
			coordinates, ok := v["coordinates"].([]interface{})
			if !ok || len(coordinates) != 2 {
				return GeoPoint{}, errors.New("invalid point")
			}
			lon, _ := zutils.ToFloat64(coordinates[0])
			lat, _ := zutils.ToFloat64(coordinates[1])
			return GeoPoint{
				Lat: lat,
				Lon: lon,
			}, nil
		}
		// {"lat": 41.12, "lon": -71.34}
		lat, _ := zutils.ToFloat64(v["lat"])
		lon, _ := zutils.ToFloat64(v["lon"])
		return GeoPoint{
			Lat: lat,
			Lon: lon,
		}, nil
	default:
		return GeoPoint{}, errors.New("invalid point")
	}
}
