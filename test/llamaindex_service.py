#!/usr/bin/env python3
"""启动 LlamaIndex 服务用于测试扫描器 - 使用本地模型"""
import sys
import time

try:
    from llama_index.core import Document, VectorStoreIndex, Settings
    from llama_index.embeddings.huggingface import HuggingFaceEmbedding
    
    print("LlamaIndex loaded")
    
    # 使用本地嵌入模型
    Settings.embed_model = HuggingFaceEmbedding(model_name="BAAI/bge-small-en-v1.5")
    
    # 创建简单的索引
    documents = [Document(text="Hello world")]
    index = VectorStoreIndex.from_documents(documents)
    print("Index created with local embeddings, running...")
    
    # 保持进程运行
    while True:
        time.sleep(10)
        
except ImportError as e:
    print(f"Import error: {e}")
    print("Installing required packages...")
    time.sleep(60)
    sys.exit(1)
except Exception as e:
    print(f"Error: {e}")
    time.sleep(60)
    sys.exit(1)
