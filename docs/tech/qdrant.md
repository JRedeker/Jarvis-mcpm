# Qdrant Vector Store Documentation

**GitHub**: https://github.com/qdrant/qdrant
**Stars**: ~20k
**Status**: Production-Ready
**Purpose**: High-performance vector similarity search engine

---

## Overview

Qdrant is the vector store used by Cipher for:
- 🧠 Semantic search in knowledge memory
- 🔍 Vector similarity matching
- 📊 High-dimensional embeddings storage
- ⚡ Fast retrieval (< 10ms typical)

---

## Research Required

**⚠️ Documentation to be completed using Context7 + GPT-Researcher**

### Context7 Query Plan
```
Library: qdrant/qdrant
Focus Areas:
- Cloud vs self-hosted deployment
- Docker setup for development
- Production configuration
- Performance tuning
- Collection management
- Integration with embedding models
- Backup and recovery
```

### GPT-Researcher Query Plan
```json
{
  "task": "Research Qdrant vector database deployment and integration with Node.js applications",
  "report_type": "research_report",
  "sources": [
    "Qdrant official documentation",
    "Qdrant Cloud setup guide",
    "Docker deployment best practices",
    "Performance optimization guides"
  ]
}
```

---

## Cipher Integration

**Current Setup** (from cipher-aggregator.md):
```yaml
memory:
  vectorStore:
    type: qdrant
    url: ${QDRANT_URL}
```

**Configuration Options** (from existing docs):
```bash
# Qdrant Cloud
VECTOR_STORE_TYPE=qdrant
VECTOR_STORE_URL=https://your-cluster.qdrant.io
VECTOR_STORE_API_KEY=your-qdrant-api-key

# Local Docker
VECTOR_STORE_TYPE=qdrant
VECTOR_STORE_HOST=localhost
VECTOR_STORE_PORT=6333
VECTOR_STORE_URL=http://localhost:6333
```

---

## MetaMCP Hybrid Considerations

**Questions for Research**:
- Does Cipher default mode need its own Qdrant instance?
- Or can MetaMCP-managed Cipher share Qdrant with other services?
- Performance impact of multiple Cipher instances accessing same Qdrant
- Collection naming strategy to avoid conflicts

---

## Research Deliverables

- [ ] Qdrant Cloud vs self-hosted comparison
- [ ] Docker compose configuration
- [ ] Production deployment checklist
- [ ] Backup/restore procedures
- [ ] Performance benchmarks for Cipher workload
- [ ] Multi-tenant access patterns