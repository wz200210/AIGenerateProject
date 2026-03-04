#!/usr/bin/env python3
"""
AI 组件模拟服务启动脚本
用于测试扫描器的检测能力
"""

import http.server
import socketserver
import json
import sys
import os
from threading import Thread

class QdrantHandler(http.server.SimpleHTTPRequestHandler):
    """Qdrant 向量数据库模拟"""
    def do_GET(self):
        if self.path == '/':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "title": "Qdrant - Vector Database",
                "version": "1.7.4",
                "service": "qdrant"
            }
            self.wfile.write(json.dumps(response).encode())
        else:
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b'OK')
    
    def log_message(self, format, *args):
        pass

class MilvusHandler(http.server.SimpleHTTPRequestHandler):
    """Milvus 向量数据库模拟"""
    def do_GET(self):
        if self.path == '/v1/health':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "status": "healthy",
                "version": "2.3.0"
            }
            self.wfile.write(json.dumps(response).encode())
        else:
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b'Milvus Mock')
    
    def log_message(self, format, *args):
        pass

class JupyterHandler(http.server.SimpleHTTPRequestHandler):
    """Jupyter Notebook 模拟"""
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(b'<html><title>Jupyter Notebook</title><body>Jupyter Server v7.0.0</body></html>')
    
    def log_message(self, format, *args):
        pass

class MLflowHandler(http.server.SimpleHTTPRequestHandler):
    """MLflow 模拟"""
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        response = {
            "mlflow_version": "2.8.0",
            "status": "running"
        }
        self.wfile.write(json.dumps(response).encode())
    
    def log_message(self, format, *args):
        pass

def start_server(handler_class, port):
    """启动单个服务"""
    with socketserver.TCPServer(('0.0.0.0', port), handler_class) as httpd:
        print(f"Started {handler_class.__name__} on port {port}")
        httpd.serve_forever()

def main():
    # 启动多个 AI 组件模拟服务
    services = [
        (QdrantHandler, 6333, "Qdrant"),
        (MilvusHandler, 19530, "Milvus"),
        (JupyterHandler, 8888, "Jupyter"),
        (MLflowHandler, 5000, "MLflow"),
    ]
    
    print("🚀 Starting AI Component Mock Services...")
    print("=" * 50)
    
    threads = []
    for handler, port, name in services:
        t = Thread(target=start_server, args=(handler, port))
        t.daemon = True
        t.start()
        threads.append((t, name, port))
    
    print("\n✅ All services started:")
    for _, name, port in threads:
        print(f"  • {name}: http://localhost:{port}")
    
    print("\nPress Ctrl+C to stop all services")
    
    try:
        while True:
            import time
            time.sleep(1)
    except KeyboardInterrupt:
        print("\n\n🛑 Stopping all services...")
        sys.exit(0)

if __name__ == '__main__':
    main()