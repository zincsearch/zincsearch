package api

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestApiBase(t *testing.T) {
	Convey("test zinc api", t, func() {
		// r := test.Server()
		Convey("POST /api/login", func() {
			Convey("with username and password", func() {
			})
			Convey("with error username or password", func() {
			})
		})

		Convey("PUT /api/user", func() {
			Convey("create user with payload", func() {
			})
			Convey("create user with error input", func() {
			})
		})

		Convey("DELETE /api/user/:userID", func() {
			Convey("delete user with exist userid", func() {
			})
			Convey("delete user with not exist userid", func() {
			})
			Convey("delete user with error input", func() {
			})
		})

		Convey("GET /api/users", func() {
		})

		Convey("PUT /api/index", func() {
			Convey("create index with payload", func() {
			})
			Convey("create index with error input", func() {
			})
		})

		Convey("GET /api/index", func() {
		})

		Convey("DELETE /api/index/:indexName", func() {
			Convey("delete index with exist indexName", func() {
			})
			Convey("delete index with not exist indexName", func() {
			})
			Convey("delete index with error input", func() {
			})
		})

		Convey("POST /api/_bulk", func() {
			Convey("bulk create documents without indexName", func() {
			})
			Convey("bulk create documents with indexName", func() {
			})
			Convey("bulk with error input", func() {
			})
		})

		Convey("POST /api/:target/_bulk", func() {
			Convey("bulk create documents with not exist indexName", func() {
			})
			Convey("bulk create documents with exist indexName", func() {
			})
			Convey("bulk with error input", func() {
			})
		})

		Convey("PUT /api/:target/document", func() {
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

		Convey("POST /api/:target/_doc", func() {
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

		Convey("PUT /api/:target/_doc/:id", func() {
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

		Convey("DELETE /api/:target/_doc/:id", func() {
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

		Convey("POST /api/:target/_search", func() {
			Convey("search document with not exist indexName", func() {
			})
			Convey("search document with exist indexName", func() {
			})
			Convey("search document with not exist term", func() {
			})
			Convey("search document with exist term", func() {
			})
			Convey("search document type: alldocuments", func() {
			})
			Convey("search document type: wildcard", func() {
			})
			Convey("search document type: fuzzy", func() {
			})
			Convey("search document type: term", func() {
			})
			Convey("search document type: daterange", func() {
			})
			Convey("search document type: matchall", func() {
			})
			Convey("search document type: match", func() {
			})
			Convey("search document type: matchphrase", func() {
			})
			Convey("search document type: multiphrase", func() {
			})
			Convey("search document type: prefix", func() {
			})
			Convey("search document type: querystring", func() {
			})
			Convey("search with error input", func() {
			})
		})

	})
}
