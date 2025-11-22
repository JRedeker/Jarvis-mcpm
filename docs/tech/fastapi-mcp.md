# FastAPI MCP Server

**GitHub**: https://github.com/tadata-org/fastapi_mcp
**Package**: fastapi-mcp
**Language**: Python
**Transport**: stdio
**Status**: ✅ Active Community Server

---

## Overview

A Model Context Protocol server that provides FastAPI development and API management capabilities. This server enables LLMs to create, modify, and manage FastAPI applications with intelligent code generation and API documentation.

**Key Features:**
- FastAPI application generation
- API endpoint creation and management
- Automatic OpenAPI documentation
- Pydantic model integration
- Database integration support
- Middleware and dependency injection

---

## Installation

### Using pip
```bash
pip install fastapi-mcp
```

### Using uv (recommended)
```bash
# Install uv if not already installed
curl -LsSf https://astral.sh/uv/install.sh | sh

# Run directly with uvx
uvx fastapi-mcp
```

---

## Configuration

The FastAPI MCP server requires no environment variables for basic functionality. It runs with default settings.

**Optional Configuration:**
- `FASTAPI_MCP_API_KEY`: For enhanced features
- `FASTAPI_BASE_URL`: Custom API base URL
- `DATABASE_URL`: For database integration features

---

## Available Tools

### Application Generation
1. `fastapi_create_app`
   - Generate complete FastAPI applications
   - Inputs: `app_name` (required), `features` (optional), `database` (optional)
   - Returns: Generated application structure

2. `fastapi_create_endpoint`
   - Create API endpoints with proper structure
   - Inputs: `path` (required), `method` (required), `response_model` (optional)
   - Returns: Endpoint implementation

### Model and Schema Management
3. `fastapi_create_model`
   - Create Pydantic models for request/response
   - Inputs: `model_name` (required), `fields` (required)
   - Returns: Pydantic model code

4. `fastapi_create_schema`
   - Generate OpenAPI schemas
   - Inputs: `schema_type` (required), `fields` (required)
   - Returns: Schema definition

### Database Integration
5. `fastapi_setup_database`
   - Configure database connections
   - Inputs: `database_type` (required), `connection_string` (required)
   - Returns: Database configuration code

6. `fastapi_create_crud`
   - Generate CRUD operations for models
   - Inputs: `model_name` (required), `operations` (optional)
   - Returns: CRUD endpoint implementations

### Middleware and Security
7. `fastapi_add_middleware`
   - Add middleware to applications
   - Inputs: `middleware_type` (required), `config` (optional)
   - Returns: Middleware implementation

8. `fastapi_add_authentication`
   - Implement authentication systems
   - Inputs: `auth_type` (required), `config` (optional)
   - Returns: Authentication setup code

---

## Usage Examples

### Application Creation
```json
{
  "tool": "fastapi_create_app",
  "arguments": {
    "app_name": "user_management_api",
    "features": ["authentication", "database", "validation"],
    "database": "postgresql"
  }
}
```

### Endpoint Creation
```json
{
  "tool": "fastapi_create_endpoint",
  "arguments": {
    "path": "/users/{user_id}",
    "method": "GET",
    "response_model": "UserResponse"
  }
}
```

### Model Creation
```json
{
  "tool": "fastapi_create_model",
  "arguments": {
    "model_name": "UserCreate",
    "fields": {
      "username": "str",
      "email": "EmailStr",
      "password": "str"
    }
  }
}
```

### CRUD Generation
```json
{
  "tool": "fastapi_create_crud",
  "arguments": {
    "model_name": "User",
    "operations": ["create", "read", "update", "delete"]
  }
}
```

---

## Supported Features

### Framework Integration
- **FastAPI**: Full support for all FastAPI features
- **Pydantic**: Model validation and serialization
- **SQLAlchemy**: Database ORM integration
- **Alembic**: Database migrations
- **Uvicorn**: ASGI server configuration

### Database Support
- **PostgreSQL**: Full-featured relational database
- **MySQL**: Popular relational database
- **SQLite**: Lightweight file-based database
- **MongoDB**: NoSQL document database
- **Redis**: In-memory data store

### Authentication Methods
- **JWT**: JSON Web Token authentication
- **OAuth2**: OAuth2 flow implementation
- **API Keys**: Simple API key authentication
- **Basic Auth**: HTTP basic authentication

---

## Project Structure

Generated FastAPI applications follow best practices:
```
my_api/
├── app/
│   ├── __init__.py
│   ├── main.py
│   ├── models/
│   ├── routes/
│   ├── schemas/
│   └── database.py
├── tests/
├── requirements.txt
└── README.md
```

---

## Known Limitations

1. **Python Focus**: Limited to Python/FastAPI development
2. **Complex Architectures**: May not handle highly complex microservice patterns
3. **Performance Tuning**: Generated code may need optimization for high-load scenarios
4. **Security**: Generated authentication may need security review
5. **Database Complexity**: Limited support for complex database relationships

---

## Testing

```bash
# Test server startup
fastapi-mcp

# Test tool availability
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | fastapi-mcp

# Test API generation (requires MCP client)
```

---

## Use Cases

1. **API Development**: Rapid API prototyping and development
2. **Microservices**: Create microservice components
3. **CRUD Applications**: Generate standard CRUD operations
4. **Authentication Systems**: Implement secure authentication
5. **Database Integration**: Connect APIs to various databases

---

## Best Practices

### API Design
- Follow RESTful principles
- Use appropriate HTTP status codes
- Implement proper error handling
- Document endpoints with descriptions

### Security
- Validate all input data
- Implement rate limiting
- Use HTTPS in production
- Keep dependencies updated

### Performance
- Use database connection pooling
- Implement caching where appropriate
- Optimize database queries
- Monitor API performance

---

## Related Documentation

- [FastAPI Documentation](https://fastapi.tiangolo.com/)
- [FastAPI MCP Server](https://github.com/tadata-org/fastapi_mcp)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [Pydantic Documentation](https://docs.pydantic.dev/)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 1A)
**Next Steps**: Register with jarvis and test functionality
**Category**: Backend Development Tool
