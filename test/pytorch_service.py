#!/usr/bin/env python3
"""启动 PyTorch 模型服务用于测试扫描器"""
import sys
import time

try:
    import torch
    print(f"PyTorch {torch.__version__} loaded")
    
    # 创建一个简单的模型
    model = torch.nn.Linear(10, 1)
    print("Model created, running inference loop...")
    
    # 保持进程运行
    while True:
        x = torch.randn(1, 10)
        y = model(x)
        time.sleep(1)
        
except ImportError:
    print("PyTorch not installed yet, waiting...")
    time.sleep(60)
    sys.exit(1)
