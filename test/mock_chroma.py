#!/usr/bin/env python3
"""
模拟 Chroma 向量数据库服务 - 用于扫描器测试
监听端口 8000
"""
import http.server
import socketserver
import json
import sys

PORT = 8000
VERSION = "0.4.15"

# 支持命令行版本查询
if len(sys.argv) > 1 and sys.argv[1] in ["--version", "-v", "-version"]:
    print(f"{VERSION}")
    sys.exit(0)

class ChromaHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/api/v1/heartbeat":
            self.send_response(200)
            self.send_header("Content-type", "application/json")
            self.end_headers()
            response = {"status": "ok", "version": VERSION}
            self.wfile.write(json.dumps(response).encode())
        elif self.path == "/":
            self.send_response(200)
            self.send_header("Content-type", "application/json")
            self.end_headers()
            response = {"version": VERSION}
            self.wfile.write(json.dumps(response).encode())
        else:
            self.send_response(200)
            self.send_header("Content-type", "text/plain")
            self.end_headers()
            self.wfile.write(f"ChromaDB Mock Server v{VERSION}".encode())
    
    def log_message(self, format, *args):
        pass

print(f"Starting Chroma mock server on port {PORT}...")
print(f"Version: {VERSION}")

with socketserver.TCPServer(("", PORT), ChromaHandler) as httpd:
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nShutting down...")
