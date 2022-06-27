## Test Search Document with Endpoint /es/:target/_search

### Init Data for Search 
This function will Index indexData to be used by search API's unit test,`PUT` call to the endpoint `/api/"+indexName+"/_doc` with indexData in body. 

1. ### Test Search Document with no existing Index
    This unit test will search documents with no existing Index:
    1. `POST` call to the endpoint `/es/notExistSearch/_search` with `{}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusBadRequest (i.e 400) , that means test is OK.

2. ### Test Search Document with existing Index
    This unit test will search documents with existing Index:
    1. `POST` call to the endpoint `/es/"+indexName+"/_search` with `{"query": {"match_all":{}}, "size":10}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.

3. ### Test Search Document with no existing Term
    This unit test will search document with no existing Term:
    1. `POST` call to the endpoint `/es/"+indexName+"/_search` with `{"query": {"match": {"_all": "xxxx"}}, "size":10}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.Equal to check that total hit values are euqal to zero.

4. ### Test Search Document with existing Term
    This unit test will search document with existing Term:
    1. `POST` call to the endpoint `/es/"+indexName+"/_search` with `{"query": {"match": {"_all": "DEMTSCHENKO"}}, "size":10}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.

5. ### Test Search Document with search type match all
    This unit test will search document with search type match all:
    1. `POST` call to the endpoint `/es/"+indexName+"/_search` with `{"query": {"match_all": {}}, "size":10}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.  
   
6. ### Test Search Document with search type wildcard
    This unit test will search document with search type wildcard:
    1. `POST` call to the endpoint `/es/"+indexName+"/_search` with `{"query": {"wildcard": {"_all": "dem*"}}, "size":10}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.     

7. ### Test Search Document with search type fuzzy
    This unit test will search document with search type fuzzy:
    1. `POST` call to the endpoint `/es/"+indexName+"/_search` with `{"query": {"fuzzy": {"Athlete": "demtschenk"}}, "size":10}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.

8. ### Test Search Document with search type term
    This unit test will search document with search type term:
    1. `POST` call to the endpoint `/es/"+indexName+"/_search` with `{"query": {"term": {"City": "turin"}}, "size":10}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.    

9. ### Test Search Document with search type data range
    This unit test will search document with search type data range:   
    1. `POST` call to the endpoint `/es/"+indexName+"/_search` with `{"query": {"range": {"@timestamp": { "gte": "%s", "lt": "%s"}}}, "size":10}` in body with `lt` as now() and `gte` 24 hours before the `lt`.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1. 

10. ### Test Search Document with search type match
    This unit test will search document with search type match:
    1. `POST` call to the endpoint `/es/"+indexName+"/_search` with `{"query": {"match": {"_all": "DEMTSCHENKO"}}, "size":10}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1. 

11. ### Test Search Document with search type match phrase
    This unit test will search document with search type match phrase:
    1. `POST` call to the endpoint `/es/"+indexName+"/_search` with `{"query": {"match_phrase": {"_all": "DEMTSCHENKO"}}, "size":10}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1. 


12. ### Test Search Document with search type prefix
    This unit test will search document with search type prefix:
    1. `POST` call to the endpoint `/es/"+indexName+"/_search` with `{"query": {"prefix": {"_all": "dem"}}, "size":10}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.

13. ### Test Search Document with search type query string
    This unit test will search document with search type prefix:
    1. `POST` call to the endpoint `/es/"+indexName+"/_search` with `{"query": {"query_string": {"query": "DEMTSCHENKO"}}, "size":10}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.