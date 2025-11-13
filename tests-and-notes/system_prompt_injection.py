#!/usr/bin/env python3
"""
System Prompt Injection System for Cipher-Aggregator

This module provides system to inject routing enforcement rules into agent
system prompts, ensuring Cline/KiloCodes receive proper routing guidance
from cipher.yml configuration.
"""

import yaml
import json
import re
import logging
from typing import Dict, List, Optional, Any
from pathlib import Path
from dataclasses import dataclass


@dataclass
class SystemPromptConfig:
    """Configuration for system prompt injection"""
    base_system_prompt: str
    routing_rules_section: str
    performance_constraints_section: str
    task_categorization_section: str
    conflict_resolution_section: str


class SystemPromptInjector:
    """System prompt injection engine for routing enforcement"""

    def __init__(self, cipher_yml_path: str = "/home/jrede/dev/MCP/cipher.yml"):
        self.cipher_yml_path = Path(cipher_yml_path)
        self.logger = self._setup_logging()

        # Load cipher configuration
        self.cipher_config = self._load_cipher_config()

        # Initialize routing rules
        self.routing_rules = self._extract_routing_rules()

        # Create injection templates
        self._setup_injection_templates()

    def _setup_logging(self) -> logging.Logger:
        """Set up logging for system prompt injection"""
        logger = logging.getLogger('system_prompt_injection')
        logger.setLevel(logging.INFO)

        if not logger.handlers:
            handler = logging.StreamHandler()
            formatter = logging.Formatter(
                '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
            )
            handler.setFormatter(formatter)
            logger.addHandler(handler)

        return logger

    def _load_cipher_config(self) -> Dict[str, Any]:
        """Load cipher.yml configuration"""
        try:
            with open(self.cipher_yml_path, 'r') as f:
                config = yaml.safe_load(f)
            return config
        except FileNotFoundError:
            self.logger.warning(f"cipher.yml not found at {self.cipher_yml_path}")
            return self._get_default_config()
        except yaml.YAMLError as e:
            self.logger.error(f"Failed to parse cipher.yml: {e}")
            return self._get_default_config()

    def _get_default_config(self) -> Dict[str, Any]:
        """Get default configuration if cipher.yml is not available"""
        return {
            'systemPrompt': self._get_default_system_prompt(),
            'routing': {
                'rules': {
                    'domain_specific_first': True,
                    'max_calls_per_task': 8,
                    'serial_execution_only': True
                }
            }
        }

    def _get_default_system_prompt(self) -> str:
        """Get default system prompt with routing rules"""
        return """
You are an AI agent operating through the Cipher Aggregator system.
Always follow these routing rules for tool selection:

## DOMAIN-SPECIFIC FIRST (CRITICAL)
Always prioritize domain-specific tools over generic ones:
- GitHub operations → `github` MCP (NEVER use fetch/curl)
- Web scraping → `firecrawl` MCP (NEVER use fetch)
- Code analysis → `code-index` MCP (NEVER manually scan files)
- API testing → `schemathesis` MCP
- Test execution → `pytest` MCP
- File operations → `filesystem` MCP (not manual bash)
- Web search → `brave-search` MCP (not fetch/curl)
- HTTP requests → `httpie` MCP (for testing APIs)

## PERFORMANCE CONSTRAINTS
- Serial execution only (maxParallelCalls: 1)
- Minimize tool calls per task (max 8)
- Use batch operations when available

## TASK CATEGORIZATION
**Development Tasks**: code-index → filesystem → github → memory-bank
**Web Research**: brave-search → firecrawl → memory-bank
**API Testing**: schemathesis → httpie → pytest
**File Management**: filesystem → file-batch → memory-bank
**Documentation**: context7 → filesystem → memory-bank

When making tool selection decisions, always consider the domain and use the most appropriate specialized tool.
"""

    def _setup_injection_templates(self):
        """Set up string templates for system prompt injection"""

        # Main routing rules template
        self.routing_template = """
## ROUTING ENFORCEMENT RULES

You MUST follow these routing rules for all tool selection decisions:

### DOMAIN-SPECIFIC FIRST (CRITICAL)
Always prioritize domain-specific tools over generic ones:

{% for domain, rules in routing_rules.items() %}
**{{ domain|replace('_', ' ')|title }} Operations**:
- Preferred: {{ rules.primary_tools|join(', ') }}
- Forbidden: {{ rules.forbidden_tools|join(', ') }}
- Fallback: {{ rules.fallback_tools|join(', ') }}
  Keywords: {{ rules.keywords|join(', ') }}
{% endfor %}

### TASK CATEGORIZATION & ROUTING
{% for category, workflow in task_workflows.items() %}
**{{ category }} Tasks**: {{ workflow|join(' → ') }}
{% endfor %}

### PERFORMANCE CONSTRAINTS
- Serial execution only (maxParallelCalls: 1)
- Minimize tool calls per task (max {{ max_calls_per_task }})
- Use batch operations when available
- Cache results in memory-bank for reuse

### CONFLICT RESOLUTION
When multiple tools could work, use this priority:
1. **Specialized tools** (domain-specific MCP servers)
2. **Reliable tools** (well-tested, stable interfaces)
3. **Simple tools** (basic functionality, minimal complexity)

### MEMORY INTEGRATION
- Search memory-bank FIRST for known solutions before routing
- Store important routing decisions for future reference
- Use routing patterns location: `/home/jrede/dev/MCP/data/memory-bank/routing-patterns/`
"""

        # Performance constraints template
        self.performance_template = """
## PERFORMANCE CONSTRAINTS

**Execution Model**:
- Serial execution only (maxParallelCalls: 1)
- Maximum {{ max_calls_per_task }} tool calls per task
- {{ call_timeout }}ms timeout per tool call

**Optimization Guidelines**:
- Use batch operations when available
- Minimize redundant tool calls
- Cache frequently accessed data
- Monitor tool performance and adjust routing accordingly

**Domain-Specific Timeouts**:
{% for domain, timeout in domain_timeouts.items() %}
- {{ domain }}: {{ timeout }}ms
{% endfor %}
"""

        # Task categorization template
        self.categorization_template = """
## TASK CATEGORIZATION

Classify tasks using these patterns:

{% for domain, rules in routing_rules.items() %}
**{{ domain|replace('_', ' ')|title }}**:
  Keywords: {{ rules.keywords|join(', ') }}
  Tools: {{ rules.primary_tools|join(', ') }} (primary), {{ rules.fallback_tools|join(', ') }} (fallback)
{% endfor %}

### WORKFLOW PATTERNS

**Development Workflow**:
1. Code Analysis → File Operations → Version Control → Knowledge Storage
2. Tools: code-index → filesystem → github → memory-bank

**Web Research Workflow**:
1. Web Search → Content Extraction → Analysis → Storage
2. Tools: brave-search → firecrawl → memory-bank

**API Testing Workflow**:
1. Schema Validation → Endpoint Testing → Documentation → Integration
2. Tools: schemathesis → httpie → pytest → context7

**File Management Workflow**:
1. File Operations → Batch Processing → Documentation → Storage
2. Tools: filesystem → file-batch → memory-bank

**Documentation Workflow**:
1. Documentation Retrieval → File Operations → Knowledge Storage
2. Tools: context7 → filesystem → memory-bank
"""

    def generate_system_prompt(self, agent_type: str = "default") -> str:
        """Generate complete system prompt with routing enforcement"""

        # Base system prompt from cipher.yml
        base_prompt = self.cipher_config.get('systemPrompt', self._get_default_system_prompt())

        # Generate routing sections
        routing_section = self._generate_routing_section()
        performance_section = self._generate_performance_section()
        categorization_section = self._generate_categorization_section()

        # Combine all sections
        agent_specific_instructions = self._get_agent_specific_instructions(agent_type)
        agent_type_title = agent_type.upper() if agent_type in ['cline', 'kilocode'] else agent_type.title()
        complete_prompt = f"""
{base_prompt}

{routing_section}

{performance_section}

{categorization_section}

## AGENT-SPECIFIC INSTRUCTIONS

{agent_type}
## AGENT-SPECIFIC INSTRUCTIONS

{agent_specific_instructions}

## ENFORCEMENT COMPLIANCE
"""

        return complete_prompt.strip()

    def _get_task_workflows(self) -> Dict[str, List[str]]:
        """Get task workflow mappings for routing templates"""
        return {
            'development': ['code-index', 'filesystem', 'github', 'memory-bank'],
            'web research': ['brave-search', 'firecrawl', 'memory-bank'],
            'api testing': ['schemathesis', 'httpie', 'pytest'],
            'file management': ['filesystem', 'file-batch', 'memory-bank'],
            'documentation': ['context7', 'filesystem', 'memory-bank']
        }

    def _generate_routing_section(self) -> str:
        """Generate routing rules section"""
        try:
            # Replace template rendering with Python string formatting
            routing_rules_content = ""
            for category, rules in self.routing_rules.items():
                routing_rules_content += f"### {category.title()} Routing\n"
                for rule in rules:
                    routing_rules_content += f"- {rule}\n"
                routing_rules_content += "\n"

            workflows_content = ""
            for task_type, workflow in self._get_task_workflows().items():
                workflows_content += f"**{task_type.title()} Tasks**: {' → '.join(workflow)}\n"

            return f"""## ROUTING RULES

### Domain-Specific Priority
{routing_rules_content}### Task Workflows
{workflows_content}

### Performance Constraints
- Maximum {8} calls per task
- Serial execution only (maxParallelCalls: 1)
- Use batch operations when available
- Cache results in memory-bank for reuse

### Conflict Resolution Hierarchy
1. Specialized tools first (domain-specific MCP servers)
2. Reliable tools second (well-tested generic tools)
3. Simple tools last (basic functionality)

### Memory Integration Requirements
- Search memory-bank FIRST for known solutions
- Store important decisions via memory-bank
- Use memory_bank_search("routing patterns") before decisions
- Auto-capture session summaries in memory-bank
"""

        except Exception as e:
            self.logger.error(f"Failed to generate routing section: {e}")
            return self._get_fallback_routing_section()

    def _generate_performance_section(self) -> str:
        """Generate performance constraints section"""
        try:
            # Replace template rendering with Python string formatting
            domain_timeouts = {
                'github': 30000,
                'brave-search': 25000,
                'firecrawl': 60000,
                'schemathesis': 90000,
                'code-index': 120000,
                'filesystem': 15000,
                'memory-bank': 20000
            }

            timeouts_content = ""
            for domain, timeout in domain_timeouts.items():
                timeouts_content += f"- {domain}: {timeout}ms\n"

            return f"""## PERFORMANCE CONSTRAINTS

### Call Limits
- Maximum calls per task: 8
- Serial execution only (maxParallelCalls: 1)
- Call timeout: 45000ms

### Domain-Specific Timeouts
{timeouts_content}

### Optimization Requirements
- Minimize tool calls per task
- Use batch operations when available
- Cache results in memory-bank for reuse
- Monitor response times and adjust accordingly

### Error Handling
- Implement exponential backoff for retries
- Handle SSE connection drops gracefully
- Check for rate limiting responses
- Monitor server resources and connections
"""

        except Exception as e:
            self.logger.error(f"Failed to generate performance section: {e}")
            return self._get_fallback_performance_section()

    def _generate_categorization_section(self) -> str:
        """Generate task categorization section"""
        try:
            # Replace template rendering with Python string formatting
            return f"""## TASK CATEGORIZATION

### Development Tasks
**Workflow**: code-index → filesystem → github → memory-bank
**Priority**: Code analysis, file operations, version control, knowledge storage

### Web Research Tasks
**Workflow**: brave-search → firecrawl → memory-bank
**Priority**: Information discovery, content extraction, research patterns

### API Testing Tasks
**Workflow**: schemathesis → httpie → pytest
**Priority**: Schema validation, endpoint testing, automated testing

### File Management Tasks
**Workflow**: filesystem → file-batch → memory-bank
**Priority**: File operations, batch processing, storage management

### Documentation Tasks
**Workflow**: context7 → filesystem → memory-bank
**Priority**: Documentation retrieval, file operations, knowledge storage

### Domain-Specific Routing Rules
- GitHub operations → github MCP (NEVER use fetch/curl)
- Web scraping → firecrawl MCP (NEVER use fetch)
- Code analysis → code-index MCP (NEVER manually scan files)
- API testing → schemathesis MCP
- Test execution → pytest MCP
- File operations → filesystem MCP (not manual bash)
- Web search → brave-search MCP (not fetch/curl)
- HTTP requests → httpie MCP (for testing APIs)
"""

        except Exception as e:
            self.logger.error(f"Failed to generate categorization section: {e}")
            return self._get_fallback_categorization_section()

    def _get_agent_specific_instructions(self, agent_type: str) -> str:
        """Get agent-specific instructions based on agent type"""
        if agent_type == "cline":
            return """As a Cline agent, you have access to specialized MCP servers through cipher-aggregator.
Always use domain-specific MCP tools rather than generic alternatives.
Your routing decisions are monitored for compliance with these rules."""
        elif agent_type == "kilocode":
            return """As a Kilo Code agent, you are a highly skilled software engineer with routing expertise.
You must enforce these routing rules and guide other agents toward optimal tool selection.
Monitor routing performance and suggest improvements when patterns emerge."""
        elif agent_type == "researcher":
            return """As a researcher agent, focus on web research and data extraction workflows.
Use brave-search for information discovery and firecrawl for content extraction.
Document research patterns for future reference in memory-bank."""
        elif agent_type == "developer":
            return """As a developer agent, prioritize code analysis and development workflows.
Use code-index for code discovery and github for version control operations.
Maintain development knowledge in memory-bank for team collaboration."""
        else:
            return """As a general agent, follow all routing rules and adapt your workflow to the task domain.
Use specialized tools when available and escalate to memory-bank for complex decisions."""

    def _extract_routing_rules(self) -> Dict[str, Any]:
        """Extract routing rules from cipher configuration"""
        routing_config = self.cipher_config.get('routing', {})

        # Extract domain-specific rules
        domain_rules = routing_config.get('domain_rules', {})

        # Default domain mappings if not specified
        default_domains = {
            'github': {
                'primary_tools': ['github'],
                'forbidden_tools': ['fetch', 'curl'],
                'fallback_tools': ['filesystem'],
                'keywords': ['repository', 'github', 'pull request', 'issue', 'commit']
            },
            'web_search': {
                'primary_tools': ['brave-search'],
                'forbidden_tools': ['fetch', 'firecrawl_search'],
                'fallback_tools': ['firecrawl'],
                'keywords': ['search for', 'find information', 'google', 'web search']
            },
            'web_scraping': {
                'primary_tools': ['firecrawl'],
                'forbidden_tools': ['fetch', 'curl'],
                'fallback_tools': ['brave-search'],
                'keywords': ['scrape', 'extract data', 'crawl website', 'web scraping']
            },
            'code_analysis': {
                'primary_tools': ['code-index'],
                'forbidden_tools': ['filesystem', 'manual scanning'],
                'fallback_tools': ['filesystem'],
                'keywords': ['find code', 'search code', 'analyze code', 'code search']
            },
            'api_testing': {
                'primary_tools': ['schemathesis'],
                'forbidden_tools': ['fetch', 'httpie'],
                'fallback_tools': ['httpie', 'pytest'],
                'keywords': ['api test', 'schema test', 'openapi', 'api validation']
            }
        }

        # Merge with defaults
        for domain, config in default_domains.items():
            if domain not in domain_rules:
                domain_rules[domain] = config

        return domain_rules

    def _get_fallback_routing_section(self) -> str:
        """Fallback routing section if template generation fails"""
        return """
## ROUTING ENFORCEMENT RULES

**Domain-Specific Priority**:
- GitHub operations → github MCP (NEVER fetch/curl)
- Web scraping → firecrawl MCP (NEVER fetch)
- Code analysis → code-index MCP (NEVER manual scanning)
- Web search → brave-search MCP (NEVER fetch/curl)

**Performance**: Max 8 calls per task, serial execution only
**Resolution**: Specialized → Reliable → Simple
"""

    def _get_fallback_performance_section(self) -> str:
        """Fallback performance section"""
        return """
## PERFORMANCE CONSTRAINTS

**Execution**: Serial only (maxParallelCalls: 1)
**Limits**: Max 8 calls per task, 45s timeout
**Optimization**: Use batch operations, cache results
"""

    def _get_fallback_categorization_section(self) -> str:
        """Fallback categorization section"""
        return """
## TASK CATEGORIZATION

**Development**: code-index → filesystem → github → memory-bank
**Research**: brave-search → firecrawl → memory-bank
**Testing**: schemathesis → httpie → pytest
"""

    def inject_routing_rules(self, base_prompt: str, agent_config: Dict[str, Any]) -> str:
        """Inject routing rules into existing system prompt"""

        # Extract agent-specific routing rules
        agent_rules = agent_config.get('routing_rules', {})
        agent_workflows = agent_config.get('workflows', {})

        # Generate agent-specific routing section
        agent_routing = self._generate_agent_specific_rules(agent_rules, agent_workflows)

        # Check if prompt already has routing rules
        if self._has_routing_rules(base_prompt):
            # Update existing routing rules
            updated_prompt = self._update_existing_rules(base_prompt, agent_routing)
        else:
            # Add routing rules to end
            updated_prompt = f"{base_prompt}\n\n{agent_routing}"

        return updated_prompt

    def _has_routing_rules(self, prompt: str) -> bool:
        """Check if prompt already contains routing rules"""
        routing_indicators = [
            'DOMAIN-SPECIFIC',
            'routing rules',
            'tool selection',
            'performance constraints'
        ]
        return any(indicator.lower() in prompt.lower() for indicator in routing_indicators)

    def _update_existing_rules(self, prompt: str, new_rules: str) -> str:
        """Update existing routing rules in prompt"""
        # Simple approach: replace the routing section if found
        routing_pattern = r'(## ROUTING ENFORCEMENT RULES.*?)(?=\n\n##|\Z)'
        if re.search(routing_pattern, prompt, re.DOTALL):
            return re.sub(routing_pattern, new_rules, prompt, flags=re.DOTALL)
        else:
            return f"{prompt}\n\n{new_rules}"

    def _generate_agent_specific_rules(self, agent_rules: Dict, agent_workflows: Dict) -> str:
        """Generate agent-specific routing rules"""
        sections = []

        if agent_rules:
            sections.append("## AGENT-SPECIFIC ROUTING RULES")
            for domain, tools in agent_rules.items():
                sections.append(f"- {domain}: {', '.join(tools)}")

        if agent_workflows:
            sections.append("\n## AGENT WORKFLOWS")
            for workflow_name, steps in agent_workflows.items():
                sections.append(f"- {workflow_name}: {' → '.join(steps)}")

        return '\n'.join(sections)

    def validate_system_prompt(self, prompt: str) -> Dict[str, Any]:
        """Validate system prompt for routing rule completeness"""
        issues = []
        recommendations = []

        # Check for required sections
        required_sections = [
            'DOMAIN-SPECIFIC',
            'performance constraints',
            'task categorization',
            'conflict resolution'
        ]

        for section in required_sections:
            if section.lower() not in prompt.lower():
                issues.append(f"Missing {section} section")
                recommendations.append(f"Add {section} to system prompt")

        # Check for forbidden tools mentions
        forbidden_patterns = ['fetch', 'curl']
        for pattern in forbidden_patterns:
            if pattern in prompt.lower():
                issues.append(f"Found potentially problematic tool: {pattern}")
                recommendations.append(f"Specify when {pattern} should NOT be used")

        # Check performance constraints
        if 'maxParallelCalls: 1' not in prompt and 'serial execution' not in prompt.lower():
            issues.append("Missing serial execution requirement")
            recommendations.append("Add serial execution constraint")

        return {
            'is_valid': len(issues) == 0,
            'issues': issues,
            'recommendations': recommendations,
            'completeness_score': max(0, 100 - len(issues) * 20)
        }

    def export_agent_config(self, agent_type: str, output_path: str):
        """Export complete agent configuration with routing rules"""
        agent_config = {
            'agent_type': agent_type,
            'system_prompt': self.generate_system_prompt(agent_type),
            'routing_rules': self.routing_rules,
            'performance_constraints': {
                'max_calls_per_task': 8,
                'serial_execution': True,
                'timeout': 45000
            },
            'validation': self.validate_system_prompt(
                self.generate_system_prompt(agent_type)
            )
        }

        config_path = Path(output_path)
        config_path.parent.mkdir(exist_ok=True)

        with open(config_path, 'w') as f:
            json.dump(agent_config, f, indent=2)

        self.logger.info(f"Agent configuration exported to {config_path}")
        return config_path


def main():
    """Demo usage of system prompt injection"""
    injector = SystemPromptInjector()

    # Generate system prompts for different agent types
    agent_types = ['default', 'cline', 'kilocode', 'researcher', 'developer']

    for agent_type in agent_types:
        prompt = injector.generate_system_prompt(agent_type)
        print(f"\n=== {agent_type.upper()} AGENT SYSTEM PROMPT ===")
        print(prompt)
        print("\n" + "="*60)

        # Validate prompt
        validation = injector.validate_system_prompt(prompt)
        print(f"Validation: {validation}")

        # Export configuration
        config_path = f"/home/jrede/dev/MCP/data/{agent_type}_agent_config.json"
        injector.export_agent_config(agent_type, config_path)


if __name__ == "__main__":
    main()