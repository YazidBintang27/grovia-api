# API Specification - Grovia Backend

# Authentication

## POST api/auth/login

### Request
 - #### Method: `POST`
 - #### Content-type: `application/json`
 - #### Endpoint: `/auth/login`
 - #### Request Body: 
 ```json
 {
    "phoneNumber": "08771265393",
    "password": "passwordexample33"
 }
 ```
### Response
 - #### Status: `200 OK`
 - #### Body: 
 ```json
 {
    "success": true,
    "message": "Login Success",
    "data": {
      "token": {
         "accessToken": "wou8h2yr8038e2940uq0jp",
         "refreshToken": "wou8h2yr8038e2940uq0jp"
      }
    },
    "error": null
 }
 ```

 - #### Status: `400 Bad Request`
 - #### Body: 
 ```json
 {
    "success": false,
    "message": "Login Failed",
    "data": null,
    "error": {
      "code": "USER_NOT_FOUND",
      "message": "User not found"
    }
 }
 ```

 ## POST api/auth/reset-password

 ### Request
 - #### Method: `POST`
 - #### Content-type: `application/json`
 - #### Endpoint: `/auth/reset-password`
 - #### Request Body: 
 ```json
 {
    "firebaseToken": "wou8h2yr8038e2940uq0jp",
    "phoneNumber": "08771265393",
    "password": "passwordexample33",
    "confirmPassword": "passwordexample33"
 }
 ```
### Response
 - #### Status: `200 OK`
 - #### Body: 
 ```json
 {
    "success": true,
    "message": "Reset password Success",
    "data": {
      "phoneNumber": "08771265393",
      "updatedAt": "2025-08-03T00:00:00Z"
    },
    "error": null
 }
 ```

 - #### Status: `400 Bad Request`
 - #### Body: 
 ```json
 {
    "success": false,
    "message": "Reset password Failed",
    "data": null,
    "error": {
      "code": "PASSWORD_MISMATCH",
      "message": "New password and confirm password do not match"
    }
 }
 ```

# Users

## GET api/users/current

### Request
 - #### Method: `GET`
 - #### Header: `Bearer <token>`
 - #### Content-type: `application/json`
 - #### Endpoint: `/users/current`

### Response
 - #### Status: `200 OK`
 - #### Body: 
 ```json
 {
    "success": true,
    "message": "Get user Success",
    "data": {
      "name": "yazid",
      "phoneNumber": "087712354798",
      "address": "Ds. Pangadegan RT 002 RW 003",
      "nik": "3603312711040004",
      "profilePicture": "https://storage.googleapis.com/tesdemas/users/dR9RzwI2hrVE4DdXmZKkg.png",
      "role": "Kepala Posyandu",
      "createdBy": "Admin",
      "createdAt": "2025-08-03T00:00:00Z",
      "updatedAt": "2025-08-03T00:00:00Z"
    },
    "error": null
 }
 ```

 - #### Status: `401 Unauthorized`
 - #### Body: 
 ```json
 {
    "success": false,
    "message": "Get user Failed",
    "data": null,
    "error": {
      "code": "UNAUTHORIZED",
      "message": "Invalid Token"
    }
 }
 ```

 ## PUT api/users/current

### Request
 - #### Method: `PUT`
 - #### Header: `Bearer <token>`
 - #### Content-type: `application/json`
 - #### Endpoint: `/users/current`
 - #### Body: 

 ```json
 {
   "name": "yazid",
   "phoneNumber": "087712354798",
   "address": "Ds. Pangadegan RT 002 RW 003",
   "nik": "3603312711040004",
   "profilePicture": "https://storage.googleapis.com/tesdemas/users/dR9RzwI2hrVE4DdXmZKkg.png",
 }
 ```

### Response
 - #### Status: `200 OK`
 - #### Body: 
 ```json
 {
    "success": true,
    "message": "Get user Success",
    "data": {
      "name": "yazid",
      "phoneNumber": "087712354798",
      "address": "Ds. Pangadegan RT 002 RW 003",
      "nik": "3603312711040004",
      "profilePicture": "https://storage.googleapis.com/tesdemas/users/dR9RzwI2hrVE4DdXmZKkg.png",
      "role": "Kepala Posyandu",
      "createdBy": "Admin",
      "createdAt": "2025-08-03T00:00:00Z",
      "updatedAt": "2025-08-03T00:00:00Z"
    },
    "error": null
 }
 ```

 - #### Status: `401 Unauthorized`
 - #### Body: 
 ```json
 {
    "success": false,
    "message": "Get user Failed",
    "data": null,
    "error": {
      "code": "UNAUTHORIZED",
      "message": "Invalid Token"
    }
 }
 ```

 ## DELETE api/users/current

### Request
 - #### Method: `DELETE`
 - #### Header: `Bearer <token>`
 - #### Content-type: `application/json`
 - #### Endpoint: `/users/current`

### Response
 - #### Status: `200 OK`
 - #### Body: 
 ```json
 {
    "success": true,
    "message": "Delete user Success",
    "error": null
 }
 ```

 - #### Status: `401 Unauthorized`
 - #### Body: 
 ```json
 {
    "success": false,
    "message": "Delete user Failed",
    "data": null,
    "error": {
      "code": "UNAUTHORIZED",
      "message": "Invalid Token"
    }
 }
 ```

 ## GET api/users

### Request
 - #### Method: `GET`
 - #### Header: `Bearer <token>`
 - #### Content-type: `application/json`
 - #### Endpoint: `/users`

### Response
 - #### Status: `200 OK`
 - #### Body: 
 ```json
 {
    "success": true,
    "message": "Get user Success",
    "data": [
      {
         "name": "yazid",
         "phoneNumber": "087712354798",
         "address": "Ds. Pangadegan RT 002 RW 003",
         "nik": "3603312711040004",
         "profilePicture": "https://storage.googleapis.com/tesdemas/users/dR9RzwI2hrVE4DdXmZKkg.png",
         "role": "Kepala Posyandu",
         "createdBy": "Admin",
         "createdAt": "2025-08-03T00:00:00Z",
         "updatedAt": "2025-08-03T00:00:00Z"
      },
      {
         "name": "yazid",
         "phoneNumber": "087712354798",
         "address": "Ds. Pangadegan RT 002 RW 003",
         "nik": "3603312711040004",
         "profilePicture": "https://storage.googleapis.com/tesdemas/users/dR9RzwI2hrVE4DdXmZKkg.png",
         "role": "Kepala Posyandu",
         "createdBy": "Admin",
         "createdAt": "2025-08-03T00:00:00Z",
         "updatedAt": "2025-08-03T00:00:00Z"
      }
    ],
    "error": null
 }
 ```

 - #### Status: `401 Unauthorized`
 - #### Body: 
 ```json
 {
    "success": false,
    "message": "Get user Failed",
    "data": null,
    "error": {
      "code": "UNAUTHORIZED",
      "message": "Invalid Token"
    }
 }
 ```

  ## GET api/users/:id

### Request
 - #### Method: `GET`
 - #### Header: `Bearer <token>`
 - #### Content-type: `application/json`
 - #### Endpoint: `/users/:id`

### Response
 - #### Status: `200 OK`
 - #### Body: 
 ```json
 {
    "success": true,
    "message": "Get user Success",
    "data": {
      "name": "yazid",
      "phoneNumber": "087712354798",
      "address": "Ds. Pangadegan RT 002 RW 003",
      "nik": "3603312711040004",
      "profilePicture": "https://storage.googleapis.com/tesdemas/users/dR9RzwI2hrVE4DdXmZKkg.png",
      "role": "Kepala Posyandu",
      "createdBy": "Admin",
      "createdAt": "2025-08-03T00:00:00Z",
      "updatedAt": "2025-08-03T00:00:00Z"
    }
    "error": null
 }
 ```

 - #### Status: `401 Unauthorized`
 - #### Body: 
 ```json
 {
    "success": false,
    "message": "Get user Failed",
    "data": null,
    "error": {
      "code": "UNAUTHORIZED",
      "message": "Invalid Token"
    }
 }
 ```

  ## POST api/users

### Request
 - #### Method: `POST`
 - #### Header: `Bearer <token>`
 - #### Content-type: `application/json`
 - #### Endpoint: `/users`
 - #### Body:

 ```json
 {
   "name": "yazid",
   "phoneNumber": "087712354798",
   "address": "Ds. Pangadegan RT 002 RW 003",
   "nik": "3603312711040004",
   "profilePicture": "https://storage.googleapis.com/tesdemas/users/dR9RzwI2hrVE4DdXmZKkg.png",
   "role": "Kepala Posyandu"
 }
 ```

### Response
 - #### Status: `200 OK`
 - #### Body: 
 ```json
 {
    "success": true,
    "message": "Create user Success",
    "data": {
      "name": "yazid",
      "phoneNumber": "087712354798",
      "address": "Ds. Pangadegan RT 002 RW 003",
      "nik": "3603312711040004",
      "profilePicture": "https://storage.googleapis.com/tesdemas/users/dR9RzwI2hrVE4DdXmZKkg.png",
      "role": "Kepala Posyandu",
      "createdBy": "Admin",
      "createdAt": "2025-08-03T00:00:00Z",
      "updatedAt": "2025-08-03T00:00:00Z"
    },
    "error": null
 }
 ```

 - #### Status: `401 Unauthorized`
 - #### Body: 
 ```json
 {
    "success": false,
    "message": "Create user Failed",
    "data": null,
    "error": {
      "code": "UNAUTHORIZED",
      "message": "Invalid Token"
    }
 }
 ```

 ## PUT api/users/:id

### Request
 - #### Method: `PUT`
 - #### Header: `Bearer <token>`
 - #### Content-type: `application/json`
 - #### Endpoint: `/users/:id`
 - #### Body: 

 ```json
 {
   "name": "yazid",
   "phoneNumber": "087712354798",
   "address": "Ds. Pangadegan RT 002 RW 003",
   "nik": "3603312711040004",
   "profilePicture": "https://storage.googleapis.com/tesdemas/users/dR9RzwI2hrVE4DdXmZKkg.png",
 }
 ```

### Response
 - #### Status: `200 OK`
 - #### Body: 
 ```json
 {
    "success": true,
    "message": "Get user Success",
    "data": {
      "name": "yazid",
      "phoneNumber": "087712354798",
      "address": "Ds. Pangadegan RT 002 RW 003",
      "nik": "3603312711040004",
      "profilePicture": "https://storage.googleapis.com/tesdemas/users/dR9RzwI2hrVE4DdXmZKkg.png",
      "role": "Kepala Posyandu",
      "createdBy": "Admin",
      "createdAt": "2025-08-03T00:00:00Z",
      "updatedAt": "2025-08-03T00:00:00Z"
    },
    "error": null
 }
 ```

 - #### Status: `401 Unauthorized`
 - #### Body: 
 ```json
 {
    "success": false,
    "message": "Get user Failed",
    "data": null,
    "error": {
      "code": "UNAUTHORIZED",
      "message": "Invalid Token"
    }
 }
 ```

 ## DELETE api/users/:id

### Request
 - #### Method: `DELETE`
 - #### Header: `Bearer <token>`
 - #### Content-type: `application/json`
 - #### Endpoint: `/users/:id`

### Response
 - #### Status: `200 OK`
 - #### Body: 
 ```json
 {
    "success": true,
    "message": "Delete user Success",
    "error": null
 }
 ```

 - #### Status: `401 Unauthorized`
 - #### Body: 
 ```json
 {
    "success": false,
    "message": "Delete user Failed",
    "data": null,
    "error": {
      "code": "UNAUTHORIZED",
      "message": "Invalid Token"
    }
 }
 ```