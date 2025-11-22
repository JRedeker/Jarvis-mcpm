# Task Master MCP Server

**GitHub**: https://github.com/eyaltoledano/claude-task-master
**Package**: task-master-mcp
**Language**: TypeScript/JavaScript
**Transport**: stdio
**Status**: ✅ Active Community Server

---

## Overview

A Model Context Protocol server that provides task management and project organization capabilities. Task Master enables LLMs to create, track, and manage development tasks with intelligent prioritization and workflow automation.

**Key Features:**
- Task creation and management
- Project organization and planning
- Intelligent task prioritization
- Workflow automation
- Integration with development workflows
- Progress tracking and reporting

---

## Installation

### Using npm (via npx)
```bash
npx task-master-mcp
```

### Direct installation
```bash
npm install -g task-master-mcp
```

### From source
```bash
git clone https://github.com/eyaltoledano/claude-task-master.git
cd claude-task-master
npm install
npm run build
```

---

## Configuration

The Task Master MCP server requires minimal configuration for basic functionality.

**Optional Configuration:**
- `TASK_MASTER_API_KEY`: For enhanced features and cloud sync
- `TASK_MASTER_DATA_DIR`: Custom data storage location
- `TASK_MASTER_WORKFLOW`: Default workflow configuration

---

## Available Tools

### Task Management
1. `task_create`
   - Create new tasks with detailed specifications
   - Inputs: `title` (required), `description` (optional), `priority` (optional), `tags` (optional)
   - Returns: Created task with unique ID

2. `task_update`
   - Update existing task properties
   - Inputs: `task_id` (required), `updates` (required)
   - Returns: Updated task information

3. `task_delete`
   - Delete tasks from the system
   - Inputs: `task_id` (required)
   - Returns: Deletion confirmation

4. `task_list`
   - List tasks with filtering options
   - Inputs: `status` (optional), `priority` (optional), `tags` (optional)
   - Returns: Array of matching tasks

### Project Organization
5. `project_create`
   - Create new projects for task organization
   - Inputs: `name` (required), `description` (optional), `goals` (optional)
   - Returns: Created project with ID

6. `project_add_task`
   - Add tasks to specific projects
   - Inputs: `project_id` (required), `task_id` (required)
   - Returns: Association confirmation

7. `project_status`
   - Get project status and progress
   - Inputs: `project_id` (required)
   - Returns: Project status report

### Workflow and Automation
8. `workflow_create`
   - Create custom workflows for task processing
   - Inputs: `name` (required), `steps` (required), `conditions` (optional)
   - Returns: Workflow configuration

9. `task_auto_prioritize`
   - Automatically prioritize tasks based on various factors
   - Inputs: `project_id` (optional), `algorithm` (optional)
   - Returns: Prioritized task list

### Reporting and Analytics
10. `task_report`
    - Generate task completion and progress reports
    - Inputs: `date_range` (optional), `project_id` (optional)
    - Returns: Comprehensive task report

11. `task_analytics`
    - Get task analytics and insights
    - Inputs: `metric_type` (optional), `time_period` (optional)
    - Returns: Analytics data and insights

---

## Usage Examples

### Task Creation
```json
{
  "tool": "task_create",
  "arguments": {
    "title": "Implement user authentication",
    "description": "Add JWT-based authentication to the API endpoints",
    "priority": "high",
    "tags": ["backend", "security", "api"]
  }
}
```

### Project Organization
```json
{
  "tool": "project_create",
  "arguments": {
    "name": "MCP Integration Phase 1B",
    "description": "Server registration and tool group configuration",
    "goals": ["Register all 11 servers", "Test tool discovery", "Configure tool groups"]
  }
}
```

### Workflow Automation
```json
{
  "tool": "workflow_create",
  "arguments": {
    "name": "Code Review Workflow",
    "steps": ["create", "review", "test", "deploy"],
    "conditions": {"requires_review": true}
  }
}
```

### Task Prioritization
```json
{
  "tool": "task_auto_prioritize",
  "arguments": {
    "project_id": "proj_123",
    "algorithm": "deadline_weighted"
  }
}
```

### Progress Reporting
```json
{
  "tool": "task_report",
  "arguments": {
    "date_range": "last_7_days",
    "project_id": "proj_123"
  }
}
```

---

## Task Properties

### Task Status
- **pending**: Task created but not started
- **in_progress**: Currently being worked on
- **completed**: Task finished successfully
- **blocked**: Task cannot proceed due to dependencies
- **cancelled**: Task was cancelled

### Priority Levels
- **critical**: Must be completed immediately
- **high**: Important task with near-term deadline
- **medium**: Standard priority task
- **low**: Nice-to-have or future task

### Task Types
- **feature**: New feature development
- **bug**: Bug fix or issue resolution
- **improvement**: Code or process improvement
- **documentation**: Documentation tasks
- **testing**: Testing and quality assurance

---

## Workflow System

### Built-in Workflows
- **Agile**: Standard agile development workflow
- **Kanban**: Continuous flow workflow
- **Waterfall**: Sequential phase-based workflow
- **Custom**: User-defined workflow patterns

### Workflow Steps
1. **Backlog**: Initial task creation and triage
2. **Ready**: Task prepared for development
3. **In Progress**: Active development phase
4. **Review**: Code review and testing
5. **Done**: Completed and deployed

---

## Integration Capabilities

### Development Tools
- **Git Integration**: Link tasks to commits and branches
- **CI/CD**: Trigger builds and deployments based on task status
- **Issue Tracking**: Sync with GitHub Issues, Jira, etc.

### Communication
- **Slack Integration**: Send task updates to Slack channels
- **Email Notifications**: Email task assignments and updates
- **Webhook Support**: Custom webhook notifications

---

## Known Limitations

1. **Single User Focus**: Designed for individual developer use
2. **No Real-time Collaboration**: Limited multi-user support
3. **Basic Analytics**: Simple reporting capabilities
4. **File-based Storage**: Local storage only, no cloud sync in basic version
5. **Limited Integrations**: Fewer third-party tool integrations

---

## Testing

```bash
# Test server startup
npx task-master-mcp

# Test tool availability
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | npx task-master-mcp

# Test task management (requires MCP client)
```

---

## Use Cases

1. **Personal Task Management**: Organize personal development tasks
2. **Project Planning**: Plan and track software projects
3. **Sprint Management**: Manage agile development sprints
4. **Bug Tracking**: Track and resolve software bugs
5. **Feature Development**: Manage feature development lifecycle

---

## Best Practices

### Task Creation
- Write clear, actionable task titles
- Include sufficient detail in descriptions
- Set appropriate priorities
- Use tags for better organization

### Project Management
- Break large projects into smaller tasks
- Set realistic deadlines
- Regular progress reviews
- Keep task status updated

### Workflow Optimization
- Customize workflows to match your process
- Use automation to reduce manual work
- Regular analytics reviews
- Continuous process improvement

---

## Related Documentation

- [Task Master GitHub Repository](https://github.com/eyaltoledano/claude-task-master)
- [Task Management Best Practices](https://www.atlassian.com/agile/project-management)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [Agile Development Guide](https://www.agilealliance.org/agile101/)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 1A)
**Next Steps**: Register with jarvis and test functionality
**Category**: Core Development Tool
