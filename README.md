# RESTful API using Go with Repository Pattern

This project is a simple RESTful API built with Go using the Repository Pattern. The Repository Pattern is a design pattern in software development that aims to separate data access logic from business logic. By using this pattern, interactions with data sources (such as databases) are encapsulated in a separate layer known as the repository. This project demonstrates how to implement this design pattern in a RESTful API.

## Table of Contents

- [Getting Started](#getting-started)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

### Prerequisites

Before you can run the project locally, you'll need to have the following installed on your machine:

- [Go](https://golang.org/doc/install) (version 1.18+ recommended)
- [MySQL](https://dev.mysql.com/downloads/mysql/)

### Clone the Repository

```bash
git clone https://github.com/hendriam/RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern.git
cd RESTful-APIs-with-Go-and-MySQL-Using-the-Repository-Pattern
```

### Create a Database and Table
```sql
CREATE DATABASE book_db;
```
```sql
CREATE TABLE books (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Setting Environment Variables
Make sure the environment variables for database configuration and logging are set before running the application. You can do this with the following command in the terminal:
```bash
export HOST=localhost
export PORT=8080
export DB_USER=your_db_user
export DB_PASSWORD=your_db_password
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_NAME=book_db
export LOG_LEVEL=info
```

### Running the Project
```bash
go mod tidy
go run main.go
```
The server will start on http://localhost:8080.

### CRUD API Endpoints
##### Create a New Book
Endpoint: POST /books
Description: Adds a new book to the collection.
Request Body:
```json
{
	"title": "Go Programming - From Beginner to Professional",
	"author": "Samantha Coyle",
	"year": 2024
}
```
Response:
Success (201 Created):
```json
{
	"code": 201,
	"message": "Successfully created data",
	"data": {
		"id": 1,
		"title": "Go Programming - From Beginner to Professional",
		"author": "Samantha Coyle",
		"year": 2024,
		"created_at": "2024-08-09T04:51:07.795858637Z",
		"updated_at": "2024-08-09T04:51:07.795858837Z"
	}
}
```
Validation Error (422 Unprocessable Entity):
```json
{
	"code": 422,
	"message": "validation error",
	"errors": [
		{
			"field": "Title",
			"message": "This field is required"
		}
	]
}
```

##### Get All Books
Endpoint: GET /books
Description: Retrieves a list of all books.
Response:
Success (200 OK)
```json
{
	"code": 200,
	"message": "Successfully got all data",
	"data": [
		{
			"id": 1,
			"title": "Go Programming - From Beginner to Professional",
			"author": "Samantha Coyle",
			"year": 2024,
			"created_at": "2024-08-09T04:52:35Z",
			"updated_at": "2024-08-09T04:52:35Z"
		}
	]
}
```

##### Get a Book by ID
Endpoint: GET /books/:id
Description: Retrieves details of a book by its ID.
Response:
Success (200 OK):
```json
{
	"code": 200,
	"message": "Successfully got the data",
	"data": {
		"id": 1,
		"title": "Go Programming - From Beginner to Professional",
		"author": "Samantha Coyle",
		"year": 2024,
		"created_at": "2024-08-09T04:52:35Z",
		"updated_at": "2024-08-09T04:52:35Z"
	}
}
```

##### Update a Book
Endpoint: PUT /books/{id}
Description: Updates the details of an existing book by its ID.
Request Body:
```json
{
	"title": "Go Programming - From Beginner to Professional",
	"author": "Samantha Coyle",
	"year": 2023
}
```
Response:
Success (200 OK):
```json
{
	"code": 200,
	"message": "Successfully updated data.",
	"data": {
		"id": 1,
		"title": "Go Programming - From Beginner to Professional",
		"author": "Samantha Coyle",
		"year": 2023,
		"created_at": "2024-08-09T04:52:35Z",
		"updated_at": "2024-08-09T07:45:31.908214005Z"
	}
}
```
Not Found (404 Not Found):
```json
{
	"code": 404,
	"message": "Data with that ID does not exist, cannot update",
	"errors": null
}
```

##### Delete a Book
Endpoint: DELETE /books/:id
Description: Deletes a book by its ID.
Response:
```json
{
	"code": 200,
	"message": "Successfully deleted data.",
	"data": null
}
```
Not Found (404 Not Found):
```json
{
	"code": 404,
	"message": "Data with that ID does not exist, cannot deleted",
	"errors": null
}
```

## Contributing
If you find a bug or have an idea for a feature, feel free to open an issue or submit a pull request. Contributions are welcome!

## License
This project is licensed under the MIT License. See the LICENSE file for details.