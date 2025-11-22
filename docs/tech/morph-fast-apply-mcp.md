n# Morph Fast Apply MCP Server

**GitHub**: https://github.com/morph/fast-apply-mcp
**Package**: @morph/fast-apply-mcp
**Language**: TypeScript/JavaScript
**Transport**: stdio
**Status**: ✅ Active Community Server

---

## Overview

A Model Context Protocol server that provides fast and intelligent code transformation capabilities. This server enables LLMs to apply code changes, refactor code, and perform automated code modifications with high accuracy and speed.

**Key Features:**
- Intelligent code transformation and refactoring
- Fast diff application with conflict resolution
- Multi-language support
- Safe code modifications with validation
- Batch processing capabilities
- Integration with version control systems

---

## Installation

### Using npm (via npx)
```bash
npx @morph/fast-apply-mcp
```

### Direct installation
```bash
npm install -g @morph/fast-apply-mcp
```

---

## Configuration

**Required Environment Variable:**
- `MORPH_API_KEY`: Morph API key for enhanced features

**Optional Configuration:**
- `MORPH_BASE_URL`: Custom Morph API base URL
- `MORPH_TIMEOUT`: Request timeout in seconds (default: 30)

**Get a Morph API Key:**
1. Visit [Morph](https://morph.sh/)
2. Sign up for an account
3. Generate an API key from the dashboard
4. Set as environment variable: `export MORPH_API_KEY="your_api_key"`

**Note**: The server can run without an API key with basic functionality.

---

## Available Tools

### Code Transformation
1. `apply_diff`
   - Apply code changes using diff format
   - Inputs: `path` (required), `diff` (required), `dry_run` (optional)
   - Returns: Applied changes and any conflicts

2. `transform_code`
   - Apply intelligent code transformations
   - Inputs: `path` (required), `transformations` (required), `language` (optional)
   - Returns: Transformed code and change summary

### Code Analysis
3. `analyze_code`
   - Analyze code structure and dependencies
   - Inputs: `path` (required), `analysis_type` (optional)
   - Returns: Code analysis report

4. `find_and_replace`
   - Find and replace text with regex support
   - Inputs: `path` (required), `find` (required), `replace` (required), `regex` (optional)
   - Returns: Number of replacements and changed content

### Batch Operations
5. `batch_transform`
   - Apply transformations to multiple files
   - Inputs: `paths` (required), `transformation` (required), `options` (optional)
   - Returns: Batch transformation results

---

## Usage Examples

### Diff Application
```json
{
  "tool": "apply_diff",
  "arguments": {
    "path": "src/utils.js",
    "diff": "<<<<<<< SEARCH\nfunction oldFunction() {\n  return 'old';\n}\n=======\nfunction newFunction() {\n  return 'new';\n}\n>>>>>>> REPLACE",
    "dry_run": false
  }
}
```

### Code Transformation
```json
{
  "tool": "transform_code",
  "arguments": {
    "path": "src/api.js",
    "language": "javascript",
    "transformations": [
      {
        "type": "add_import",
        "import": "import { logger } from './logger';"
      },
      {
        "type": "wrap_function",
        "function_name": "apiCall",
        "wrapper": "logger.wrap"
      }
    ]
  }
}
```

### Find and Replace
```json
{
  "tool": "find_and_replace",
  "arguments": {
    "path": "src/components/*.js",
    "find": "console\\.log\\((.*)\\)",
    "replace": "logger.debug($1)",
    "regex": true
  }
}
```

### Code Analysis
```json
{
  "tool": "analyze_code",
  "arguments": {
    "path": "src/",
    "analysis_type": "dependencies"
  }
}
```

### Batch Transformation
```json
{
  "tool": "batch_transform",
  "arguments": {
    "paths": ["src/file1.js", "src/file2.js", "src/file3.js"],
    "transformation": {
      "type": "add_header",
      "header": "// Copyright 2024 - Auto-generated header"
    }
  }
}
```

---

## Supported Languages

### Primary Languages
- **JavaScript/TypeScript**: Full support with AST parsing
- **Python**: Full support with syntax analysis
- **Java**: Structural transformations
- **C/C++**: Basic transformations
- **Go**: Basic transformations

### Web Technologies
- **HTML/CSS**: DOM-based transformations
- **JSON/YAML**: Structured data transformations
- **Markdown**: Text-based transformations

### Configuration Files
- **XML**: Basic transformations
- **INI/Config**: Text-based transformations

---

## Transformation Types

### Structural Transformations
- `add_import`: Add import statements
- `remove_import`: Remove import statements
- `add_function`: Add new functions
- `remove_function`: Remove functions
- `rename_function`: Rename functions with references
- `add_class`: Add new classes
- `wrap_function`: Wrap functions with decorators/middleware

### Text Transformations
- `replace_text`: Simple text replacement
- `regex_replace`: Regex-based replacement
- `add_header`: Add file headers
- `add_footer`: Add file footers
- `indent`: Fix code indentation

### Advanced Transformations
- `refactor_variable`: Rename variables safely
- `extract_function`: Extract code blocks into functions
- `inline_function`: Inline function calls
- `move_code`: Move code between files

---

## Known Limitations

1. **API Dependency**: Some advanced features require API key
2. **Complex Refactoring**: Limited support for complex architectural changes
3. **Language Support**: Not all languages have full AST support
4. **Conflict Resolution**: May require manual intervention for complex conflicts
5. **Performance**: Large codebases may process slowly
6. **Validation**: Code validation is limited to syntax checking

---

## Testing

```bash
# Test server startup
npx @morph/fast-apply-mcp

# Test tool availability
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | npx @morph/fast-apply-mcp

# Test transformation functionality (requires MCP client)
```

---

## Use Cases

1. **Code Refactoring**: Safely refactor large codebases
2. **Migration**: Migrate code between frameworks or versions
3. **Standardization**: Enforce coding standards across projects
4. **Automation**: Automate repetitive code modifications
5. **Code Generation**: Generate boilerplate code
6. **Legacy Modernization**: Update legacy code to modern patterns

---

## Performance Tips

### Efficient Transformations
- Use batch operations for multiple files
- Specify language for better parsing
- Use dry_run to preview changes
- Apply transformations incrementally

### Large Codebases
- Process files in smaller batches
- Use specific file patterns
- Monitor memory usage
- Implement change validation

---

## Integration with Development Workflows

### CI/CD Integration
```yaml
# GitHub Actions example
- name: Apply Code Transformations
  run: |
    echo '{"tool": "batch_transform", "arguments": {"paths": ["src/"], "transformation": {"type": "add_header", "header": "// Auto-generated"}}}' | \
    npx @morph/fast-apply-mcp
```

### Pre-commit Hooks
```bash
# Apply formatting and standards
morph-transform --config .morphrc.json --path src/
```

---

## Related Documentation

- [Morph Documentation](https://morph.sh/)
- [Morph Fast Apply MCP Server](https://github.com/morph-artifacts/fast-apply-mcp)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [Code Transformation Best Practices](https://morph.sh/docs/transformations)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 1A)
**Next Steps**: Register with jarvis and test functionality
