## Testing  Auth API
1. ### Testing Authentication with Authorization
    This unit test will check authentication with user authorization:
    1. `GET` call to endpoint `/api/index/`.
    2. Setting Basic Auth that will set the request's Authorization header to use HTTP Basic Authentication with the provided username and password.
    3. Send a request to the server and response would be recorded. 
    4. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK and user has the authorization. 

2. ### Testing Authorization with Error Password
    This unit test will check the user authorization with error password:
    1. `GET` call to endpoint `/api/index/`.
    2. Setting Basic Auth that will set the request's Authorization header to use HTTP Basic Authentication with the provided username and password {'xxx'}.
    3. Send a request to the server and response would be recorded.
    4. Using assert.Equal will test if the response code is equal to http.StatusUnauthorized (i.e 401),that means test is OK and authorization with error password is working fine.

3. ### Testing Authentication without Authorization 
   This unit test will check the authentication without user authorization:
   1. `GET` call to endpoint `/api/index/`.
   2. Send a request to the server and response would be recorded.
   3. Using assert.Equal will test if the response code is equal to http.StatusUnauthorized (i.e 401),that means test is OK and authentication without authorization is working fine.
