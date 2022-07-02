## Base API Test
1. ### Test Base API 
   This unit test will check the Zinc base API:
   1. `GET ` call to the endpoint `/`.
   2. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.
2. ### Test Version API
   This unit test will check the Zinc version:
   1. `GET` call to the endpoint `/version`.
   2. Send a request to the server and response would be recorded.
   3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.
   4. Unmarshal the json recieved in the response body and will check if there is no error using assert.NoError and using assert.True to check if the data has the field "Version".

3. ### Test Health API
   This unit test will check the Zinc Server health:
   1. `GET` call to the endpoint `/healthz`.
   2. Send a request to the server and response would be recorded.
   3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.
   4. Unmarshal the json recieved in the response body and will check if there is no error using assert.NoError and using assert.True to check if the data has the field "Status"and using asset.Equal to check if the status is 'ok'.

4. ### Test User Interface API
   This unit test will check the Zinc UI API:
    1. `GET` call to the endpoint `/ui`.
   2. Send a request to the server and response would be recorded.
   3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.