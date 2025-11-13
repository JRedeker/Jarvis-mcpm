# Universal Code Quality Standards

*Applies to ALL projects - ensuring maintainable, readable, and reliable code*

## 🎯 Quality Gates

Every code change must pass these fundamental quality checks before merging:

### ✅ Mandatory Quality Standards
- [ ] **Auto-commit enforcement** (if project uses auto-commit)
- [ ] **Type hints** everywhere (Python projects)
- [ ] **Linting** passes completely
- [ ] **Type checking** passes completely
- [ ] **Test coverage** meets project requirements
- [ ] **Documentation** updated for API changes

### 🔄 Automated Quality Gates
- **Pre-commit hooks** configured and working
- **CI/CD pipeline** includes quality checks
- **Test automation** runs on every change
- **Security scanning** automated (no secrets committed)

## 🐍 Python-Specific Standards

### Type Hints (MANDATORY)
```python
# ✅ GOOD - Complete type hints
def process_data(user_id: str, options: dict[str, Any]) -> tuple[User, bool]:
    pass

# ❌ BAD - Missing type hints
def process_data(user_id, options):
    pass
```

### Linting Requirements
- **mypy** for type checking
- **black** for code formatting
- **isort** for import sorting
- **ruff** for additional linting
- **No linting errors** tolerated

### Package Management
- **Workspace .venv** always active
- **uv** preferred for package management
- **requirements.txt** for pinning
- **pyproject.toml** for configuration
- **NO global installs** - everything in .venv

## 📁 Code Organization

### File Structure
```
project/
├── src/
│   ├── models/         # Data models
│   ├── services/       # Business logic
│   ├── api/           # HTTP endpoints
│   ├── ui/            # User interface
│   └── utils/         # Shared utilities
├── tests/
│   ├── unit/          # Unit tests
│   ├── integration/   # Integration tests
│   └── e2e/          # End-to-end tests
├── docs/              # Documentation
└── scripts/           # Automation scripts
```

### Naming Conventions
- **snake_case** for functions, variables, methods
- **PascalCase** for classes, exceptions
- **UPPER_CASE** for constants
- **Descriptive names** - avoid abbreviations
- **No generic names** like `data`, `temp`, `result`

## 📝 Documentation Standards

### Code Documentation
```python
def calculate_price(item: Item, quantity: int) -> float:
    """Calculate the total price with quantity discounts.

    Args:
        item: The item being purchased
        quantity: Number of items to purchase

    Returns:
        Total price in dollars with discount applied

    Raises:
        ValueError: If quantity is less than 1
        InvalidItemError: If item is not purchasable
    """
    pass
```

### API Documentation
- **OpenAPI/Swagger** for HTTP APIs
- **README.md** for setup and usage
- **CHANGELOG.md** for version history
- **API documentation** automatically generated
- **Examples** included in documentation

## 🧪 Testing Standards

### Test Organization
- **Unit tests**: 80%+ coverage for business logic
- **Integration tests**: API endpoints and database
- **E2E tests**: Critical user workflows
- **Test isolation**: Each test independent
- **Fast tests**: Unit tests < 1 second each

### Test Quality
```python
# ✅ GOOD - Clear test with proper assertions
def test_calculate_discount_for_bulk_purchase():
    item = Item(price=100.0, discount_tiers=[(10, 0.05), (50, 0.1)])
    total = calculate_price(item, quantity=20)
    assert total == 1800.0  # 20 * 100 * 0.9

# ❌ BAD - Vague test without clear assertions
def test_discount():
    result = calculate_price(item, 20)
    assert result > 0
```

### Test Data Management
- **Factory patterns** for test data
- **No hardcoded values** in tests
- **Test fixtures** for complex setup
- **Mock external services** consistently

## 🔧 Code Review Standards

### Review Checklist
- [ ] **Functionality**: Does it solve the problem?
- [ ] **Readability**: Is the code easy to understand?
- [ ] **Testing**: Are tests comprehensive?
- [ ] **Documentation**: Is documentation complete?
- [ ] **Security**: No security issues introduced?
- [ ] **Performance**: No performance regressions?
- [ ] **Maintainability**: Easy to modify in future?

### Code Review Process
- **Self-review** first - catch obvious issues
- **Automated checks** pass before human review
- **Specific feedback** - not just "looks good"
- **Constructive criticism** - focus on improvement

## 📊 Architecture Quality

### Code Complexity
- **Cyclomatic complexity** < 10 per function
- **Function length** < 50 lines
- **Class length** < 200 lines
- **File length** < 500 lines
- **Method complexity** - one clear purpose

### Design Patterns
- **DRY** (Don't Repeat Yourself)
- **SOLID** principles applied
- **Dependency injection** for testability
- **Strategy pattern** for algorithm variations
- **Observer pattern** for event handling

## 🏗️ Performance Standards

### Performance Budgets
- **API responses** < 500ms (95th percentile)
- **Page loads** < 2 seconds
- **Database queries** < 100ms
- **Background tasks** < 5 minutes
- **Test execution** < 5 minutes total

### Performance Monitoring
- **Response time tracking**
- **Memory usage monitoring**
- **Database query optimization**
- **Caching strategies** where appropriate
- **Load testing** before deployment

## 🛠️ Development Environment

### Standard Development Setup
```bash
# ✅ REQUIRED setup for Python projects
python -m venv .venv
source .venv/bin/activate  # or .venv\Scripts\activate on Windows
uv pip install -e .
pre-commit install
```

### Required Tools
- **Code editor** with linting integration
- **Git hooks** for quality checks
- **Testing framework** configured
- **Debugging tools** setup
- **Version control** with proper branching

## 🚀 Deployment Quality

### Release Process
- [ ] **Version tagging** semantic
- [ ] **Changelog** updated
- [ ] **Tests pass** in production-like environment
- [ ] **Performance benchmarks** met
- [ ] **Security scan** clean
- [ ] **Rollback plan** documented

### Production Readiness
- **Error handling** graceful degradation
- **Monitoring** and alerting configured
- **Backup strategy** implemented
- **Documentation** complete for operations
- **Incident response** procedures documented

---

**Remember**: Quality is not negotiable. It's easier to write quality code from the start than to fix it later. These standards protect both current and future developers.
