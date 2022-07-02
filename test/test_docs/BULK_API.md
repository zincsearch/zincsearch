## Test Bulk Document with Endpoint /api/_bulk

1. ### Test Bulk Document Index 
   This unit test will index documents in bulk:
   1. `POST` call to the endpoint `/api/_bulk` with bulkData in body.
   2. The **bulkData** will be in ndjson form defined in init.go.
   3. Send a request to the server and response would be recorded.
   4. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK and documents are indexed.
2. ### Test Bulk Document Delete
   This unit test will delete documents in bulk:
   1. `POST` call to the endpoint `/api/_bulk` with bulkDataWithDelete in body.
   2. The **bulkDataWithDelete** will be in ndjson form defined in init.go.
   3. Send a request to the server and response would be recorded.
   4. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK and documents are deleted.
3. ### Test Bulk Document with Error Input
   This unit test will index documents with error input:
   1. `POST` call to the endpoint `/api/_bulk` with `{"index": {}}` in body .
   2. Send a request to the server and response would be recorded.
   3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.  