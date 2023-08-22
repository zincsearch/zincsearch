## Testing  User API
1. ### Testing Login with username and password
   This unit test will check the login using username and password:
   1. `POST` call to endpoint `/api/login/` with username and password in the body.
   2. Send a request to the server and response would be recorded.
   3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK. 
   4. Unmarshal the json received in the response body and will check if there is no error using assert.NoError and using assert.True to check the validated field in the response body is true.

2. ### Testing Login with bad username or password
   This unit test will check the login using bad username and password:
   1. `POST` call to endpoint `/api/login/` with username and password {xxx} in the body.
   2.  Send a request to the server and response would be recorded.
   3. Using assert.Equal will test if the response code is equal to http.StatusUnauthorized (i.e 401),that means test is OK.
   4. Unmarshal the json received in the response body and will check if there is no error using assert.NoError and using assert.False to check the validated field in the response body is false.

3. ### Testing Create User API
   This unit test will check the API that will create a new user:
   1. Creating user with payload (i.e _id: username , name: username , password : password , role: admin).
   2. `PUT` call to the endpoint `/api/user`.
   3. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK. 
   4. Test new user login using the above username and password in body to the endpoint `/api/login`.
   5. Response will be recorded and using assert.Equal will test if the response code is equal to http.StatusUnauthorized (i.e 401),that means test is OK.
   6. Unmarshal the json received in the response body and will check if there is no error using assert.NoError and using assert.True to check the validated field in the response body is true.

4. ### Test Update User
   This unit test will update the existing user : 
   1. Update the existing username by using paylod (i.e _id: username , name: username-updated , password : password , role: admin).
   2. `PUT` call to the endpoint `/api/user`.
   3.  Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.
   4. Get the user name using auth.Getuser method and using assert.Equal to match that the user Name has been successfully updated.

5. ### Test Create User with error input
   This unit test will create user with error input: 
   1. `PUT` call to the endpoint `/api/user` with `xxx` in body.
   2. Request to the server and response will be recorded.
   3. Using assert.Equal will test if the response code is equal to http.StatusBadRequest (i.e 400),that means test is OK.

6. ### Test Delete User with existing UserID
   This unit test will delete the existing userid:
   1. `DELETE` call to the endpoint `/api/user/username` with empty body.
   2. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK and user is deleted.
   3. Using auth.Getuser method to check user existence and using assert.False to verify that the user doest not exists.
   
7. ### Test Delete User with no existing UserID
    This unit test will test delete the user with no existing userid:
    1. `DELETE` call to the endpoint `/api/user/userNotExist` with empty body.
    2. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK ,will get a null response.

8. ### Test GET User 
    This unit test will enlist all the users:
    1. `GET` call to the endpoint `/api/user` with empty body.
    2. Using assert.Equal will test if the response code is equal to http.StatusOK (i.e 200),that means test is OK.
    3. Unmarshal the json received in the response body and will check if there is no error using assert.NoError and using assert.GreaterorEqual to check the length of data (i.e number of users) are greater than or equal to 1 , finally using assert.Equal to verify that the first user is admin.
