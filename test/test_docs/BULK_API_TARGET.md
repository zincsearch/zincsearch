## Test Bulk Document with Enpoint /api/:target/_bulk

1. ### Test Create Document with not existing Index
   This unit test will index documents with not existing index:
   1. `POST` call to the endpoint `/api/notExistIndex/_bulk` with bulkData in body.
   2. Send a request to the server and response would be recorded.
   3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.
2. ### Test Create Document with Existing Index
   This unit test will index documents with existing index:
   1. `PUT` call to the endpoint `/api/index/` with `{"name": indexName", "storage_type": "disk"}` in body indexName defined in init.go.
   2. Send a request to the server and response would be recorded and using assert.Equal will test if the response code is equal to http.StatusBadRequest (i.e 400).
   3. Unmarshal the response and check there is no error using assert.Error method and using assert.Equal to check that the response is equal to  **index ["+indexName+"] already exists**.
   4. `POST` call to the endpoint `/api/indexName/_bulk` with bulkData in body.
   5. Send a request to the server and response would be recorded.
   6. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.
3. ### Test Bulk Document with error input
   This unit test will index bulk document with error input:
   1. `POST` call to the endpoint `/api/indexName/_bulk` with `{"index": {}}` in body.
   2. Send a request to the server and response would be recorded.
   3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.