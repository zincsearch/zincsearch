package query

import (
	"github.com/blugelabs/bluge"
	"github.com/prabhatsharma/zinc/pkg/errors"
)

func GeoBoundingBoxQuery(query map[string]interface{}) (bluge.Query, error) {
	return nil, errors.New(errors.ErrorTypeNotImplemented, "[geo_bounding_box] query doesn't support")
}

func GeoDistanceQuery(query map[string]interface{}) (bluge.Query, error) {
	return nil, errors.New(errors.ErrorTypeNotImplemented, "[geo_distance] query doesn't support")
}

func GeoPolygonQuery(query map[string]interface{}) (bluge.Query, error) {
	return nil, errors.New(errors.ErrorTypeNotImplemented, "[geo_polygon] query doesn't support")
}

func GeoShapeQuery(query map[string]interface{}) (bluge.Query, error) {
	return nil, errors.New(errors.ErrorTypeNotImplemented, "[geo_shape] query doesn't support")
}
