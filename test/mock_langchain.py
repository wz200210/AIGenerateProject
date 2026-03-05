#!/usr/bin/env python3
"""
模拟 LangChain 服务进程 - 用于扫描器测试
"""
import time
import sys

# 模拟 LangChain 进程
print("Starting LangChain mock service...")
print("Version: 0.2.0")

# 保持进程运行
try:
    while True:
        time.sleep(60)
except KeyboardInterrupt:
    print("Shutting down...")
    sys.exit(0)
