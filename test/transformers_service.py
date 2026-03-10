#!/usr/bin/env python3
"""启动 Transformers 服务用于测试扫描器"""
import sys
import time

try:
    from transformers import AutoTokenizer, AutoModel
    print("Transformers loaded")
    
    # 加载一个小模型
    tokenizer = AutoTokenizer.from_pretrained("bert-base-uncased")
    model = AutoModel.from_pretrained("bert-base-uncased")
    print("BERT model loaded, running...")
    
    # 保持进程运行
    while True:
        time.sleep(10)
        
except ImportError as e:
    print(f"Import error: {e}")
    time.sleep(60)
    sys.exit(1)
except Exception as e:
    print(f"Error: {e}")
    time.sleep(60)
    sys.exit(1)
