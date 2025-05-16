# Todo App API

A simple yet powerful RESTful API for managing todo tasks, built with Go and MongoDB. This project serves as a boilerplate for building complete backend APIs with proper architecture and security.

## Features

- User authentication with JWT
- Todo task management (create, read, update, delete)
- Swagger documentation
- MongoDB integration
- Secure password handling with bcrypt and pepper
- Well-structured codebase following clean architecture principles
- Middleware for authentication and error handling
- Comprehensive validation

## Tech Stack

- [Go](https://golang.org/) - Programming language
- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [MongoDB](https://www.mongodb.com/) - Database
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation
- [Swagger](https://swagger.io/) - API documentation
- [godotenv](https://github.com/joho/godotenv) - Environment variable loading

## Project Structure

```
todo-app/
├── cmd/
│   └── api/          # Application entry point
├── internal/
│   ├── auth/         # Authentication logic
│   ├── config/       # Application configuration
│   ├── controller/   # API controllers
│   ├── repository/   # Data access layer
│   ├── routes/       # API routes
│   └── service/      # Business logic
├── pkg/
│   └── database/     # Database connections
├── docs/             # Swagger documentation
└── .env.example      # Example environment variables
```

## Getting Started

### Prerequisites

- Go 1.16+
- MongoDB
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ovansa/todo-app-go.git
   cd todo-app-go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   ```bash
   cp .env.example .env
   ```
   Edit the `.env` file with your configuration:
   ```
   PORT=8080
   MONGO_URI=mongodb://localhost:27017
   DATABASE_NAME=todo_app
   JWT_SECRET=your_secret_key
   JWT_EXPIRATION=24
   PASSWORD_PEPPER=your_password_pepper
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

## API Documentation

API documentation is available via Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

## Authentication

The API uses JWT for authentication. To access protected endpoints:

1. Register a new user at `/api/auth/register`
2. Login at `/api/auth/login` to obtain a JWT token
3. Include the token in the Authorization header as `Bearer {token}`

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login and get JWT token

### Todo Operations

- `GET /api/todos` - Get all todos for authenticated user
- `GET /api/todos/:id` - Get a specific todo
- `POST /api/todos` - Create a new todo
- `PUT /api/todos/:id` - Update a todo
- `DELETE /api/todos/:id` - Delete a todo

## Configuration

The application can be configured using environment variables:

- `MONGO_URI` - MongoDB connection string
- `DATABASE_NAME` - MongoDB database name
- `JWT_SECRET` - Secret for JWT signing
- `JWT_EXPIRATION` - JWT token expiration time in hours
- `PASSWORD_PEPPER` - Additional security for password hashing
- `TEST_MODE` - Enable test mode

## Development

### Running Tests

```bash
go test ./...
```

### Generating Swagger Documentation

```bash
swag init
```

## Project Architecture

This project follows clean architecture principles:

1. **Controller Layer** - Handles HTTP requests/responses and input validation
2. **Service Layer** - Contains the business logic
3. **Repository Layer** - Data access abstraction
4. **Entity Layer** - Domain models and DTOs

### Benefits of This Architecture

- Clear separation of concerns
- Improved testability
- Maintainable and scalable codebase
- Decoupled from external dependencies

## Security Features

- Password hashing with bcrypt and additional "pepper"
- JWT authentication with expiration
- Input validation to prevent injection attacks
- Secure HTTP headers

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

Muhammed Ibrahim - aminmuhammad18@gmail.com

Project Link: [https://github.com/ovansa/todo-app-go](https://github.com/ovansa/todo-app-go)
