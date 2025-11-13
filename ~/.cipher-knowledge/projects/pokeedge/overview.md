# PokeEdge-Specific Knowledge

*PokeEdge Web Application - extracted knowledge from constitution and .clinerules*

## 📋 Project Overview

PokeEdge is a modern web application focused on Pokemon data with:
- **Technology Stack**: React-based web application with Node.js backend
- **Testing Framework**: Playwright MCP for end-to-end testing
- **Quality Focus**: Production-grade quality, no shortcuts or technical debt
- **Architecture**: Simplified patterns with abstraction over complexity

## 🎯 Core PokeEdge Principles

### 1. **Quality-First Engineering (NON-NEGOTIABLE)**
- **No shortcuts or hacks** - implement the correct solution
- **Zero technical debt by default** - no TODOs without tracking
- **Immediate refactor** - complete migrations in same change set
- **Production-grade quality** - avoid placeholder code

### 2. **Simplification Cascades**
- **One powerful abstraction** over many special-cased implementations
- **Optimize for deletions** - measure success by components removed
- **Unifying insights** - find the single concept that collapses complexity
- **No parallel "old + new"** - complete migration promptly

### 3. **Test-First Development**
- **Tests written first** - TDD mandatory for all features
- **Red-Green-Refactor** cycle strictly enforced
- **Integration testing focus** - new library contracts, service communication
- **Playwright MCP for E2E** - browser automation via MCP server

## 🏗️ PokeEdge Architecture Patterns

### Simplification Principle
```javascript
// ❌ BAD - Multiple similar implementations
function handleAmazonPricing(item) { /* Amazon-specific logic */ }
function handleEbayPricing(item) { /* eBay-specific logic */ }
function handleWalmartPricing(item) { /* Walmart-specific logic */ }

// ✅ GOOD - One unifying abstraction
function calculatePricing(item, sources) {
  return sources.reduce((price, source) => {
    const discount = getDiscountForSource(source, item.quantity);
    return price * (1 - discount);
  }, item.basePrice);
}
```

### Quality-First Implementation
```javascript
// ✅ GOOD - Production-grade implementation
async function fetchPokemonData(pokemonId) {
  // Input validation
  if (!pokemonId || typeof pokemonId !== 'string') {
    throw new ValidationError('Invalid pokemon ID');
  }

  // Proper error handling
  try {
    const response = await fetch(`/api/pokemon/${pokemonId}`);
    if (!response.ok) {
      throw new APIError(`Failed to fetch: ${response.status}`);
    }

    const data = await response.json();

    // Business rule validation
    if (!data.id || !data.name) {
      throw new DataIntegrityError('Invalid pokemon data structure');
    }

    return data;
  } catch (error) {
    logger.error('Pokemon fetch failed', { pokemonId, error: error.message });
    throw error;
  }
}

// ❌ BAD - Quick hack that becomes technical debt
async function getPokemon(id) {
  return await fetch(`/api/pokemon/${id}`); // No validation, error handling, logging
}
```

## 🧪 PokeEdge Testing Standards

### Playwright MCP Usage (REQUIRED)
```javascript
// Always use Playwright MCP server, never install Playwright locally
use_mcp_tool({
  server_name: "github.com/executeautomation/mcp-playwright",
  tool_name: "playwright_navigate",
  arguments: {
    url: "http://localhost:5176/",
    timeout: 30000,
  },
});

use_mcp_tool({
  server_name: "github.com/executeautomation/mcp-playwright",
  tool_name: "playwright_screenshot",
  arguments: {
    name: "pokemon-search-results",
    fullPage: true,
  },
});
```

### Test-First Example
```javascript
// ✅ GOOD - Test written before implementation
describe('Pokemon Search', () => {
  it('should return results for valid search term', async () => {
    // Test setup
    const searchTerm = 'pikachu';

    // Execute search
    const results = await searchPokemon(searchTerm);

    // Assertions
    expect(results).toHaveLength(expect.any(Number));
    expect(results[0]).toHaveProperty('name');
    expect(results[0]).toHaveProperty('id');
  });

  it('should handle empty search results gracefully', async () => {
    const results = await searchPokemon('nonexistent');
    expect(results).toEqual([]);
  });
});
```

## 🔧 PokeEdge Development Workflow

### Definition of Done
For every PokeEdge change:
- [ ] **Tests updated/added** and passing (unit/integration/e2e as appropriate)
- [ ] **Linters/formatters pass** - no new warnings
- [ ] **Documentation updated** - README/ADRs/migration steps
- [ ] **No temporary placeholders** - unless explicitly approved with owner/deadline
- [ ] **No unused symbols/imports** - clean up during change
- [ ] **Refactors completed** - obsolete paths removed

### Code Review Checklist
- [ ] **Quality standards met** - production-grade implementation
- [ ] **Simplification applied** - no unnecessary complexity
- [ ] **Tests comprehensive** - covers edge cases and error conditions
- [ ] **Error handling robust** - proper logging and user feedback
- [ ] **No shortcuts** - correct solution implemented
- [ ] **Migration plan** - if refactoring legacy code

## 🎨 PokeEdge UI/UX Patterns

### React Component Standards
```jsx
// ✅ GOOD - Clean, testable component
function PokemonCard({ pokemon, onSelect }) {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSelect = async () => {
    setIsLoading(true);
    setError(null);

    try {
      await onSelect(pokemon);
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="pokemon-card">
      {error && <div className="error">{error}</div>}
      <img src={pokemon.image} alt={pokemon.name} />
      <h3>{pokemon.name}</h3>
      <button
        onClick={handleSelect}
        disabled={isLoading}
      >
        {isLoading ? 'Selecting...' : 'Select'}
      </button>
    </div>
  );
}
```

## 🔄 PokeEdge Simplification Examples

### Before: Complex Special Cases
```javascript
// Multiple similar functions that could be unified
function getAmazonPrice(item) { /* Amazon API logic */ }
function getEbayPrice(item) { /* eBay API logic */ }
function getWalmartPrice(item) { /* Walmart API logic */ }
function getBestPrice(item) {
  const prices = [
    getAmazonPrice(item),
    getEbayPrice(item),
    getWalmartPrice(item)
  ].filter(p => p !== null);

  return Math.min(...prices);
}
```

### After: Unified Abstraction
```javascript
// One function that handles all sources via configuration
const PRICE_SOURCES = {
  amazon: { api: AmazonAPI, discount: 0.05 },
  ebay: { api: EbayAPI, discount: 0.03 },
  walmart: { api: WalmartAPI, discount: 0.02 }
};

async function getPrice(item, sources = Object.keys(PRICE_SOURCES)) {
  const prices = await Promise.all(
    sources.map(async (source) => {
      const config = PRICE_SOURCES[source];
      const basePrice = await config.api.getPrice(item);
      return basePrice * (1 - config.discount);
    })
  );

  return Math.min(...prices.filter(p => p !== null));
}
```

## 🚨 PokeEdge Red Flags

### Quality Issues (Stop and Fix)
- **Commented-out code** in production
- **TODO comments** without tracking
- **Quick workarounds** instead of proper solutions
- **Duplicate logic** across components
- **Missing error handling** in user-facing code

### Complexity Issues (Simplify)
- **Growing list of special cases** in pricing logic
- **Multiple implementations** of similar patterns
- **"Just one more case"** keeps being added
- **Parallel old vs new** systems running indefinitely
- **Configuration options proliferating** instead of sensible defaults

## 📁 PokeEdge File Organization

```
web-app-pokeedge/
├── src/
│   ├── components/        # React components
│   │   ├── common/        # Shared components
│   │   ├── pokemon/       # Pokemon-specific components
│   │   └── search/        # Search-related components
│   ├── services/          # API and business logic
│   │   ├── pokemon/       # Pokemon data services
│   │   ├── pricing/       # Price calculation services
│   │   └── api/           # HTTP clients
│   ├── utils/             # Shared utilities
│   ├── hooks/             # Custom React hooks
│   └── types/             # TypeScript definitions
├── tests/
│   ├── unit/              # Unit tests
│   ├── integration/       # Service integration tests
│   └── e2e/               # Playwright MCP tests
├── .clinerules/           # Development guidelines
└── docs/                  # Project documentation
```

## 🎯 PokeEdge Success Metrics

### Code Quality Metrics
- **Zero TODO comments** in production code
- **100% test coverage** for business logic
- **< 3 seconds** page load times
- **Zero critical** linting warnings
- **Single abstraction** for each concept

### Simplification Metrics
- **Net negative LOC** after refactoring
- **Reduced branching** complexity
- **Eliminated parallel systems**
- **Simplified configuration** surface
- **Measurable deletion** of complexity

---

**Remember**: PokeEdge prioritizes production-grade quality over speed. Every change should make the codebase simpler, more maintainable, and more reliable. If you find yourself adding complexity, stop and look for the unifying abstraction that could eliminate it entirely.
