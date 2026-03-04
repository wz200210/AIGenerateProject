# AI Component Scanner

一个用于**运行时检测** AI 组件的 Go 工具，专注于发现和识别运行中的 LLM 服务、向量数据库、ML 框架等 AI 基础设施。

⚠️ **v0.3.0 重大变更**：已废弃静态文件扫描，改为纯运行时检测架构。

## 🚀 核心特性

### 🔍 运行时进程检测
- 扫描 `/proc` 目录发现运行中的 AI 进程
- 分析进程父子关系，构建服务拓扑
- 通过 inode 映射端口到进程

### 🎯 智能语义分析
- **Python ML 分析器**：识别 PyTorch、TensorFlow、Transformers 等
- **Node.js 分析器**：检测 OpenAI SDK、LangChain.js 等
- **Docker 分析器**：识别容器化 AI 服务
- **Service Mesh 分析器**：分析服务间调用关系

### 📊 自动版本探测
从多维度获取 AI 服务版本：
1. 命令行参数 (`--version=1.2.3`)
2. 环境变量 (`VERSION`, `APP_VERSION`)
3. 执行探测 (`ollama --version`)
4. HTTP API 探测 (`GET /version`)
5. Docker 镜像标签

### 🌐 网络服务发现
- 检测常见 AI 服务端口（11434 Ollama, 8000 vLLM, 6333 Qdrant 等）
- HTTP 健康/版本端点主动探测
- 端口到进程关联映射

### 🔐 安全检测
- 扫描进程环境变量中的 API Key
- 检测潜在敏感信息泄露
- 标记关键安全风险

## 📦 安装

```bash
# 从源码安装
go install github.com/wz200210/AIGenerateProject/cmd/scanner@latest

# 或克隆后本地编译
git clone git@github.com:wz200210/AIGenerateProject.git
cd AIGenerateProject
go mod tidy
go build -o scanner ./cmd/scanner
```

## 🎮 使用

### 基本扫描（运行时检测）
```bash
# 扫描当前系统运行的 AI 组件
scanner scan

# 输出 JSON 格式（含详细元数据）
scanner scan -o json
```

### 详细扫描
```bash
# 包含进程树和网络连接的详细分析
scanner detail
```

### 版本检查（开发中）
```bash
# 主动探测所有 AI 服务的版本信息
scanner version-check
```

## 🎯 检测能力

### LLM 推理服务
| 服务 | 检测方式 | 版本探测 |
|------|---------|---------|
| Ollama | 进程 + 端口 11434 | ✅ `--version` |
| vLLM | 进程 + 端口 8000 | ✅ `--version` |
| TGI | 进程 + 端口 8080 | ✅ API |
| OpenAI API 代理 | 进程匹配 | ✅ 环境变量 |

### 向量数据库
| 服务 | 默认端口 | 检测方式 |
|------|---------|---------|
| Milvus | 19530 | 进程 + HTTP API |
| Chroma | 8000 | 进程匹配 |
| Weaviate | 8080 | 进程 + HTTP |
| Qdrant | 6333/6334 | 进程 + HTTP API |
| Pinecone | - | 进程匹配 |
| pgvector | 5432 | 进程 + 端口 |

### ML 框架服务
| 框架 | 检测特征 |
|------|---------|
| Hugging Face Transformers | Python 进程 + import |
| PyTorch | `torch` 进程匹配 |
| TensorFlow | `tensorflow` 进程 |
| ONNX Runtime | `onnxruntime` 进程 |

### Agent/RAG 框架
| 框架 | 检测方式 |
|------|---------|
| LangChain | Python 进程分析 |
| LlamaIndex | Python 进程分析 |
| AutoGPT | 进程匹配 |
| CrewAI | 进程匹配 |

## 📊 示例输出

### 控制台报告
```
🔍 AI Component Runtime Scanner v0.3.0
═══════════════════════════════════════════════════════

🔍 Scanning process tree...
🔍 Scanning network services...
🔍 Scanning Docker containers...

╔══════════════════════════════════════════════════════╗
║     Runtime AI Component Scan Report                 ║
╚══════════════════════════════════════════════════════╝

⏱️  Scan Time: 2026-03-04T17:52:05+08:00

📊 Scan Summary:
  • Processes scanned: 3
  • Ports scanned: 2
  • Containers scanned: 1
  • Total components found: 3

Running AI Components:
────────────────────────────────────────────────────────────

[deployment]
  • Ollama [medium] v0.1.33
    Source: /proc/1234 (port 11434)
    PID: 1234 | Ports: 11434 | Exe: /usr/local/bin/ollama

[vector_database]
  • Milvus [medium] v2.3.0
    Source: /proc/2345 (port 19530)
    Docker container: milvus-standalone | Status: Up 3 days

[llm_framework]
  • LangChain Service [low] v0.2.0
    Source: /proc/3456
    Python process using LangChain (PID: 3456)
```

### JSON 输出
```json
{
  "scan_time": "2026-03-04T17:52:05+08:00",
  "scan_duration": "45.6ms",
  "process_count": 3,
  "port_count": 2,
  "container_count": 1,
  "components": [
    {
      "name": "Ollama",
      "type": "deployment",
      "version": "0.1.33",
      "file_path": "/proc/1234 (port 11434)",
      "confidence": 0.95,
      "severity": "medium",
      "description": "PID: 1234 | Ports: 11434 | Exe: /usr/local/bin/ollama"
    }
  ]
}
```

## 🏗️ 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                    Runtime Scanner v0.3.0                    │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ Process Tree │  │ Network Scan │  │ Docker Scan  │     │
│  │   Scanner    │  │   (ss/net)   │  │  (docker ps) │     │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘     │
│         │                 │                 │              │
│         └────────┬────────┴────────┬────────┘              │
│                  │                 │                       │
│         ┌────────▼─────────────────▼────────┐             │
│         │      Semantic Analyzers           │             │
│         │  • Python ML Analyzer             │             │
│         │  • Node.js Analyzer               │             │
│         │  • Docker Analyzer                │             │
│         │  • Service Mesh Analyzer          │             │
│         └────────┬──────────────────────────┘             │
│                  │                                         │
│         ┌────────▼────────┐                              │
│         │ Version Prober  │                              │
│         │  • CLI args     │                              │
│         │  • Environment  │                              │
│         │  • Exec probe   │                              │
│         │  • HTTP API     │                              │
│         └────────┬────────┘                              │
│                  │                                         │
│         ┌────────▼────────┐                              │
│         │ Report Generator│                              │
│         │  • Console      │                              │
│         │  • JSON         │                              │
│         └─────────────────┘                              │
│                                                           │
└─────────────────────────────────────────────────────────────┘
```

## 🔧 技术细节

### 进程检测
- 读取 `/proc/<pid>/cmdline` 获取命令行
- 读取 `/proc/<pid>/environ` 获取环境变量
- 读取 `/proc/<pid>/status` 获取 PPID
- 通过 `/proc/<pid>/exe` 获取可执行文件路径

### 端口映射
- 解析 `/proc/net/tcp` 和 `/proc/net/tcp6`
- 通过 inode 查找绑定端口的进程
- 使用 `/proc/<pid>/fd/` 目录匹配 socket

### 版本探测策略
```
1. 命令行参数解析 → --version, -v
2. 环境变量读取 → VERSION, APP_VERSION
3. 执行版本命令 → <exe> --version
4. HTTP 端点探测 → GET localhost:<port>/version
5. Docker 标签提取 → image:tag
```

## 📁 项目结构

```
.
├── cmd/scanner/              # 命令行入口
│   └── main.go
├── internal/
│   ├── runtime/              # 运行时扫描器（核心）
│   │   └── scanner.go        # 进程/网络/容器扫描
│   └── scanner/              # 报告生成
│       └── scanner.go        # 控制台/JSON 输出
├── pkg/ai/types/             # 类型定义
│   └── types.go
├── go.mod
├── Makefile
└── README.md
```

## 🛠️ 开发

```bash
# 编译
go build -o scanner ./cmd/scanner

# 运行测试
make test

# 交叉编译
make build-all
```

## ⚠️ 已知限制

- 需要 Linux 环境（依赖 `/proc` 文件系统）
- 需要 root 权限或 CAP_SYS_PTRACE 以读取其他用户进程
- Docker 检测需要 docker CLI 和权限
- 某些版本探测可能需要服务响应 HTTP 请求

## 📈 性能

- 扫描 1000+ 进程：~100ms
- 网络端口扫描：~50ms
- Docker 容器扫描：~100ms
- 总扫描时间：通常 < 1秒

## 📝 更新日志

### v0.3.0 (2026-03-04)
- **重大重构**：废弃静态文件扫描，改为纯运行时检测
- **新增**：进程树分析，支持父子进程关系
- **新增**：语义分析器框架，支持 Python/Node.js/Docker
- **新增**：自动版本探测（5 种策略）
- **新增**：端口到进程关联映射
- **新增**：HTTP 服务版本探测

### v0.2.0 (2026-03-04)
- 新增运行时扫描功能（进程、端口、容器）
- 新增 API Key 泄露检测
- 支持运行时 + 静态扫描混合模式

### v0.1.0 (2026-03-04)
- 初始版本
- 支持 50+ AI 组件静态检测
- 控制台和 JSON 输出

## 📄 License

MIT

---

**🔗 GitHub**: https://github.com/wz200210/AIGenerateProject