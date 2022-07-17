package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zinclabs/zinc/test/utils"
)

func TestNewPage(t *testing.T) {
	c, _ := utils.NewGinContext()
	params := map[string]string{
		"page_num":  "0",
		"page_size": "20",
	}
	utils.SetGinRequestURL(c, "/", params)
	page := NewPage(c)
	assert.Equal(t, page.PageNum, int64(0))
	assert.Equal(t, page.PageSize, int64(20))
}

func TestGetStartEndIndex(t *testing.T) {
	page := Page{
		PageNum:  0,
		PageSize: 20,
	}
	page.Total = 25
	startIndex, endIndex := page.GetStartEndIndex()
	assert.Equal(t, startIndex, int64(0))
	assert.Equal(t, endIndex, int64(20))
}
