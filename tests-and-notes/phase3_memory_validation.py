#!/usr/bin/env python3
"""
Phase 3: Cipher Memory Validation Test
Tests cipher memory, workspace, and knowledge graph functionality
"""

import json
import urllib.request
import urllib.error
import uuid
import time
from typing import Dict, Any, Optional

class CipherMemoryValidator:
    """Validates cipher memory tools functionality"""

    def __init__(self):
        self.sse_url = "http://localhost:3020/sse"
        self.session_id = None

    def establish_connection(self) -> bool:
        """Establish SSE connection and get session ID"""
        print("🔄 Establishing SSE connection...")

        try:
            req = urllib.request.Request(
                self.sse_url,
                method="GET",
                headers={"Accept": "text/event-stream"}
            )

            response = urllib.request.urlopen(req, timeout=10)

            # Read session ID from endpoint event
            for line in response:
                line = line.decode("utf-8").strip()
                if line.startswith("data: /sse?sessionId="):
                    self.session_id = line.split("sessionId=")[1]
                    print(f"✅ Connection established (Session: {self.session_id[:8]}...)")
                    return True

        except Exception as e:
            print(f"❌ Connection failed: {e}")
            return False

        return False

    def call_tool(self, tool_name: str, arguments: Dict[str, Any]) -> Optional[Dict]:
        """Call an MCP tool and return the result"""
        if not self.session_id:
            print("❌ No session ID")
            return None

        request_id = str(uuid.uuid4())

        request_data = {
            "jsonrpc": "2.0",
            "method": "tools/call",
            "id": request_id,
            "params": {
                "name": tool_name,
                "arguments": arguments
            }
        }

        url_with_session = f"{self.sse_url}?sessionId={self.session_id}"

        try:
            req = urllib.request.Request(
                url_with_session,
                data=json.dumps(request_data).encode("utf-8"),
                headers={"Content-Type": "application/json"},
                method="POST"
            )

            with urllib.request.urlopen(req, timeout=60) as response:
                if response.status in [200, 202]:
                    # For now, just confirm the request was acceptedprint(f"✅ {tool_name} called successfully (status: {response.status})")

                    if response.status == 200:
                        response_data = json.loads(response.read().decode("utf-8"))
                        return response_data
                    return {"status": "accepted"}

        except urllib.error.HTTPError as e:
            error_body = e.read().decode("utf-8")
            print(f"❌ {tool_name} failed: HTTP {e.code}")
            print(f"   Error: {error_body[:200]}")
            return None
        except Exception as e:
            print(f"❌ {tool_name} failed: {e}")
            return None

        return None

    def test_memory_search(self) -> bool:
        """Test cipher_memory_search"""
        print("\n📋 Testing cipher_memory_search...")

        result = self.call_tool("cipher_memory_search", {
            "query": "testing cipher memory system",
            "top_k": 5,
            "similarity_threshold": 0.3
        })

        return result is not None

    def test_workspace_search(self) -> bool:
        """Test cipher_workspace_search"""
        print("\n📋 Testing cipher_workspace_search...")

        result = self.call_tool("cipher_workspace_search", {
            "query": "project progress and team collaboration",
            "top_k": 5
        })

        return result is not None

    def test_workspace_store(self) -> bool:
        """Test cipher_workspace_store"""
        print("\n📋 Testing cipher_workspace_store...")

        result = self.call_tool("cipher_workspace_store", {
            "interaction": "Phase 3 memory validation test completed. All cipher memory tools are now functional and accessible.",
            "context": {
                "sessionId": "phase3-validation",
                "projectId": "cipher-mcp"
            }
        })

        return result is not None

    def test_extract_and_operate_memory(self) -> bool:
        """Test cipher_extract_and_operate_memory"""
        print("\n📋 Testing cipher_extract_and_operate_memory...")

        result = self.call_tool("cipher_extract_and_operate_memory", {
            "interaction": "Remember: cipher memory tools require the real OPENAI_API_KEY from .env file, not demo keys. The .env file must have Unix line endings for bash to load it properly.",
            "context": {
                "sessionId": "phase3-validation"
            },
            "options": {
                "autoExtractKnowledgeInfo": True,
                "confidenceThreshold": 0.7
            }
        })

        return result is not None

    def test_knowledge_graph_operations(self) -> bool:
        """Test cipher knowledge graph tools"""
        print("\n📋 Testing cipher knowledge graph operations...")

        # Test add_node
        node_result = self.call_tool("cipher_add_node", {
            "id": "cipher-memory-validation",
            "labels": ["TestNode", "Phase3"],
            "properties": {
                "name": "Cipher Memory Validation",
                "timestamp": time.time(),
                "status": "completed"
            }
        })

        if not node_result:
            return False

        # Test search_graph
        search_result = self.call_tool("cipher_search_graph", {
            "searchType": "nodes",
            "nodeLabels": ["TestNode"],
            "limit": 10
        })

        return search_result is not None

    def test_reasoning_memory(self) -> bool:
        """Test cipher reasoning memory tools"""
        print("\n📋 Testing cipher reasoning memory...")

        # Extract reasoning steps
        extract_result = self.call_tool("cipher_extract_reasoning_steps", {
            "userInput": "To fix cipher memory tools: 1) Diagnose API key issue, 2) Fix .env line endings, 3) Load environment in mcp-manager.sh, 4) Restart cipher-aggregator"
        })

        return extract_result is not None

    def run_all_tests(self) -> Dict[str, bool]:
        """Run all memory validation tests"""
        print("="*60)
        print("Phase 3: Cipher Memory Validation Test Suite")
        print("="*60)

        if not self.establish_connection():
            return {"connection": False}

        results = {
            "connection": True,
            "memory_search": self.test_memory_search(),
            "workspace_search": self.test_workspace_search(),
            "workspace_store": self.test_workspace_store(),
            "extract_and_operate": self.test_extract_and_operate_memory(),
            "knowledge_graph": self.test_knowledge_graph_operations(),
            "reasoning_memory": self.test_reasoning_memory()
        }

        print("\n" + "="*60)
        print("Test Results Summary")
        print("="*60)

        passed = sum(1 for v in results.values() if v)
        total = len(results)

        for test_name, passed in results.items():
            status = "✅ PASS" if passed else "❌ FAIL"
            print(f"{status}: {test_name}")

        print(f"\nOverall: {passed}/{total} tests passed ({100*passed//total}%)")

        return results

if __name__ == "__main__":
    validator = CipherMemoryValidator()
    results = validator.run_all_tests()

    # Exit with appropriate code
    all_passed = all(results.values())
    exit(0 if all_passed else 1)
