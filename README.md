# AI Component Scanner

一个用于**运行时检测** AI 组件的 Go 工具，专注于发现和识别运行中的 LLM 服务、向量数据库、ML 框架等 AI 基础设施。

⚠️ **v0.3.1 重要更新**：网络服务检测已优化为三重验证机制（端口+进程名+版本号），有效解决自定义端口配置导致的误报问题。

## 🚀 核心特性

### 🔍 运行时进程检测（跨平台）
- **Linux**：扫描 `/proc` 目录发现运行中的 AI 进程
- **Windows**：使用 PowerShell/wmic/netstat 获取进程信息
- 分析进程父子关系，构建服务拓扑
- 通过端口映射关联进程与服务

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

### 🌐 网络服务发现（三重验证机制）
- **Linux**：通过 inode 映射端口到进程 (`/proc/net/tcp`)
- **Windows**：使用 `netstat -ano` 获取端口-PID 映射
- **第一重**：检测端口是否有进程监听
- **第二重**：验证监听进程名是否匹配组件特征
- **第三重**：必须能获取到版本号才算匹配成功
- 有效避免自定义端口配置导致的误报

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
│                    Runtime Scanner v0.3.1                    │
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
│         └────────┬──────────────────────────┘             │
│                  │                                         │
│         ┌────────▼────────┐                              │
│         │ Network Service │                              │
│         │    Validator    │                              │
│         │  ├─ Port Check  │                              │
│         │  ├─ Process Match│                             │
│         │  └─ Version Probe│                             │
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

### 网络服务检测流程

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ 端口监听？  │──NO──┤    跳过     │     │ 进程名匹配？ │──NO──┤    跳过     │
└──────┬──────┘     └─────────────┘     └──────┬──────┘     └─────────────┘
      │ YES                                   │ YES
      ▼                                       ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ 获取进程PID │────▶│ 进程名匹配？ │────▶│ 版本可获取？ │──NO──┤    跳过     │
└─────────────┘     └─────────────┘     └──────┬──────┘     └─────────────┘
                                               │ YES
                                               ▼
                                        ┌─────────────┐
                                        │ 确认AI组件  │
                                        │ +版本号     │
                                        └─────────────┘
```

## 🔧 技术细节

### 进程检测

#### Linux
- 读取 `/proc/<pid>/cmdline` 获取命令行
- 读取 `/proc/<pid>/environ` 获取环境变量
- 读取 `/proc/<pid>/status` 获取 PPID
- 通过 `/proc/<pid>/exe` 获取可执行文件路径

#### Windows
- 使用 PowerShell `Get-Process` 获取进程列表（主要方案）
- 降级到 `wmic process` 获取命令行和父进程信息
- 使用 `tasklist` 作为最后的降级方案
- 支持解析 CSV 格式输出

### 端口映射

#### Linux
- 解析 `/proc/net/tcp` 和 `/proc/net/tcp6`
- 通过 inode 查找绑定端口的进程
- 使用 `/proc/<pid>/fd/` 目录匹配 socket

#### Windows
- 执行 `netstat -ano` 获取端口和 PID 映射
- 解析状态为 LISTENING 的连接
- 关联端口与进程信息

### 版本探测策略（5种方式）
```
1. 命令行参数解析 → --version, -v
2. 环境变量读取 → VERSION, APP_VERSION
3. 执行版本命令 → <exe> --version
4. HTTP 端点探测 → GET localhost:<port>/version
5. Docker 标签提取 → image:tag
```

### 网络服务检测策略
```
三重验证机制：

1. 端口验证
   └─ 端口是否有进程监听？否 → 跳过
   
2. 进程验证  
   └─ 监听进程名是否匹配组件模式？否 → 跳过
   
3. 版本验证
   └─ 是否能获取版本号？否 → 跳过（避免误报）
   └─ 是 → 确认匹配，记录版本
```

该机制有效避免以下误报场景：
- 其他服务占用组件默认端口（如 nginx 占用 8080）
- 组件使用自定义端口配置
- 残留僵尸进程

## 🤖 OpenClaw 检测

本项目特别优化了对 **OpenClaw AI Agent Gateway** 的检测支持，可在 Linux 和 Windows 平台上识别运行中的 OpenClaw 实例：

### 检测模式
- 进程名匹配：`openclaw`, `openclaw.exe`, `openclaw.*gateway`
- Node.js 模式：`node.*openclaw`, `openclaw.*index.js`
- 路径匹配：`/usr/lib/node_modules/openclaw`, `\\openclaw`
- npm 包：`@openclaw`

### 环境变量检测
- `OPENCLAW_API_KEY`
- `OPENCLAW_GATEWAY_TOKEN`
- `OPENCLAW_CONFIG_PATH`
- `OPENCLAW_HOME`

### 默认端口
- `8787` (OpenClaw 默认端口)
- `3000` (Node.js 常用端口)
- `8080` (通用网关端口)

## 📁 项目结构

```
.
├── cmd/scanner/              # 命令行入口
│   └── main.go
├── internal/
│   ├── runtime/              # 运行时扫描器（核心）
│   │   ├── config_scanner_linux.go    # Linux 实现
│   │   └── config_scanner_windows.go  # Windows 实现
│   ├── scanner/              # 报告生成
│   │   └── scanner.go        # 控制台/JSON 输出
│   └── detector/             # 静态文件检测
│       └── detector.go
├── pkg/ai/types/             # 类型定义
│   └── types.go
├── config/
│   └── rules.yaml            # AI 组件检测规则
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

### Linux
- 需要 root 权限或 CAP_SYS_PTRACE 以读取其他用户进程
- Docker 检测需要 docker CLI 和权限

### Windows
- 部分功能需要 PowerShell 或 wmic 支持
- 进程命令行获取可能受限于系统权限
- Docker 检测需要 Docker Desktop 或 Docker Engine

### 通用
- 某些版本探测可能需要服务响应 HTTP 请求
- 环境变量扫描受限于当前用户权限

## 📈 性能

- 扫描 1000+ 进程：~100ms
- 网络端口扫描：~50ms
- Docker 容器扫描：~100ms
- 总扫描时间：通常 < 1秒

## ✅ 测试验证

### 已验证的 AI 组件 (8个)

| 组件 | 类型 | 检测方式 | 版本 |
|------|------|----------|------|
| **Ollama** | LLM推理服务 | 进程+端口+版本API | v0.17.7 |
| **PyTorch** | ML框架 | 语义分析器 | 2.10.0 |
| **TensorFlow** | ML框架 | 语义分析器 | 2.20.0 |
| **LlamaIndex** | RAG框架 | 语义分析器 | - |
| **Gradio** | 部署工具 | 语义分析器 | - |
| **Streamlit** | 部署工具 | 进程匹配 | - |
| **Jupyter** | 监控工具 | 进程匹配 | - |
| **Chroma** | 向量数据库 | 进程匹配 | - |

### 配置优化

通过添加 `semantic_analyzers` 配置，提升 Python 库的识别率：

```yaml
# config/rules.yaml 示例
- name: "Gradio"
  semantic_analyzers:
    - type: "python_import"
      patterns:
        - "import gradio"
        - "from gradio import"

- name: "LlamaIndex"
  semantic_analyzers:
    - type: "python_import"
      patterns:
        - "import llama_index"
        - "from llama_index import"
```

### 测试脚本

详见 `test/` 目录：
- `pytorch_service.py` - PyTorch 测试服务
- `tensorflow_service.py` - TensorFlow 测试服务
- `llamaindex_service.py` - LlamaIndex 测试服务
- `transformers_service.py` - Transformers 测试服务

## 📝 更新日志

### v0.3.2 (2026-03-12)
- **新增**：Windows 平台支持
  - 使用 PowerShell/wmic/tasklist 获取进程信息
  - 使用 netstat 获取端口映射
  - 支持跨平台编译 (`GOOS=windows`)
- **新增**：OpenClaw AI Gateway 检测优化
  - 添加 Windows 特定进程模式
  - 支持检测 `openclaw.exe` 等 Windows 可执行文件
  - 增强 Node.js 模式匹配

### v0.3.1 (2026-03-05)
- **优化**：网络服务检测改为三重验证机制（端口+进程名+版本号）
- **修复**：解决自定义端口配置导致的误报问题
- **改进**：只有能获取版本号的服务才确认为 AI 组件

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