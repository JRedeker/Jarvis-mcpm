# OpenSpec MCP Server

**GitHub**: https://github.com/fission-ai/openspec
**Package**: openspec-mcp
**Language**: Python
**Transport**: stdio
**Status**: ✅ Active Community Server

---

## Overview

A Model Context Protocol server that provides API specifications and documentation capabilities. OpenSpec enables LLMs to create, validate, and manage OpenAPI specifications, generate documentation, and ensure API consistency across projects.

**Key Features:**
- OpenAPI specification generation
- API documentation automation
- Specification validation and linting
- Code generation from specifications
- API consistency checking
- Multi-format support (OpenAPI 3.0, 3.1, Swagger 2.0)

---

## Installation

### Using pip
```bash
pip install openspec-mcp
```

### Using uv (recommended)
```bash
# Install uv if not already installed
curl -LsSf https://astral.sh/uv/install.sh | sh

# Run directly with uvx
uvx openspec-mcp
```

---

## Configuration

The OpenSpec MCP server requires minimal configuration for basic functionality.

**Optional Configuration:**
- `OPENAPI_MCP_API_KEY`: For enhanced features
- `OPENAPI_SPEC_DIR`: Default specifications directory
- `OPENAPI_TEMPLATE_DIR`: Custom templates directory

---

## Available Tools

### Specification Generation
1. `openspec_create_spec`
   - Create new OpenAPI specifications from descriptions
   - Inputs: `title` (required), `version` (required), `description` (optional), `endpoints` (optional)
   - Returns: Generated OpenAPI specification

2. `openspec_add_endpoint`
   - Add endpoints to existing specifications
   - Inputs: `spec_path` (required), `path` (required), `method` (required), `parameters` (optional)
   - Returns: Updated specification

### Validation and Linting
3. `openspec_validate`
   - Validate OpenAPI specifications for correctness
   - Inputs: `spec_path` (required), `strict_mode` (optional)
   - Returns: Validation results and errors

4. `openspec_lint`
   - Lint specifications for best practices
   - Inputs: `spec_path` (required), `rules` (optional)
   - Returns: Linting results and recommendations

### Documentation Generation
5. `openspec_generate_docs`
   - Generate documentation from specifications
   - Inputs: `spec_path` (required), `format` (optional), `template` (optional)
   - Returns: Generated documentation

6. `openspec_generate_client`
   - Generate client code from specifications
   - Inputs: `spec_path` (required), `language` (required), `options` (optional)
   - Returns: Generated client code

### Schema Management
7. `openspec_create_schema`
   - Create JSON schemas for request/response bodies
   - Inputs: `schema_name` (required), `properties` (required), `type` (optional)
   - Returns: JSON schema definition

8. `openspec_extract_schemas`
   - Extract schemas from existing specifications
   - Inputs: `spec_path` (required), `component_type` (optional)
   - Returns: Extracted schema definitions

### Comparison and Merging
9. `openspec_compare`
   - Compare two API specifications for differences
   - Inputs: `spec1_path` (required), `spec2_path` (required)
   - Returns: Comparison report

10. `openspec_merge`
    - Merge multiple specifications
    - Inputs: `spec_paths` (required), `merge_strategy` (optional)
    - Returns: Merged specification

---

## Usage Examples

### Specification Creation
```json
{
  "tool": "openspec_create_spec",
  "arguments": {
    "title": "User Management API",
    "version": "1.0.0",
    "description": "API for managing user accounts and profiles",
    "endpoints": [
      {
        "path": "/users",
        "method": "GET",
        "description": "List all users"
      },
      {
        "path": "/users/{userId}",
        "method": "GET",
        "description": "Get user by ID"
      }
    ]
  }
}
```

### Endpoint Addition
```json
{
  "tool": "openspec_add_endpoint",
  "arguments": {
    "spec_path": "api-spec.yaml",
    "path": "/users/{userId}/profile",
    "method": "PUT",
    "parameters": [
      {
        "name": "userId",
        "in": "path",
        "required": true,
        "schema": {"type": "string"}
      }
    ]
  }
}
```

### Validation
```json
{
  "tool": "openspec_validate",
  "arguments": {
    "spec_path": "api-spec.yaml",
    "strict_mode": true
  }
}
```

### Documentation Generation
```json
{
  "tool": "openspec_generate_docs",
  "arguments": {
    "spec_path": "api-spec.yaml",
    "format": "html",
    "template": "modern"
  }
}
```

### Client Generation
```json
{
  "tool": "openspec_generate_client",
  "arguments": {
    "spec_path": "api-spec.yaml",
    "language": "python",
    "options": {
      "async": true,
      "type_hints": true
    }
  }
}
```

---

## Supported Formats

### OpenAPI Versions
- **OpenAPI 3.0**: Industry standard specification
- **OpenAPI 3.1**: Latest version with JSON Schema support
- **Swagger 2.0**: Legacy format support

### Output Formats
- **YAML**: Human-readable specification format
- **JSON**: Machine-readable specification format
- **HTML**: Interactive documentation
- **Markdown**: GitHub-friendly documentation

### Client Languages
- **Python**: Full async/sync support
- **JavaScript/TypeScript**: Node.js and browser support
- **Go**: Native Go client generation
- **Java**: Spring Boot integration
- **C#**: .NET client support

---

## Schema Support

### Data Types
- **Primitive**: string, integer, number, boolean
- **Complex**: arrays, objects, enums
- **Advanced**: oneOf, anyOf, allOf combinations
- **Custom**: User-defined types and formats

### Validation Rules
- **Format**: email, uuid, date, uri, etc.
- **Pattern**: Regular expression validation
- **Range**: minimum/maximum values
- **Length**: string and array length constraints

---

## Known Limitations

1. **Complex Schemas**: Limited support for highly nested schemas
2. **Custom Extensions**: May not support all OpenAPI extensions
3. **Large Specifications**: Performance issues with very large specs
4. **Real-time Validation**: No live validation during editing
5. **Advanced Security**: Basic security scheme support only

---

## Testing

```bash
# Test server startup
uvx openspec-mcp

# Test tool availability
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | uvx openspec-mcp

# Test specification operations (requires MCP client)
```

---

## Use Cases

1. **API Design**: Design APIs with proper specifications
2. **Documentation**: Generate comprehensive API documentation
3. **Code Generation**: Create client SDKs from specifications
4. **API Validation**: Ensure API consistency and correctness
5. **Team Collaboration**: Share API specifications across teams

---

## Best Practices

### Specification Design
- Start with clear API requirements
- Use consistent naming conventions
- Include comprehensive descriptions
- Define proper response schemas

### Documentation
- Write clear, concise descriptions
- Include examples for all endpoints
- Document error responses
- Keep documentation synchronized

### Validation
- Validate specifications regularly
- Use linting to catch issues early
- Test generated clients thoroughly
- Maintain version compatibility

---

## Related Documentation

- [OpenAPI Specification](https://swagger.io/specification/)
- [OpenSpec MCP GitHub Repository](https://github.com/fission-ai/openspec)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [API Design Best Practices](https://swagger.io/resources/articles/best-practices-in-api-design/)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 1A)
**Next Steps**: Register with jarvis and test functionality
**Category**: Core Development Tool
