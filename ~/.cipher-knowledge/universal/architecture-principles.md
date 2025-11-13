# Universal Architecture Principles

*Applies to ALL projects - scalable and maintainable software architecture*

## 🏗️ Layered Architecture

Every project should follow clear layer separation with well-defined responsibilities:

### Layer Hierarchy (Top → Bottom)
```
┌─────────────────────────────────────┐
│            UI Layer                │  ← User Interface
│     (TUI, Web, CLI, API UI)        │
├─────────────────────────────────────┤
│           API Layer                │  ← HTTP/REST Endpoints
│    (Request/Response Handling)     │
├─────────────────────────────────────┤
│         Business Logic             │  ← Core Application Logic
│         (Services)                 │  ← Use Cases & Rules
├─────────────────────────────────────┤
│          Data Layer                │  ← Database Access
│        (Repositories)              │  ← Data Persistence
├─────────────────────────────────────┤
│        External Services           │  ← Third-party APIs
│         (Clients)                  │  ← External Systems
├─────────────────────────────────────┤
│         Models/Domain              │  ← Core Data Models
│         (Entities)                 │  ← Business Objects
└─────────────────────────────────────┘
```

## 📋 Layer Responsibilities

### 1. **UI Layer** (User Interface)
**Responsibilities:**
- Present information to users
- Handle user input and interactions
- Route requests to appropriate handlers
- Render views and components

**What belongs here:**
- Textual TUI components (PokeEdge)
- API response models
- Input validation (basic)
- Session management

**❌ What does NOT belong:**
- Business logic
- Database access
- External API calls
- Complex data processing

### 2. **API Layer** (Interface Boundaries)
**Responsibilities:**
- Define request/response contracts
- Handle HTTP concerns (status codes, headers)
- Route requests to business logic
- Input/output serialization

**What belongs here:**
- HTTP endpoint definitions
- Request/response models (Pydantic)
- HTTP status code handling
- Basic validation (format, types)

**❌ What does NOT belong:**
- Business rules
- Data access logic
- Complex transformations

### 3. **Business Logic Layer** (Core Application)
**Responsibilities:**
- Implement business rules and use cases
- Coordinate between different services
- Handle complex business operations
- Enforce business invariants

**What belongs here:**
- Service classes
- Use case implementations
- Business rule validation
- Process orchestration

**Example (PokeEdge pricing):**
```python
class PricingService:
    def calculate_total(self, items: List[Item], discounts: List[Discount]) -> float:
        """Calculate total price with all applicable discounts."""
        subtotal = sum(item.base_price for item in items)

        # Apply bulk discounts
        for discount in discounts:
            if discount.type == "bulk":
                subtotal = self._apply_bulk_discount(subtotal, discount)

        # Apply source-specific discounts
        for item in items:
            if item.source in DISCOUNT_SOURCES:
                subtotal *= DISCOUNT_SOURCES[item.source]

        return round(subtotal, 2)
```

### 4. **Data Layer** (Persistence)
**Responsibilities:**
- Abstract database operations
- Provide data access interfaces
- Handle data transformation
- Manage transactions

**What belongs here:**
- Repository classes
- Database connection management
- Query builders
- Data mappers

**Example:**
```python
class ItemRepository:
    async def get_by_id(self, item_id: str) -> Optional[Item]:
        """Retrieve item by ID with proper error handling."""
        async with self.db_pool as conn:
            result = await conn.fetchrow(
                "SELECT * FROM items WHERE id = $1", item_id
            )
            return self._map_to_model(result) if result else None
```

### 5. **External Services Layer** (Third-party Integration)
**Responsibilities:**
- Abstract external API calls
- Handle authentication and timeouts
- Transform external data formats
- Provide retry logic and error handling

**What belongs here:**
- API client classes
- Authentication handling
- Data format conversion
- Error handling and retries

**Example:**
```python
class FireCrawlClient:
    def __init__(self, api_key: str):
        self.api_key = api_key
        self.base_url = "https://api.firecrawl.dev"

    async def scrape_url(self, url: str) -> ScrapeResult:
        """Scrape a single URL with proper error handling."""
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{self.base_url}/scrape",
                headers={"Authorization": f"Bearer {self.api_key}"},
                json={"url": url, "formats": ["markdown", "html"]},
                timeout=30.0
            )
            response.raise_for_status()
            return ScrapeResult(**response.json())
```

### 6. **Models/Domain Layer** (Business Objects)
**Responsibilities:**
- Define core business entities
- Encapsulate business rules
- Provide type safety
- Maintain data integrity

**What belongs here:**
- Data models (Pydantic/SQLAlchemy)
- Value objects
- Domain exceptions
- Business rules

**Example:**
```python
@dataclass
class Item:
    id: str
    name: str
    price: float
    source: str

    def validate(self) -> None:
        """Business rule validation."""
        if self.price < 0:
            raise InvalidItemError("Price cannot be negative")

        if self.source not in VALID_SOURCES:
            raise InvalidItemError(f"Invalid source: {self.source}")
```

## 🔄 Data Flow Patterns

### Request Flow
```
User Input → UI Layer → API Layer → Business Logic → Data Layer → Database
     ↑                                                                    ↓
Response ← UI Layer ← API Layer ← Business Logic ← Data Layer ← Database
```

### Key Principles:
1. **Dependencies flow downward** - UI depends on API, API on Business Logic, etc.
2. **No circular dependencies** - strict layer ordering
3. **Interface segregation** - each layer exposes clear interfaces
4. **Single responsibility** - each layer has one clear purpose

## 🌉 Cross-Cutting Concerns

### Error Handling
- **Consistent error format** across all layers
- **Domain exceptions** for business rule violations
- **Technical exceptions** for system failures
- **Error propagation** - let exceptions bubble up appropriately

### Logging
- **Structured logging** with consistent format
- **Context preservation** across layer boundaries
- **Security considerations** - no sensitive data in logs
- **Performance tracking** at layer boundaries

### Configuration
- **Environment-based** configuration
- **Type-safe** configuration loading
- **Default values** for development
- **Validation** of configuration at startup

## 🏛️ Architectural Patterns

### Repository Pattern
```python
# Abstract interface
class ItemRepositoryProtocol(Protocol):
    async def get_by_id(self, item_id: str) -> Optional[Item]: ...
    async def create(self, item: Item) -> Item: ...
    async def update(self, item: Item) -> Item: ...
    async def delete(self, item_id: str) -> None: ...

# Concrete implementation
class PostgreSQLItemRepository(ItemRepositoryProtocol):
    async def get_by_id(self, item_id: str) -> Optional[Item]:
        # Database implementation
        pass
```

### Service Layer Pattern
```python
class ItemService:
    def __init__(self, repository: ItemRepositoryProtocol):
        self._repository = repository

    async def get_item_with_pricing(self, item_id: str) -> ItemWithPricing:
        """Business logic coordination."""
        item = await self._repository.get_by_id(item_id)
        if not item:
            raise ItemNotFoundError(f"Item {item_id} not found")

        pricing = await self._calculate_pricing(item)
        return ItemWithPricing(item=item, pricing=pricing)
```

### Factory Pattern
```python
class ServiceFactory:
    def __init__(self, repository: ItemRepositoryProtocol):
        self._repository = repository

    def create_pricing_service(self) -> PricingService:
        return PricingService(self._repository)

    def create_item_service(self) -> ItemService:
        return ItemService(self._repository)
```

## 🔒 Dependency Management

### Dependency Injection
```python
# In PokeEdge main.py or app.py
def create_app() -> FastAPI:
    """Application factory with dependency injection."""
    repository = PostgreSQLItemRepository(db_pool)
    service_factory = ServiceFactory(repository)

    app = FastAPI()

    # Register API endpoints with injected services
    app.include_router(
        items_router(service_factory),
        prefix="/api/v1",
        tags=["items"]
    )

    return app
```

### Interface Segregation
```python
# Instead of one big interface
class AllOperations(Protocol):
    async def create(self, item: Item) -> Item: ...
    async def read(self, item_id: str) -> Optional[Item]: ...
    async def update(self, item: Item) -> Item: ...
    async def delete(self, item_id: str) -> None: ...
    async def search(self, query: str) -> List[Item]: ...  # Extra method

# Use focused interfaces
class CRUDOperations(Protocol):
    async def create(self, item: Item) -> Item: ...
    async def read(self, item_id: str) -> Optional[Item]: ...
    async def update(self, item: Item) -> Item: ...
    async def delete(self, item_id: str) -> None: ...

class SearchOperations(Protocol):
    async def search(self, query: str) -> List[Item]: ...
```

## 🧪 Architecture Testing

### Layer Testing Strategy
- **Unit tests** for each layer in isolation
- **Integration tests** for layer interactions
- **Contract tests** for API boundaries
- **End-to-end tests** for complete workflows

### Mock Dependencies
```python
# Test business logic without database
def test_pricing_calculation():
    # Mock the repository
    mock_repo = Mock(spec=ItemRepositoryProtocol)
    mock_repo.get_by_id.return_value = Item(id="1", price=100.0, source="amazon")

    service = PricingService(mock_repo)
    result = service.calculate_price("1", quantity=20)

    assert result == 1800.0  # 20 * 100 * 0.9 bulk discount
```

## 📈 Scalability Considerations

### Horizontal Scaling
- **Stateless services** - no in-memory state
- **Database connection pooling**
- **Distributed caching** (Redis)
- **Message queues** for async processing

### Performance Optimization
- **Lazy loading** of data
- **Caching** at appropriate layers
- **Database query optimization**
- **Connection pooling**

### Monitoring
- **Health checks** for each layer
- **Performance metrics** per layer
- **Error tracking** and alerting
- **Business metrics** (not just technical)

---

**Remember**: Good architecture is like good plumbing - you don't notice it when it works, but you definitely notice when it doesn't. Invest in clean architecture from the start.
