# OpenAPI Implementation Plan for Todo API

## Current API Analysis

I've analyzed the existing Todo API implementation and found the following:

### API Structure
- **Framework**: Gorilla Mux for HTTP routing
- **Architecture**: Domain-Driven Design (DDD) with Clean Architecture
- **Data Storage**: In-memory repository implementation
- **Authentication**: No authentication currently implemented

### Endpoints
The API has 7 endpoints:
1. `GET /api/todos`: Get all todos
2. `GET /api/todos/{id}`: Get a specific todo
3. `POST /api/todos`: Create a new todo
4. `PUT /api/todos/{id}`: Update a todo
5. `DELETE /api/todos/{id}`: Delete a todo
6. `PATCH /api/todos/{id}/complete`: Mark a todo as complete
7. `PATCH /api/todos/{id}/incomplete`: Mark a todo as incomplete

### Request/Response Models
- **TodoResponse**: Response structure for todos
- **CreateTodoRequest**: Request structure for creating todos
- **UpdateTodoRequest**: Request structure for updating todos

### Error Handling
- HTTP status codes with error messages
- Domain-specific errors like `ErrTodoNotFound` and `ErrEmptyTitle`

## OpenAPI Specification

I've created an initial OpenAPI 3.0 specification in `api/openapi.yaml` that documents all 7 endpoints with their request/response schemas and error responses. Here's an example of one endpoint:

```yaml
  /api/todos:
    post:
      summary: Create a new todo
      description: Creates a new todo item
      operationId: createTodo
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateTodoRequest'
      responses:
        '201':
          description: Todo created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/TodoResponse'
        '400':
          description: Bad request
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
```

## Implementation Strategy

### Phase 1: OpenAPI Specification Refinement

1. **Priority Endpoints to Document**:
   - Start with the core CRUD operations:
     1. `POST /api/todos` (Create)
     2. `GET /api/todos` (Read All)
     3. `GET /api/todos/{id}` (Read One)
     4. `PUT /api/todos/{id}` (Update)
     5. `DELETE /api/todos/{id}` (Delete)
   - Then document the specialized endpoints:
     6. `PATCH /api/todos/{id}/complete`
     7. `PATCH /api/todos/{id}/incomplete`

2. **Refinement Steps**:
   - Add detailed descriptions for each endpoint
   - Include examples for request/response bodies
   - Document all possible error responses
   - Add validation rules for request parameters

### Phase 2: Code Generation Setup

1. **Install oapi-codegen**:
   ```bash
   go install github.com/oapi-codegen/oapi-codegen/cmd/oapi-codegen@latest
   ```

2. **Create Configuration File**:
   Create `api/oapi-codegen-config.yaml` with the following content:
   ```yaml
   package: api
   generate:
     models: true
     server: true
     client: false
     spec: true
   output: internal/interfaces/api/generated.go
   import-mapping:
     ./components.yaml: github.com/tomotakashimizu/devin-sandbox/todo-app/internal/interfaces/api
   ```

3. **Create Makefile**:
   Create a Makefile with a target for code generation:
   ```makefile
   .PHONY: generate-api
   generate-api:
       oapi-codegen -config api/oapi-codegen-config.yaml -package api api/openapi.yaml > internal/interfaces/api/generated.go
   ```

### Phase 3: Integration with Existing Code

1. **Integration Strategy**:
   - **Adapter Pattern**: Create adapter functions that map between generated types and domain types
   - **Interface Implementation**: Implement the generated server interface with the existing handler logic
   - **Gradual Migration**: Start with one endpoint, test thoroughly, then migrate others

2. **Implementation Steps**:
   - Create a new file `internal/interfaces/api/oapi_handler.go` that:
     - Implements the generated server interface
     - Uses the existing `TodoService` for business logic
     - Maps between generated types and domain types
   - Update `cmd/api/main.go` to use the new handler

3. **Example Integration Code**:
   ```go
   // internal/interfaces/api/oapi_handler.go
   package api

   import (
       "net/http"

       "github.com/tomotakashimizu/devin-sandbox/todo-app/internal/application"
   )

   // TodoAPIHandler implements the generated server interface
   type TodoAPIHandler struct {
       todoService *application.TodoService
   }

   // NewTodoAPIHandler creates a new TodoAPIHandler
   func NewTodoAPIHandler(todoService *application.TodoService) *TodoAPIHandler {
       return &TodoAPIHandler{
           todoService: todoService,
       }
   }

   // Implement generated interface methods...
   ```

4. **Main.go Updates**:
   ```go
   // cmd/api/main.go
   // ...
   
   // Setup HTTP handlers
   todoService := application.NewTodoService(todoRepo)
   todoHandler := api.NewTodoAPIHandler(todoService)
   
   // Setup router with generated middleware
   swagger, err := api.GetSwagger()
   if err != nil {
       log.Fatalf("Error loading swagger spec: %v", err)
   }
   swagger.Servers = nil // Clear servers to use request host
   
   router := mux.NewRouter()
   router.Use(loggingMiddleware)
   
   // Use oapi-codegen middleware for validation
   router.Use(api.OapiRequestValidator(swagger))
   
   // Register routes
   api.HandlerFromMux(todoHandler, router)
   
   // ...
   ```

### Phase 4: Testing and Validation

1. **Validation Testing**:
   - Test request validation using the generated validators
   - Ensure all error cases are properly handled
   - Verify content type validation

2. **Integration Testing**:
   - Create integration tests for each endpoint
   - Test with valid and invalid requests
   - Verify response formats match the OpenAPI spec

## Benefits of This Approach

1. **API Documentation**: Clear, standardized documentation of the API
2. **Request Validation**: Automatic validation of incoming requests
3. **Type Safety**: Generated types ensure type safety
4. **Client Generation**: Ability to generate client code for consumers
5. **Maintainability**: Single source of truth for API definition
6. **Discoverability**: OpenAPI UI for exploring the API

## Next Steps After Implementation

1. **API Documentation UI**: Set up Swagger UI for interactive documentation
2. **Client SDK Generation**: Generate client SDKs for other languages
3. **Authentication**: Add authentication to the OpenAPI spec and implementation
4. **Versioning**: Implement API versioning in the OpenAPI spec
