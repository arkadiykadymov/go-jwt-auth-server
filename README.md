# jwt-auth-server

    This is a simple authentication server.

## Installation
    This application are deployed in heroku - https://pacific-everglades-72035.herokuapp.com

## Description

  The API supports 5 endpoints.
  To get started, you need to register a user on the /register path by sending a username and password.
  After registration, you will receive a response in the form

  ```json
  {
      "data": {
          "message": "User created successfully!",
          "GUID": "93fdf0e5-00e4-4c66-8d41-369dd35cbd4a",
          "status": 201
      }
  }
  ```
  The resulting GUID can be used in further requests

 # DB
  In this project, a Mongo DB database deployed on Mongo DB Atlas is used.
  DB uri - mongodb+srv://test:1234@cluster0.5bbqz.mongodb.net/test?retryWrites=true&w=majority
  DB type - Replica set
  DB username - test
  DB passwrod - 1234
  DB name - test
  DB collection - users


# API

### HTTP

#### `GET /`

```bash
curl --request GET \
  --url http://https://pacific-everglades-72035.herokuapp.com/ \
```

```json
{"data": "Hello world!"}
```

Returns "Hello world!"


#### `POST /register`

Endpoint for registretion user with specific username and password.

```bash
curl --request POST \
  --url http://https://pacific-everglades-72035.herokuapp.com/register \
  --header 'content-type: application/json' \
  --data '{ 
    "username": "username",
    "password": "password"
}'
```

```json
{
    "data": {
        "message": "User created successfully!",
        "GUID": "93fdf0e5-00e4-4c66-8d41-369dd35cbd4a",
        "status": 201
    }
}
```

#### `POST /login/:uuid`

Authenticate and login user by user GUID

```bash
curl --request POST \
  --url https://pacific-everglades-72035.herokuapp.com/login/93fdf0e5-00e4-4c66-8d41-369dd35cbd4a\
  --header 'content-type: application/json' \
```
When the login succeeds, an access token is returned

```json
{
    "access_token":    "eyJhbGciOiJIUzI1NiIsImtpZCI6InNpZ25pbl8xIiwidHlwIjoiSldUIn0.eyJVc2VybmFtZSI6ImVyaWMiLCJleHAiOjE1NzA3NjI5NzksInN1YiI6IjVkOTNlMTFjNmY4Zjk4YzlmYjI0ZGU0NiJ9.n70EAaiY6rbH1QzpoUJhx3hER4odW8FuN2wYG1sgH7g",
"refresh_token": "eyJhbGciOiJIUzI1NiIsImtpZCI6InNpZ25pbl8yIiwidHlwIjoiSldUIn0.eyJleHAiOjE1NzA3NjM1NzksInN1YiI6IjVkOTNlMTFjNmY4Zjk4YzlmYjI0ZGU0NiJ9.zwGB1340IVMLjMf_UnFC_rEeNdD131OGPcg_S0ea8DE",
"status": 200
    }
```


#### `POST /refresh-token`

Request new access_token by using the `refresh_token`

```bash
curl --request POST \
  --url http://https://pacific-everglades-72035.herokuapp.com/refresh-token \
  --header 'content-type: application/json' \
  --data '{
    "refresh_token" : "eyJhbGciOiJIUzI1NiIsImtpZCI6InNpZ25pbl8yIiwidHlwIjoiSldUIn0.eyJleHAiOjE1NzA3NjM1NzksInN1YiI6IjVkOTNlMTFjNmY4Zjk4YzlmYjI0ZGU0NiJ9.zwGB1340IVMLjMf_UnFC_rEeNdD131OGPcg_S0ea8DE"
}'
```
When the token is valid, a new access_token is returned

```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIsImtpZCI6InNpZ25pbl8xIiwidHlwIjoiSldUIn0.eyJVc2VybmFtZSI6ImVyaWMiLCJleHAiOjE1NzA3NjMyMjksInN1YiI6IjVkOTNlMTFjNmY4Zjk4YzlmYjI0ZGU0NiJ9.wrWsDNor28aWv6huKUHAuVyROGAXqjO5luPfa5K5NQI",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsImtpZCI6InNpZ25pbl8yIiwidHlwIjoiSldUIn0.eyJleHAiOjE1NzA3NjM1NzksInN1YiI6IjVkOTNlMTFjNmY4Zjk4YzlmYjI0ZGU0NiJ9.zwGB1340IVMLjMf_UnFC_rEeNdD131OGPcg_S0ea8DE",
    "status": 200
}
```


#### `POST /delete-refresh-token`

Request are deleteing exact refresh token

```bash
curl --request POST \
  --url http://https://pacific-everglades-72035.herokuapp.com/delete-refresh-token \
  --header 'content-type: application/json' \
  --data '{
    "refresh_token" : "eyJhbGciOiJIUzI1NiIsImtpZCI6InNpZ25pbl8yIiwidHlwIjoiSldUIn0.eyJleHAiOjE1NzA3NjM1NzksInN1YiI6IjVkOTNlMTFjNmY4Zjk4YzlmYjI0ZGU0NiJ9.zwGB1340IVMLjMf_UnFC_rEeNdD131OGPcg_S0ea8DE"
}'
```

When the token is valid and is in the database, it is deleted

```json
{
    "message": "Token deleted",
    "status": 200
}
```

#### `POST /delete-refresh-tokens/:uuid`

Request are deleteing refresh token by user id


```bash
curl --request POST \
  --url http://https://pacific-everglades-72035.herokuapp.com/delete-refresh-tokens/93fdf0e5-00e4-4c66-8d41-369dd35cbd4a \
  --header 'content-type: application/json' \
```

Response 

```json
{
    "message": "2 tokens deleted",
    "status": 200
}
```