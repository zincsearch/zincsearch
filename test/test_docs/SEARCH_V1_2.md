## Test Search Document with Aggregation

1. ### Test Search Document with Term Aggregation
    This unit test will search document with term aggregation:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "matchall", "aggs": {"my-agg": {"agg_type": "terms","field": "Sport"}}}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check length of aggreagation data is greater or equal to 1.

2. ### Test Search Document with Metric Aggregation
    This unit test will search document with metric aggregation:
    1. `POST` call to the endpoint `/api/"+indexName+"/_search ` with `{"search_type": "matchall", "aggs": {"my-agg-max": {"agg_type": "max","field": "Year"},"my-agg-min": {"agg_type": "min","field": "Year"},"my-agg-avg": {"agg_type": "avg","field": "Year"}}}` in body.
    2. Send the request to the server and response would be recorded.
    3. Using assert.Equal to check the response code is euqal to http.StatusOk (i.e 200) , that means test is OK.
    4. Unmarshal the json response and using assert.NoError to check there is no error and using assert.GreaterorEqual to check length of aggreagation data is greater or equal to 1.    