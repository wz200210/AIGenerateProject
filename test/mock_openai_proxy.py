#!/usr/bin/env python3
"""
模拟 OpenAI API 代理服务 - 用于扫描器测试
监听端口 8080
"""
import http.server
import socketserver
import json
import os

PORT = 8080
os.environ["OPENAI_API_KEY"] = "sk-mock-key-for-testing"

class OpenAIHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/v1/models":
            self.send_response(200)
            self.send_header("Content-type", "application/json")
            self.end_headers()
            response = {
                "object": "list",
                "data": [{"id": "gpt-4", "object": "model"}]
            }
            self.wfile.write(json.dumps(response).encode())
        else:
            self.send_response(200)
            self.send_header("Content-type", "text/plain")
            self.end_headers()
            self.wfile.write(b"OpenAI API Proxy v1.0.0")
    
    def log_message(self, format, *args):
        pass

print(f"Starting OpenAI API proxy on port {PORT}...")
print("Version: 1.0.0")

with socketserver.TCPServer(("", PORT), OpenAIHandler) as httpd:
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nShutting down...")
