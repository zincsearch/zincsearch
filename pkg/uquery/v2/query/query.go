package query

import (
	"fmt"
	"strings"

	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/analysis"

	"github.com/prabhatsharma/zinc/pkg/errors"
	meta "github.com/prabhatsharma/zinc/pkg/meta/v2"
)

func Query(query map[string]interface{}, mappings *meta.Mappings, analyzers map[string]*analysis.Analyzer) (bluge.Query, error) {
	var subq bluge.Query
	var cmd string
	var err error
	for k, t := range query {
		if subq != nil && cmd != "" {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[%s] malformed query, excepted [END_OBJECT] but found [FIELD_NAME] %s", cmd, k))
		}
		k := strings.ToLower(k)
		cmd = k
		v, ok := t.(map[string]interface{})
		if !ok {
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[%s] query doesn't support value type %T", k, t))
		}
		switch k {
		case "bool":
			if subq, err = BoolQuery(v, mappings, analyzers); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[bool] failed to parse field").Cause(err)
			}
		case "boosting":
			if subq, err = BoostingQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[boosting] failed to parse field").Cause(err)
			}
		case "match":
			if subq, err = MatchQuery(v, mappings, analyzers); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[match] failed to parse field").Cause(err)
			}
		case "match_bool_prefix":
			if subq, err = MatchBoolPrefixQuery(v, mappings, analyzers); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[match_bool_prefix] failed to parse field").Cause(err)
			}
		case "match_phrase":
			if subq, err = MatchPhraseQuery(v, mappings, analyzers); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[match_phrase] failed to parse field").Cause(err)
			}
		case "match_phrase_prefix":
			if subq, err = MatchPhrasePrefixQuery(v, mappings, analyzers); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[match_phrase_prefix] failed to parse field").Cause(err)
			}
		case "multi_match":
			if subq, err = MultiMatchQuery(v, mappings, analyzers); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[multi_match] failed to parse field").Cause(err)
			}
		case "match_all":
			if subq, err = MatchAllQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[match_all] failed to parse field").Cause(err)
			}
		case "match_none":
			if subq, err = MatchNoneQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[match_none] failed to parse field").Cause(err)
			}
		case "combined_fields":
			if subq, err = CombinedFieldsQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[combined_fields] failed to parse field").Cause(err)
			}
		case "query_string":
			if subq, err = QueryStringQuery(v, mappings, analyzers); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[query_string] failed to parse field").Cause(err)
			}
		case "simple_query_string":
			if subq, err = SimpleQueryStringQuery(v, mappings, analyzers); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[simple_query_string] failed to parse field").Cause(err)
			}
		case "exists":
			if subq, err = ExistsQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[exists] failed to parse field").Cause(err)
			}
		case "ids":
			if subq, err = IdsQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[ids] failed to parse field").Cause(err)
			}
		case "range":
			if subq, err = RangeQuery(v, mappings); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[range] failed to parse field").Cause(err)
			}
		case "regexp":
			if subq, err = RegexpQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[regexp] failed to parse field").Cause(err)
			}
		case "prefix":
			if subq, err = PrefixQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[prefix] failed to parse field").Cause(err)
			}
		case "fuzzy":
			if subq, err = FuzzyQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[fuzzy] failed to parse field").Cause(err)
			}
		case "wildcard":
			if subq, err = WildcardQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[wildcard] failed to parse field").Cause(err)
			}
		case "term":
			if subq, err = TermQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[term] failed to parse field").Cause(err)
			}
		case "terms":
			if subq, err = TermsQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[terms] failed to parse field").Cause(err)
			}
		case "terms_set":
			if subq, err = TermsSetQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[terms_set] failed to parse field").Cause(err)
			}
		case "geo_bounding_box":
			if subq, err = GeoBoundingBoxQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[geo_bounding_box] failed to parse field").Cause(err)
			}
		case "geo_distance":
			if subq, err = GeoDistanceQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[geo_distance] failed to parse field").Cause(err)
			}
		case "geo_polygon":
			if subq, err = GeoPolygonQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[geo_polygon] failed to parse field").Cause(err)
			}
		case "geo_shape":
			if subq, err = GeoShapeQuery(v); err != nil {
				return nil, errors.New(errors.ErrorTypeXContentParseException, "[geo_shape] failed to parse field").Cause(err)
			}
		default:
			return nil, errors.New(errors.ErrorTypeParsingException, fmt.Sprintf("[%s] query doesn't support", k))
		}
	}

	return subq, nil
}
