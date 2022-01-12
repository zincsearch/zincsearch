package es

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestApiBase(t *testing.T) {
	Convey("test es api", t, func() {
		// r := test.Server()
		Convey("POST /es/_bulk", func() {
			Convey("bulk create documents without indexName", func() {
			})
			Convey("bulk create documents with indexName", func() {
			})
			Convey("bulk with error input", func() {
			})
		})

		Convey("POST /es/:target/_bulk", func() {
			Convey("bulk create documents with not exist indexName", func() {
			})
			Convey("bulk create documents with exist indexName", func() {
			})
			Convey("bulk with error input", func() {
			})
		})

		Convey("POST /es/:target/_doc", func() {
			Convey("create document with not exist indexName", func() {
			})
			Convey("create document with exist indexName", func() {
			})
			Convey("create document with exist indexName not exist id", func() {
			})
			Convey("create document with exist indexName and exist id", func() {
			})
			Convey("create document with error input", func() {
			})
		})

		Convey("PUT /es/:target/_doc/:id", func() {
			Convey("create document with not exist indexName", func() {
			})
			Convey("create document with exist indexName", func() {
			})
			Convey("create document with exist indexName not exist id", func() {
			})
			Convey("create document with exist indexName and exist id", func() {
			})
			Convey("create document with error input", func() {
			})
		})

		Convey("DELETE /es/:target/_doc/:id", func() {
			Convey("delete document with not exist indexName", func() {
			})
			Convey("delete document with exist indexName", func() {
			})
			Convey("delete document with exist indexName not exist id", func() {
			})
			Convey("delete document with exist indexName and exist id", func() {
			})
			Convey("delete document with error input", func() {
			})
		})

		Convey("PUT /es/:target/_create/:id", func() {
			Convey("create document with not exist indexName", func() {
			})
			Convey("create document with exist indexName", func() {
			})
			Convey("create document with exist indexName not exist id", func() {
			})
			Convey("create document with exist indexName and exist id", func() {
			})
			Convey("create document with error input", func() {
			})
		})

		Convey("POST /es/:target/_create/:id", func() {
			Convey("create document with not exist indexName", func() {
			})
			Convey("create document with exist indexName", func() {
			})
			Convey("create document with exist indexName not exist id", func() {
			})
			Convey("create document with exist indexName and exist id", func() {
			})
			Convey("create document with error input", func() {
			})
		})

		Convey("POST /es/:target/_update/:id", func() {
			Convey("update document with not exist indexName", func() {
			})
			Convey("update document with exist indexName", func() {
			})
			Convey("update document with exist indexName not exist id", func() {
			})
			Convey("update document with exist indexName and exist id", func() {
			})
			Convey("update document with error input", func() {
			})
		})

	})
}
