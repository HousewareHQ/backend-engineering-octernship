# Authentication and Authorization service in GoLang using JWT

## Tech Stack
- GoLang Framework: [`Gin`](https://gin-gonic.com/)
- [MongoDB](https://www.mongodb.com/)
- [Mongo Go Driver](https://www.mongodb.com/docs/drivers/go/current/)

### Reason for choosing Tech Stack: 
- Gin is a great choice for building web applications in Go, especially scalability. Gin has a small footprint, making it easy to build and deploy applications quickly. Gin has a simple and intuitive API, making it easy to learn and use even for beginners. The reason for choosing a framework was to make the code more faster and easier to read and understand and also that gorilla mux is outdated and not maintained anymore.

- Mongodb: as it is a NoSQL document database, it is very easy to scale and also it is very fast and efficient. It is also very easy to use and understand. It is also very easy to integrate with GoLang.
- mongo-go-driver: it is the official Go driver for MongoDB. It is very easy to use and understand and also it is very fast and efficient.

## Features
- User/Admin Login
- User/Admin Logout
- Admin User adds a new User account(by providing the username & password)
- Admin User deletes an existing User account from their organization
- List all Users in their organization
- Passwords are hashed using bcrypt


## Installation
This app requires go 1.19+ to run.

Clone the repository
and run this command 
```sh
make run 
```
this will start the server on port 8080 (the .env file is already present in the repo for the simplicity to run and check the program and can alter the .env file according to your needs)


## API Endpoints Testing
- [Postman Collection](https://elements.getpostman.com/redirect?entityId=16228609-6328e8f5-b08d-49e6-9a77-3e6aabd24e37&entityType=collection)

<br>

#### Admin Signup
```http
  POST /signup
```
This endpoint takes username, password and orgid of the admin and returns inserts the data in User Collection and usertoken collection and returns the insertedID of both the collections.

<br>


#### Admin can Add user to the Organisation using username and password
```http
  POST /users
```
This endpoint takes username, password, token of the admin and returns inserts the data in User Collection and usertoken collection and returns the insertedID of both the collections. using the token orgid is extracted to save the new user in the specified orgid

<br>

#### Admin can Delete the User from the organization using the Userid
```http
  DELETE /users/:userid
```
This endpoint takes the token and deletes the user of the orgid which is extracted from the token. the userid is extracted from the query parameters.

<br>

#### Logout the user 
```http
  POST /logout
```
This endpoint takes the accesstoken from the user and removes the accesstoken and refreshtoken from the database so that the user cant access the protected routes.

<br>

#### Login
```http
  POST /login
```
Endpoint takes in username and password as input and and returns the accesstoken and refreshtoken and all the user details.

<br>

#### Refresh access token 
```http
  POST /refresh
```
This endpoint takes in the refreshtoken and updates the accesstoken if the refreshtokne hasnt expired yet then.

<br>

#### Get all users in the organization by users or admin 
```http
  GET /users
```
This endpoint takes the token and extracts the orgID from the token and get all users according to the page and uses mongodb aggregation feature. 


<br>

## DB Design
The Database is designed in such a way that it is very easy to scale and also it is very fast and efficient. The Database has two Collections called Users and UserTokens and the schema is as follows:

<br>

### Users Table:

| Column Name | Data Type | Description                  |
| :---------- | :-------- | :----------------------      |
| `id`        | `ObejectId`| Unique ID for each user     |
| `userid`    | `string`   | User ID for each user       |
| `username`  | `string`   | username                    |
| `password`  | `string`   | Password of the user        |
| `usertype`  | `string`   | Is the user an admin or user|
| `orgid`     | `string`   | ID of the organisation      |
| `createdat` | `Date`     | When was the Document created|
| `updatedat` | `Date`     | When was the Document updated|

### UserTokens Table:

| Column Name     | Data Type | Description                  |
| :----------     | :-------- | :----------------------      |
| `id`            | `ObejectId`| Unique ID for each user     |
| `userid`        | `string`   | User ID for each user       |
| `accesstoken`   | `string`   | Access token with 1hr validity|
| `refreshtoken`  | `string`   | Access token with 24hrs validity|
| `createdat`     | `Date`     | When was the document created|
| `updatedat`     | `Date`     | When was the document updated|

<br>

This schema makes the database scalable and seperates the tokens collection from the user collection.

## Authors

- [@shreykhandelwal](https://github.com/HawkingRadiation42)
