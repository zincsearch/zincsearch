## Test Search with Endpoint /api/:target/_search

### Init Data for Search 
This function will Index indexData to be used by search API's unit test,`PUT` call to the endpoint `/api/"+indexName+"/_doc` with indexData in body. 

1. ### Test Search Document with no existing Index
    This unit test will search documents with no existing Index:
    1. `POST` call to the endpoint `/api/notExistSearch/_search` with `{}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusBadRequest (i.e 400) , that means test is OK.

2. ### Test Search Document with existing Index
    This unit test will search documents with existing Index:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search` with `{"search_type": "alldocuments"}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.

3. ### Test Search Document with no existing Term
    This unit test will search document with no existing Term:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "match", "query": {"term": "hello"}}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.Equal to check that total hit values are euqal to zero.

4. ### Test Search Document with existing Term
    This unit test will search document with existing Term:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "match", "query": {"term": "DEMTSCHENKO"}}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.

5. ### Test Search Document with search type all documents
    This unit test will search document with search type all documents:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "alldocuments", "query": {}}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.    

6. ### Test Search Document with search type wildcard
    This unit test will search document with search type wildcard:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "wildcard", "query": {"term": "dem*"}}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.     

7. ### Test Search Document with search type fuzzy
    This unit test will search document with search type fuzzy:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "fuzzy", "query": {"term": "demtschenk"}}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.

8. ### Test Search Document with search type term
    This unit test will search document with search type term:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "term","query": {"term": "turin","field":"City"}}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.    

9. ### Test Search Document with search type data range
    This unit test will search document with search type data range:   
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "daterange","query": {"start_time": "%s","end_time": "%s"}}` in body with `end_time` as now() and `start_time` 24 hours before the `end_time`.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1. 

10. ### Test Search Document with search type match all
    This unit test will search document with search type match all:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "matchall", "query": {"term": "demtschenk"}}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.  


11. ### Test Search Document with search type match
    This unit test will search document with search type match:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "match", "query": {"term": "DEMTSCHENKO"}}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1. 

12. ### Test Search Document with search type match phrase
    This unit test will search document with search type match phrase:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "matchphrase", "query": {"term": "DEMTSCHENKO"}}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1. 

13. ### Test Search Document with search type multiphrase
    This unit test will search document with search type multi phrase:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "multiphrase","query": {"terms": [["demtschenko"],["albert"]]}}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.

14. ### Test Search Document with search type prefix
    This unit test will search document with search type prefix:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "prefix", "query": {"term": "dem"}}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.

15. ### Test Search Document with search type query string
    This unit test will search document with search type prefix:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "querystring", "query": {"term": "DEMTSCHENKO"}}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check that total hit values are greater or equal to 1.

