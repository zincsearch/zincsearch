## Test ES API with Endpoint /es/:target/_bulk

1. ### Test Create Document with no existing Index
    This unit test will create document with not existing Index:
    1. `POST` call to the endpoint `/es/notExistIndex/_bulk` with bulkData in body.
    2. The **bulkData** in ndjson form defined in init.go file.
    3. Using Replaceall to replace bulkData `_index` value with an empty string. 
    4. Send a request to the server and response would be recorded.
    5. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK and bulk document is indexed with Index notExistIndex.

2. ### Test Create Document with existing Index
    This unit test will create document with existing Index:
    1. Create Index request to check that Index already exists by `PUT` call to the endpoint `/api/index` with `{"name": "` + indexName + `", "storage_type": "disk"}` in  body.
    2. Using assert.Equal will test if the response code is equal to http.StatusBadRequest (i.e 400), that means Index already exists.
    3. Umarshal the response recieved and using assert.Equal will check that error string is equal to index `["+indexName+"] already exists`.
    4. `POST` call to the endpoint `/es/"+indexName+"/_bulk` with bulkData in body.
    5. Send a request to the server and response would be recorded.
    6.  Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK and bulk document is indexed.

3. ### Test Bulk with error Input
   This unit test will index documents with error input:
   1. `POST` call to the endpoint `/es/"+indexName+"/_bulk` with `{"index": {}}` in body .
   2. Send a request to the server and response would be recorded.
   3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.  
