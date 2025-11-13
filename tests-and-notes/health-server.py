#!/usr/bin/env python3
"""
MCP Health Monitoring Server
Provides HTTP health endpoint for cipher-aggregator monitoring
"""

import json
import os
import socket
import subprocess
import sys
import time
from datetime import datetime, timezone
from http.server import BaseHTTPRequestHandler, HTTPServer
from pathlib import Path
import signal

class HealthCheckHandler(BaseHTTPRequestHandler):
    """Handler for health check HTTP requests"""

    def __init__(self, *args, **kwargs):
        self.config = kwargs.pop('config', {})
        super().__init__(*args, **kwargs)

    def log_message(self, format, *args):
        """Override to reduce logging noise"""
        if os.getenv('MCP_HEALTH_DEBUG', 'false').lower() == 'true':
            sys.stderr.write(f"[{datetime.now().isoformat()}] {format % args}\n")

    def do_GET(self):
        """Handle GET requests to health endpoint"""
        if self.path == '/health':
            self.send_health_response()
        elif self.path == '/ready':
            self.send_readiness_response()
        elif self.path == '/':
            self.send_status_page()
        else:
            self.send_error(404, "Not Found")

    def send_health_response(self):
        """Send comprehensive health check response"""
        try:
            health_data = self.get_health_data()
            response = json.dumps(health_data, indent=2)

            self.send_response(200)
            self.send_header('Content-Type', 'application/json')
            self.send_header('Cache-Control', 'no-cache')
            self.send_header('Access-Control-Allow-Origin', '*')
            self.end_headers()
            self.wfile.write(response.encode())

        except Exception as e:
            self.send_error(500, f"Health check failed: {str(e)}")

    def send_readiness_response(self):
        """Send simple readiness check"""
        try:
            cipher_running = self.is_cipher_running()
            sse_responsive = self.test_sse_connectivity() if cipher_running else False

            if cipher_running and sse_responsive:
                self.send_response(200)
                status = "ready"
            else:
                self.send_response(503)
                status = "not ready"

            self.send_header('Content-Type', 'application/json')
            self.send_header('Cache-Control', 'no-cache')
            self.end_headers()

            response = {
                "status": status,
                "timestamp": datetime.now(timezone.utc).isoformat(),
                "ready": status == "ready"
            }
            self.wfile.write(json.dumps(response, indent=2).encode())

        except Exception as e:
            self.send_error(500, f"Readiness check failed: {str(e)}")

    def send_status_page(self):
        """Send simple HTML status page"""
        try:
            health_data = self.get_health_data()

            html = f"""
            <!DOCTYPE html>
            <html>
            <head>
                <title>MCP Health Monitor</title>
                <style>
                    body {{ font-family: Arial, sans-serif; margin: 40px; }}
                    .status {{ padding: 10px; border-radius: 5px; margin: 10px 0; }}
                    .healthy {{ background-color: #d4edda; color: #155724; }}
                    .degraded {{ background-color: #fff3cd; color: #856404; }}
                    .unhealthy {{ background-color: #f8d7da; color: #721c24; }}
                    .check {{ margin: 5px 0; }}
                </style>
                <meta http-equiv="refresh" content="30">
            </head>
            <body>
                <h1>MCP Health Monitor</h1>
                <div class="status {health_data['status']}">
                    <strong>Overall Status: {health_data['status'].upper()}</strong><br>
                    <small>Last updated: {health_data['timestamp']}</small>
                </div>

                <h2>Health Checks</h2>
                {self.format_checks_html(health_data['checks'])}

                <h2>System Info</h2>
                <p><strong>Uptime:</strong> {self.format_uptime(health_data.get('uptime_seconds', 0))}</p>
                <p><strong>Version:</strong> {health_data.get('version', 'unknown')}</p>
            </body>
            </html>
            """

            self.send_response(200)
            self.send_header('Content-Type', 'text/html')
            self.send_header('Cache-Control', 'no-cache')
            self.end_headers()
            self.wfile.write(html.encode())

        except Exception as e:
            self.send_error(500, f"Status page failed: {str(e)}")

    def format_checks_html(self, checks):
        """Format health checks as HTML"""
        html = ""
        for check_name, check_data in checks.items():
            status = check_data.get('status', 'unknown')
            icon = "✅" if status == "pass" else "❌" if status == "fail" else "⚠️"
            html += f'<div class="check">{icon} <strong>{check_name}:</strong> {status}'
            if 'message' in check_data:
                html += f' - {check_data["message"]}'
            html += '</div>'
        return html

    def format_uptime(self, seconds):
        """Format uptime in human readable format"""
        if seconds < 60:
            return f"{seconds}s"
        elif seconds < 3600:
            return f"{seconds//60}m {seconds%60}s"
        else:
            hours = seconds // 3600
            minutes = (seconds % 3600) // 60
            return f"{hours}h {minutes}m"

    def get_health_data(self):
        """Get comprehensive health data"""
        start_time = time.time()

        # Check if cipher is running
        cipher_running = self.is_cipher_running()
        cipher_pid = self.get_cipher_pid()

        # Test SSE connectivity
        sse_responsive = self.test_sse_connectivity() if cipher_running else False

        # Check port availability
        port_status = self.check_sse_port()

        # Check for conflicts
        conflicts = self.detect_conflicts()

        # Calculate uptime
        uptime_seconds = self.get_uptime_seconds()

        # Determine overall status
        failed_checks = [
            not cipher_running,
            not sse_responsive and cipher_running,
            port_status.get('status') != 'available' if cipher_running else False,
            conflicts.get('count', 0) > 0
        ]

        failed_count = sum(failed_checks)
        if failed_count == 0:
            status = "healthy"
        elif failed_count <= 2:
            status = "degraded"
        else:
            status = "unhealthy"

        # Build checks data
        checks = {
            "cipher_running": {
                "status": "pass" if cipher_running else "fail",
                "message": f"PID {cipher_pid}" if cipher_running else "Not running"
            },
            "sse_responsive": {
                "status": "pass" if sse_responsive else "fail",
                "message": "Responding" if sse_responsive else "Not responding"
            },
            "port_available": {
                "status": port_status.get('status', 'unknown'),
                "message": port_status.get('message', 'Unknown status')
            },
            "conflicts": {
                "status": "pass" if conflicts.get('count', 0) == 0 else "fail",
                "message": f"{conflicts.get('count', 0)} conflicts found"
            }
        }

        return {
            "status": status,
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "checks": checks,
            "uptime_seconds": uptime_seconds,
            "version": "2.0",
            "response_time_ms": round((time.time() - start_time) * 1000, 2)
        }

    def is_cipher_running(self):
        """Check if cipher-aggregator process is running"""
        pid_file = self.config.get('pid_file', '../MCP/cipher-aggregator.pid')

        if not os.path.exists(pid_file):
            return False

        try:
            with open(pid_file, 'r') as f:
                pid = int(f.read().strip())

            # Check if process exists
            try:
                os.kill(pid, 0)  # Signal 0 just checks if process exists
                return True
            except OSError:
                return False

        except (ValueError, FileNotFoundError, PermissionError):
            return False

    def get_cipher_pid(self):
        """Get cipher-aggregator PID"""
        pid_file = self.config.get('pid_file', '../MCP/cipher-aggregator.pid')

        try:
            with open(pid_file, 'r') as f:
                return int(f.read().strip())
        except:
            return None

    def test_sse_connectivity(self, timeout=3):
        """Test if SSE server is responding"""
        sse_host = self.config.get('sse_host', '127.0.0.1')
        sse_port = self.config.get('sse_port', 3020)

        try:
            # Try to connect to SSE endpoint
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.settimeout(timeout)
            sock.connect((sse_host, sse_port))
            sock.close()
            return True
        except (socket.timeout, socket.error, ConnectionRefusedError):
            return False

    def check_sse_port(self):
        """Check if SSE port status"""
        sse_port = self.config.get('sse_port', 3020)

        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.settimeout(1)
            result = sock.connect_ex(('127.0.0.1', sse_port))
            sock.close()

            if result == 0:
                return {"status": "in_use", "message": f"Port {sse_port} in use"}
            else:
                return {"status": "available", "message": f"Port {sse_port} available"}

        except Exception as e:
            return {"status": "unknown", "message": f"Error checking port: {str(e)}"}

    def detect_conflicts(self):
        """Detect conflicting MCP processes"""
        # This is a simplified version - could be enhanced to match mcp-manager.sh logic
        conflicts = []

        try:
            # Look for MCP processes
            result = subprocess.run(
                ['ps', 'aux'],
                capture_output=True,
                text=True,
                timeout=5
            )

            if result.returncode == 0:
                for line in result.stdout.split('\n'):
                    if any(pattern in line.lower() for pattern in ['mcp', 'cipher', 'firecrawl']):
                        if 'health-server.py' not in line:  # Exclude ourselves
                            conflicts.append(line.strip())

            return {
                "count": len(conflicts),
                "processes": conflicts
            }

        except Exception as e:
            return {"count": 0, "error": str(e)}

    def get_uptime_seconds(self):
        """Get health server uptime in seconds"""
        pid = os.getpid()

        try:
            result = subprocess.run(
                ['ps', '-o', 'etime', '-p', str(pid)],
                capture_output=True,
                text=True,
                timeout=5
            )

            if result.returncode == 0:
                # Parse elapsed time format (can be various formats)
                lines = result.stdout.strip().split('\n')
                if len(lines) >= 2:
                    elapsed = lines[1].strip()
                    return self.parse_elapsed_time(elapsed)

        except:
            pass

        return 0

    def parse_elapsed_time(self, elapsed_str):
        """Parse elapsed time string from ps command"""
        # Handle formats like "1-02:03:04", "02:03:04", "3:04", etc.
        try:
            parts = elapsed_str.split(':')

            if '-' in parts[0]:  # Format: days-HH:MM:SS
                days_part, time_part = parts[0].split('-')
                days = int(days_part)
                hours, minutes, seconds = map(int, time_part.split(':'))
                return days * 86400 + hours * 3600 + minutes * 60 + seconds
            elif len(parts) == 3:  # Format: HH:MM:SS
                hours, minutes, seconds = map(int, parts)
                return hours * 3600 + minutes * 60 + seconds
            elif len(parts) == 2:  # Format: MM:SS
                minutes, seconds = map(int, parts)
                return minutes * 60 + seconds
            else:  # Format: SS
                return int(parts[0])

        except (ValueError, IndexError):
            return 0

def load_config():
    """Load configuration from environment and files"""
    config = {
        'port': int(os.getenv('MCP_HEALTH_PORT', '3021')),
        'sse_host': os.getenv('SSE_HOST', '127.0.0.1'),
        'sse_port': int(os.getenv('SSE_PORT', '3020')),
        'pid_file': os.getenv('MCP_PID_FILE', '../MCP/cipher-aggregator.pid'),
        'log_file': os.getenv('MCP_HEALTH_LOG', '../MCP/logs/health-server.log')
    }

    # Try to load from health-config.env
    env_file = os.getenv('MCP_HEALTH_CONFIG', '../MCP/health-config.env')
    if os.path.exists(env_file):
        with open(env_file, 'r') as f:
            for line in f:
                line = line.strip()
                if line and not line.startswith('#') and '=' in line:
                    key, value = line.split('=', 1)
                    os.environ[key.strip()] = value.strip().strip('"')

    return config

def setup_logging(config):
    """Setup logging configuration"""
    import logging
    from logging.handlers import RotatingFileHandler

    # Create logs directory if it doesn't exist
    log_dir = os.path.dirname(config['log_file'])
    os.makedirs(log_dir, exist_ok=True)

    # Setup file handler with rotation
    file_handler = RotatingFileHandler(
        config['log_file'],
        maxBytes=10*1024*1024,  # 10MB
        backupCount=5
    )

    # Setup console handler
    console_handler = logging.StreamHandler()

    # Create formatter
    formatter = logging.Formatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )

    file_handler.setFormatter(formatter)
    console_handler.setFormatter(formatter)

    # Configure root logger
    root_logger = logging.getLogger()
    root_logger.setLevel(logging.INFO)
    root_logger.addHandler(file_handler)
    root_logger.addHandler(console_handler)

    return logging.getLogger('health-server')

def create_handler_class(config):
    """Create a handler class with config bound"""
    class HandlerWithConfig(HealthCheckHandler):
        def __init__(self, *args, **kwargs):
            super().__init__(*args, config=config, **kwargs)

    return HandlerWithConfig

def signal_handler(signum, frame):
    """Handle shutdown signals gracefully"""
    sys.stderr.write(f"\nReceived signal {signum}, shutting down...\n")
    sys.exit(0)

def main():
    """Main entry point"""
    # Setup signal handlers for graceful shutdown
    signal.signal(signal.SIGTERM, signal_handler)
    signal.signal(signal.SIGINT, signal_handler)

    # Load configuration
    config = load_config()
    logger = setup_logging(config)

    logger.info(f"Starting MCP Health Server on port {config['port']}")
    logger.info(f"SSE Target: {config['sse_host']}:{config['sse_port']}")
    logger.info(f"PID File: {config['pid_file']}")

    # Create handler class with config
    handler_class = create_handler_class(config)

    try:
        # Create and start server
        server = HTTPServer(('0.0.0.0', config['port']), handler_class)

        logger.info(f"Health server listening on http://0.0.0.0:{config['port']}")
        logger.info(f"Health endpoint: http://localhost:{config['port']}/health")
        logger.info(f"Status page: http://localhost:{config['port']}/")

        server.serve_forever()

    except KeyboardInterrupt:
        logger.info("Shutdown requested by user")
    except Exception as e:
        logger.error(f"Server error: {e}")
    finally:
        logger.info("Health server stopped")

if __name__ == '__main__':
    main()
