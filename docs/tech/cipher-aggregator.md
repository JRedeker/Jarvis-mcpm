# Cipher

### Set Up Cipher Development Environment

Source: https://github.com/campfirein/cipher/blob/main/CONTRIBUTING.md

Installs dependencies using pnpm, sets up environment variables by copying an example file and requiring API keys, builds the project, and verifies the setup by running tests.

```bash
# Install dependencies
pnpm install

# Set up environment variables
cp .env.example .env
# Edit .env and add your API keys (at least one required):
# OPENAI_API_KEY=your_openai_key
# ANTHROPIC_API_KEY=your_anthropic_key

# Build the project
pnpm run build

# Verify setup by running tests
pnpm test
```

--------------------------------

### Cipher CLI: Development Setup Examples

Source: https://github.com/campfirein/cipher/blob/main/docs/cli-reference.md

Provides examples for setting up Cipher in a development environment, including using in-memory storage, starting an API server, and testing MCP integration with different transport types.

```bash
# Start with in-memory storage for testing
VECTOR_STORE_TYPE=in-memory cipher

# Start API server for development
cipher --mode api --port 8080 --cors

# Test MCP integration
cipher --mode mcp

# Test MCP with SSE transport
cipher --mode mcp --mcp-transport-type sse --mcp-port 4000

# Test MCP with streamable-HTTP transport
cipher --mode mcp --mcp-transport-type streamable-http --mcp-port 5000
```

--------------------------------

### Install and Use Cipher CLI

Source: https://github.com/campfirein/cipher/blob/main/docs/cli-reference.md

Instructions for installing the Cipher CLI globally using npm or executing it directly with npx, including a help command example.

```bash
# Install globally
npm install -g @byterover/cipher

# Or use npx
npx @byterover/cipher --help
```

--------------------------------

### Install Cipher via Smithery for Claude

Source: https://github.com/campfirein/cipher/blob/main/README.md

Installs the Cipher project for use with Claude via the Smithery CLI. This command automates the setup process for integrating Cipher with Claude Desktop.

```bash
npx -y @smithery/cli install @campfirein/cipher --client claude
```

--------------------------------

### Environment Variables Setup

Source: https://github.com/campfirein/cipher/blob/main/docs/llm-providers.md

Example of setting up environment variables in a `.env` file for various LLM providers, including API keys.

```bash
# OpenAI
OPENAI_API_KEY=sk-your-openai-key

# Anthropic
ANTHROPIC_API_KEY=sk-ant-your-anthropic-key

# OpenRouter
OPENROUTER_API_KEY=sk-or-your-openrouter-key
```

--------------------------------

### Setup and Run Byterover Cipher with Docker

Source: https://github.com/campfirein/cipher/blob/main/README.md

Steps to clone the Byterover Cipher repository, configure environment variables, and start the application using Docker Compose. Includes instructions for testing the API and building the UI.

```bash
# Clone and setup
git clone https://github.com/campfirein/cipher.git
cd cipher

# Configure environment
cp .env.example .env
# Edit .env with your API keys

# Start with Docker
docker-compose up --build -d

# Test
curl http://localhost:3000/health

# To include the UI in the Docker build:
docker build --build-arg BUILD_UI=true .
```

--------------------------------

### Cipher CLI: Production Deployment Examples

Source: https://github.com/campfirein/cipher/blob/main/docs/cli-reference.md

Illustrates examples for deploying Cipher in a production environment, covering API servers, MCP servers with various transports, and running the Web UI.

```bash
# Production API server
cipher --mode api --port 3000 --host 0.0.0.0 --agent /etc/cipher/production.yml

# MCP server with custom timeout (stdio)
cipher --mode mcp --timeout 120000 --agent /etc/cipher/mcp-config.yml

# Production MCP server with SSE transport
cipher --mode mcp --mcp-transport-type sse --mcp-port 3000 --agent /etc/cipher/mcp-config.yml

# Production MCP server with streamable-HTTP transport
cipher --mode mcp --mcp-transport-type streamable-http --mcp-port 3001 --agent /etc/cipher/mcp-config.yml

# Web UI for team access
cipher --mode ui --ui-port 80 --ui-host 0.0.0.0
```

--------------------------------

### Start Cipher Server (Bash)

Source: https://github.com/campfirein/cipher/blob/main/src/app/ui/README.md

This command sequence starts the Cipher backend API server. It first builds the project, links it globally, and then launches the server on port 3001.

```Bash
pnpm run build && npm link && cipher --mode server
```

--------------------------------

### Copy Example Environment File

Source: https://github.com/campfirein/cipher/blob/main/README.md

This command copies the example environment file to a new file named .env, which is typically used for local development and configuration.

```bash
cp .env.example .env
```

--------------------------------

### Install Byterover Cipher via NPM

Source: https://github.com/campfirein/cipher/blob/main/README.md

Instructions for installing the Byterover Cipher package globally or locally within a project using npm.

```bash
# Install globally
npm install -g @byterover/cipher

# Or install locally in your project
npm install @byterover/cipher
```

--------------------------------

### Start and Manage Cipher Sessions

Source: https://github.com/campfirein/cipher/blob/main/docs/cli-reference.md

Demonstrates how to start a new workspace memory session and switch between different team sessions using the Cipher CLI.

```bash
USE_WORKSPACE_MEMORY=true cipher --new-session team-project

cipher
> /session switch team-frontend
> /session switch team-backend
```

--------------------------------

### Start UI Development Server (Bash)

Source: https://github.com/campfirein/cipher/blob/main/src/app/ui/README.md

This command starts the Next.js development server for the Cipher UI, typically running on port 3000. It allows for real-time updates during development.

```Bash
pnpm run dev
```

--------------------------------

### Launch Cipher Assistant

Source: https://github.com/campfirein/cipher/blob/main/examples/01-kimi-k2-coding-assistant/README.md

This bash command navigates to the example directory and launches the Cipher coding assistant using a specific configuration file (`cipher.yml`). Ensure you are in the correct directory before executing.

```bash
# Navigate to the example directory
cd examples/01-kimi-k2-coding-assistant

# Start Cipher with this configuration
cipher --agent ./cipher.yml
```

--------------------------------

### Configure Qdrant Cloud for Cipher

Source: https://github.com/campfirein/cipher/blob/main/docs/vector-stores.md

Set up Qdrant Cloud by providing the cluster URL and API key in your environment variables for Cipher integration. This is the easiest way to get started with Qdrant.

```env
VECTOR_STORE_TYPE=qdrant
VECTOR_STORE_URL=https://your-cluster.qdrant.io
VECTOR_STORE_API_KEY=your-qdrant-api-key
```

--------------------------------

### Configure Ollama Embeddings (Local)

Source: https://github.com/campfirein/cipher/blob/main/docs/embedding-configuration.md

Example configuration for running Ollama embeddings locally. Specifies the model and optionally the base URL.

```yaml
embedding:
  type: ollama
  model: nomic-embed-text
  baseUrl: http://localhost:11434  # Optional, defaults to this
```

--------------------------------

### Build and Link Byterover Cipher from Source

Source: https://github.com/campfirein/cipher/blob/main/README.md

Commands to install dependencies, build the project, and link it locally for development when running Byterover Cipher from its source code.

```bash
pnpm i && pnpm run build && npm link
```

--------------------------------

### Configure LM Studio Embeddings (Local)

Source: https://github.com/campfirein/cipher/blob/main/docs/embedding-configuration.md

Example configuration for LM Studio embeddings, allowing local model selection and optional dimension specification.

```yaml
embedding:
  type: lmstudio
  model: nomic-embed-text-v1.5  # or bge-large, bge-base, bge-small
  baseUrl: http://localhost:1234/v1  # Optional, defaults to this
  # dimensions: 768  # Optional, auto-detected based on model
```

--------------------------------

### Configure Qdrant Vector Store (Remote)

Source: https://github.com/campfirein/cipher/blob/main/docs/configuration.md

Bash configuration for connecting to Qdrant Cloud, specifying the type, URL, and API key.

```bash
# Remote (Qdrant Cloud)
VECTOR_STORE_TYPE=qdrant
VECTOR_STORE_URL=your-qdrant-endpoint
VECTOR_STORE_API_KEY=your-qdrant-api-key
```

--------------------------------

### Cipher CLI Execution Modes Examples

Source: https://github.com/campfirein/cipher/blob/main/docs/cli-reference.md

Provides examples for running Cipher in its different execution modes: CLI for interactive chat, API for RESTful services, MCP for Model Context Protocol, and UI for a web interface.

```bash
# CLI Mode Examples
cipher
cipher "Add this to memory: CORS issues in Vite are usually solved by configuring the proxy"
cipher --agent ./my-config.yml "What do I know about authentication?"

# API Mode Examples
cipher --mode api
cipher --mode api --port 8080 --host 0.0.0.0 --cors
cipher --mode api --agent ./production-config.yml

# MCP Mode Examples
cipher --mode mcp
cipher --mode mcp --mcp-transport-type stdio
cipher --mode mcp --mcp-transport-type sse --mcp-port 3000
cipher --mode mcp --mcp-transport-type streamable-http --mcp-port 3001
cipher --mode mcp --strict

# UI Mode Examples
cipher --mode ui
cipher --mode ui --ui-port 8080
cipher --mode ui --ui-host 0.0.0.0 --ui-port 3001
```

--------------------------------

### Configure Qdrant Vector Store (Local)

Source: https://github.com/campfirein/cipher/blob/main/docs/configuration.md

Bash configuration for connecting to a local Qdrant instance via Docker, specifying the type, host, and port.

```bash
# Local (Docker)
VECTOR_STORE_TYPE=qdrant
VECTOR_STORE_HOST=localhost
VECTOR_STORE_PORT=6333
VECTOR_STORE_URL=http://localhost:6333
```

--------------------------------

### Automatic Embedding Fallback Example

Source: https://github.com/campfirein/cipher/blob/main/docs/embedding-configuration.md

Demonstrates Cipher's automatic embedding provider selection when no explicit embedding configuration is present, falling back to Voyage AI.

```yaml
# Example: Only LLM configured, embedding auto-selected
llm:
  provider: anthropic
  model: claude-3-5-sonnet-20241022
  apiKey: $ANTHROPIC_API_KEY
# No embedding config = auto-fallback to Voyage
```

--------------------------------

### Configure Azure OpenAI Embeddings

Source: https://github.com/campfirein/cipher/blob/main/docs/embedding-configuration.md

Example configuration for Azure OpenAI embeddings, similar to OpenAI but requires an Azure endpoint and API key.

```yaml
embedding:
  type: openai
  model: text-embedding-3-small
  apiKey: $AZURE_OPENAI_API_KEY
  baseUrl: $AZURE_OPENAI_ENDPOINT
```

--------------------------------

### ChromaDB Basic and SSL Configuration

Source: https://github.com/campfirein/cipher/blob/main/docs/vector-stores.md

Environment variables for configuring ChromaDB connection, including basic setup with URL and advanced setup with SSL/TLS enabled.

```bash
# Basic setup
VECTOR_STORE_TYPE=chroma
VECTOR_STORE_URL=http://localhost:8000

# With SSL/TLS
VECTOR_STORE_TYPE=chroma
VECTOR_STORE_HOST=localhost
VECTOR_STORE_PORT=8000
VECTOR_STORE_SSL=true
```

--------------------------------

### Test Cipher MCP Connection

Source: https://github.com/campfirein/cipher/blob/main/docs/mcp-integration.md

Demonstrates how to test the Cipher MCP connection by starting the server manually and then using the MCP client tools.

```bash
# Start Cipher MCP server manually
export OPENAI_API_KEY="sk-your-key"
cipher --mode mcp

# Test with MCP client tools
npx @modelcontextprotocol/inspector cipher --mode mcp
```

--------------------------------

### MCP Aggregator Hub Configuration

Source: https://github.com/campfirein/cipher/blob/main/docs/examples.md

Configuration for Cipher's MCP aggregator hub, demonstrating integration with multiple MCP servers. It shows how to set up the filesystem MCP server with its command and arguments.

```yaml
mcpServers:
  filesystem:
    type: stdio
    command: npx
    args: ['-y', '@modelcontextprotocol/server-filesystem', '.']
```

--------------------------------

### Configure Milvus Vector Store (Local)

Source: https://github.com/campfirein/cipher/blob/main/docs/configuration.md

Bash configuration for connecting to a local Milvus instance via Docker, specifying the type, host, and port.

```bash
# Local (Docker)
VECTOR_STORE_TYPE=milvus
VECTOR_STORE_HOST=localhost
VECTOR_STORE_PORT=19530
```

--------------------------------

### Configure OpenAI Embeddings

Source: https://github.com/campfirein/cipher/blob/main/docs/embedding-configuration.md

Example configuration for using OpenAI embeddings with Cipher. Requires an OpenAI API key and specifies the model and type.

```yaml
embedding:
  type: openai
  model: text-embedding-3-small
  apiKey: $OPENAI_API_KEY
```

--------------------------------

### Run Qdrant Locally with Docker for Cipher

Source: https://github.com/campfirein/cipher/blob/main/docs/vector-stores.md

Deploy Qdrant locally using Docker for development and testing with Cipher. Supports basic setup and persistent storage configurations.

```bash
# Basic setup (data lost on removing the container)
docker run -d --name qdrant-basic -p 6333:6333 qdrant/qdrant

# With persistent storage
docker run -d --name qdrant-storage -v ./qdrant-data:/qdrant/storage -p 6333:6333 qdrant/qdrant
```

```env
# .env configuration
VECTOR_STORE_TYPE=qdrant
VECTOR_STORE_HOST=localhost
VECTOR_STORE_PORT=6333
VECTOR_STORE_URL=http://localhost:6333
```

--------------------------------

### Use Byterover Cipher CLI Commands

Source: https://github.com/campfirein/cipher/blob/main/README.md

Examples of how to use the Byterover Cipher command-line interface for interactive mode, one-shot commands, API server mode, MCP server mode, and Web UI mode. Includes notes on environment variables and session management.

```bash
# Interactive mode
cipher

# One-shot command
cipher "Add this to memory as common causes of 'CORS error' in local dev with Vite + Express."

# API server mode
cipher --mode api

# MCP server mode
cipher --mode mcp

# Web UI mode
cipher --mode ui

# Start a fresh session:
/session new <session-name>
```

--------------------------------

### Configure Gemini Embeddings

Source: https://github.com/campfirein/cipher/blob/main/docs/embedding-configuration.md

Example configuration for using Gemini embeddings with Cipher. Requires a Gemini API key and specifies the model and type.

```yaml
embedding:
  type: gemini
  model: gemini-embedding-001
  apiKey: $GEMINI_API_KEY
```

--------------------------------

### Cipher CLI Global Options Examples

Source: https://github.com/campfirein/cipher/blob/main/docs/cli-reference.md

Illustrates the use of core and mode-specific global options for the Cipher CLI, such as setting the execution mode, specifying configuration files, managing sessions, and configuring server ports and hosts.

```bash
# Using core options
cipher --mode api
cipher --agent /path/to/config.yml
cipher --new-session project-alpha
cipher --mode mcp --strict
cipher --help
cipher --version

# API mode options
cipher --mode api --port 8080 --host 0.0.0.0 --cors

# MCP mode options
cipher --mode mcp --mcp-transport-type sse --mcp-port 4000 --mcp-host 0.0.0.0 --mcp-dns-rebinding-protection
cipher --mode mcp --timeout 120000

# UI mode options
cipher --mode ui --ui-port 8080 --ui-host 0.0.0.0
```

--------------------------------

### Configure Qwen Embeddings with Fixed Dimensions

Source: https://github.com/campfirein/cipher/blob/main/docs/embedding-configuration.md

Example configuration for Qwen embeddings, which requires specifying fixed dimensions. Supports 1024, 768, or 512 dimensions.

```yaml
embedding:
  type: qwen
  model: text-embedding-v3
  apiKey: $QWEN_API_KEY
  dimensions: 1024  # Required: 1024, 768, or 512
```

--------------------------------

### Kimi K2 Coding Assistant Configuration

Source: https://github.com/campfirein/cipher/blob/main/docs/examples.md

Configuration for integrating Cipher with the Kimi K2 AI coding assistant. It specifies the LLM provider, model, API key, and a system prompt for AI-powered code assistance with persistent memory.

```yaml
llm:
  provider: openai
  model: gpt-4-turbo
  apiKey: $OPENAI_API_KEY

systemPrompt: 'You are a coding assistant with memory capabilities.'
```

--------------------------------

### Configure Milvus Vector Store (Remote)

Source: https://github.com/campfirein/cipher/blob/main/docs/configuration.md

Bash configuration for connecting to Zilliz Cloud (Milvus), specifying the type, endpoint, username, and password.

```bash
# Remote (Zilliz Cloud)
VECTOR_STORE_TYPE=milvus
VECTOR_STORE_URL=your-milvus-cluster-endpoint
VECTOR_STORE_USERNAME=your-zilliz-username
VECTOR_STORE_PASSWORD=your-zilliz-password
```

--------------------------------

### Commit Changes with Conventional Commits

Source: https://github.com/campfirein/cipher/blob/main/CONTRIBUTING.md

Example of committing changes using the conventional commit format, with examples for features, fixes, documentation, and tests.

```bash
git add .
git commit -m "feat: add persistent memory layer for agent sessions"
# Other examples:
# git commit -m "fix: resolve MCP connection timeout issues"
# git commit -m "docs: update configuration examples"
# git commit -m "test: add integration tests for memory persistence"
```

--------------------------------

### Configure Cipher Agent and LLM

Source: https://github.com/campfirein/cipher/blob/main/docs/configuration.md

This snippet shows the basic structure for configuring the Cipher agent, including LLM provider, model, API key, system prompt, and optional MCP server settings.

```yaml
# LLM Configuration
llm:
  provider: openai # openai, anthropic, openrouter, ollama, qwen
  model: gpt-4-turbo
  apiKey: $OPENAI_API_KEY

# System Prompt
systemPrompt: 'You are a helpful AI assistant with memory capabilities.'

# MCP Servers (optional)
mcpServers:
  filesystem:
    type: stdio
    command: npx
    args: ['-y', '@modelcontextprotocol/server-filesystem', '.']
```

--------------------------------

### Configure Additional Vector Store Settings

Source: https://github.com/campfirein/cipher/blob/main/docs/configuration.md

Bash configuration for general vector store settings, including collection name, dimension, and distance metric.

```bash
# Collection configuration
VECTOR_STORE_COLLECTION=knowledge_memory
VECTOR_STORE_DIMENSION=1536
VECTOR_STORE_DISTANCE=Cosine
```

--------------------------------

### Configure Qwen Embeddings with Dimensions

Source: https://github.com/campfirein/cipher/blob/main/docs/configuration.md

Configuration for Qwen embeddings, requiring the model, API key, and specifying dimensions.

```yaml
embedding:
  type: qwen
  model: text-embedding-v3
  apiKey: $QWEN_API_KEY
  dimensions: 1024  # Required: 1024, 768, or 512
```

--------------------------------

### Cipher Project Environment Variables

Source: https://github.com/campfirein/cipher/blob/main/README.md

Example environment variables for configuring the Cipher project. This includes settings for vector stores, chat history, AWS integration, and advanced options like logging and memory search.

```env
VECTOR_STORE_TYPE=qdrant
VECTOR_STORE_URL=https://your-cluster.qdrant.io
VECTOR_STORE_API_KEY=your-qdrant-api-key

CIPHER_PG_URL=postgresql://user:pass@localhost:5432/cipher_db

USE_WORKSPACE_MEMORY=true
WORKSPACE_VECTOR_STORE_COLLECTION=workspace_memory

AWS_ACCESS_KEY_ID=your-aws-access-key
AWS_SECRET_ACCESS_KEY=your-aws-secret-key
AWS_DEFAULT_REGION=us-east-1

CIPHER_LOG_LEVEL=info
REDACT_SECRETS=true

VECTOR_STORE_DIMENSION=1536
VECTOR_STORE_DISTANCE=Cosine
VECTOR_STORE_MAX_VECTORS=10000

SEARCH_MEMORY_TYPE=knowledge
DISABLE_REFLECTION_MEMORY=true
```

--------------------------------

### Configure OpenAI Embeddings

Source: https://github.com/campfirein/cipher/blob/main/docs/configuration.md

Configuration for using OpenAI embeddings, specifying the model and API key.

```yaml
embedding:
  type: openai
  model: text-embedding-3-small
  apiKey: $OPENAI_API_KEY
```

--------------------------------

### Run ChromaDB Locally with Docker for Cipher

Source: https://github.com/campfirein/cipher/blob/main/docs/vector-stores.md

Deploy ChromaDB locally using Docker for development and testing with Cipher. Supports both basic setup and configurations with persistent storage.

```bash
# Basic setup (data lost on removing the container)
docker run -d --name chroma-basic -p 8000:8000 chromadb/chroma

# With persistent storage
docker run -d --name chroma-storage -v ./chroma-data:/data -p 8000:8000 chromadb/chroma
```

```env
# .env configuration
VECTOR_STORE_TYPE=chroma
VECTOR_STORE_HOST=localhost
VECTOR_STORE_PORT=8000
VECTOR_STORE_URL=http://localhost:8000
```

--------------------------------

### Workspace Memory Team Progress Configuration

Source: https://github.com/campfirein/cipher/blob/main/docs/examples.md

Configuration for Cipher's workspace memory feature, enabling team collaboration and project progress tracking. It sets the LLM provider, model, API key, and a system prompt for a team-aware AI assistant.

```yaml
llm:
  provider: openai
  model: gpt-4-turbo
  apiKey: $OPENAI_API_KEY

# Workspace memory enabled
systemPrompt: 'You are a team-aware AI assistant with workspace memory.'
```

--------------------------------

### Run Milvus Locally with Docker for Cipher

Source: https://github.com/campfirein/cipher/blob/main/docs/vector-stores.md

Set up Milvus locally for development with Cipher by downloading and running the official standalone installation script. This script manages the Milvus Docker container and its dependencies.

```bash
# Download the official installation script
curl -sfL https://raw.githubusercontent.com/milvus-io/milvus/master/scripts/standalone_embed.sh -o standalone_embed.sh

# Start the Docker container
bash standalone_embed.sh start
```

```env
# .env configuration
VECTOR_STORE_TYPE=milvus
VECTOR_STORE_HOST=localhost
VECTOR_STORE_PORT=19530
```

--------------------------------

### Basic Cipher CLI Usage

Source: https://github.com/campfirein/cipher/blob/main/docs/cli-reference.md

Demonstrates the fundamental ways to interact with the Cipher CLI, including interactive mode, one-shot commands, and starting different server modes (API, MCP, UI).

```bash
# Interactive CLI mode
cipher

# One-shot command
cipher "Your prompt here"

# Server modes
cipher --mode api        # REST API server
cipher --mode mcp        # MCP server
cipher --mode ui         # Web UI server
```

--------------------------------

### Configure AWS Bedrock Embeddings with Fixed Dimensions

Source: https://github.com/campfirein/cipher/blob/main/docs/embedding-configuration.md

Example configuration for AWS Bedrock embeddings, requiring specific dimensions (1024, 512, or 256) and AWS credentials.

```yaml
embedding:
  type: aws-bedrock
  model: amazon.titan-embed-text-v2:0
  region: $AWS_REGION
  accessKeyId: $AWS_ACCESS_KEY_ID
  secretAccessKey: $AWS_SECRET_ACCESS_KEY
  dimensions: 1024  # Required: 1024, 512, or 256
```

--------------------------------

### Configure Voyage AI Embeddings (Fixed Dimensions)

Source: https://github.com/campfirein/cipher/blob/main/docs/embedding-configuration.md

Example configuration for Voyage AI embeddings, which uses fixed 1024 dimensions. Requires a Voyage API key.

```yaml
embedding:
  type: voyage
  model: voyage-3-large
  apiKey: $VOYAGE_API_KEY
  # Note: Voyage models use fixed 1024 dimensions
```

--------------------------------

### PostgreSQL Database Setup SQL

Source: https://github.com/campfirein/cipher/blob/main/docs/chat-history.md

SQL commands to create the PostgreSQL database and user, and grant necessary privileges for Cipher to operate.

```sql
-- Connect to PostgreSQL as superuser
CREATE DATABASE cipher_db;
CREATE USER cipher_user WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE cipher_db TO cipher_user;
```

```sql
-- Connect to cipher_db
GRANT USAGE, CREATE ON SCHEMA public TO cipher_user;
GRANT ALL ON ALL TABLES IN SCHEMA public TO cipher_user;
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO cipher_user;
```

--------------------------------

### Workspace Memory Team Progress Environment Variables

Source: https://github.com/campfirein/cipher/blob/main/docs/examples.md

Environment variables for enabling and configuring workspace memory for team collaboration in Cipher. It includes settings for enabling workspace memory, specifying the collection name, and disabling default memory.

```bash
# .env
USE_WORKSPACE_MEMORY=true
WORKSPACE_VECTOR_STORE_COLLECTION=team_progress
DISABLE_DEFAULT_MEMORY=false
```

--------------------------------

### Configure LM Studio Embeddings

Source: https://github.com/campfirein/cipher/blob/main/docs/configuration.md

Configuration for LM Studio embeddings, specifying the model and optional base URL. Dimensions are auto-detected.

```yaml
embedding:
  type: lmstudio
  model: nomic-embed-text-v1.5  # or bge-large, bge-base, bge-small
  baseUrl: http://localhost:1234/v1  # Optional, defaults to this
  # dimensions: 768  # Optional, auto-detected based on model
```

--------------------------------

### Configure Azure OpenAI Embeddings

Source: https://github.com/campfirein/cipher/blob/main/docs/configuration.md

Configuration for Azure OpenAI embeddings, specifying the model, API key, and endpoint URL.

```yaml
embedding:
  type: openai
  model: text-embedding-3-small
  apiKey: $AZURE_OPENAI_API_KEY
  baseUrl: $AZURE_OPENAI_ENDPOINT
```

--------------------------------

### Querying Team Progress and Displaying Results (Bash)

Source: https://github.com/campfirein/cipher/blob/main/examples/05-workspace-memory-team-progress/README.md

Shows an example of a natural language query for team progress and the expected structured output. The results are presented in a table format, detailing team member contributions, task status, and completion details.

```bash
# Query
"Show me what the backend team accomplished this week"

# Returns organized table with:
| Team Member | Feature/Task   | Status      | Completion | Notes               |
| ----------- | -------------- | ----------- | ---------- | ------------------- |
| Sarah       | API Gateway    | Completed   | 100%       | Deployed to staging |
| Mike        | Authentication | In Progress | 75%        | Testing phase       |
```

--------------------------------

### Disable Embeddings in Campfirein Cipher

Source: https://github.com/campfirein/cipher/blob/main/docs/embedding-configuration.md

This configuration disables all memory functionality, forcing the Cipher to operate in a chat-only mode. This eliminates the need for a vector database connection and simplifies the setup for basic chat assistant functionality.

```yaml
embedding:
  disabled: true
```

--------------------------------

### Start Cipher with SSE Transport

Source: https://github.com/campfirein/cipher/blob/main/docs/mcp-integration.md

This bash command starts the Cipher application in MCP mode, specifically configured to use the Server-Sent Events (SSE) transport protocol on port 4000.

```bash
cipher --mode mcp --mcp-transport-type sse --mcp-port 4000
```

--------------------------------

### Pinecone Cloud Configuration

Source: https://github.com/campfirein/cipher/blob/main/docs/vector-stores.md

Configuration for connecting to Pinecone's managed cloud service. Requires API key and collection name for basic setup.

```bash
# Basic configuration
VECTOR_STORE_TYPE=pinecone
VECTOR_STORE_API_KEY=your-pinecone-api-key
VECTOR_STORE_COLLECTION=your-index-name # Collection names are used as indexes in Pinecone
```

--------------------------------

### Configure AWS Bedrock Embeddings with Dimensions

Source: https://github.com/campfirein/cipher/blob/main/docs/configuration.md

Configuration for AWS Bedrock embeddings, including model, region, credentials, and required dimensions.

```yaml
embedding:
  type: aws-bedrock
  model: amazon.titan-embed-text-v2:0
  region: $AWS_REGION
  accessKeyId: $AWS_ACCESS_KEY_ID
  secretAccessKey: $AWS_SECRET_ACCESS_KEY
  dimensions: 1024  # Required: 1024, 512, or 256
```

--------------------------------

### PostgreSQL Connection Pooling Configuration

Source: https://github.com/campfirein/cipher/blob/main/docs/chat-history.md

Example of how to configure connection pooling for PostgreSQL in a .env file by setting the 'max_connections' parameter in the CIPHER_PG_URL.

```bash
# .env - Connection pooling
CIPHER_PG_URL="postgresql://user:pass@localhost:5432/cipher_db?max_connections=10"
```

--------------------------------

### Set Environment Variables for LLM Providers

Source: https://github.com/campfirein/cipher/blob/main/docs/embedding-configuration.md

Configure your environment by setting API keys and necessary credentials for various Large Language Model (LLM) providers. This includes OpenAI, Gemini, Qwen, Voyage AI, AWS Bedrock, and Azure OpenAI, enabling the Cipher to connect to and utilize these services.

```bash
# OpenAI
OPENAI_API_KEY=sk-your-openai-key

# Gemini
GEMINI_API_KEY=your-gemini-api-key

# Qwen
QWEN_API_KEY=your-qwen-api-key

# Voyage AI
VOYAGE_API_KEY=your-voyage-key

# AWS Bedrock
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key

# Azure OpenAI
AZURE_OPENAI_API_KEY=your-azure-key
AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com
```

--------------------------------

### Strict Memory Layer Configuration

Source: https://github.com/campfirein/cipher/blob/main/docs/examples.md

Configuration for Cipher's strict memory layer, demonstrating controlled access and precise memory operations. It includes LLM settings and embedding configurations for memory integrity checks.

```yaml
llm:
  provider: anthropic
  model: claude-3-5-sonnet-20241022
  apiKey: $ANTHROPIC_API_KEY

embedding:
  type: openai
  model: text-embedding-3-small
  apiKey: $OPENAI_API_KEY
```

--------------------------------

### Disable Cipher Embeddings

Source: https://github.com/campfirein/cipher/blob/main/docs/configuration.md

Disables all memory-related tools and operates Cipher in chat-only mode.

```yaml
embedding:
  disabled: true
```

--------------------------------

### AI Agent Rules for Workspace Memory Integration (Bash)

Source: https://github.com/campfirein/cipher/blob/main/examples/05-workspace-memory-team-progress/README.md

These rules guide AI agents on how to effectively use Cipher's workspace memory tools. They specify when to search for project updates, store progress information, and query team member activities, ensuring efficient team progress tracking.

```bash
+ When users are about to implement a new feature in the project, please use `cipher_workspace_search` to get the latest updates on the project's progress across the team.

+ When users finished a new feature or implemented it partially, ideally at the end of that generation, store the progress using `cipher_workspace_memory`. Remember to provide sufficient information for the storage payload

+ When users ask for team members' progress, provide enough context for the `cipher_workspace_search` tool, such as project name, repository, and present the results in a table format.
```

--------------------------------

### Cipher Memory Operations

Source: https://github.com/campfirein/cipher/blob/main/docs/cli-reference.md

Provides examples for searching specific memories and adding structured knowledge to the Cipher memory. It also shows how to view memory statistics.

```bash
# Search specific memories
cipher "What do I know about React hooks?"

# Add structured knowledge
cipher "Remember: Next.js 13+ uses app directory structure with layout.tsx and page.tsx files"

# Memory statistics
cipher
> /memory stats
> /stats
```

--------------------------------

### Launch Cipher as MCP Server

Source: https://github.com/campfirein/cipher/blob/main/examples/03-strict-memory-layer/README.md

Start the Cipher service in MCP (Multi-Client Protocol) mode, specifying the agent configuration file. This command initiates the memory service.

```bash
cipher --mode mcp --agent ./examples/03-strict-memory-layer/cipher.yml
```

--------------------------------

### Troubleshoot MCP Connection Issues

Source: https://github.com/campfirein/cipher/blob/main/docs/cli-reference.md

Provides commands to test MCP server connections, including setting custom timeouts and specifying transport types (SSE, Streamable-HTTP). It also includes curl examples for checking HTTP endpoints.

```bash
# Test MCP server with custom timeout
cipher --mode mcp --timeout 120000

# Test specific transport types
cipher --mode mcp --mcp-transport-type sse --mcp-port 3000
cipher --mode mcp --mcp-transport-type streamable-http --mcp-port 3001

# Check HTTP endpoints
# SSE (establish stream)
curl -N -H "Accept: text/event-stream" http://localhost:3000/sse
# Streamable-HTTP (client must send Accept headers)
curl -sS -X POST \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":"1","method":"ping","params":{}}' \
  http://localhost:3001/http
```

--------------------------------

### Start Cipher with Streamable-HTTP Transport

Source: https://github.com/campfirein/cipher/blob/main/docs/mcp-integration.md

This bash command initiates the Cipher process in MCP mode, configured to use the Streamable-HTTP transport protocol on port 4000.

```bash
cipher --mode mcp --mcp-transport-type streamable-http --mcp-port 4000
```

--------------------------------

### Configure Voyage Embeddings

Source: https://github.com/campfirein/cipher/blob/main/docs/configuration.md

Configuration for Voyage embeddings, specifying the model and API key. Note that Voyage models use fixed 1024 dimensions.

```yaml
embedding:
  type: voyage
  model: voyage-3-large
  apiKey: $VOYAGE_API_KEY
  # Note: Voyage models use fixed 1024 dimensions
```

--------------------------------

### MCP Client Configuration for Cipher

Source: https://github.com/campfirein/cipher/blob/main/examples/03-strict-memory-layer/README.md

Example JSON configuration for an MCP client to connect to the Cipher server. It includes server command, arguments, and environment variables for API keys, vector store type, URL, and collection names.

```json
{
  "mcpServers": {
    "cipher": {
      "command": "cipher",
      "args": [
        "--mode", "mcp",
        "--agent", "./examples/03-strict-memory-layer/cipher.yml"
      ],
      "env": {
        "MCP_SERVER_MODE": "default",
        "OPENAI_API_KEY": "sk-வுகளை",
        "VECTOR_STORE_TYPE": "milvus",
        "VECTOR_STORE_URL": "...",
        "VECTOR_STORE_API_KEY": "...",
        "VECTOR_STORE_USERNAME": "...",
        "VECTOR_STORE_PASSWORD": "...",
        "VECTOR_STORE_COLLECTION": "knowledge_memory",
        "REFLECTION_VECTOR_STORE_COLLECTION": "reflection_memory",
        "DISABLE_REFLECTION_MEMORY": "true"  # default: true
      }
    }
  }
}
```

--------------------------------

### Configure Reflection Memory

Source: https://github.com/campfirein/cipher/blob/main/docs/configuration.md

Sets the collection name for the reflection memory vector store and controls whether reflection memory is disabled. The default for disabling reflection memory is true.

```shell
REFLECTION_VECTOR_STORE_COLLECTION=reflection_memory
DISABLE_REFLECTION_MEMORY=true  # default: true
```

--------------------------------

### Extracting Progress Updates from Natural Language (Bash)

Source: https://github.com/campfirein/cipher/blob/main/examples/05-workspace-memory-team-progress/README.md

Demonstrates how Cipher automatically extracts structured information from natural language input regarding task completion. This example shows the extraction of team member, feature, and status from a simple sentence.

```bash
# Input
"Alice completed the payment integration yesterday"

# Automatically extracts
- teamMember: "Alice"
- feature: "payment integration"
- status: "completed"
```

--------------------------------

### Clear Ports (Bash)

Source: https://github.com/campfirein/cipher/blob/main/src/app/ui/README.md

This bash command clears specified ports (3000-3001) by finding and killing processes listening on them. It's a common utility for freeing up ports before starting development servers.

```Bash
lsof -ti:3000-3001 | xargs kill -9
```

--------------------------------

### Basic Byterover Cipher Configuration (YAML)

Source: https://github.com/campfirein/cipher/blob/main/README.md

A sample YAML configuration file for Byterover Cipher, demonstrating settings for LLM provider, model, API key, system prompt, and optional MCP servers.

```yaml
# LLM Configuration
llm:
  provider: openai # openai, anthropic, openrouter, ollama, qwen
  model: gpt-4-turbo
  apiKey: $OPENAI_API_KEY

# System Prompt
systemPrompt: 'You are a helpful AI assistant with memory capabilities.'

# MCP Servers (optional)
mcpServers:
  filesystem:
    type: stdio
    command: npx
    args: ['-y', '@modelcontextprotocol/server-filesystem', '.']
```

--------------------------------

### Cipher Configuration: File System Access

Source: https://github.com/campfirein/cipher/blob/main/examples/01-kimi-k2-coding-assistant/README.md

This YAML snippet configures the file system access for Cipher. The `args` parameter specifies the root directory that the assistant can access, which should be customized to your project's workspace.

```yaml
filesystem:
  args:
    - /Users  # Change to your workspace directory (e.g., /home/user/projects)
```

--------------------------------

### Add File-Based Prompt Provider

Source: https://github.com/campfirein/cipher/blob/main/memAgent/prompt-provider-commands-examples.md

Adds or updates a file-based prompt provider, like 'project-guidelines', with options such as summarizing content. It indicates that an LLM summary has been generated and cached for the provider.

```bash
/prompt-providers add-file project-guidelines --summarize true
```

--------------------------------

### List Prompt Providers

Source: https://github.com/campfirein/cipher/blob/main/memAgent/prompt-provider-commands-examples.md

Lists active and available prompt providers in the Cipher CLI. It shows which providers are currently active and which are available but not yet loaded, with a hint on how to activate more.

```bash
/prompt-providers list
```

--------------------------------

### Verify MCP Server Help

Source: https://github.com/campfirein/cipher/blob/main/examples/01-kimi-k2-coding-assistant/README.md

These bash commands are used to verify the help information for the MCP (Model Context Protocol) servers, specifically `server-filesystem` and `firecrawl-mcp`. This is useful for troubleshooting connection issues.

```bash
# Verify MCP servers
npx -y @modelcontextprotocol/server-filesystem --help
npx -y firecrawl-mcp --help
```

--------------------------------

### Cipher CLI: System Commands

Source: https://github.com/campfirein/cipher/blob/main/docs/cli-reference.md

Lists essential system commands for the Cipher CLI, including displaying configuration, viewing statistics, accessing help, clearing the screen, and exiting the application.

```bash
/config
/stats
/help
/clear
/exit
```

--------------------------------

### Show All Prompt Providers

Source: https://github.com/campfirein/cipher/blob/main/memAgent/prompt-provider-commands-examples.md

Displays all prompt providers, including those that are enabled (active or available) and disabled. This command helps in managing the status of different providers.

```bash
/prompt-providers show-all
```

--------------------------------

### Byterover Cipher Environment Variables (.env)

Source: https://github.com/campfirein/cipher/blob/main/README.md

Template for the .env file used by Byterover Cipher, listing essential API keys for various LLM providers like OpenAI, Anthropic, Gemini, and Qwen.

```bash
# ====================
# API Keys (At least one required)
# ====================
OPENAI_API_KEY=sk-your-openai-api-key
ANTHROPIC_API_KEY=sk-ant-your-anthropic-key
GEMINI_API_KEY=your-gemini-api-key
QWEN_API_KEY=your-qwen-api-key
```

--------------------------------

### Enable Prompt Provider

Source: https://github.com/campfirein/cipher/blob/main/memAgent/prompt-provider-commands-examples.md

Enables a disabled prompt provider, such as 'project-guidelines', making it active or available for use. It confirms that the provider has been successfully enabled.

```bash
/prompt-providers enable project-guidelines
```

--------------------------------

### Clone Cipher Repository

Source: https://github.com/campfirein/cipher/blob/main/CONTRIBUTING.md

Steps to fork and clone the Cipher repository from GitHub. This involves forking the project to your account and then cloning your fork locally.

```bash
git clone https://github.com/YOUR_USERNAME/cipher.git
cd cipher
```

--------------------------------

### Cipher Configuration: LLM Model

Source: https://github.com/campfirein/cipher/blob/main/examples/01-kimi-k2-coding-assistant/README.md

This YAML snippet defines the Language Model (LLM) settings for Cipher. It specifies the provider as 'openrouter' and the model as 'moonshotai/kimi-k2'. The `maxIterations` parameter controls the number of attempts for complex tasks.

```yaml
llm:
  provider: openrouter
  model: moonshotai/kimi-k2  # Can change to other OpenRouter models
  maxIterations: 75          # Increase for complex tasks, decrease for speed
```

--------------------------------

### Show Current System Prompt

Source: https://github.com/campfirein/cipher/blob/main/memAgent/prompt-provider-commands-examples.md

Displays the current system prompt being used by the Cipher CLI, including its content and metadata like length and line count. This is useful for understanding the AI's current instructions.

```bash
/prompt
```

--------------------------------

### Configure Qdrant with Docker Compose for Cipher

Source: https://github.com/campfirein/cipher/blob/main/docs/vector-stores.md

Integrate Qdrant into your project using Docker Compose by defining the Qdrant service in your `docker-compose.yml` file.

```yaml
services:
  qdrant:
    image: qdrant/qdrant:latest
    ports:
      - "6333:6333"
    volumes:
      - qdrant_data:/qdrant/storage
    environment:
      - QDRANT__SERVICE__HTTP_PORT=6333

volumes:
  qdrant_data:
```

--------------------------------

### Cipher Configuration and Logging

Source: https://github.com/campfirein/cipher/blob/main/docs/cli-reference.md

Shows how to check the Cipher configuration and enable verbose logging for debugging purposes using environment variables.

```bash
# Check configuration
cipher --config

# Use verbose logging
DEBUG=cipher:* cipher --mode api
```

--------------------------------

### Create Feature Branch for Cipher

Source: https://github.com/campfirein/cipher/blob/main/CONTRIBUTING.md

Demonstrates how to create a new Git branch for feature development in the Cipher project, following a descriptive naming convention.

```bash
git checkout -b feature/your-descriptive-branch-name
```

--------------------------------

### Run Cipher Tests

Source: https://github.com/campfirein/cipher/blob/main/CONTRIBUTING.md

Command to execute all tests within the Cipher project to ensure code functionality and stability.

```bash
# Run all tests
pnpm test
```

--------------------------------

### Cipher Aggregator Server Configuration (cipher.yml)

Source: https://github.com/campfirein/cipher/blob/main/examples/04-mcp-aggregator-hub/README.md

This YAML configuration defines various MCP servers for the Cipher aggregator, including Exa Search (stdio), Context7 (streamable-http), Semgrep (streamable-http), and TaskMaster (stdio, disabled by default). It specifies transport types, commands, URLs, and environment variables for each server.

```yaml
# LLM Configuration
llm:
  provider: openai
  model: gpt-4o-mini
  apiKey: $OPENAI_API_KEY

mcpServers:
  # Exa Search (stdio transport)
  exa:
    type: stdio
    command: npx
    args: ["-y", "exa-mcp-server"]
    env:
      EXA_API_KEY: $EXA_API_KEY

  # Context7 (streamable-http transport)
  context7:
    type: "streamable-http"
    url: "https://mcp.context7.com/mcp"
    enabled: true

  # Semgrep (streamable-http transport)
  semgrep:
    type: "streamable-http"
    url: "https://mcp.semgrep.ai/mcp/"
    enabled: true

  # TaskMaster (disabled by default)
  taskmaster:
    type: stdio
    command: npx
    args: ["-y", "--package=task-master-ai", "task-master-ai"]
    enabled: false
    env:
      OPENAI_API_KEY: $OPENAI_API_KEY
```

--------------------------------

### Cipher CLI: Core Environment Variables

Source: https://github.com/campfirein/cipher/blob/main/docs/cli-reference.md

Details core environment variables for configuring Cipher, including API keys for language models (OpenAI, Anthropic), memory settings, vector store configurations, and chat history database URLs.

```bash
# LLM Configuration
OPENAI_API_KEY=sk-your-key
ANTHROPIC_API_KEY=sk-ant-your-key

# Memory Settings
USE_WORKSPACE_MEMORY=true
DISABLE_REFLECTION_MEMORY=true  # default: true

# Vector Store
VECTOR_STORE_TYPE=qdrant
VECTOR_STORE_URL=your-endpoint

# Chat History
CIPHER_PG_URL=postgresql://user:pass@host:5432/db
```

--------------------------------

### Show Prompt Statistics

Source: https://github.com/campfirein/cipher/blob/main/memAgent/prompt-provider-commands-examples.md

Displays performance statistics for the system prompt, including the number of providers used, total prompt length, generation time, and success status. This helps in evaluating prompt efficiency.

```bash
/prompt-stats
```

--------------------------------

### Update Prompt Provider Configuration

Source: https://github.com/campfirein/cipher/blob/main/memAgent/prompt-provider-commands-examples.md

Updates the configuration of an existing prompt provider, for instance, changing the 'summarize' setting for 'project-guidelines'. It confirms that the provider has been updated.

```bash
/prompt-providers update project-guidelines --summarize false
```

--------------------------------

### Add Dynamic Prompt Provider

Source: https://github.com/campfirein/cipher/blob/main/memAgent/prompt-provider-commands-examples.md

Adds or updates a dynamic prompt provider, such as 'summary', with specified parameters like conversation history. It confirms the addition and provides a generated summary for the provider.

```bash
/prompt-providers add-dynamic summary --history 10
```

--------------------------------

### Cipher MCP Server Environment Variables

Source: https://github.com/campfirein/cipher/blob/main/docs/cli-reference.md

Shows how to configure the behavior of the Cipher MCP server using environment variables, including setting the server mode, conflict resolution for aggregators, and timeouts.

```bash
# MCP server behavior
export MCP_SERVER_MODE=aggregator  # or 'default'
export AGGREGATOR_CONFLICT_RESOLUTION=prefix  # 'first-wins', 'error'
export AGGREGATOR_TIMEOUT=60000
```

--------------------------------

### Configure Multiple Cipher Instances

Source: https://github.com/campfirein/cipher/blob/main/docs/mcp-integration.md

Sets up multiple Cipher MCP server instances for different projects, each with its own configuration file and memory collection.

```json
{
	"mcpServers": {
		"cipher-frontend": {
			"type": "stdio",
			"command": "cipher",
			"args": ["--mode", "mcp", "--agent", "/path/to/frontend-config.yml"],
			"env": {
				"OPENAI_API_KEY": "sk-your-key",
				"VECTOR_STORE_COLLECTION": "frontend_memory"
			}
		},
		"cipher-backend": {
			"type": "stdio",
			"command": "cipher",
			"args": ["--mode", "mcp", "--agent", "/path/to/backend-config.yml"],
			"env": {
				"OPENAI_API_KEY": "sk-your-key",
				"VECTOR_STORE_COLLECTION": "backend_memory"
			}
		}
	}
}
```

--------------------------------

### MCP Server Configuration for Cipher

Source: https://github.com/campfirein/cipher/blob/main/README.md

This JSON configuration snippet defines how to set up Cipher as an MCP server within an MCP client. It specifies the server type, command, arguments, and environment variables required for operation.

```json
{
	"mcpServers": {
		"cipher": {
			"type": "stdio",
			"command": "cipher",
			"args": ["--mode", "mcp"],
			"env": {
				"MCP_SERVER_MODE": "aggregator",
				"OPENAI_API_KEY": "your_openai_api_key",
				"ANTHROPIC_API_KEY": "your_anthropic_api_key"
			}
		}
	}
}
```