## Test Create Document with Endpoint /api/:target/_doc
1. ### Test Create Document with not existing Index
    This unit test will create document with not existing index:
    1. `PUT` call to the endpoint `/api/notExistIndex/_doc` with indexData in body.
    2. The **indexData** in json form defined in init.go file.
    3. Send a request to the server and response would be recorded.
    4.  Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK and document is indexed.

2. ### Test Create Document with existing Index
    This unit test will create document with exisiting index:
    1. `PUT` call to the endpoint `/api/"+indexName+"/_doc` with indexData in body.
    2. Send a request to the server and response would be recorded.
    3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK and document is indexed.

3. ### Test Create Document with existing Index and not existing Document ID
   This unit test will index document with existing Index and not existing ID:
    1. `PUT` call to the endpoint `/api/"+indexName+"/_doc` with indexData in body.
    2. Send a request to the server and response would be recorded.
    3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK and document is indexed.
    4. Unmarshal the json recieved in the response body and will check if there is no error using assert.NoError and using assert.NotEqual to check that "id" is not an empty value.
    5. Storing the value of id recieved in the json response in a variable _id to be used in the next test.
    
4. ### Test Update Document with existing Index and existing Document ID
   This unit test will update document with existing Index and exisitng ID:
   1. `PUT` call to the endpoint `/api/"+indexName+"/_doc` with indexData in body.
   2. Updating indexData by adding a field _id with the value that was stored in the last test using string.Replace.
   3. Send a request to the server and response would be recorded.
   4. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK and document is updated with existing Document ID.
5. ### Test Create Document with error Input
   THis unit test will create document with error Input:
   1. `PUT` call to the endpoint `/api/"+indexName+"/_doc` with `data` in the body.
   2. Send a request to the server and response would be recorded.
   3. Using assert.Equal will test if the response code is equal to http.StatusBadRequest (i.e 400),that means test is OK.
