# Filesystem MCP Server

**GitHub**: https://github.com/modelcontextprotocol/servers/tree/main/src/filesystem
**Package**: @modelcontextprotocol/server-filesystem
**Language**: TypeScript/JavaScript
**Transport**: stdio
**Status**: ✅ Official MCP Reference Server

---

## Overview

A Model Context Protocol server for filesystem operations, providing secure file access with configurable controls. This server enables LLMs to read, write, and manage files and directories.

**Key Features:**
- Read/write files with encoding support
- Directory creation and listing
- File moving and searching
- File metadata access
- Dynamic directory access control via MCP Roots
- Secure access controls

---

## Installation

### Using npm (via npx)
```bash
npx @modelcontextprotocol/server-filesystem
```

### Direct installation
```bash
npm install -g @modelcontextprotocol/server-filesystem
```

---

## Configuration

**Directory Access Control**: The server uses a flexible directory access control system with two methods:

### Method 1: Command-line Arguments (Server-side)
Specify allowed directories when starting the server:
```bash
mcp-server-filesystem /path/to/dir1 /path/to/dir2
```

### Method 2: MCP Roots (Recommended - Client-side)
MCP clients that support [Roots](https://modelcontextprotocol.io/docs/learn/client-concepts#roots) can dynamically update the allowed directories.

**Important**: If server starts without command-line arguments AND client doesn't support roots protocol (or provides empty roots), the server will throw an error during initialization.

---

## Available Tools

### File Operations
1. `read_file`
   - Read file contents with optional line range
   - Inputs: `path` (required), `encoding` (optional), `start_line` (optional), `end_line` (optional)
   - Returns: File contents with line numbers

2. `write_file`
   - Write content to files (create or overwrite)
   - Inputs: `path` (required), `content` (required), `encoding` (optional)
   - Returns: Success confirmation

3. `edit_file`
   - Apply edits to existing files
   - Inputs: `path` (required), `edits` (array of edit objects)
   - Returns: Applied changes confirmation

### Directory Operations
4. `list_directory`
   - List directory contents with metadata
   - Inputs: `path` (required), `recursive` (optional)
   - Returns: Directory listing with file types, sizes, modification times

5. `create_directory`
   - Create new directories
   - Inputs: `path` (required)
   - Returns: Success confirmation

6. `move_file`
   - Move files or directories
   - Inputs: `source` (required), `destination` (required)
   - Returns: Success confirmation

### Search and Metadata
7. `search_files`
   - Search for files by pattern
   - Inputs: `path` (required), `pattern` (required), `recursive` (optional)
   - Returns: Matching file paths

8. `get_file_info`
   - Get detailed file metadata
   - Inputs: `path` (required)
   - Returns: File size, permissions, modification time, etc.

---

## Usage Examples

### Reading Files
```json
{
  "tool": "read_file",
  "arguments": {
    "path": "README.md"
  }
}
```

```json
{
  "tool": "read_file",
  "arguments": {
    "path": "src/main.js",
    "start_line": 10,
    "end_line": 50
  }
}
```

### Writing Files
```json
{
  "tool": "write_file",
  "arguments": {
    "path": "config.json",
    "content": "{\n  \"version\": \"1.0.0\",\n  \"name\": \"my-app\"\n}"
  }
}
```

### Directory Operations
```json
{
  "tool": "list_directory",
  "arguments": {
    "path": "src",
    "recursive": true
  }
}
```

```json
{
  "tool": "create_directory",
  "arguments": {
    "path": "new-folder"
  }
}
```

### File Search
```json
{
  "tool": "search_files",
  "arguments": {
    "path": ".",
    "pattern": "*.js",
    "recursive": true
  }
}
```

### File Editing
```json
{
  "tool": "edit_file",
  "arguments": {
    "path": "config.js",
    "edits": [
      {
        "oldText": "const port = 3000;",
        "newText": "const port = process.env.PORT || 3000;"
      }
    ]
  }
}
```

---

## Security Features

### Access Control
- **Directory Restrictions**: Only access allowed directories
- **Path Traversal Protection**: Prevents `../` attacks
- **Symbolic Link Handling**: Safe handling of symlinks
- **Permission Checks**: Respects file system permissions

### Safe Defaults
- **No Root Access**: Cannot access system directories by default
- **Encoding Validation**: Validates file encodings
- **Size Limits**: Prevents reading extremely large files
- **Error Sanitization**: Safe error messages

---

## Known Limitations

1. **Directory Restrictions**: Can only access explicitly allowed directories
2. **No Binary Files**: Limited support for binary file operations
3. **Single Encoding**: Files must be readable as text
4. **No Concurrent Access**: No file locking mechanisms
5. **Path Resolution**: Relative paths resolved against allowed directories
6. **No Network Files**: Local filesystem only, no network drives

---

## Testing

```bash
# Test server startup with allowed directory
mcp-server-filesystem /home/user/projects

# Test tool availability
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | mcp-server-filesystem /home/user/projects

# Test file operations (requires MCP client)
```

---

## Use Cases

1. **Code Analysis**: Read and analyze source code files
2. **Configuration Management**: Manage application configs
3. **Documentation**: Create and update documentation files
4. **Project Navigation**: Explore project structure
5. **File Transformation**: Batch file operations and edits

---

## Docker Configuration

When running in Docker, mount host directories:
```bash
docker run -v /host/path:/container/path mcp-server-filesystem /container/path
```

**Docker Compose Example:**
```yaml
services:
  mcp-filesystem:
    image: mcp-server-filesystem
    volumes:
      - ./project:/workspace
    command: ["/workspace"]
```

---

## Integration with MCP Roots

**Client-side Roots Configuration:**
```json
{
  "mcpServers": {
    "filesystem": {
      "command": "mcp-server-filesystem",
      "args": [],
      "roots": [
        {
          "uri": "file:///home/user/project1",
          "name": "Project 1"
        },
        {
          "uri": "file:///home/user/project2",
          "name": "Project 2"
        }
      ]
    }
  }
}
```

---

## Error Handling

Common errors and solutions:
- **"No allowed directories"**: Configure allowed directories or MCP roots
- **"File not found"**: Check file path and permissions
- **"Access denied"**: Verify directory access permissions
- **"Invalid encoding"**: Use supported text encodings (UTF-8, ASCII, etc.)

---

## Related Documentation

- [Official MCP Servers Repository](https://github.com/modelcontextprotocol/servers)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [Filesystem Server Source Code](https://github.com/modelcontextprotocol/servers/tree/main/src/filesystem)
- [MCP Roots Documentation](https://modelcontextprotocol.io/docs/learn/client-concepts#roots)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 1A)
**Next Steps**: Register with jarvis and test functionality
