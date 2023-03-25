# Authorization + Authentication service in Golang using JWT

<summary><a href="#important-change-to-be-made-before-running">Please check this out before proceeding to installation.</a></summary>

<summary><a href="#important-change-to-be-made-before-running">Did you follow this or not ?</a></summary>

<summary><a href="#important-change-to-be-made-before-running">Error connecting to DB. Follow this.</a></summary>

## Where is the source code ?

> `backend/cmd` - contains the main.go file and all the other necessary folders.

> `backend/test` - contains the test code.

## Features

- User Login
- User Logout
- All passwords are hashed before storing in the database.
- Admin User adds a new User account(by providing the username & password)
- Admin User deletes an existing User account from their organization
- List all Users in their organization
- Admin User adds a new Organization(by providing the organization name & head)

## Architecture / Stack

- Golang framework - [Gofiber](https://pkg.go.dev/github.com/gofiber/fiber@v1.14.6)
- Database - [MySQL](https://www.mysql.com/)
- ORM - [GORM](https://pkg.go.dev/gorm.io/gorm@v1.24.6)
- JWT - [jwt-go](https://pkg.go.dev/github.com/dgrijalva/jwt-go@v3.2.0)
- Password Hashing - [bcrypt](https://pkg.go.dev/golang.org/x/crypto@v0.7.0/bcrypt)
- Unit Testing: [Testing package](https://pkg.go.dev/testing)

## Installation

This app requires [Go](https://go.dev/doc/install) v1.20+ to run.  
Also make sure to install [docker](https://www.docker.com/products/docker-desktop/) and [docker-compose](https://docs.docker.com/compose/install/) for your operating system.

Clone the repository

```bash
  git clone https://github.com/HousewareHQ/houseware---backend-engineering-octernship-RohitShah1706.git
```

## Run using go locally.

Install the dependencies:

```bash
  cd houseware---backend-engineering-octernship-RohitShah1706/backend
  go mod download
```

Run the application

```bash
  go run cmd/main.go
```

This will start the dev server on port 8080.

## Or run using docker in dev environment

We will use the following docker images along with our own go-app image  
| Plugin | Links |
| ------ | ------ |
| Nginx | https://hub.docker.com/_/nginx |
| MYSQL | https://hub.docker.com/_/mysql |
| Golang | https://hub.docker.com/_/golang |

Build required docker images and start the container.

```bash
  docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build
```

This will automatically start the go-app on port 8080 if in dev mode.  
If in prod mod it will start the server on port 80 TCP.

This docker-compose file sets up a basic web application with three services:  
| Container | Role |
| ------ | ------ |
| go-app | Handles the main login for running the app. |
| nginx | Acts as a reverse proxy to the go-app service, so incoming traffic to the application can be efficiently directed to the appropriate service. Can also act as a load-balancer to distribute requests across multiple instances of the go-app service. |
| mysql | Runs MySQL database and creates a persistent volume for data storage. This way your data persists even after tearing down the containers. |

This setup provides a scalable and easily deployable web application architecture that can be customized further according to specific requirements.

## Environment Variables

To run this application, you will need to add the following environment variables to your .env file.  
For development mode the `.env` file is present in backend directory.

`SECRET_KEY`

`NGINX_PORT`

`DB_USERNAME`

`DB_PASSWORD`

`DB_HOST`

`DB_PORT`

`DB_NAME`

`ADMIN_ORG_NAME`

`ADMIN_ORG_HEAD`

`ADMIN_USERNAME`

`ADMIN_PASSWORD`

## Important change to be made before running

> `DB_HOST` in `.env` file should be localhost if running the app using go command.  
> `DB_HOST` in `.env` file should be `mysql` or whatever is the name of your MySQL service so that docker's internal dns will replace with the appropriate IP Address.

## Postman API Reference

> [Postman API access](https://elements.getpostman.com/redirect?entityId=20239076-994ff5e1-5715-49b9-a59e-d3f84e2f3c78&entityType=collection) - all the endpoints are here for your direct testing.

### Admin endpoints - protected with isAdminCheck middleware

#### Create Organisation

```http
  POST /api/org
```

| Body   | Type     | Description                                   |
| :----- | :------- | :-------------------------------------------- |
| `name` | `string` | **Required** Name of the organisation         |
| `head` | `string` | **Required** Name of the head of organisation |

#### Delete Organisation

```http
  DELETE /api/org
```

| Body   | Type     | Description                           |
| :----- | :------- | :------------------------------------ |
| `name` | `string` | **Required** Name of the organisation |

#### Update Organisation details

```http
  PATCH /api/org
```

| Body   | Type     | Description                                       |
| :----- | :------- | :------------------------------------------------ |
| `name` | `string` | **Required** Name of the organisation             |
| `head` | `string` | **Required** Name of the new head of organisation |

#### Create User Account

```http
  POST /api/auth
```

| Body       | Type     | Description                                                       |
| :--------- | :------- | :---------------------------------------------------------------- |
| `username` | `string` | **Required** Name of the user                                     |
| `password` | `string` | **Required** Password that will be hashed while storing in the DB |
| `org_id`   | `string` | **Required** Organisation the user belongs to                     |

#### Delete User Account

```http
  DELETE /api/auth
```

| Body       | Type     | Description                   |
| :--------- | :------- | :---------------------------- |
| `username` | `string` | **Required** Name of the user |

### User endpoints - public

#### Login to account - will create a new JWT token with expiry time of 1 hour

```http
  POST /api/auth/login
```

| Body       | Type     | Description                       |
| :--------- | :------- | :-------------------------------- |
| `username` | `string` | **Required** Name of the user     |
| `password` | `string` | **Required** Password of the user |

#### Refresh access token - grants a new acces token with an expiry time of 24 hours

```http
  POST /api/auth/refresh
```

| Cookie | Type           | Description                                                                                       |
| :----- | :------------- | :------------------------------------------------------------------------------------------------ |
| `jwt`  | `access token` | **Required** Users should be logged in and have a valid access token to refresh the current token |

#### Get user account

```http
  GET /api/auth/user
```

| Cookie | Type           | Description                                                                                       |
| :----- | :------------- | :------------------------------------------------------------------------------------------------ |
| `jwt`  | `access token` | **Required** Users should be logged in and have a valid access token to get their account details |

#### Logout

```http
  GET /api/auth/logout
```

| Cookie | Type           | Description                                                                    |
| :----- | :------------- | :----------------------------------------------------------------------------- |
| `jwt`  | `access token` | **Required** Users should be logged in and have a valid access token to logout |

#### Get people in same organisation

```http
  GET /api/org
```

| Cookie | Type           | Description                                                                                                     |
| :----- | :------------- | :-------------------------------------------------------------------------------------------------------------- |
| `jwt`  | `access token` | **Required** Users should be logged in and have a valid access token to get list of users in their organisation |

## Running Unit Tests

To run tests, run the following commands.  
Flush the whole DB before running the tests so that mock data doesn't intersect with the previous data.

Make sure you are in the root directory of the project where the docker-compose files are present.  
Start the MySQL server service only if you are using Docker.

```bash
  docker-compose -f .\docker-compose.yml -f .\docker-compose.dev.yml up -d mysql
```

Now run unit tests.  
Make sure you are in the backend directory where the go.mod file is located.

```bash
    cd backend
    go test ./test -v
```

## Rationale

- **Golang Framework**: GoFiber is a lightweight and fast web framework that is perfect for building APIs. It has a simple and intuitive API that makes it easy to build and maintain APIs. It also has excellent performance and scalability, making it ideal for high-traffic applications.

- **MySQL Database**: MySQL is a popular open-source relational database management system. It is known for its robustness, scalability, and performance. It also has excellent support for ACID transactions, making it ideal for applications that require data consistency and reliability.

- **Gorm**: GORM is a powerful and easy-to-use ORM for Golang. It provides a simple and intuitive API for interacting with databases, making it easy to build and maintain database-driven applications. It also has excellent support for MySQL, making it an ideal choice for this project.

## DB Design

The database will have two tables called "users" and "organisations" that will store the user and organisation details respectively. The tables will have the following columns:

### Users Table

| Column Name | Data Type | Description             |
| :---------- | :-------- | :---------------------- |
| `id`        | `int`     | Unique ID for each user |
| `username`  | `string`  | Name of the user        |
| `password`  | `string`  | Password of the user    |
| `is_admin`  | `bool`    | Is the user an admin    |
| `org_id`    | `int`     | ID of the organisation  |

### Organisations Table

| Column Name | Data Type | Description                          |
| :---------- | :-------- | :----------------------------------- |
| `id`        | `int`     | Unique ID for each org               |
| `name`      | `string`  | Name of the organisation             |
| `head`      | `string`  | Name of the head of the organisation |

The "id" column will be the primary key for both tables. The "org_id" column in the "users" table will be a foreign key that references the "id" column in the "organisations" table.

The password will be hashed using bcrypt before being stored in the database to ensure security.

## Acknowledgements

- [Dev.to article](https://dev.to/koddr/go-fiber-by-examples-testing-the-application-1ldf)
- [Golang docs](https://go.dev/doc/)
- [Golang packages](https://pkg.go.dev/fmt)

## Authors

- [@RohitShah1706](https://github.com/RohitShah1706)

## License

[MIT](https://choosealicense.com/licenses/mit/)
