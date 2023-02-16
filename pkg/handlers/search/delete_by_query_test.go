package search

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zincsearch/pkg/core"
	"github.com/zinclabs/zincsearch/pkg/ider"
	"github.com/zinclabs/zincsearch/pkg/meta"
	"github.com/zinclabs/zincsearch/test/utils"
)

type arg struct {
	doc    map[string]interface{}
	query  string
	params map[string]string
}

type body struct {
	is       string
	contains string
}

type success struct {
	outcome    bool
	statusCode int
	body       body
}

type failure struct {
	outcome    bool
	statusCode int
	body       body
}

type want struct {
	success success
	failure failure
}

func TestDeleteByQuery(t *testing.T) {
	tests := []struct {
		name string
		arg  arg
		want want
	}{
		{
			name: "should delete matched documents",
			arg: arg{
				doc: map[string]interface{}{
					"name": "zinc",
				},
				query: `{"query":{"match":{"name":"zinc"}},"size":10}`,
				params: map[string]string{
					"target": "TestDeleteByQuery.index",
				},
			},
			want: want{
				success: success{
					outcome:    true,
					statusCode: 200,
					body: body{
						contains: `"time_out":false,"total":1,"deleted":1,"batches":0,"version_conflicts":0,"noops":0,"failures":[],"retries":{"bulk":0,"search":0},"throttled_millis":0,"requests_per_second":-1,"throttled_until_millis":0}`,
					},
				},
			},
		},
		{
			name: "should return bad request with invalid json body",
			arg: arg{
				doc: map[string]interface{}{
					"name": "zinc",
				},
				query: `invalid { json }`,
				params: map[string]string{
					"target": "TestDeleteByQuery.index",
				},
			},
			want: want{
				failure: failure{
					statusCode: 400,
					body: body{
						is: `{"error":"invalid character 'i' looking for beginning of value"}`,
					},
				},
			},
		},
		{
			name: "should return bad request when no matching indices are found",
			arg: arg{
				doc: map[string]interface{}{
					"name": "zinc",
				},
				query: `{"query":{"match":{"name":"zinc"}},"size":10}`,
				params: map[string]string{
					"target": "noneMatchingIndex",
				},
			},
			want: want{
				failure: failure{
					statusCode: 400,
					body: body{
						is: `{"error":"index noneMatchingIndex does not exists"}`,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			index, err := core.NewIndex("TestDeleteByQuery.index", "disk", 2)
			assert.NoError(t, err)
			assert.NoError(t, core.StoreIndex(index))
			id := ider.Generate()
			assert.NoError(t, index.CreateDocument(id, test.arg.doc, false))
			time.Sleep(time.Second)

			c, w := utils.NewGinContext()
			utils.SetGinRequestData(c, test.arg.query)
			utils.SetGinRequestParams(c, test.arg.params)
			DeleteByQuery(c)

			if test.want.success.outcome {
				time.Sleep(time.Second)
				assertHTTPResponse(t, w, test.want.success.statusCode, test.want.success.body)
				assertZeruResultQuery(t, index, test.arg.query)
			} else {
				assertHTTPResponse(t, w, test.want.failure.statusCode, test.want.failure.body)
			}

			assert.NoError(t, core.DeleteIndex(index.GetName()))
		})

	}
}

func assertHTTPResponse(t *testing.T, w *httptest.ResponseRecorder, statusCode int, body body) {
	assert.Equal(t, w.Code, statusCode)
	if body.is != "" {
		assert.Equal(t, w.Body.String(), body.is)
	} else {
		assert.Contains(t, w.Body.String(), body.contains)
	}
}

func assertZeruResultQuery(t *testing.T, index *core.Index, query interface{}) {
	jsonQuery, err := json.Marshal(&query)
	assert.NoError(t, err)
	var search, serr = index.Search(&meta.ZincQuery{
		Query: &meta.Query{
			Match: map[string]*meta.MatchQuery{
				"_all": {
					Query: string(jsonQuery),
				},
			},
		},
		Size: 10,
	})
	assert.NoError(t, serr)
	assert.Equal(t, 0, search.Hits.Total.Value)
}
