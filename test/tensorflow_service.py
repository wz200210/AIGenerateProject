#!/usr/bin/env python3
"""启动 TensorFlow 模型服务用于测试扫描器"""
import sys
import time

try:
    import tensorflow as tf
    print(f"TensorFlow {tf.__version__} loaded")
    
    # 创建一个简单的模型
    model = tf.keras.Sequential([
        tf.keras.layers.Dense(1, input_shape=(10,))
    ])
    print("Model created, running...")
    
    # 保持进程运行
    while True:
        time.sleep(10)
        
except ImportError:
    print("TensorFlow not installed yet, waiting...")
    time.sleep(60)
    sys.exit(1)
