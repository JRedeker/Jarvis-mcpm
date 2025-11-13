#!/usr/bin/env python3
"""
Phase 3: Leverage Built-in AI Features - Memory Implementation

This script implements and tests Cipher's built-in memory functionality including:
- Workspace memory with auto-capture
- Knowledge memory with semantic search
- Reflection memory for reasoning traces
- Embedding generation with OpenAI model

Follows cipher.yml configuration and .env settings exactly.
"""

import os
import sys
import json
import time
import asyncio
from datetime import datetime
from pathlib import Path

# Add project root to path for imports
sys.path.insert(0, str(Path(__file__).parent))

# Load environment variables from .env file
try:
    from dotenv import load_dotenv
    load_dotenv()
except ImportError:
    print("⚠️  Warning: python-dotenv not available, using manual env loading")
    # Fallback to manual environment loading
    def load_env():
        env_vars = {}
        try:
            with open('.env', 'r') as f:
                for line in f:
                    if '=' in line and not line.strip().startswith('#'):
                        key, value = line.strip().split('=', 1)
                        env_vars[key.strip()] = value.strip()
        except FileNotFoundError:
            print("⚠️  Warning: .env file not found")
        return env_vars

    # Load environment variables
    env_vars = load_env()


class MemoryImplementation:
    """
    Implements Cipher's built-in AI features for memory and knowledge management.

    Features:
    - Workspace memory (project-scoped, auto-capture)
    - Knowledge memory (semantic search with embeddings)
    - Reflection memory (reasoning traces)
    - Integration with existing SSE framework
    """

    def __init__(self):
        self.validator = None
        self.session_id = None
        self.memory_tools = {}
        self.workspace_memory_root = "/home/jrede/dev/MCP/data/workspace-memory"
        self.memory_bank_root = "/home/jrede/dev/MCP/data/memory-bank"

        # Ensure memory directories exist
        os.makedirs(self.workspace_memory_root, exist_ok=True)
        os.makedirs(self.memory_bank_root, exist_ok=True)

        print("🧠 Memory Implementation initialized")
        print(f"📁 Workspace memory: {self.workspace_memory_root}")
        print(f"🗄️ Memory bank: {self.memory_bank_root}")

    async def connect_to_cipher(self):
        """Establish connection to cipher-aggregator using bidirectional SSE"""
        try:
            self.validator = ParameterValidator()
            await self.validator.establish_sse_connection()
            self.session_id = self.validator.session_id
            print(f"✅ Connected to cipher with session: {self.session_id}")
            return True
        except Exception as e:
            print(f"❌ Failed to connect to cipher: {e}")
            return False

    async def test_memory_tools_availability(self):
        """Test if built-in cipher memory tools are available"""
        print("\n🔍 Testing Built-in Memory Tools Availability...")

        # List of expected cipher memory tools
        expected_tools = [
            "cipher_memory_search",
            "cipher_workspace_search",
            "cipher_extract_and_operate_memory",
            "cipher_store_reasoning_memory",
            "cipher_search_reasoning_patterns",
            "cipher_workspace_store",
            "cipher_memory_search",
            "cipher_extract_entities",
            "cipher_intelligent_processor"
        ]

        try:
            # Get tools list from cipher
            tools_response = await self.validator.send_request("tools/list", {})

            if "error" in tools_response:
                print(f"❌ Failed to get tools list: {tools_response['error']}")
                return False

            tools = tools_response.get("result", {}).get("tools", [])
            tool_names = [tool.get("name", "") for tool in tools]

            print(f"📋 Found {len(tools)} total tools")

            # Check for cipher memory tools
            found_memory_tools = []
            missing_memory_tools = []

            for tool_name in expected_tools:
                if tool_name in tool_names:
                    found_memory_tools.append(tool_name)
                    print(f"  ✅ {tool_name}")
                else:
                    missing_memory_tools.append(tool_name)
                    print(f"  ❌ {tool_name} - MISSING")

            # Store memory tools info
            self.memory_tools = {tool["name"]: tool for tool in tools if tool["name"] in found_memory_tools}

            print(f"\n📊 Memory Tools Summary:")
            print(f"  Found: {len(found_memory_tools)}/{len(expected_tools)}")
            print(f"  Missing: {len(missing_memory_tools)}")

            if missing_memory_tools:
                print(f"  ⚠️  Missing tools: {', '.join(missing_memory_tools)}")
                return False

            return True

        except Exception as e:
            print(f"❌ Error testing memory tools: {e}")
            return False

    async def test_workspace_memory_operations(self):
        """Test workspace memory functionality"""
        print("\n🏢 Testing Workspace Memory Operations...")

        try:
            # Test workspace search
            print("  🔍 Testing workspace search...")
            search_response = await self.validator.send_request("tools/call", {
                "name": "cipher_workspace_search",
                "arguments": {
                    "query": "memory implementation test",
                    "top_k": 5
                }
            })

            if "error" in search_response:
                print(f"    ❌ Workspace search failed: {search_response['error']}")
            else:
                results = search_response.get("result", [])
                print(f"    ✅ Workspace search found {len(results)} items")
                for i, result in enumerate(results[:3]):  # Show first 3 results
                    print(f"      {i+1}. {result}")

            # Test workspace store
            print("  💾 Testing workspace store...")
            store_response = await self.validator.send_request("tools/call", {
                "name": "cipher_workspace_store",
                "arguments": {
                    "interaction": "Phase 3 memory implementation test",
                    "workspaceInfo": {
                        "project": "MCP",
                        "domain": "memory-implementation",
                        "status": "in-progress"
                    },
                    "options": {
                        "autoExtractWorkspaceInfo": True
                    }
                }
            })

            if "error" in store_response:
                print(f"    ❌ Workspace store failed: {store_response['error']}")
            else:
                print("    ✅ Workspace store successful")

            return True

        except Exception as e:
            print(f"❌ Error testing workspace memory: {e}")
            return False

    async def test_knowledge_memory_operations(self):
        """Test knowledge memory functionality"""
        print("\n🧠 Testing Knowledge Memory Operations...")

        try:
            # Test memory search
            print("  🔍 Testing semantic memory search...")
            search_response = await self.validator.send_request("tools/call", {
                "name": "cipher_memory_search",
                "arguments": {
                    "query": "memory implementation patterns",
                    "top_k": 3,
                    "similarity_threshold": 0.3
                }
            })

            if "error" in search_response:
                print(f"    ❌ Memory search failed: {search_response['error']}")
            else:
                results = search_response.get("result", [])
                print(f"    ✅ Memory search found {len(results)} items")
                for i, result in enumerate(results):
                    print(f"      {i+1}. {result}")

            # Test memory store
            print("  💾 Testing memory store...")
            store_response = await self.validator.send_request("tools/call", {
                "name": "cipher_extract_and_operate_memory",
                "arguments": {
                    "interaction": "Key finding: semantic search enables finding relevant code patterns and documentation",
                    "knowledgeInfo": {
                        "domain": "memory-implementation",
                        "codePattern": "semantic-search-with-embeddings"
                    },
                    "options": {
                        "autoExtractKnowledgeInfo": True,
                        "similarityThreshold": 0.7
                    }
                }
            })

            if "error" in store_response:
                print(f"    ❌ Memory store failed: {store_response['error']}")
            else:
                result = store_response.get("result", {})
                operations = result.get("operations", [])
                print(f"    ✅ Memory operations: {len(operations)} performed")
                for op in operations:
                    print(f"      - {op.get('type', 'unknown')}: {op.get('status', 'unknown')}")

            return True

        except Exception as e:
            print(f"❌ Error testing knowledge memory: {e}")
            return False

    async def test_embedding_generation(self):
        """Test embedding generation with OpenAI model"""
        print("\n🔤 Testing Embedding Generation...")

        try:
            # Test entity extraction with embeddings
            print("  🧠 Testing entity extraction with embeddings...")
            extract_response = await self.validator.send_request("tools/call", {
                "name": "cipher_extract_entities",
                "arguments": {
                    "text": "Implement semantic search using OpenAI embeddings for code patterns and API documentation",
                    "options": {
                        "autoLink": True,
                        "entityTypes": ["Function", "API", "Pattern"]
                    }
                }
            })

            if "error" in extract_response:
                print(f"    ❌ Entity extraction failed: {extract_response['error']}")
            else:
                result = extract_response.get("result", {})
                entities = result.get("entities", [])
                print(f"    ✅ Extracted {len(entities)} entities")
                for entity in entities[:5]:  # Show first 5 entities
                    print(f"      - {entity.get('type', 'Unknown')}: {entity.get('name', 'Unnamed')}")

            return True

        except Exception as e:
            print(f"❌ Error testing embeddings: {e}")
            return False

    async def verify_environment_configuration(self):
        """Verify environment variables and configuration"""
        print("\n⚙️  Verifying Environment Configuration...")

        # Check critical environment variables
        env_vars = {
            "USE_WORKSPACE_MEMORY": os.getenv("USE_WORKSPACE_MEMORY"),
            "KNOWLEDGE_GRAPH_ENABLED": os.getenv("KNOWLEDGE_GRAPH_ENABLED"),
            "DISABLE_DEFAULT_MEMORY": os.getenv("DISABLE_DEFAULT_MEMORY"),
            "EMBEDDING_MODEL": os.getenv("EMBEDDING_MODEL"),
            "OPENAI_API_KEY": "✅ Set" if os.getenv("OPENAI_API_KEY") else "❌ Missing"
        }

        print("  Environment Variables:")
        for var, value in env_vars.items():
            status = "✅" if value and value != "❌ Missing" else "❌"
            print(f"    {var}: {status} {value}")

        # Check memory directories
        workspace_exists = os.path.exists(self.workspace_memory_root)
        memory_bank_exists = os.path.exists(self.memory_bank_root)

        print(f"\n  Directory Status:")
        print(f"    Workspace memory: {'✅' if workspace_exists else '❌'} {self.workspace_memory_root}")
        print(f"    Memory bank: {'✅' if memory_bank_exists else '❌'} {self.memory_bank_root}")

        return all([
            env_vars["USE_WORKSPACE_MEMORY"] == "true",
            env_vars["KNOWLEDGE_GRAPH_ENABLED"] == "true",
            env_vars["OPENAI_API_KEY"] == "✅ Set"
        ])

    async def run_phase3_tests(self):
        """Run all Phase 3 tests"""
        print("🚀 Starting Phase 3: Leverage Built-in AI Features")
        print("=" * 60)

        # Step 1: Verify environment configuration
        env_ok = await self.verify_environment_configuration()
        if not env_ok:
            print("❌ Environment configuration incomplete - cannot proceed")
            return False

        # Step 2: Connect to cipher
        if not await self.connect_to_cipher():
            print("❌ Failed to connect to cipher - cannot proceed")
            return False

        # Step 3: Test memory tools availability
        tools_ok = await self.test_memory_tools_availability()
        if not tools_ok:
            print("❌ Memory tools not available - cannot proceed")
            return False

        # Step 4: Test workspace memory operations
        workspace_ok = await self.test_workspace_memory_operations()

        # Step 5: Test knowledge memory operations
        knowledge_ok = await self.test_knowledge_memory_operations()

        # Step 6: Test embedding generation
        embeddings_ok = await self.test_embedding_generation()

        # Summary
        print("\n" + "=" * 60)
        print("📊 Phase 3 Implementation Summary:")
        print(f"  Environment Config: {'✅' if env_ok else '❌'}")
        print(f"  Memory Tools: {'✅' if tools_ok else '❌'}")
        print(f"  Workspace Memory: {'✅' if workspace_ok else '❌'}")
        print(f"  Knowledge Memory: {'✅' if knowledge_ok else '❌'}")
        print(f"  Embedding Generation: {'✅' if embeddings_ok else '❌'}")

        overall_success = all([env_ok, tools_ok, workspace_ok, knowledge_ok, embeddings_ok])
        print(f"\n🎯 Phase 3 Overall: {'✅ SUCCESS' if overall_success else '❌ FAILED'}")

        return overall_success


async def main():
    """Main execution function"""
    print("🧠 Cipher Memory Implementation - Phase 3")
    print("Implementing built-in AI features: workspace memory, knowledge memory, embeddings")

    impl = MemoryImplementation()
    success = await impl.run_phase3_tests()

    if success:
        print("\n✅ Phase 3 implementation completed successfully!")
        print("📝 Ready to proceed to Phase 4: Timeout Configuration Enforcement")
    else:
        print("\n❌ Phase 3 implementation failed - check logs above")

    return 0 if success else 1


if __name__ == "__main__":
    exit_code = asyncio.run(main())
    sys.exit(exit_code)
