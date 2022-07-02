## Test Delete Document with Endpoint /api/:target/_doc/:id

1. ### Test Delete Document with not existing Index
    This unit test will delete document with not existing Index:
    1. `DELETE` call to the endpoint `/api/notExistIndexDelete/_doc/1111`.
    2. Using assert.Equal will test if the response code is equal to http.StatusBadRequest (i.e 400),that means test is OK.

2. ### Test Delete Document with existing Index and not existing ID
    This unit test will delete document with existing Index and not existing ID:
    1. `DELETE` call to the endpoint `/api/"+indexName+"/_doc/notexist`.
    2. Using assert.Equal will test if the response code is equal to http.OK (i.e 200),that means test is OK and document is deleted.

3. ### Test Delete Document with existing Index and existing ID
    This unit test will delete document with existing Index and existing ID:
    1. `DELETE`  call to the endpoint `/api/"+indexName+"/_doc/1111`.
    2. Using assert.Equal will test if the response code is equal to http.OK (i.e 200),that means test is OK and document is deleted.

