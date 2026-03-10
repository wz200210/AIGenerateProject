#!/usr/bin/env python3
"""启动 LlamaIndex 服务用于测试扫描器"""
import sys
import time

try:
    from llama_index.core import Document, VectorStoreIndex
    print("LlamaIndex loaded")
    
    # 创建简单的索引
    documents = [Document(text="Hello world")]
    index = VectorStoreIndex.from_documents(documents)
    print("Index created, running...")
    
    # 保持进程运行
    while True:
        time.sleep(10)
        
except ImportError:
    print("LlamaIndex not installed yet, waiting...")
    time.sleep(60)
    sys.exit(1)
