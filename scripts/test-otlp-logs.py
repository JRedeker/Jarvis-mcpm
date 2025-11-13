#!/usr/bin/env python3
"""
OTLP Logging Pipeline Test Script

Tests the end-to-end OTLP logging pipeline:
1. Initializes OTel logging
2. Emits test logs with various severities
3. Verifies logs reach Loki
4. Tests structured logging with attributes

Usage:
    python3 scripts/test-otlp-logs.py
"""

import asyncio
import logging
import sys
import os
import time
from datetime import datetime

# Add project root to path
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

import httpx
from servers.otel_logs import setup_otel_logging, emit_structured_log

# Test configuration
LOKI_BASE_URL = os.getenv("LOKI_BASE_URL", "http://localhost:3100")
TEST_SERVICE_NAME = "test-otlp-pipeline"

def print_banner(text):
    """Print a formatted banner"""
    print(f"\n{'='*70}")
    print(f"  {text}")
    print(f"{'='*70}\n")

def print_step(num, text):
    """Print a formatted step"""
    print(f"[Step {num}] {text}")

def print_result(success, message):
    """Print a formatted result"""
    status = "✓" if success else "✗"
    print(f"  {status} {message}")

async def test_otel_initialization():
    """Test 1: Initialize OTel logging"""
    print_step(1, "Testing OTel initialization")

    try:
        provider = setup_otel_logging(TEST_SERVICE_NAME, enable_console_bridge=True)

        if provider:
            print_result(True, "OTel logger provider initialized successfully")
            return provider
        else:
            print_result(False, "OTel disabled (OTEL_LOGS_EXPORTER=none)")
            return None
    except Exception as e:
        print_result(False, f"Failed to initialize OTel: {e}")
        return None

async def test_structured_logging(provider):
    """Test 2: Emit structured logs"""
    print_step(2, "Testing structured log emission")

    test_cases = [
        ("INFO", logging.INFO, {"test_type": "info_test", "value": 123}),
        ("WARN", logging.WARNING, {"test_type": "warn_test", "alert": True}),
        ("ERROR", logging.ERROR, {"test_type": "error_test", "is_error": True}),
    ]

    for name, level, attrs in test_cases:
        try:
            emit_structured_log(
                provider,
                TEST_SERVICE_NAME,
                f"test_event_{name.lower()}",
                level=level,
                **attrs
            )
            print_result(True, f"Emitted {name} level log with attributes")
        except Exception as e:
            print_result(False, f"Failed to emit {name} log: {e}")

    # Allow time for batching
    print(f"\n  Waiting 3s for batch export...")
    await asyncio.sleep(3)

async def test_cost_log_simulation(provider):
    """Test 3: Simulate cost logging"""
    print_step(3, "Testing cost log simulation (OpenRouter style)")

    try:
        emit_structured_log(
            provider,
            TEST_SERVICE_NAME,
            "openrouter_cost",
            level=logging.INFO,
            tier_id="m2",
            model="minimax/minimax-01",
            task_summary="Test LLM inference task",
            tokens_prompt=1500,
            tokens_completion=500,
            tokens_total=2000,
            cost_input=0.001575,
            cost_output=0.000525,
            cost_total=0.002100,
            pricing_input_per_1m=1.05,
            pricing_output_per_1m=1.05,
        )
        print_result(True, "Emitted cost tracking log with full attributes")
    except Exception as e:
        print_result(False, f"Failed to emit cost log: {e}")

    await asyncio.sleep(2)

async def test_http_instrumentation_simulation(provider):
    """Test 4: Simulate HTTP instrumentation"""
    print_step(4, "Testing HTTP request instrumentation simulation")

    try:
        emit_structured_log(
            provider,
            TEST_SERVICE_NAME,
            "openrouter_request",
            level=logging.INFO,
            tier_id="m3",
            model="openai/gpt-5",
            latency_ms=1250.5,
            http_status=200,
            request_id="test-req-12345",
            is_error=False,
        )
        print_result(True, "Emitted HTTP request instrumentation log")
    except Exception as e:
        print_result(False, f"Failed to emit HTTP log: {e}")

    await asyncio.sleep(2)

async def query_loki(query, limit=10):
    """Query Loki for test logs"""
    try:
        async with httpx.AsyncClient(timeout=10.0) as client:
            # Calculate time range (last 30 seconds)
            end_ns = int(time.time() * 1e9)
            start_ns = end_ns - (30 * int(1e9))

            params = {
                "query": query,
                "start": str(start_ns),
                "end": str(end_ns),
                "limit": str(limit),
                "direction": "backward"
            }

            url = f"{LOKI_BASE_URL}/loki/api/v1/query_range"
            response = await client.get(url, params=params)
            response.raise_for_status()

            data = response.json()
            streams = data.get("data", {}).get("result", [])

            total_entries = sum(len(stream.get("values", [])) for stream in streams)
            return total_entries, streams

    except Exception as e:
        return 0, []

async def test_loki_ingestion():
    """Test 5: Verify logs in Loki"""
    print_step(5, "Verifying logs reached Loki")

    # Wait for final batch export
    print("  Waiting 5s for final batch export...")
    await asyncio.sleep(5)

    # Query for our test logs
    query = f'{{service_name="{TEST_SERVICE_NAME}"}}'

    try:
        count, streams = await query_loki(query)

        if count > 0:
            print_result(True, f"Found {count} log entries in Loki")

            # Show sample logs
            print("\n  Sample log entries:")
            shown = 0
            for stream in streams[:2]:
                labels = stream.get("stream", {})
                for ts_ns, line in stream.get("values", [])[:3]:
                    if shown < 5:
                        print(f"    [{labels.get('level', 'INFO')}] {line[:100]}")
                        shown += 1
        else:
            print_result(False, "No logs found in Loki - check if stack is running")
            print(f"  Query: {query}")
            print(f"  Loki URL: {LOKI_BASE_URL}")

    except httpx.ConnectError:
        print_result(False, f"Cannot connect to Loki at {LOKI_BASE_URL}")
        print("  Make sure to start the stack: ./scripts/deploy-otlp-stack.sh start")
    except Exception as e:
        print_result(False, f"Failed to query Loki: {e}")

async def test_specific_events():
    """Test 6: Query specific event types"""
    print_step(6, "Querying specific event types")

    event_types = [
        "openrouter_cost",
        "openrouter_request",
        "test_event_info"
    ]

    for event in event_types:
        query = f'{{service_name="{TEST_SERVICE_NAME}"}} |= "{event}"'
        count, _ = await query_loki(query, limit=5)

        if count > 0:
            print_result(True, f"Found {count} '{event}' events")
        else:
            print_result(False, f"No '{event}' events found")

async def main():
    """Main test suite"""
    print_banner("OTLP Logging Pipeline End-to-End Test")

    print("Configuration:")
    print(f"  OTEL_EXPORTER_OTLP_ENDPOINT: {os.getenv('OTEL_EXPORTER_OTLP_ENDPOINT', 'http://localhost:4317')}")
    print(f"  LOKI_BASE_URL: {LOKI_BASE_URL}")
    print(f"  Test Service: {TEST_SERVICE_NAME}")

    # Run tests
    provider = await test_otel_initialization()

    if provider is None:
        print("\nOTel is disabled. Set OTEL_LOGS_EXPORTER=otlp to enable.")
        print("Exiting test suite.")
        return

    await test_structured_logging(provider)
    await test_cost_log_simulation(provider)
    await test_http_instrumentation_simulation(provider)
    await test_loki_ingestion()
    await test_specific_events()

    print_banner("Test Summary")

    print("""
Next steps:
  1. View logs in Grafana: http://localhost:3000
  2. Use LogQL queries in Explore tab:
     - {service_name="test-otlp-pipeline"}
     - {service_name="test-otlp-pipeline",event="openrouter_cost"}
     - {service_name="llm-inference-mcp",tier_id="m3"}

  3. Test with real servers:
     - Enable OTEL_* env vars in cipher.yml
     - Restart cipher-aggregator
     - Make some LLM inference calls
     - Query logs via logs-mcp server or Grafana

Run deployment:
  ./scripts/deploy-otlp-stack.sh start

View logs:
  ./scripts/deploy-otlp-stack.sh logs [otel-collector|loki|grafana]
""")

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("\n\nTest interrupted by user")
    except Exception as e:
        print(f"\n\nFatal error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)