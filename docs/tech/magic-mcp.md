# Magic MCP Server

**GitHub**: https://github.com/21st-dev/magic-mcp
**Package**: @21st-dev/magic-mcp
**Language**: TypeScript/JavaScript
**Transport**: stdio
**Status**: ✅ Active Community Server

---

## Overview

A Model Context Protocol server that provides frontend development tooling and UI component generation capabilities. Magic MCP enables LLMs to create, modify, and optimize frontend components with intelligent design suggestions and code generation.

**Key Features:**
- Frontend component generation
- UI/UX design assistance
- CSS and styling optimization
- Component library integration
- Responsive design automation
- Accessibility improvements

---

## Installation

### Using npm (via npx)
```bash
npx @21st-dev/magic-mcp
```

### Direct installation
```bash
npm install -g @21st-dev/magic-mcp
```

---

## Configuration

The Magic MCP server requires no environment variables for basic functionality. It runs with default settings.

**Optional Configuration:**
- `MAGIC_API_KEY`: For enhanced features and premium capabilities
- `MAGIC_BASE_URL`: Custom API endpoint (if using self-hosted version)

---

## Available Tools

### Component Generation
1. `magic_create_component`
   - Generate frontend components from descriptions
   - Inputs: `description` (required), `framework` (optional), `style` (optional)
   - Returns: Generated component code

2. `magic_modify_component`
   - Modify existing components with intelligent suggestions
   - Inputs: `component_code` (required), `modification_request` (required)
   - Returns: Modified component code

### Styling and CSS
3. `magic_optimize_css`
   - Optimize CSS for performance and maintainability
   - Inputs: `css_code` (required), `optimization_type` (optional)
   - Returns: Optimized CSS code

4. `magic_generate_styles`
   - Generate CSS/styling for components
   - Inputs: `component_description` (required), `theme` (optional)
   - Returns: Generated styles

### Accessibility and UX
5. `magic_improve_accessibility`
   - Enhance component accessibility
   - Inputs: `component_code` (required), `accessibility_focus` (optional)
   - Returns: Accessibility-improved code

6. `magic_responsive_design`
   - Make components responsive
   - Inputs: `component_code` (required), `breakpoints` (optional)
   - Returns: Responsive component code

---

## Usage Examples

### Component Creation
```json
{
  "tool": "magic_create_component",
  "arguments": {
    "description": "A modern navigation bar with dropdown menus and mobile hamburger menu",
    "framework": "react",
    "style": "tailwind"
  }
}
```

### Component Modification
```json
{
  "tool": "magic_modify_component",
  "arguments": {
    "component_code": "<div className='card'>...</div>",
    "modification_request": "Add hover effects and smooth transitions"
  }
}
```

### CSS Optimization
```json
{
  "tool": "magic_optimize_css",
  "arguments": {
    "css_code": ".button { padding: 10px 20px; margin: 5px; }",
    "optimization_type": "performance"
  }
}
```

### Accessibility Enhancement
```json
{
  "tool": "magic_improve_accessibility",
  "arguments": {
    "component_code": "<button>Click me</button>",
    "accessibility_focus": "screen_reader_support"
  }
}
```

---

## Supported Frameworks

### Frontend Frameworks
- **React**: Full support with hooks and modern patterns
- **Vue**: Vue 3 composition API support
- **Angular**: Component and service generation
- **Svelte**: Svelte and SvelteKit support

### Styling Systems
- **Tailwind CSS**: Utility-first CSS framework
- **Styled Components**: CSS-in-JS solutions
- **CSS Modules**: Scoped CSS support
- **Sass/SCSS**: Advanced CSS preprocessing

### Component Libraries
- **Material-UI**: Google's Material Design
- **Ant Design**: Enterprise-class UI design
- **Chakra UI**: Simple, modular and accessible components
- **Bootstrap**: Popular CSS framework

---

## Design Patterns

### Component Architecture
- **Atomic Design**: Atoms, molecules, organisms structure
- **Compound Components**: Flexible component composition
- **Render Props**: Component behavior sharing
- **Higher-Order Components**: Component enhancement patterns

### State Management
- **Local State**: useState, useReducer patterns
- **Context API**: React Context for state sharing
- **Prop Drilling**: Alternative to context for simple cases

---

## Known Limitations

1. **Frontend Focus**: Limited to frontend development scenarios
2. **Design Subjectivity**: AI-generated designs may need human refinement
3. **Complex Interactions**: May struggle with highly complex UI logic
4. **Browser Compatibility**: Generated code may need browser-specific adjustments
5. **Performance**: Generated components may need optimization for production

---

## Testing

```bash
# Test server startup
npx @21st-dev/magic-mcp

# Test tool availability
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | npx @21st-dev/magic-mcp

# Test component generation (requires MCP client)
```

---

## Use Cases

1. **Rapid Prototyping**: Quickly generate UI components
2. **Design System**: Create consistent component libraries
3. **Accessibility**: Ensure components meet accessibility standards
4. **Responsive Design**: Create mobile-friendly interfaces
5. **Code Refactoring**: Modernize and optimize existing components

---

## Best Practices

### Component Generation
- Provide clear, specific descriptions
- Specify framework and styling preferences
- Include accessibility requirements upfront
- Test generated components across different screen sizes

### Code Quality
- Review generated code for security best practices
- Optimize for performance when needed
- Ensure proper error handling
- Add appropriate comments and documentation

---

## Related Documentation

- [Magic MCP GitHub Repository](https://github.com/21st-dev/magic-mcp)
- [21st.dev Magic Platform](https://21st.dev/)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [Frontend Development Best Practices](https://web.dev/)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 1A)
**Next Steps**: Register with jarvis and test functionality
**Category**: Frontend Development Tool
