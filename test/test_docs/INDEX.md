## Test Index with Endpoint /api/index

1. ### Test Create Index with Payload
    This unit test will create an Index with Payload:
    1. `PUT` call to the endpoint `/api/index` with `{"name":"newindex","storage_type":"disk"}` in body.
    2. Send a request to the server and response would be recorded.
    3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.
    4. Using assert.Equal will test the response message is euqal to the `{"index":"newindex","message":"ok","storage_type":"disk"}` , that means test is OK and index is created successfully.  

2. ### Test Create Index with error Input
    This unit test will create an Index with error Input:
    1. `PUT` call to the endpoint `/api/index` with `{"name":"","storage_type":"disk"}` in body.
    2. Send a request to the server and response would be recorded.
    3. Using assert.Equal will test if the response code is equal to http.StatusBadRequest (i.e 400),that means test is OK.

3. ### Test GET Index 
    This unit test will get all the Indexes:
    1. `GET` call to the endpoint `/api/index`.
    2. Send a request to the server and response would be recorded.
    3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.
    4. Unmarshal the json response from the server and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that the length of the response is greater or equal to 1.

4. ### Test Delete with existing Index
    This unit test will delete Index with existing Index:
    1. `DELETE` call to the endpoint `/api/index/newindex`.
    2. Send a request to the server and response would be recorded.
    3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK and index is deleted.

5. ### Test Delete with not existing Index 
    This unit test will delete Index with no existing Index:
    1. `DELETE` call to the endpoint `/api/index/newindex`. 
    2. As Index was already deleted in the above test , request is sent to the server and response is recorded.
    3. Using assert.Equal will test if the response code is equal to http.StatusBadRequest (i.e 400),that means test is OK. 