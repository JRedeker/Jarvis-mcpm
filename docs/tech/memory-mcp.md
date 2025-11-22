# Memory MCP Server

**GitHub**: https://github.com/modelcontextprotocol/servers/tree/main/src/memory
**Package**: @modelcontextprotocol/server-memory
**Language**: TypeScript/JavaScript
**Transport**: stdio
**Status**: ✅ Official MCP Reference Server

---

## Overview

A basic implementation of persistent memory using a local knowledge graph. This lets Claude remember information about the user across chats.

**Key Features:**
- Knowledge graph-based persistent memory system
- Entity and relation management
- Cross-session memory persistence
- Simple graph-based data model

---

## Installation

### Using npm (via npx)
```bash
npx @modelcontextprotocol/server-memory
```

### Direct installation
```bash
npm install -g @modelcontextprotocol/server-memory
```

---

## Configuration

The memory server requires no environment variables or API keys. It runs with default settings and stores data locally.

**Data Storage**: Memory data is persisted locally in a knowledge graph format.

---

## Core Concepts

### Entities
Entities are the primary nodes in the knowledge graph. Each entity has:
- A unique name (identifier)
- An entity type (e.g., "person", "organization", "event")
- A list of observations

**Example Entity:**
```json
{
  "name": "John_Smith",
  "entityType": "person",
  "observations": ["Speaks fluent Spanish", "Works at Anthropic"]
}
```

### Relations
Relations define directed connections between entities. They are always stored in active voice and describe how entities interact or relate to each other.

**Example Relation:**
```json
{
  "from": "John_Smith",
  "to": "Anthropic",
  "relationType": "works_at"
}
```

---

## Available Tools

### Entity Management
1. `memory_create_entity`
   - Create a new entity in the knowledge graph
   - Inputs: name, entityType, observations (array)

2. `memory_get_entity`
   - Retrieve an entity by name
   - Inputs: name

3. `memory_update_entity`
   - Update an existing entity
   - Inputs: name, observations (array to add)

4. `memory_delete_entity`
   - Delete an entity and its relations
   - Inputs: name

### Relation Management
5. `memory_create_relation`
   - Create a relation between entities
   - Inputs: from, to, relationType

6. `memory_get_relations`
   - Get relations for an entity
   - Inputs: name, relationType (optional)

7. `memory_delete_relation`
   - Delete a specific relation
   - Inputs: from, to, relationType

### Search and Discovery
8. `memory_search`
   - Search for entities and relations
   - Inputs: query

---

## Usage Examples

### Creating Entities
```json
{
  "tool": "memory_create_entity",
  "arguments": {
    "name": "Alice_Developer",
    "entityType": "person",
    "observations": ["Senior software engineer", "Specializes in Python", "Likes coffee"]
  }
}
```

### Creating Relations
```json
{
  "tool": "memory_create_relation",
  "arguments": {
    "from": "Alice_Developer",
    "to": "Python_Programming",
    "relationType": "specializes_in"
  }
}
```

### Searching Memory
```json
{
  "tool": "memory_search",
  "arguments": {
    "query": "python developer"
  }
}
```

### Updating Entity Observations
```json
{
  "tool": "memory_update_entity",
  "arguments": {
    "name": "Alice_Developer",
    "observations": ["Recently learned TypeScript", "Working on MCP project"]
  }
}
```

---

## Knowledge Graph Structure

The memory system organizes information as a directed graph:

```
[Alice_Developer] --specializes_in--> [Python_Programming]
     |
     | --works_at--> [Tech_Company]
     |
     | --likes--> [Coffee]
```

This structure allows for:
- Flexible relationship modeling
- Efficient traversal and search
- Natural language query support
- Persistent storage across sessions

---

## Known Limitations

1. **Local Storage**: Data is stored locally, not cloud-synced
2. **No Schema Validation**: Flexible but less structured than database systems
3. **Simple Search**: Basic text search, no advanced querying
4. **Memory Management**: No automatic cleanup of old/irrelevant data
5. **Single User**: Designed for single-user scenarios, not multi-user

---

## Testing

```bash
# Test server startup
npx @modelcontextprotocol/server-memory

# Test tool availability
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | npx @modelcontextprotocol/server-memory

# Test memory operations (requires MCP client)
```

---

## Use Cases

1. **User Preference Tracking**: Remember user preferences and settings
2. **Project Context**: Maintain context about ongoing projects
3. **Learning Progress**: Track what user has learned or accomplished
4. **Relationship Mapping**: Map connections between people, organizations, topics
5. **Session Continuity**: Maintain continuity across multiple chat sessions

---

## Related Documentation

- [Official MCP Servers Repository](https://github.com/modelcontextprotocol/servers)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [Memory Server Source Code](https://github.com/modelcontextprotocol/servers/tree/main/src/memory)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 1A)
**Next Steps**: Register with jarvis and test functionality
