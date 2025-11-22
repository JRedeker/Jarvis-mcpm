# Kilo Code Documentation (Context7 Fetch)

- **Source**: Context7 `get-library-docs` (`/kilo-org/docs`, tokens=3000)
- **Fetched**: 2025-11-16

Kilo Code is an open-source VS Code extension that accelerates development through AI-driven code generation and task automation. The extension integrates directly into VS Code to provide intelligent assistance for writing, refactoring, debugging, and explaining code through natural language interaction. This documentation site is built with Docusaurus 3.8.1 and provides comprehensive guides for installation, usage, and customization.

The project serves as a complete technical reference for Kilo Code, covering everything from basic setup to advanced customization. It includes documentation for multiple operational modes, custom tool creation via Model Context Protocol (MCP), and extensive API provider integration. The documentation architecture supports internationalization (English and Chinese) with local search capabilities, PostHog analytics integration, and a community-driven approach to feature development.

## Installation and Setup

Install Kilo Code extension into VS Code:

```bash
# Install from VS Code Marketplace
code --install-extension kilocode.kilo-code

# Alternative: Install from .vsix file
code --install-extension kilocode-kilo-code-*.vsix

# For VSCodium or Windsurf - use Open VSX Registry
# Navigate to https://open-vsx.org/extension/kilocode/Kilo-Code

# Start documentation development server
npm start

# Build documentation for production
npm run build
```

## Tool System - read_file

Read file contents with line numbers for analysis and reference.

```javascript
// Reading entire file
{
  tool: "read_file",
  parameters: {
    path: "src/app.js"
  }
}

// Expected output:
// 1 | const express = require('express');
// 2 | const app = express();
// 3 | const port = 3000;

// Reading specific line range
{
  tool: "read_file",
  parameters: {
    path: "src/app.js",
    start_line: 46,
    end_line: 68
  }
}

// Expected output:
// 46 | function processData(input) {
// 47 |   return input.map(x => x * 2);
// 48 | }
// ...
// 68 | module.exports = { processData };

// Reading large file with auto-truncation
{
  tool: "read_file",
  parameters: {
    path: "src/large-module.js",
    auto_truncate: true
  }
}

// Output truncates at configured limit (typically 500-1000 lines):
// 1 | import React from 'react';
// ...
// 1000 | export default Component;
// [... truncated 500 lines ...]
```

## Tool System - apply_diff

Apply surgical code changes with fuzzy matching and line number targeting.

```javascript
// Making targeted changes to existing code
{
  tool: "apply_diff",
  parameters: {
    path: "src/calculator.js",
    diff: `<<<<<<< SEARCH
:start_line:10
:end_line:12
-------
    // Old calculation logic
    const result = value * 0.9;
    return result;
=======
    // Updated calculation logic with logging
    console.log(\`Calculating for value: \${value}\`);
    const result = value * 0.95; // Adjusted factor
    return result;
>>>>>>> REPLACE

<<<<<<< SEARCH
:start_line:25
:end_line:25
-------
    const defaultTimeout = 5000;
=======
    const defaultTimeout = 10000; // Increased timeout
>>>>>>> REPLACE`
  }
}

// Result: Shows diff view for user approval, then applies changes
// Line 10-12: Replaced with new calculation logic
// Line 25: Updated timeout value
// All changes preserve indentation and formatting
```

## Tool System - write_to_file

Create new files or replace entire file contents with approval workflow.

```javascript
// Creating a new configuration file
{
  tool: "write_to_file",
  parameters: {
    path: "config/api.json",
    content: `{
  "endpoint": "https://api.example.com",
  "timeout": 5000,
  "retries": 3,
  "headers": {
    "Content-Type": "application/json"
  }
}`,
    line_count: 8
  }
}

// Creating a JavaScript utility module
{
  tool: "write_to_file",
  parameters: {
    path: "src/utils/helpers.js",
    content: `export function formatDate(date) {
  return new Date(date).toLocaleDateString();
}

export function calculateTotal(items) {
  return items.reduce((sum, item) => sum + item.price, 0);
}

export function debounce(func, delay) {
  let timeout;
  return function(...args) {
    clearTimeout(timeout);
    timeout = setTimeout(() => func.apply(this, args), delay);
  };
}`,
    line_count: 14
  }
}

// Process:
// 1. Opens diff view showing changes
// 2. User reviews and can edit proposed content
// 3. User approves or rejects
// 4. File written with approved content
```

## Tool System - execute_command

Execute terminal commands with output capture and security validation.

```bash
# Install project dependencies
{
  tool: "execute_command",
  parameters: {
    command: "npm install express mongodb dotenv"
  }
}

# Build and start development server
{
  tool: "execute_command",
  parameters: {
    command: "npm run build && npm start"
  }
}
```
