# GPT Researcher MCP Server

**GitHub**: https://github.com/assafelovic/gptr-mcp

---

## Overview

A custom Model Context Protocol server that provides comprehensive research capabilities using GPT-based analysis. This server enables LLMs to conduct deep research on topics, generate detailed reports, and analyze complex information.

**Key Features:**
- Deep research with multiple sources
- Comprehensive report generation
- Multi-step research process
- Source citation and verification
- Topic analysis and summarization
- Customizable research depth

---

## Installation

### Local Development
```bash
# Navigate to the gptr-mcp directory
cd gptr-mcp

# Install dependencies
pip install -r requirements.txt

# Run the server
python server.py
```

### Using the wrapper
```bash
# From the project root
python -m gptr-mcp.server
```

---

## Configuration

**Required Environment Variables:**
- `OPENAI_API_KEY`: OpenAI API key for GPT models
- `TAVILY_API_KEY`: Tavily API key for web search (optional but recommended)

**Optional Environment Variables:**
- `OPENAI_BASE_URL`: Custom OpenAI API base URL
- `DEFAULT_MODEL`: GPT model to use (default: gpt-4-turbo-preview)
- `MAX_ITERATIONS`: Maximum research iterations (default: 3)
- `RESEARCH_DEPTH`: Research depth level (default: "comprehensive")

**API Keys Setup:**
1. **OpenAI API Key**: Get from [OpenAI Platform](https://platform.openai.com/)
2. **Tavily API Key**: Get from [Tavily](https://tavily.com/)

---

## Available Tools

### Research Tools
1. `research_topic`
   - Conduct comprehensive research on a topic
   - Inputs: `topic` (required), `depth` (optional), `format` (optional)
   - Returns: Detailed research report with sources

2. `generate_report`
   - Generate structured research report
   - Inputs: `topic` (required), `sections` (optional), `style` (optional)
   - Returns: Formatted research report

### Analysis Tools
3. `analyze_sources`
   - Analyze and verify research sources
   - Inputs: `sources` (required), `criteria` (optional)
   - Returns: Source analysis and credibility assessment

4. `summarize_research`
   - Summarize existing research
   - Inputs: `content` (required), `length` (optional), `focus` (optional)
   - Returns: Condensed summary

### Specialized Research
5. `research_competitors`
   - Research competitors or alternatives
   - Inputs: `target` (required), `industry` (optional), `depth` (optional)
   - Returns: Competitive analysis report

6. `research_trends`
   - Research current trends and developments
   - Inputs: `field` (required), `timeframe` (optional), `region` (optional)
   - Returns: Trend analysis report

---

## Usage Examples

### Basic Research
```json
{
  "tool": "research_topic",
  "arguments": {
    "topic": "Model Context Protocol adoption in enterprise",
    "depth": "comprehensive",
    "format": "detailed"
  }
}
```

### Report Generation
```json
{
  "tool": "generate_report",
  "arguments": {
    "topic": "AI safety in 2024",
    "sections": ["executive_summary", "key_findings", "recommendations", "sources"],
    "style": "academic"
  }
}
```

### Competitor Research
```json
{
  "tool": "research_competitors",
  "arguments": {
    "target": "Anthropic Claude",
    "industry": "AI assistants",
    "depth": "detailed"
  }
}
```

### Trend Analysis
```json
{
  "tool": "research_trends",
  "arguments": {
    "field": "machine learning",
    "timeframe": "last_6_months",
    "region": "global"
  }
}
```

### Source Analysis
```json
{
  "tool": "analyze_sources",
  "arguments": {
    "sources": [
      "https://arxiv.org/abs/2024.12345",
      "https://techcrunch.com/2024/01/15/ai-news"
    ],
    "criteria": ["credibility", "recency", "relevance"]
  }
}
```

---

## Research Process

### Multi-step Research
1. **Topic Analysis**: Break down complex topics into subtopics
2. **Source Discovery**: Find relevant and credible sources
3. **Content Extraction**: Extract key information from sources
4. **Synthesis**: Combine information into coherent analysis
5. **Verification**: Cross-check facts and claims
6. **Report Generation**: Create structured final report

### Quality Assurance
- Source credibility assessment
- Fact-checking and verification
- Bias detection and mitigation
- Citation and reference management
- Content originality checking

---

## Report Formats

### Available Formats
- **Detailed**: Comprehensive analysis with full sources
- **Executive**: Concise summary for decision makers
- **Academic**: Formal research paper style
- **Technical**: Technical specifications and details
- **Casual**: Easy-to-read blog post style

### Report Sections
- Executive Summary
- Key Findings
- Detailed Analysis
- Supporting Evidence
- Sources and References
- Recommendations
- Future Research Directions

---

## Known Limitations

1. **API Costs**: Research can consume significant API tokens
2. **Research Time**: Deep research may take several minutes
3. **Source Availability**: Limited by available online sources
4. **Real-time Information**: May not include very recent developments
5. **Bias Potential**: Research quality depends on source diversity
6. **Complex Topics**: Very specialized topics may have limited sources

---

## Testing

```bash
# Set environment variables
export OPENAI_API_KEY="your_openai_key"
export TAVILY_API_KEY="your_tavily_key"  # Optional

# Test server startup
python gptr-mcp/server.py

# Test tool availability
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | python gptr-mcp/server.py

# Test research functionality (requires MCP client)
```

---

## Use Cases

1. **Market Research**: Analyze markets, competitors, and trends
2. **Academic Research**: Conduct literature reviews and analysis
3. **Technology Assessment**: Evaluate new technologies and tools
4. **Business Intelligence**: Gather business insights and data
5. **Content Creation**: Research topics for articles and reports
6. **Due Diligence**: Research companies, products, and investments

---

## Performance Tips

### Efficient Research
- Use specific and focused research questions
- Set appropriate depth levels
- Leverage existing research when possible
- Use source analysis to verify quality

### Cost Management
- Monitor API usage and costs
- Use shallower research for preliminary analysis
- Cache research results when appropriate
- Batch similar research topics

---

## Integration with Workflows

### Research Pipeline
```python
# Example research pipeline
research_question = "Impact of AI on software development"
initial_research = research_topic(research_question, depth="basic")
if initial_research.quality > 0.8:
    detailed_research = research_topic(research_question, depth="comprehensive")
    report = generate_report(detailed_research)
```

### Automated Research
```yaml
# GitHub Actions for automated research
- name: Conduct Weekly Research
  run: |
    research_output=$(echo '{"tool": "research_trends", "arguments": {"field": "AI", "timeframe": "last_week"}}' | \
    python gptr-mcp/server.py)
    echo "$research_output" > research_report.md
```

---

## Related Documentation

- [GPTR MCP Server GitHub](https://github.com/assafelovic/gptr-mcp)
- [Custom MCP Server Implementation](gptr-mcp/)
- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [Tavily API Documentation](https://docs.tavily.com/)

---

**Last Updated**: 2025-11-18
**Research Status**: ✅ Complete (Phase 1A)
**Next Steps**: Register with jarvis and test functionality
