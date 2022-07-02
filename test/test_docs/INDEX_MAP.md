## Test Index Mapping with Endpoint /api/:target/_mapping

1. ### Test Update mapping for Index
    This unit test will update mappings for Index:
    1. `PUT` call to the endpoint `/api/"+indexName+"-mapping/_mapping` with `{"properties":{"Athlete": {"type": "keyword"}}}` in body.
    2. Send the request to server and response would be recorded.
    3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.
    4. Using assert.Contains to test if the response body contains a string `"message":"ok"` , that means test is OK.

2. ### Test GET mappings from Index
    This unit test will GET mapping from Index:
    1. `GET` call to the endpoint `/api/"+indexName+"/_mapping`.
    2. Send the request to server and response would be recorded.
    3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.NotNil to check if the data[indexName] is not nil. 
    5. Creating map from data[indexName] and using assert.NotNil to check if the v["mappings] is not nil, that means the test is OK.