# Todo API with Go and DDD

A simple Todo API built with Go following Domain-Driven Design principles.

## Project Structure

- `cmd/api`: Application entry point
- `internal/domain`: Domain models and repository interfaces
- `internal/application`: Application services and use cases
- `internal/infrastructure`: Repository implementations and other infrastructure concerns
- `internal/interfaces`: API handlers and controllers

## Getting Started

To run the application:

```bash
cd cmd/api
go run main.go
```

The API will be available at http://localhost:8080

## API Endpoints

- `GET /api/todos`: Get all todos
- `GET /api/todos/{id}`: Get a specific todo
- `POST /api/todos`: Create a new todo
- `PUT /api/todos/{id}`: Update a todo
- `DELETE /api/todos/{id}`: Delete a todo
- `PATCH /api/todos/{id}/complete`: Mark a todo as completed
- `PATCH /api/todos/{id}/incomplete`: Mark a todo as incomplete

## Example API Usage

### Create a Todo

```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Go", "description": "Learn Go programming language"}'
```

### Get All Todos

```bash
curl http://localhost:8080/api/todos
```

### Mark Todo as Completed

```bash
curl -X PATCH http://localhost:8080/api/todos/{id}/complete
```
