# AI Component Scanner

一个用于扫描和识别项目中 AI 组件的 Go 工具，支持 LLM 框架、RAG 工具、向量数据库、Agent 框架等多种组件类型。

## 功能

- 🔍 扫描代码库中的 AI 相关依赖和引用
- 🤖 识别 **50+** 种 AI/ML 框架和库
- 💾 检测 AI 模型文件（.gguf, .onnx 等）
- 🔐 检测潜在硬编码的 API Key
- 📊 生成彩色控制台报告或 JSON 输出

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

### 🌐 国际 LLM 框架
- OpenAI (GPT-4, GPT-3.5)
- Anthropic Claude
- Google Gemini / Vertex AI
- Cohere
- Mistral AI

### 🇨🇳 国产大模型
- 百度文心一言 (ERNIE Bot)
- 阿里通义千问 (Qwen)
- 讯飞星火
- 智谱 GLM
- 月之暗面 Kimi
- MiniMax

### 🔗 LLM 编排框架
- LangChain / LangGraph
- LlamaIndex
- LiteLLM
- Vercel AI SDK
- Ollama

### 🤖 Agent 框架
- AutoGPT
- AutoGen (Microsoft)
- CrewAI

### 🧠 ML 框架
- Hugging Face Transformers
- TensorFlow / Keras
- PyTorch / Lightning
- ONNX Runtime

### 💾 向量数据库 (RAG)
- Pinecone
- Chroma
- Milvus
- Weaviate
- Qdrant
- FAISS
- pgvector
- Redis Vector
- Elasticsearch / OpenSearch
- Vespa

### 🔤 Embedding 模型
- OpenAI Embeddings
- Sentence Transformers
- BGE (中文 Embedding)

### 🚀 部署/推理
- vLLM
- Text Generation Inference (TGI)
- Triton Inference Server
- BentoML

### 📊 监控/可观测性
- LangSmith
- Langfuse
- Weights & Biases
- MLflow

### 📄 文档处理 (RAG)
- Unstructured
- PyPDF
- BeautifulSoup

### 📦 模型文件
- `.gguf` / `.ggml` - Llama 等本地模型
- `.pt` / `.pth` - PyTorch 模型
- `.onnx` - ONNX 模型
- `.safetensors` - HuggingFace SafeTensors
- `.pb` / `.h5` - TensorFlow 模型
- `.tflite` - TensorFlow Lite
- `.bin` - 通用二进制模型
- `.mlmodel` - Apple CoreML
- `.ckpt` - Checkpoint 文件

## 项目结构

```
.
├── cmd/scanner/           # 命令行入口
├── internal/
│   ├── detector/          # AI 组件检测逻辑
│   └── scanner/           # 文件扫描引擎
├── pkg/ai/types/          # 类型定义和组件配置
├── go.mod
├── Makefile
└── README.md
```

## 开发

```bash
# 克隆仓库
git clone git@github.com:wz200210/AIGenerateProject.git
cd AIGenerateProject

# 编译
go mod tidy
go build -o scanner ./cmd/scanner

# 运行测试
make test

# 交叉编译
make build-all
```

## 示例输出

```
╔════════════════════════════════════════╗
║     AI Component Scan Report           ║
╚════════════════════════════════════════╝

📁 Project: /path/to/your/project
📊 Files scanned: 1,234
🤖 AI components found: 8

[llm_framework]
  • OpenAI [medium]
    File: /path/to/project/main.py:15
    OpenAI GPT API

[vector_database]
  • Chroma [low]
    File: /path/to/project/requirements.txt:3
    Chroma 开源向量数据库

[api_key]
  • OpenAI API Key [critical]
    File: /path/to/project/.env:2
    Potential hardcoded API key detected
```

## License

MIT
