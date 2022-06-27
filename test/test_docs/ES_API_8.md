## Test ES API with Endpoint /es/:target/_update/:id 

1. ### Test Update Document with no existing Index
    This unit test will update document with no existing Index:
    1. `POST` call to the endpoint `/es/notExistIndexCreate/_update/1111` with indexData in body.
    2. indexData is in json form and defined in init.go file.
    3. Send a request to the server and response would be recorded.
    4. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK and document is updated in new Index.

2. ### Test Update Document with existing Index
    This unit test will update document with existing Index:
    1. `POST` call to the endpoint `/es/"+indexName+"/_update/1111` with indexData in body.
    2. Send a request to the server and response would be recorded.
    3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.

3. ### Test Update Document with existing Index and not existing ID
    This unit test will update document with existing Index and not existing ID:
    1. `POST` call to the endpoint `/es/"+indexName+"/_update/notexistCreate` with indexData in body.
    2. Send a request to the server and response would be recorded.
    3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.

4. ### Test Update Document with existing Index and existing ID
    This unit test will update document with existing Index and existing ID:
    1. `POST` call to the endpoint `/es/"+indexName+"/_update/1111` with indexData in body.
    2. Send a request to the server and response would be recorded.
    3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.  

5. ### Test Update Document with error Input
    This unit test will update document with error Input:
    1. `POST` call to the endpoint `/es/"+indexName+"/_update/1111` with `xxx` in body.
    2. Send a request to the server and response would be recorded.
    3. Using assert.Equal will test if the response code is equal to http.StatusBadRequest (i.e 400),that means test is OK. 