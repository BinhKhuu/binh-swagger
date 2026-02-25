# Binh Swagger

This project focuses on learning Golang by building a command-line (CMD) application. Initially started as a project to build an API server, it grew into incorporating some elements of OpenAPI 2.0, though it does not fully adhere to the specification.

## Features
- Command-line (CMD) project structure in Golang.
- Generates a Golang API server project using the Gin framework from an API configuration.
- Swagger documentation for all API endpoints.
- Easy-to-use interface for testing API requests.

## Getting Started
1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/binh-swagger.git
   cd binh-swagger
   ```
2. Install the required dependencies:
   ```bash
   go mod tidy
   ```

## Debugging
For debugging information, see [TESTING.md](TESTING.md).

## Purpose
The purpose of this project is to learn and practice Golang by building a CMD project and integrating Swagger for API documentation and testing.

## Technologies Used
- Golang
- Swagger

## Future Enhancements
- Migrate to OpenAPI 3.0.
- Refactor tests in the `generate` package:
  - Add more table-driven tests.
  - Improve robustness of mock helpers.
- Add more detailed help text to the main menu.
- Create additional commands.

## Running the Tool
1. In the target directory, initialize the Go module:
   ```bash
   go mod init
   ```
2. Run the tool with the following command:
   ```bash
   go run ./cmd/swagger/swagger.go config --api /DirectoryOfAPIConfiguration APIConfigurationFileName.yaml
   ```

## API Configuration YAML Example
Below is an example of the API configuration YAML file:
``` yaml
version: "1.0"

definitions:
  User:
    name: "User"
    fields:
      - name: "Id"
        type: "string"
        json: "id"
      - name: "Email"
        type: "string"
        json: "email"
      - name: "Created_at"
        type: "time.Time"
        json: "created_at"

  Product: 
    name: "Product"
    fields:
      - name: "Id"
        type: "string"
        json: "id"
      - name: "Title"
        type: "string"
        json: "title"
      - name: "Price"
        type: "float64"
        json: "price"

paths:
    /users:
      get:
        summary: List all users
        operationId: listUsers
        produces:
          - application/json
        responses:
          200:
            description: A list of users
            schema:
              type: array
              items:
                $ref: "#/definitions/User"
      post:
        summary: Create a new user
        operationId: createUser
        produces:
          - application/json
        responses:
          201:
            description: The created user
            schema:
              type: object
              $ref: "#/definitions/User"
````