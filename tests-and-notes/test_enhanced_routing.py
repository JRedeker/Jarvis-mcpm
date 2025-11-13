#!/usr/bin/env python3
"""
Test Enhanced Routing Logic
Verify that the improved tier selection works correctly
"""

import os
import sys

# Add servers directory to path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), "..", "servers"))

# Import the routing function directly from the file
import importlib.util

spec = importlib.util.spec_from_file_location(
    "llm_inference_mcp", "/home/jrede/dev/MCP/servers/llm-inference-mcp.py"
)
llm_module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(llm_module)
select_tier_from_task = llm_module.select_tier_from_task

# Test cases with expected tier
test_cases = [
    # l0 - The Router (quick operations)
    ("classify this data into categories", "l0"),
    ("check the status of the API", "l0"),
    ("validate this input", "l0"),
    ("quick search for user records", "l0"),
    # m1 - The Organizer (data transformation)
    ("format this JSON data", "m1"),
    ("organize these files by date", "m1"),
    ("summarize the research findings", "m1"),
    ("parse this CSV file", "m1"),
    # m2 - The Builder (code generation)
    ("write a Python function to sort a list", "m2"),
    ("generate a React component for user profile", "m2"),
    ("create an API endpoint for user authentication", "m2"),
    ("implement a database migration", "m2"),
    ("code a sorting algorithm", "m2"),
    # m3 - The Thinker (reasoning & analysis)
    ("analyze our system architecture and recommend improvements", "m3"),
    ("should we use PostgreSQL or MongoDB?", "m3"),
    ("optimize our database queries for better performance", "m3"),
    ("design a scalable microservices architecture", "m3"),
    ("evaluate the security risks in our deployment", "m3"),
    ("compare these two approaches and recommend the best", "m3"),
    # m4 - The Writer (large output)
    ("write comprehensive documentation for all API endpoints", "m4"),
    ("create a detailed migration guide covering all edge cases", "m4"),
    ("comprehensive review of our entire deployment strategy", "m4"),
    ("document the complete system architecture with examples", "m4"),
    ("write a detailed guide explaining the entire codebase", "m4"),
]

print("Testing Enhanced Routing Logic")
print("=" * 60)

passed = 0
failed = 0

for task, expected_tier in test_cases:
    result = select_tier_from_task(task)
    status = "✓" if result == expected_tier else "✗"

    if result == expected_tier:
        passed += 1
        print(f"{status} PASS: '{task[:50]}...'")
        print(f"   Expected: {expected_tier}, Got: {result}")
    else:
        failed += 1
        print(f"{status} FAIL: '{task[:50]}...'")
        print(f"   Expected: {expected_tier}, Got: {result}")
    print()

print("=" * 60)
print(f"Results: {passed} passed, {failed} failed out of {len(test_cases)} tests")
print(f"Success rate: {(passed / len(test_cases) * 100):.1f}%")
