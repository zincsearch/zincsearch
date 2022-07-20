package meta

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Page struct {
	PageNum  int64 `json:"page_num"`
	PageSize int64 `json:"page_size"`
	Total    int64 `json:"total"`
}

func NewPage(c *gin.Context) *Page {
	pageNum, _ := strconv.ParseInt(c.DefaultQuery("page_num", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "0"), 10, 64)
	if pageNum < 1 {
		pageNum = 1
	}
	return &Page{
		PageNum:  pageNum,
		PageSize: pageSize,
	}
}

func (e *Page) GetStartEndIndex() (startIndex, endIndex int64) {
	if e.PageSize == 0 {
		return 0, e.Total
	}
	startIndex = (e.PageNum - 1) * e.PageSize
	endIndex = e.PageNum * e.PageSize

	if e.Total < endIndex {
		startIndex = e.Total - e.Total%e.PageSize
		endIndex = e.Total
	}
	return
}
