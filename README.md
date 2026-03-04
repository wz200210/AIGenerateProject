# AI Generate Project

一个用于扫描和识别项目中常见 AI 组件的 Go 工具。

## 功能

- 扫描代码库中的 AI 相关依赖和引用
- 识别常见的 AI/ML 框架和库
- 检测 AI 模型文件
- 生成扫描报告

## 安装

```bash
go install github.com/wz200210/AIGenerateProject/cmd/scanner@latest
```

## 使用

```bash
# 扫描当前目录
scanner scan

# 扫描指定目录
scanner scan -p /path/to/project

# 输出 JSON 格式报告
scanner scan -o json
```

## 支持的 AI 组件

### LLM 框架
- OpenAI GPT
- Anthropic Claude
- Google Gemini
- LangChain
- LlamaIndex

### ML 框架
- TensorFlow
- PyTorch
- Hugging Face Transformers
- ONNX Runtime

### 模型文件
- GGUF (Llama, etc.)
- PyTorch (.pt, .pth)
- TensorFlow (.pb, .h5)
- ONNX (.onnx)

## 项目结构

```
.
├── cmd/scanner/        # 命令行入口
├── internal/
│   ├── detector/       # AI 组件检测逻辑
│   └── scanner/        # 文件扫描引擎
├── pkg/ai/types/       # 类型定义
└── README.md
```