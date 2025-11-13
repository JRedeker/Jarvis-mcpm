#!/usr/bin/env python3
"""
System Prompt Routing Enforcement System for Cipher-Aggregator

This system ensures Cline/KiloCodes agents receive proper routing instructions
from cipher.yml and monitors/validates tool selection decisions according to
Phase 5 requirements.

Key Features:
- Dynamic system prompt injection with routing rules from cipher.yml
- Real-time routing rule validation and enforcement
- Tool selection tracking with SQLite persistence
- Performance constraint monitoring (max 8 calls, serial execution)
- Agent configuration verification
- Routing pattern analysis and optimization recommendations
"""

import yaml
import sqlite3
import json
import logging
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Any, Tuple, Union
from pathlib import Path
from dataclasses import dataclass, asdict
from enum import Enum
import threading
import time
import re
from collections import defaultdict, Counter

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class TaskCategory(Enum):
    """Task categories for routing classification."""
    DEVELOPMENT = "development"
    WEB_RESEARCH = "web_research"
    API_TESTING = "api_testing"
    FILE_MANAGEMENT = "file_management"
    DOCUMENTATION = "documentation"
    SEARCH_OPERATIONS = "search_operations"
    MONITORING = "monitoring"
    DATA_MANAGEMENT = "data_management"
    UI_DEVELOPMENT = "ui_development"
    INFRASTRUCTURE = "infrastructure"


class ToolSelectionStatus(Enum):
    """Status of tool selection decisions."""
    COMPLIANT = "compliant"
    NON_COMPLIANT = "non_compliant"
    CONFLICT = "conflict"
    TIMEOUT = "timeout"
    ERROR = "error"


class Domain(Enum):
    """Domain categories for domain-specific routing."""
    GITHUB = "github"
    WEB_SCRAPING = "web_scraping"
    CODE_ANALYSIS = "code_analysis"
    API_TESTING = "api_testing"
    TEST_EXECUTION = "test_execution"
    FILE_OPERATIONS = "file_operations"
    WEB_SEARCH = "web_search"
    HTTP_REQUESTS = "http_requests"
    KNOWLEDGE_STORAGE = "knowledge_storage"
    DATABASE_OPERATIONS = "database_operations"
    MONITORING = "monitoring"
    UI_DEVELOPMENT = "ui_development"
    INFRASTRUCTURE = "infrastructure"
    UNKNOWN = "unknown"


@dataclass
class RoutingRule:
    """Represents a single routing rule from cipher.yml."""
    domain: Domain
    preferred_tool: str
    alternative_tools: List[str]
    forbidden_tools: List[str]
    task_categories: List[TaskCategory]
    priority: int
    description: str


@dataclass
class ToolSelection:
    """Represents a single tool selection decision."""
    session_id: str
    agent_type: str  # 'cline', 'kilocode', etc.
    task_description: str
    task_category: TaskCategory
    detected_domain: Domain
    recommended_tool: str
    selected_tool: str
    selection_status: ToolSelectionStatus
    selection_timestamp: datetime
    execution_time_ms: Optional[int] = None
    success: bool = True
    error_message: Optional[str] = None
    performance_score: float = 0.0


@dataclass
class PerformanceMetrics:
    """Performance metrics for tool selections."""
    session_id: str
    total_calls: int
    calls_remaining: int
    execution_mode: str  # 'serial', 'parallel'
    start_time: datetime
    last_update: datetime
    average_execution_time: float
    success_rate: float


class CipherRoutingEngine:
    """Core routing enforcement engine that reads from cipher.yml and validates decisions."""

    def __init__(self, config_path: str = "/home/jrede/dev/MCP/cipher.yml"):
        """Initialize the routing engine with cipher.yml configuration."""
        self.config_path = config_path
        self.config = self._load_cipher_config()
        self.routing_rules = self._parse_routing_rules()
        self.session_metrics: Dict[str, PerformanceMetrics] = {}
        self.session_lock = threading.Lock()

        # Initialize SQLite database for tracking
        self.db_path = "/home/jrede/dev/MCP/data/routing_decisions.db"
        self._init_database()

        logger.info(f"Cipher Routing Engine initialized with {len(self.routing_rules)} routing rules")

    def _load_cipher_config(self) -> Dict[str, Any]:
        """Load and parse the cipher.yml configuration."""
        try:
            with open(self.config_path, 'r') as f:
                config = yaml.safe_load(f)
            logger.info(f"Loaded cipher configuration from {self.config_path}")
            return config
        except Exception as e:
            logger.error(f"Failed to load cipher configuration: {e}")
            raise

    def _parse_routing_rules(self) -> List[RoutingRule]:
        """Parse routing rules from cipher.yml systemPrompt."""
        rules = []
        system_prompt = self.config.get('systemPrompt', '')

        # Extract domain-specific routing rules from system prompt
        domain_patterns = {
            Domain.GITHUB: [r'github', r'repository', r'repo', r'git'],
            Domain.WEB_SCRAPING: [r'firecrawl', r'scraping', r'extract.*web', r'crawl.*web'],
            Domain.CODE_ANALYSIS: [r'code-index', r'code.*analysis', r'search.*code'],
            Domain.API_TESTING: [r'schemathesis', r'api.*test', r'openapi'],
            Domain.TEST_EXECUTION: [r'pytest', r'test.*execution'],
            Domain.FILE_OPERATIONS: [r'filesystem', r'file.*operation', r'local.*file'],
            Domain.WEB_SEARCH: [r'brave-search', r'web.*search', r'online.*search'],
            Domain.HTTP_REQUESTS: [r'httpie', r'http.*request'],
            Domain.KNOWLEDGE_STORAGE: [r'memory-bank', r'knowledge.*storage'],
            Domain.DATABASE_OPERATIONS: [r'sql', r'database.*query'],
            Domain.MONITORING: [r'prometheus', r'metrics', r'monitoring'],
            Domain.UI_DEVELOPMENT: [r'svelte', r'textual', r'ui.*development'],
            Domain.INFRASTRUCTURE: [r'docker', r'infrastructure']
        }

        # Parse task categorization patterns
        task_categories = {
            'development': TaskCategory.DEVELOPMENT,
            'web_research': TaskCategory.WEB_RESEARCH,
            'api_testing': TaskCategory.API_TESTING,
            'file_management': TaskCategory.FILE_MANAGEMENT,
            'documentation': TaskCategory.DOCUMENTATION,
            'search_operations': TaskCategory.SEARCH_OPERATIONS
        }

        # Create routing rules based on cipher.yml systemPrompt
        for domain, patterns in domain_patterns.items():
            if any(re.search(pattern, system_prompt, re.IGNORECASE) for pattern in patterns):
                # Extract preferred tool for this domain
                preferred_tool = self._extract_preferred_tool(domain, system_prompt)
                alternative_tools = self._extract_alternative_tools(domain, system_prompt)
                forbidden_tools = self._extract_forbidden_tools(domain, system_prompt)

                # Map task categories found in system prompt
                matched_categories = []
                for cat_name, category in task_categories.items():
                    if re.search(cat_name.replace('_', '.*'), system_prompt, re.IGNORECASE):
                        matched_categories.append(category)

                rule = RoutingRule(
                    domain=domain,
                    preferred_tool=preferred_tool,
                    alternative_tools=alternative_tools,
                    forbidden_tools=forbidden_tools,
                    task_categories=matched_categories if matched_categories else list(TaskCategory),
                    priority=self._get_domain_priority(domain),
                    description=f"Routing rule for {domain.value} operations"
                )
                rules.append(rule)

        # Add default rules for unknown domains
        default_rule = RoutingRule(
            domain=Domain.UNKNOWN,
            preferred_tool="filesystem",
            alternative_tools=["fetch", "httpie"],
            forbidden_tools=[],
            task_categories=list(TaskCategory),
            priority=999,
            description="Default routing rule for unknown domains"
        )
        rules.append(default_rule)

        logger.info(f"Parsed {len(rules)} routing rules from system prompt")
        return rules

    def _extract_preferred_tool(self, domain: Domain, system_prompt: str) -> str:
        """Extract the preferred tool for a domain from system prompt."""
        domain_tool_mapping = {
            Domain.GITHUB: ["github", "github-mcp"],
            Domain.WEB_SCRAPING: ["firecrawl"],
            Domain.CODE_ANALYSIS: ["code-index"],
            Domain.API_TESTING: ["schemathesis"],
            Domain.TEST_EXECUTION: ["pytest"],
            Domain.FILE_OPERATIONS: ["filesystem"],
            Domain.WEB_SEARCH: ["brave-search"],
            Domain.HTTP_REQUESTS: ["httpie"],
            Domain.KNOWLEDGE_STORAGE: ["memory-bank"],
            Domain.DATABASE_OPERATIONS: ["sql"],
            Domain.MONITORING: ["prometheus"],
            Domain.UI_DEVELOPMENT: ["svelte", "textual-devtools"],
            Domain.INFRASTRUCTURE: ["docker"]
        }

        tools = domain_tool_mapping.get(domain, ["filesystem"])
        for tool in tools:
            if re.search(tool, system_prompt, re.IGNORECASE):
                return tool
        return tools[0]

    def _extract_alternative_tools(self, domain: Domain, system_prompt: str) -> List[str]:
        """Extract alternative tools for a domain from system prompt."""
        # Define fallback hierarchies based on domain
        fallback_hierarchies = {
            Domain.GITHUB: ["fetch"],
            Domain.WEB_SCRAPING: ["fetch"],
            Domain.CODE_ANALYSIS: ["filesystem"],
            Domain.API_TESTING: ["httpie", "fetch"],
            Domain.TEST_EXECUTION: [],
            Domain.FILE_OPERATIONS: ["file-batch"],
            Domain.WEB_SEARCH: ["firecrawl"],
            Domain.HTTP_REQUESTS: ["fetch"],
            Domain.KNOWLEDGE_STORAGE: ["filesystem"],
            Domain.DATABASE_OPERATIONS: [],
            Domain.MONITORING: ["fetch"],
            Domain.UI_DEVELOPMENT: ["filesystem"],
            Domain.INFRASTRUCTURE: []
        }

        return fallback_hierarchies.get(domain, ["fetch"])

    def _extract_forbidden_tools(self, domain: Domain, system_prompt: str) -> List[str]:
        """Extract forbidden tools for a domain from system prompt."""
        forbidden_mappings = {
            Domain.GITHUB: ["fetch", "curl"],  # Never use fetch for GitHub
            Domain.WEB_SCRAPING: ["fetch"],    # Use firecrawl instead
            Domain.CODE_ANALYSIS: ["fetch", "manual file scanning"],
            Domain.WEB_SEARCH: ["fetch", "firecrawl_search"],  # Use brave-search
            Domain.HTTP_REQUESTS: [],  # httpie and fetch both acceptable
        }

        return forbidden_mappings.get(domain, [])

    def _get_domain_priority(self, domain: Domain) -> int:
        """Get priority for domain (lower number = higher priority)."""
        priority_mapping = {
            Domain.GITHUB: 1,
            Domain.WEB_SEARCH: 2,
            Domain.API_TESTING: 3,
            Domain.DATABASE_OPERATIONS: 4,
            Domain.CODE_ANALYSIS: 5,
            Domain.WEB_SCRAPING: 6,
            Domain.KNOWLEDGE_STORAGE: 7,
            Domain.TEST_EXECUTION: 8,
            Domain.FILE_OPERATIONS: 9,
            Domain.MONITORING: 10,
            Domain.HTTP_REQUESTS: 11,
            Domain.UI_DEVELOPMENT: 12,
            Domain.INFRASTRUCTURE: 13
        }
        return priority_mapping.get(domain, 999)

    def _init_database(self):
        """Initialize SQLite database for tracking routing decisions."""
        Path(self.db_path).parent.mkdir(parents=True, exist_ok=True)

        with sqlite3.connect(self.db_path) as conn:
            conn.execute("""
                CREATE TABLE IF NOT EXISTS routing_decisions (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    session_id TEXT NOT NULL,
                    agent_type TEXT NOT NULL,
                    task_description TEXT,
                    task_category TEXT NOT NULL,
                    detected_domain TEXT NOT NULL,
                    recommended_tool TEXT NOT NULL,
                    selected_tool TEXT NOT NULL,
                    selection_status TEXT NOT NULL,
                    selection_timestamp DATETIME NOT NULL,
                    execution_time_ms INTEGER,
                    success BOOLEAN DEFAULT 1,
                    error_message TEXT,
                    performance_score REAL DEFAULT 0.0,
                    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
                )
            """)

            conn.execute("""
                CREATE INDEX IF NOT EXISTS idx_session_id
                ON routing_decisions(session_id)
            """)

            conn.execute("""
                CREATE INDEX IF NOT EXISTS idx_domain
                ON routing_decisions(detected_domain)
            """)

            conn.execute("""
                CREATE INDEX IF NOT EXISTS idx_timestamp
                ON routing_decisions(selection_timestamp)
            """)

            conn.execute("""
                CREATE TABLE IF NOT EXISTS session_metrics (
                    session_id TEXT PRIMARY KEY,
                    total_calls INTEGER NOT NULL,
                    calls_remaining INTEGER NOT NULL,
                    execution_mode TEXT NOT NULL,
                    start_time DATETIME NOT NULL,
                    last_update DATETIME NOT NULL,
                    average_execution_time REAL DEFAULT 0.0,
                    success_rate REAL DEFAULT 0.0
                )
            """)

        logger.info(f"Initialized routing decisions database at {self.db_path}")

    def detect_domain(self, task_description: str, context: Optional[Dict] = None) -> Domain:
        """Detect the domain for a given task description."""
        task_lower = task_description.lower()
        context = context or {}

        # Domain detection patterns
        domain_patterns = {
            Domain.GITHUB: [
                r'github', r'repository', r'repo', r'git', r'branch', r'commit',
                r'pull.*request', r'pr', r'issue', r'fork', r'star'
            ],
            Domain.WEB_SCRAPING: [
                r'scrape', r'extract.*web', r'crawl.*web', r'parse.*html',
                r'scrape.*site', r'extract.*data.*from.*web'
            ],
            Domain.CODE_ANALYSIS: [
                r'code.*analysis', r'analyze.*code', r'search.*code',
                r'find.*function', r'find.*class', r'code.*search',
                r'indexed.*search', r'code.*index'
            ],
            Domain.API_TESTING: [
                r'api.*test', r'schema.*validation', r'openapi', r'swagger',
                r'endpoint.*test', r'http.*test', r'api.*spec'
            ],
            Domain.TEST_EXECUTION: [
                r'run.*test', r'test.*execution', r'pytest', r'unit.*test',
                r'integration.*test', r'test.*suite'
            ],
            Domain.FILE_OPERATIONS: [
                r'read.*file', r'write.*file', r'create.*file', r'edit.*file',
                r'list.*files', r'directory', r'local.*file', r'workspace.*file'
            ],
            Domain.WEB_SEARCH: [
                r'search.*web', r'google', r'bing', r'online.*search',
                r'find.*information.*online', r'web.*query'
            ],
            Domain.HTTP_REQUESTS: [
                r'http.*request', r'curl', r'wget', r'request.*url',
                r'fetch.*url', r'get.*http', r'post.*http'
            ],
            Domain.KNOWLEDGE_STORAGE: [
                r'store.*knowledge', r'save.*information', r'remember',
                r'knowledge.*base', r'save.*for.*later', r'document'
            ],
            Domain.DATABASE_OPERATIONS: [
                r'sql.*query', r'database', r'select.*from', r'insert.*into',
                r'update.*table', r'run.*query', r'db.*query'
            ],
            Domain.MONITORING: [
                r'metrics', r'monitoring', r'prometheus', r'grafana',
                r'performance.*monitor', r'system.*metrics'
            ],
            Domain.UI_DEVELOPMENT: [
                r'svelte', r'ui.*development', r'frontend', r'interface',
                r'gui', r'web.*component'
            ],
            Domain.INFRASTRUCTURE: [
                r'docker', r'container', r'kubernetes', r'deployment',
                r'infrastructure', r'devops'
            ]
        }

        # Score domains based on pattern matches
        domain_scores = {}
        for domain, patterns in domain_patterns.items():
            score = 0
            for pattern in patterns:
                matches = len(re.findall(pattern, task_lower, re.IGNORECASE))
                score += matches

            # Check context for additional clues
            if context:
                context_str = str(context).lower()
                for pattern in patterns:
                    score += len(re.findall(pattern, context_str, re.IGNORECASE))

            domain_scores[domain] = score

        # Return domain with highest score, or UNKNOWN if no matches
        if not domain_scores or max(domain_scores.values()) == 0:
            return Domain.UNKNOWN

        return max(domain_scores.items(), key=lambda x: x[1])[0]

    def classify_task_category(self, task_description: str) -> TaskCategory:
        """Classify task into a category for routing."""
        task_lower = task_description.lower()

        category_patterns = {
            TaskCategory.DEVELOPMENT: [
                r'develop', r'code', r'program', r'implement', r'build',
                r'create.*function', r'create.*class', r'bug.*fix'
            ],
            TaskCategory.WEB_RESEARCH: [
                r'research', r'investigate', r'explore.*web', r'find.*information',
                r'gather.*data', r'web.*analysis'
            ],
            TaskCategory.API_TESTING: [
                r'test.*api', r'validate.*schema', r'api.*testing',
                r'endpoint.*test', r'swagger.*test'
            ],
            TaskCategory.FILE_MANAGEMENT: [
                r'manage.*file', r'organize.*file', r'file.*operation',
                r'directory.*manage', r'file.*system'
            ],
            TaskCategory.DOCUMENTATION: [
                r'document', r'docs', r'readme', r'write.*documentation',
                r'api.*doc', r'generate.*doc'
            ],
            TaskCategory.SEARCH_OPERATIONS: [
                r'search', r'find', r'locate', r'query', r'look.*for'
            ]
        }

        category_scores = {}
        for category, patterns in category_patterns.items():
            score = sum(len(re.findall(pattern, task_lower, re.IGNORECASE))
                       for pattern in patterns)
            category_scores[category] = score

        if not category_scores or max(category_scores.values()) == 0:
            return TaskCategory.DEVELOPMENT  # Default category

        return max(category_scores.items(), key=lambda x: x[1])[0]

    def get_routing_recommendation(self, task_description: str, context: Optional[Dict] = None) -> Tuple[Domain, str]:
        """Get routing recommendation for a task."""
        detected_domain = self.detect_domain(task_description, context)

        # Find the most appropriate routing rule
        best_rule = None
        for rule in self.routing_rules:
            if rule.domain == detected_domain:
                if best_rule is None or rule.priority < best_rule.priority:
                    best_rule = rule

        if not best_rule:
            # Fall back to default rule
            best_rule = next((r for r in self.routing_rules if r.domain == Domain.UNKNOWN), None)

        recommended_tool = best_rule.preferred_tool if best_rule else "filesystem"

        return detected_domain, recommended_tool

    def validate_tool_selection(self, session_id: str, agent_type: str, task_description: str,
                              selected_tool: str, context: Optional[Dict] = None) -> ToolSelection:
        """Validate a tool selection decision and track it."""
        detected_domain, recommended_tool = self.get_routing_recommendation(task_description, context)
        task_category = self.classify_task_category(task_description)

        # Determine selection status
        selection_status = ToolSelectionStatus.COMPLIANT
        if selected_tool != recommended_tool:
            # Check if selected tool is in alternatives
            rule = next((r for r in self.routing_rules if r.domain == detected_domain), None)
            if rule and selected_tool in rule.alternative_tools:
                selection_status = ToolSelectionStatus.COMPLIANT
            else:
                selection_status = ToolSelectionStatus.NON_COMPLIANT

        # Check if tool is forbidden
        rule = next((r for r in self.routing_rules if r.domain == detected_domain), None)
        if rule and selected_tool in rule.forbidden_tools:
            selection_status = ToolSelectionStatus.NON_COMPLIANT

        # Create selection record
        selection = ToolSelection(
            session_id=session_id,
            agent_type=agent_type,
            task_description=task_description,
            task_category=task_category,
            detected_domain=detected_domain,
            recommended_tool=recommended_tool,
            selected_tool=selected_tool,
            selection_status=selection_status,
            selection_timestamp=datetime.now()
        )

        # Store in database
        self._store_selection(selection)

        return selection

    def _store_selection(self, selection: ToolSelection):
        """Store tool selection in SQLite database."""
        with sqlite3.connect(self.db_path) as conn:
            conn.execute("""
                INSERT INTO routing_decisions (
                    session_id, agent_type, task_description, task_category,
                    detected_domain, recommended_tool, selected_tool, selection_status,
                    selection_timestamp, execution_time_ms, success, error_message,
                    performance_score
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """, (
                selection.session_id, selection.agent_type, selection.task_description,
                selection.task_category.value, selection.detected_domain.value,
                selection.recommended_tool, selection.selected_tool, selection.selection_status.value,
                selection.selection_timestamp, selection.execution_time_ms,
                selection.success, selection.error_message, selection.performance_score
            ))

    def track_session_performance(self, session_id: str, mode: str = "serial", max_calls: int = 8) -> PerformanceMetrics:
        """Track performance metrics for a session."""
        with self.session_lock:
            if session_id not in self.session_metrics:
                metrics = PerformanceMetrics(
                    session_id=session_id,
                    total_calls=max_calls,
                    calls_remaining=max_calls,
                    execution_mode=mode,
                    start_time=datetime.now(),
                    last_update=datetime.now(),
                    average_execution_time=0.0,
                    success_rate=100.0
                )
                self.session_metrics[session_id] = metrics

                # Store in database
                with sqlite3.connect(self.db_path) as conn:
                    conn.execute("""
                        INSERT OR REPLACE INTO session_metrics (
                            session_id, total_calls, calls_remaining, execution_mode,
                            start_time, last_update, average_execution_time, success_rate
                        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                    """, (
                        session_id, max_calls, max_calls, mode,
                        metrics.start_time, metrics.last_update,
                        metrics.average_execution_time, metrics.success_rate
                    ))
            return self.session_metrics[session_id]

    def update_session_call(self, session_id: str, execution_time_ms: int, success: bool):
        """Update session metrics after a tool call."""
        with self.session_lock:
            if session_id in self.session_metrics:
                metrics = self.session_metrics[session_id]
                metrics.calls_remaining = max(0, metrics.calls_remaining - 1)
                metrics.last_update = datetime.now()

                # Update average execution time (simple moving average)
                if metrics.average_execution_time == 0:
                    metrics.average_execution_time = execution_time_ms
                else:
                    metrics.average_execution_time = (metrics.average_execution_time + execution_time_ms) / 2

                # Update success rate
                total_calls_made = metrics.total_calls - metrics.calls_remaining
                if total_calls_made > 0:
                    current_success_count = metrics.success_rate * total_calls_made / 100
                    if success:
                        current_success_count += 1
                    metrics.success_rate = (current_success_count / total_calls_made) * 100

    def check_performance_constraints(self, session_id: str) -> Dict[str, Any]:
        """Check if session meets performance constraints."""
        if session_id not in self.session_metrics:
            return {"violation": False, "message": "Session not tracked"}

        metrics = self.session_metrics[session_id]
        violations = []

        # Check call limit
        if metrics.calls_remaining < 0:
            violations.append(f"Exceeded call limit: {metrics.calls_remaining} calls remaining")

        # Check execution mode (should be serial)
        if metrics.execution_mode != "serial":
            violations.append(f"Non-serial execution mode: {metrics.execution_mode}")

        return {
            "violation": len(violations) > 0,
            "violations": violations,
            "calls_remaining": metrics.calls_remaining,
            "execution_mode": metrics.execution_mode,
            "average_execution_time": metrics.average_execution_time,
            "success_rate": metrics.success_rate
        }

    def generate_system_prompt(self, agent_type: str = "general", include_rules: bool = True) -> str:
        """Generate system prompt with routing rules for specific agent type."""
        base_prompt = self.config.get('systemPrompt', '')

        if not include_rules:
            return base_prompt

        # Add Phase 5 enforcement instructions
        enforcement_addition = f"""

## PHASE 5 ROUTING ENFORCEMENT ({agent_type.upper()} AGENT)

### MANDATORY COMPLIANCE RULES
1. **DOMAIN-SPECIFIC FIRST**: Always use domain-specific tools before generic ones
   - GitHub operations → github MCP (NEVER use fetch/curl)
   - Web scraping → firecrawl MCP (NEVER use fetch)
   - Code analysis → code-index MCP (NEVER manually scan files)
   - API testing → schemathesis MCP
   - Web search → brave-search MCP (NEVER use fetch)
   - File operations → filesystem MCP (NEVER use fetch for local files)

2. **PERFORMANCE CONSTRAINTS**:
   - Max 8 tool calls per task
   - Serial execution only (maxParallelCalls: 1)
   - Minimize tool calls by batching when possible

3. **ROUTING VALIDATION**:
   - Your tool selections are monitored and validated
   - Non-compliant selections will be flagged
   - Follow the recommended routing patterns

4. **AGENT ACCOUNTABILITY**:
   - Always choose the most appropriate tool for the domain
   - Document your tool selection reasoning
   - Report violations of these routing rules

### CURRENT ROUTING RULES SUMMARY
"""

        # Add specific rules for this agent type
        agent_specific_rules = self._get_agent_specific_rules(agent_type)
        enforcement_addition += agent_specific_rules

        enforcement_addition += """

### ENFORCEMENT MONITORING
Your tool selections will be tracked and analyzed for compliance with these routing rules.
Non-compliant selections may result in task failure or forced tool rerouting.

Choose tools wisely and follow domain-specific routing patterns for optimal performance.
"""

        return base_prompt + enforcement_addition

    def _get_agent_specific_rules(self, agent_type: str) -> str:
        """Get agent-specific routing rules."""
        agent_rules = {
            "cline": """
**CLINE AGENT SPECIFIC RULES**:
- Prioritize GitHub operations for repository management
- Use code-index for all code analysis and search tasks
- Follow development workflow: code-index → filesystem → github → memory-bank
- Never use fetch for GitHub API operations
""",
            "kilocode": """
**KILOCODE AGENT SPECIFIC RULES**:
- Focus on development and code analysis tasks
- Use code-index for code understanding and search
- Follow development patterns: analyze → implement → test → document
- Prioritize filesystem for local file operations
""",
            "general": """
**GENERAL AGENT RULES**:
- Use domain-specific tools as primary choice
- Fall back to alternative tools only when primary is unavailable
- Minimize tool calls and use batch operations
- Report any routing conflicts or issues
"""
        }

        return agent_rules.get(agent_type.lower(), agent_rules["general"])

    def analyze_routing_patterns(self, days_back: int = 30) -> Dict[str, Any]:
        """Analyze routing patterns for optimization recommendations."""
        try:
            with sqlite3.connect(self.db_path) as conn:
                # Get decisions from specified period
                cutoff_date = datetime.now() - timedelta(days=days_back)
                cursor = conn.execute("""
                    SELECT
                        detected_domain, recommended_tool, selected_tool,
                        selection_status, success, execution_time_ms
                    FROM routing_decisions
                    WHERE selection_timestamp >= ?
                """, (cutoff_date,))

                decisions = cursor.fetchall()

                if not decisions:
                    return {"error": "No routing decisions found for analysis period"}

                # Analyze patterns
                total_decisions = len(decisions)
                compliant_selections = sum(1 for d in decisions if d[3] == 'compliant')
                successful_tasks = sum(1 for d in decisions if d[4])

                # Domain analysis
                domain_stats = defaultdict(lambda: {"total": 0, "compliant": 0, "successful": 0})
                # Properly typed tool effectiveness tracking
                tool_effectiveness: Dict[str, Dict[str, Union[int, List[int]]]] = {}

                for decision in decisions:
                    domain, recommended, selected, status, success, exec_time = decision

                    domain_stats[domain]["total"] += 1
                    if status == 'compliant':
                        domain_stats[domain]["compliant"] += 1
                    if success:
                        domain_stats[domain]["successful"] += 1

                    # Initialize tool if not seen before
                    if selected not in tool_effectiveness:
                        tool_effectiveness[selected] = {"used": 0, "successful": 0, "avg_time": []}

                    tool_effectiveness[selected]["used"] += 1
                    if success:
                        tool_effectiveness[selected]["successful"] += 1
                    if exec_time:
                        tool_effectiveness[selected]["avg_time"].append(exec_time)

                # Calculate metrics
                analysis = {
                    "analysis_period_days": days_back,
                    "total_decisions": total_decisions,
                    "compliance_rate": (compliant_selections / total_decisions) * 100,
                    "success_rate": (successful_tasks / total_decisions) * 100,
                    "domain_statistics": dict(domain_stats),
                    "tool_effectiveness": {
                        tool: {
                            "usage_count": stats["used"],
                            "success_rate": (stats["successful"] / stats["used"]) * 100 if stats["used"] > 0 else 0,
                            "average_execution_time": sum(stats["avg_time"]) / len(stats["avg_time"]) if stats["avg_time"] else 0
                        }
                        for tool, stats in tool_effectiveness.items()
                    },
                    "recommendations": self._generate_optimization_recommendations(
                        domain_stats, tool_effectiveness
                    )
                }

                return analysis

        except Exception as e:
            logger.error(f"Error analyzing routing patterns: {e}")
            return {"error": f"Analysis failed: {str(e)}"}

    def _generate_optimization_recommendations(self, domain_stats: Dict, tool_effectiveness: Dict) -> List[Dict[str, Any]]:
        """Generate optimization recommendations based on analysis."""
        recommendations = []

        # Check domain compliance rates
        for domain, stats in domain_stats.items():
            if stats["total"] > 0:
                compliance_rate = (stats["compliant"] / stats["total"]) * 100
                success_rate = (stats["successful"] / stats["total"]) * 100

                if compliance_rate < 80:
                    recommendations.append({
                        "type": "low_compliance",
                        "domain": domain,
                        "issue": f"Low compliance rate: {compliance_rate:.1f}%",
                        "recommendation": "Improve agent training on domain-specific routing rules",
                        "priority": "high"
                    })

                if success_rate < 70:
                    recommendations.append({
                        "type": "low_success_rate",
                        "domain": domain,
                        "issue": f"Low success rate: {success_rate:.1f}%",
                        "recommendation": "Review tool effectiveness and consider tool replacement",
                        "priority": "high"
                    })

        # Check tool performance
        for tool, stats in tool_effectiveness.items():
            if stats["used"] > 5:  # Only consider tools used frequently
                success_rate = (stats["successful"] / stats["used"]) * 100
                avg_time = stats["average_execution_time"]

                if success_rate < 60:
                    recommendations.append({
                        "type": "poor_tool_performance",
                        "tool": tool,
                        "issue": f"Poor performance: {success_rate:.1f}% success rate",
                        "recommendation": f"Consider alternative tools or improve {tool} configuration",
                        "priority": "medium"
                    })

                if avg_time > 10000:  # More than 10 seconds
                    recommendations.append({
                        "type": "slow_tool",
                        "tool": tool,
                        "issue": f"Slow execution: {avg_time:.0f}ms average",
                        "recommendation": f"Optimize {tool} performance or implement caching",
                        "priority": "medium"
                    })

        return recommendations


def main():
    """Example usage of the routing enforcement system."""
    # Initialize routing engine
    engine = CipherRoutingEngine()

    # Test domain detection
    test_tasks = [
        "Search for Python tutorials on GitHub",
        "Scrape product data from e-commerce website",
        "Analyze code for security vulnerabilities",
        "Test API endpoint with OpenAPI schema",
        "Read local Python file for errors"
    ]

    print("=== DOMAIN DETECTION TEST ===")
    for task in test_tasks:
        domain, recommended_tool = engine.get_routing_recommendation(task)
        print(f"Task: {task}")
        print(f"Detected Domain: {domain.value}")
        print(f"Recommended Tool: {recommended_tool}")
        print()

    # Test tool validation
    print("=== TOOL VALIDATION TEST ===")
    selections = [
        ("session-123", "cline", "Search GitHub for repositories", "github", "fetch"),
        ("session-124", "kilocode", "Scrape website data", "firecrawl", "firecrawl_extract"),
        ("session-125", "general", "Read local file", "filesystem", "filesystem_read"),
    ]

    for session_id, agent_type, task, expected_tool, selected_tool in selections:
        selection = engine.validate_tool_selection(
            session_id, agent_type, task, selected_tool
        )
        print(f"Session: {session_id}")
        print(f"Task: {task}")
        print(f"Selected Tool: {selected_tool}")
        print(f"Recommended Tool: {selection.recommended_tool}")
        print(f"Status: {selection.selection_status.value}")
        print(f"Domain: {selection.detected_domain.value}")
        print()

    # Test performance tracking
    print("=== PERFORMANCE TRACKING TEST ===")
    engine.track_session_performance("session-123", "serial", 8)
    engine.update_session_call("session-123", 2500, True)
    engine.update_session_call("session-123", 1800, True)

    constraints = engine.check_performance_constraints("session-123")
    print(f"Performance Constraints Check: {constraints}")

    # Generate system prompt
    print("=== SYSTEM PROMPT GENERATION TEST ===")
    system_prompt = engine.generate_system_prompt("cline", include_rules=True)
    print(f"Generated prompt length: {len(system_prompt)} characters")
    print(f"First 500 characters: {system_prompt[:500]}...")

    print("\n=== ROUTING ENFORCEMENT SYSTEM READY ===")


if __name__ == "__main__":
    main()