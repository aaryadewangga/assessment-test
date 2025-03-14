# Assessment Test

This project is a Go-based application designed to manage user authentication, product management, and transaction processing. It utilizes the Echo framework for HTTP routing and middleware support. The application uses PostgreSQL as its database.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Routes](#routes)
- [Middleware](#middleware)
- [Controllers](#controllers)
- [Database](#database)

## Features
- User registration and authentication
- JWT-based authorization (EdDSA - Ed25519)
- Product management (CRUD operations)
- Transaction processing (with multiple products per transaction)
- Export transaction details to PDF and Excel formats
- Uses PostgreSQL as the database

## Installation

1. **Clone the Repository**  
   `git clone https://github.com/aaryadewangga/assessment-test.git`  
   `cd assessment-test`

2. **Install Dependencies**  
   Make sure you have Go installed (version 1.18+). Then, install dependencies:  
   `go mod tidy`

3. **Install and Configure PostgreSQL**  
   Ensure you have PostgreSQL installed. You can install it via:

   **Ubuntu/Debian**  
   `sudo apt update`  
   `sudo apt install postgresql postgresql-contrib`

   **Mac (Homebrew)**  
   `brew install postgresql`

   **Windows (Chocolatey)**  
   `choco install postgresql`

   Once installed, start PostgreSQL and create a database:  
   `CREATE DATABASE assessment_test;`  
   `CREATE USER your_user WITH ENCRYPTED PASSWORD 'your_password';`  
   `GRANT ALL PRIVILEGES ON DATABASE assessment_test TO your_user;`

4. **Setup Environment Variables**  
   Create a `.env.sh` file and configure the following variables:
   ```
        export APP_NAME=your-app-name
        export APP_PORT=your-app-port
        export DB_PG_ADDR=your-db-address
        export DB_PG_USER=your-db-username
        export DB_PG_PASS=your-db-password
        export DB_PG_DBNAME=your-db-name
        export DATABASE_URL=your-db-url
        export JWT_PRIVATE_KEY=your-jwt-private-key
        export JWT_PUBLIC_KEY=your-jwt-public-key
   ```

5. **Run Database Migrations**  
This project uses Dbmate for database migrations. Install it and run:  
`dbmate up`

6. **Start the Server**  
`go run main.go`  
The server will start at `http://localhost:8080`.

## Configuration  
The application uses PostgreSQL as the database and reads configurations from environment variables stored in the `.env.sh` file.

## Usage  
Once the server is running, you can use an API client like Postman or cURL to interact with the endpoints.

### Register a User  
To register a new user, you can use the following cURL command:  
```
curl -X POST http://localhost:8080/register -d '{"username":"testuser", "password":"password"}' -H "Content-Type: application/json"
```
### Login to Get JWT Token  
To log in and obtain a JWT token, you can use the following cURL command:  
```
curl -X POST http://localhost:8080/login -d '{"username":"testuser", "password":"password"}' -H "Content-Type: application/json"
```

Upon successful login, the server will return a JWT token. This token must be included in the Authorization header for accessing protected routes. 

Example of using the JWT token in a request:  

## Routes

### Authentication  
- **POST** `/register`: Register a new user  
- **POST** `/login`: Authenticate and get JWT token  

### Admin Routes (Protected)  
- **POST** `/admin/products`: Add new product  
- **PUT** `/admin/products`: Update product by ID  
- **DELETE** `/admin/products`: Delete product by ID  

### Cashier & Admin Routes  
- **GET** `/products`: Get all products  

### Transactions (Protected)  
- **POST** `/transactions/create`: Create a new transaction  
- **GET** `/transactions/list`: Get all transactions  
- **GET** `/transactions/details`: Get transaction details by ID  
- **GET** `/transactions/generate/pdf`: Export transaction details as PDF  
- **GET** `/transactions/generate/excel`: Export transaction details as Excel  

## Middleware  
The project uses Echo's middleware for security and logging:  
- **JWT Middleware**: Protects routes using JWT authentication  
- **CORS Middleware**: Enables cross-origin requests  
- **Logging Middleware**: Logs all incoming requests  

## Controllers  
The controllers handle business logic for different API endpoints:  
- **AuthController** → Handles user authentication  
- **User Controller** → Manages user registration  
- **ProductController** → Handles product CRUD operations  
- **TransactionController** → Manages transactions and report generation  

## Database  
The project uses PostgreSQL as the main database.  
- Database migrations are handled using Dbmate.  
- All database interactions are performed using go-pg as the ORM.  
