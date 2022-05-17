package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

func NewGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.ReleaseMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := &http.Request{
		URL:    &url.URL{},
		Header: make(http.Header),
	}
	c.Request = req
	return c, w
}

func SetGinRequestData(c *gin.Context, data interface{}) {
	c.Request.Header.Set("Content-Type", "application/json;charset=utf-8")
	switch v := data.(type) {
	case string:
		c.Request.Body = ioutil.NopCloser(bytes.NewBufferString(v))
	case map[string]interface{}:
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return
		}
		buf := bytes.NewBuffer(jsonBytes)
		c.Request.Body = ioutil.NopCloser(buf)
	default:
	}
}

func SetGinRequestURL(c *gin.Context, path string, params map[string]string) {
	q := c.Request.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	c.Request.URL.Path = path
	c.Request.URL.RawQuery = q.Encode()
}

func SetGinRequestParams(c *gin.Context, params map[string]string) {
	p := gin.Params{}
	for k, v := range params {
		p = append(p, gin.Param{Key: k, Value: v})
	}
	c.Params = p
}
