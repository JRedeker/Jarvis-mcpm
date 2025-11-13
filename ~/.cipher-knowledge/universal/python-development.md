# Python Development Standards

*Applies to ALL Python projects - modern Python development best practices*

## 🐍 Environment Management

### Virtual Environment (MANDATORY)
```bash
# ✅ ALWAYS use virtual environment
python -m venv .venv
source .venv/bin/activate  # Linux/Mac
# .venv\Scripts\activate   # Windows

# ✅ Install packages in development mode
uv pip install -e .
```

### Package Management
- **uv** preferred for package management
- **pip** acceptable for simple projects
- **poetry** acceptable if team familiar
- **requirements.txt** for production dependencies
- **requirements-dev.txt** for development dependencies
- **pyproject.toml** for project configuration

### Environment Files
```bash
# ✅ .env (never commit)
OPENAI_API_KEY=sk-...
DATABASE_URL=postgresql://...

# ✅ .env.example (safe to commit)
OPENAI_API_KEY=your_api_key_here
DATABASE_URL=your_database_url_here
```

## 📦 Dependency Management

### Development Dependencies
```python
# pyproject.toml
[project.optional-dependencies]
dev = [
    "pytest>=7.0.0",
    "pytest-cov>=4.0.0",
    "black>=23.0.0",
    "isort>=5.12.0",
    "mypy>=1.0.0",
    "ruff>=0.1.0",
    "pre-commit>=3.0.0",
]
```

### Version Pinning
- **Specific versions** for production dependencies
- **Upper bounds** for development dependencies
- **Compatible versions** tested together
- **Regular updates** scheduled

## 🧪 Testing Standards

### Test Framework
- **pytest** preferred for testing framework
- **unittest** acceptable for simple projects
- **hypothesis** for property-based testing
- **pytest-cov** for coverage reporting

### Test Organization
```
tests/
├── conftest.py           # Shared fixtures
├── unit/                 # Unit tests
│   ├── test_models.py
│   ├── test_services.py
│   └── test_utils.py
├── integration/          # Integration tests
│   ├── test_api.py
│   └── test_database.py
└── e2e/                  # End-to-end tests
    ├── test_user_flows.py
    └── test_apis.py
```

### Test Best Practices
```python
# ✅ GOOD - Clear test structure
def test_calculate_discount_bulk_purchase():
    """Test that bulk purchases get proper discount."""
    # Arrange
    item = create_test_item(price=100.0, discount_tiers=[(10, 0.1)])

    # Act
    result = calculate_price(item, quantity=20)

    # Assert
    assert result == 1800.0  # 20 * 100 * 0.9
    assert item.discount_applied == 0.1

# ❌ BAD - Vague test
def test_discount():
    result = calculate_price(item, 20)
    assert result > 0
```

### Fixtures and Factories
```python
# ✅ GOOD - Reusable test data
@pytest.fixture
def test_user():
    return UserFactory(email="test@example.com", status="active")

@pytest.fixture
def sample_items():
    return [
        ItemFactory(price=100.0, category="electronics"),
        ItemFactory(price=50.0, category="books"),
    ]
```

## 🔍 Type Hints (MANDATORY)

### Type Annotation Rules
```python
# ✅ GOOD - Complete type hints
from typing import Optional, Dict, List, Union, Any
from dataclasses import dataclass
from pathlib import Path

def process_data(
    user_id: str,
    data: Dict[str, Any],
    max_items: Optional[int] = None
) -> tuple[List[Item], bool]:
    pass

@dataclass
class User:
    id: str
    email: str
    preferences: Dict[str, Union[str, int, bool]] = None

# ❌ BAD - Missing type hints
def process_data(user_id, data, max_items=None):
    pass
```

### Type Checking
- **mypy** for static type checking
- **No Any** unless absolutely necessary
- **Generic types** for collections
- **Union** for optional/multiple types
- **Type aliases** for complex types

## 🏗️ Code Organization

### Package Structure
```
src/project_name/
├── __init__.py           # Package initialization
├── models/               # Data models
│   ├── __init__.py
│   ├── user.py
│   └── item.py
├── services/             # Business logic
│   ├── __init__.py
│   ├── user_service.py
│   └── item_service.py
├── api/                  # HTTP endpoints
│   ├── __init__.py
│   ├── v1/
│   │   ├── __init__.py
│   │   ├── users.py
│   │   └── items.py
│   └── dependencies.py
├── utils/                # Utilities
│   ├── __init__.py
│   ├── database.py
│   └── logging.py
└── config.py             # Configuration

tests/                    # Tests mirror src structure
```

### Import Organization
```python
# ✅ GOOD - Organized imports
# Standard library
import json
import logging
from pathlib import Path
from typing import Optional, Dict, List

# Third-party
import requests
from fastapi import FastAPI
from pydantic import BaseModel

# Local imports
from .models import User, Item
from .services import UserService
from .utils import setup_logging

# ❌ BAD - Disorganized imports
import json
from models import User
import requests
from services import UserService
import logging
from typing import Optional
```

## ⚡ Performance Standards

### Performance Guidelines
- **Profiling before optimization**
- **List comprehensions** over loops when appropriate
- **Generators** for large datasets
- **Caching** for expensive operations
- **Database query optimization**

### Async/Await Usage
```python
# ✅ GOOD - Proper async usage
async def fetch_user_data(user_id: str) -> User:
    """Fetch user data from database asynchronously."""
    async with get_database_connection() as conn:
        result = await conn.fetchone(
            "SELECT * FROM users WHERE id = $1", user_id
        )
        return User(**result) if result else None

# Use with proper async context
async def main():
    user = await fetch_user_data("123")
    if user:
        print(f"User: {user.email}")
```

### Database Access
- **AsyncPG** for PostgreSQL
- **SQLAlchemy** with async support
- **Connection pooling** for performance
- **Prepared statements** for security
- **Migration tools** (Alembic)

## 🛠️ Development Tools

### Code Formatting
```bash
# ✅ Automated formatting
black src/ tests/ --line-length 88
isort src/ tests/ --profile black
ruff check src/ tests/ --fix
mypy src/ --strict
```

### Pre-commit Configuration
```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.4.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files

  - repo: https://github.com/psf/black
    rev: 23.3.0
    hooks:
      - id: black

  - repo: https://github.com/pycqa/isort
    rev: 5.12.0
    hooks:
      - id: isort

  - repo: https://github.com/charliermarsh/ruff-pre-commit
    rev: v0.0.270
    hooks:
      - id: ruff
```

### Linting Configuration
```toml
# pyproject.toml
[tool.ruff]
line-length = 88
target-version = "py311"
select = [
    "E",   # pycodestyle errors
    "W",   # pycodestyle warnings
    "F",   # pyflakes
    "I",   # isort
    "C",   # flake8-comprehensions
    "B",   # flake8-bugbear
]

[tool.mypy]
python_version = "3.11"
warn_return_any = true
warn_unused_configs = true
disallow_untyped_defs = true
```

## 🚀 Deployment Standards

### Production Dependencies
- **Exact versions** in production
- **Security updates** tracked and applied
- **Dependency scanning** automated
- **License compatibility** verified

### Environment Configuration
```python
# config.py
import os
from pydantic import BaseSettings
from typing import Optional

class Settings(BaseSettings):
    database_url: str
    openai_api_key: str
    debug: bool = False
    log_level: str = "INFO"

    class Config:
        env_file = ".env"

settings = Settings()
```

### Health Checks
```python
# health.py
from fastapi import FastAPI
import asyncpg

async def health_check():
    """Check if all services are healthy."""
    checks = {
        "database": await check_database(),
        "external_api": await check_external_api(),
    }

    all_healthy = all(checks.values())
    status_code = 200 if all_healthy else 503

    return {"status": "healthy" if all_healthy else "unhealthy", "checks": checks}
```

## 📊 Monitoring and Logging

### Logging Configuration
```python
# logging.py
import logging
import structlog
from pythonjsonlogger import jsonlogger

def setup_logging():
    """Configure structured logging."""
    structlog.configure(
        processors=[
            structlog.processors.TimeStamper(fmt="iso"),
            structlog.processors.add_log_level,
            structlog.processors.JSONRenderer()
        ],
        context_class=dict,
        logger_factory=structlog.PrintLoggerFactory(),
        wrapper_class=structlog.make_filtering_bound_logger(20),  # INFO
        cache_logger_on_first_use=True,
    )
```

### Metrics and Monitoring
- **Prometheus** for metrics
- **Sentry** for error tracking
- **Health check endpoints**
- **Performance monitoring**
- **Resource usage tracking**

---

**Remember**: Python development is about clarity, efficiency, and maintainability. These standards ensure that Python code is not just functional, but elegant and robust.
