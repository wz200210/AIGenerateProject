package types

// AIComponentType 表示 AI 组件的类型
type AIComponentType string

const (
	TypeLLMFramework     AIComponentType = "llm_framework"
	TypeMLFramework      AIComponentType = "ml_framework"
	TypeVectorDB         AIComponentType = "vector_database"
	TypeRAGTool          AIComponentType = "rag_tool"
	TypeAgentFramework   AIComponentType = "agent_framework"
	TypeModelFile        AIComponentType = "model_file"
	TypeAPIKey           AIComponentType = "api_key"
	TypeConfig           AIComponentType = "config"
	TypeDependency       AIComponentType = "dependency"
	TypeEmbedding        AIComponentType = "embedding"
	TypeMonitoring       AIComponentType = "monitoring"
	TypeDeployment       AIComponentType = "deployment"
)

// Severity 表示风险等级
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// AIComponent 表示检测到的 AI 组件
type AIComponent struct {
	Name        string          `json:"name"`
	Type        AIComponentType `json:"type"`
	Version     string          `json:"version,omitempty"`
	FilePath    string          `json:"file_path"`
	LineNumber  int             `json:"line_number,omitempty"`
	Confidence  float64         `json:"confidence"`
	Severity    Severity        `json:"severity"`
	Description string          `json:"description"`
	RawContent  string          `json:"raw_content,omitempty"`
}

// ScanResult 表示扫描结果
type ScanResult struct {
	ProjectPath  string         `json:"project_path"`
	TotalFiles   int            `json:"total_files"`
	Components   []AIComponent  `json:"components"`
	ScanTime     string         `json:"scan_time"`
	Errors       []string       `json:"errors,omitempty"`
}

// FrameworkInfo 包含框架的检测模式
type FrameworkInfo struct {
	Name         string
	Type         AIComponentType
	Patterns     []string
	FilePatterns []string
	Severity     Severity
	Description  string
}

// ModelFileInfo 模型文件信息
type ModelFileInfo struct {
	Extension   string
	Name        string
	Type        AIComponentType
	Severity    Severity
	Description string
}

// CommonAIFrameworks 常见 AI 框架定义
var CommonAIFrameworks = []FrameworkInfo{
	// ========== 国际 LLM 框架 ==========
	{
		Name:         "OpenAI",
		Type:         TypeLLMFramework,
		Patterns:     []string{"openai", "gpt-", "chatgpt"},
		FilePatterns: []string{"go.mod", "requirements.txt", "package.json", "Cargo.toml"},
		Severity:     SeverityMedium,
		Description:  "OpenAI GPT API",
	},
	{
		Name:         "Anthropic Claude",
		Type:         TypeLLMFramework,
		Patterns:     []string{"anthropic", "claude"},
		FilePatterns: []string{"go.mod", "requirements.txt", "package.json"},
		Severity:     SeverityMedium,
		Description:  "Anthropic Claude API",
	},
	{
		Name:         "Google Gemini",
		Type:         TypeLLMFramework,
		Patterns:     []string{"google.golang.org/genai", "gemini", "palm", "vertexai"},
		FilePatterns: []string{"go.mod", "requirements.txt"},
		Severity:     SeverityMedium,
		Description:  "Google Gemini/Vertex AI API",
	},
	{
		Name:         "Cohere",
		Type:         TypeLLMFramework,
		Patterns:     []string{"cohere"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityMedium,
		Description:  "Cohere LLM API",
	},
	{
		Name:         "Mistral AI",
		Type:         TypeLLMFramework,
		Patterns:     []string{"mistralai", "mistral"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityMedium,
		Description:  "Mistral AI API",
	},

	// ========== 国产 LLM 框架 ==========
	{
		Name:         "百度文心一言",
		Type:         TypeLLMFramework,
		Patterns:     []string{"baidu", "ernie", "wenxin", "文心"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod", "pom.xml"},
		Severity:     SeverityMedium,
		Description:  "百度文心一言/ERNIE Bot API",
	},
	{
		Name:         "阿里通义千问",
		Type:         TypeLLMFramework,
		Patterns:     []string{"dashscope", "tongyi", "qianwen", "通义千问", "qwen"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityMedium,
		Description:  "阿里通义千问/DashScope API",
	},
	{
		Name:         "讯飞星火",
		Type:         TypeLLMFramework,
		Patterns:     []string{"xunfei", "spark", "xinghuo", "讯飞", "星火"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityMedium,
		Description:  "讯飞星火认知大模型 API",
	},
	{
		Name:         "智谱 GLM",
		Type:         TypeLLMFramework,
		Patterns:     []string{"zhipu", "glm-", "chatglm", "智谱"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityMedium,
		Description:  "智谱 AI GLM 大模型 API",
	},
	{
		Name:         "月之暗面 Kimi",
		Type:         TypeLLMFramework,
		Patterns:     []string{"moonshot", "kimi", "月之暗面"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityMedium,
		Description:  "月之暗面 Kimi API",
	},
	{
		Name:         "MiniMax",
		Type:         TypeLLMFramework,
		Patterns:     []string{"minimax"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityMedium,
		Description:  "MiniMax 大模型 API",
	},
	{
		Name:         "DeepSeek",
		Type:         TypeLLMFramework,
		Patterns:     []string{"deepseek"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityMedium,
		Description:  "DeepSeek 大模型 API",
	},
	{
		Name:         "xAI Grok",
		Type:         TypeLLMFramework,
		Patterns:     []string{"xai-sdk", "grok"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityMedium,
		Description:  "xAI Grok API",
	},
	{
		Name:         "OpenRouter",
		Type:         TypeLLMFramework,
		Patterns:     []string{"openrouter"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityMedium,
		Description:  "OpenRouter 统一 LLM API 路由",
	},
	{
		Name:         "SiliconFlow",
		Type:         TypeLLMFramework,
		Patterns:     []string{"siliconflow"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityMedium,
		Description:  "SiliconFlow 硅基流动模型服务",
	},
	{
		Name:         "Volcengine 豆包",
		Type:         TypeLLMFramework,
		Patterns:     []string{"volcengine", "volc-engine", "doubao", "ark-sdk"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityMedium,
		Description:  "字节火山引擎豆包/Ark API",
	},
	{
		Name:         "Groq",
		Type:         TypeLLMFramework,
		Patterns:     []string{"groq"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityMedium,
		Description:  "Groq 高速推理 API",
	},
	{
		Name:         "Together AI",
		Type:         TypeLLMFramework,
		Patterns:     []string{"togetherai", "together-ai", "together"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityMedium,
		Description:  "Together AI 模型推理平台",
	},

	// ========== 编排框架 ==========
	{
		Name:         "LangChain",
		Type:         TypeLLMFramework,
		Patterns:     []string{"langchain", "langchain-"},
		FilePatterns: []string{"go.mod", "requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "LangChain LLM 应用框架",
	},
	{
		Name:         "LlamaIndex",
		Type:         TypeRAGTool,
		Patterns:     []string{"llamaindex", "llama_index"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "LlamaIndex RAG 数据框架",
	},
	{
		Name:         "LangGraph",
		Type:         TypeAgentFramework,
		Patterns:     []string{"langgraph"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "LangGraph 工作流框架",
	},
	{
		Name:         "CrewAI",
		Type:         TypeAgentFramework,
		Patterns:     []string{"crewai"},
		FilePatterns: []string{"requirements.txt"},
		Severity:     SeverityLow,
		Description:  "CrewAI 多智能体框架",
	},
	{
		Name:         "AutoGen",
		Type:         TypeAgentFramework,
		Patterns:     []string{"autogen", "pyautogen"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "Microsoft AutoGen 多智能体框架",
	},
	{
		Name:         "AutoGPT",
		Type:         TypeAgentFramework,
		Patterns:     []string{"autogpt"},
		FilePatterns: []string{"requirements.txt"},
		Severity:     SeverityLow,
		Description:  "AutoGPT 自主智能体",
	},
	{
		Name:         "OpenAI Agents SDK",
		Type:         TypeAgentFramework,
		Patterns:     []string{"openai-agents", "openai_agents"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "OpenAI Agents SDK 智能体框架",
	},
	{
		Name:         "Google ADK",
		Type:         TypeAgentFramework,
		Patterns:     []string{"google-adk", "google_adk"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "Google Agent Development Kit",
	},
	{
		Name:         "Pydantic AI",
		Type:         TypeAgentFramework,
		Patterns:     []string{"pydantic-ai", "pydantic_ai"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Pydantic AI 类型安全 Agent 框架",
	},
	{
		Name:         "Agno",
		Type:         TypeAgentFramework,
		Patterns:     []string{"agno"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Agno 轻量级多智能体框架(原 Phidata)",
	},
	{
		Name:         "smolagents",
		Type:         TypeAgentFramework,
		Patterns:     []string{"smolagents"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "HuggingFace smolagents 轻量 Agent",
	},
	{
		Name:         "Camel AI",
		Type:         TypeAgentFramework,
		Patterns:     []string{"camel-ai", "camelai"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Camel-AI 多智能体协作框架",
	},
	{
		Name:         "Strands Agents",
		Type:         TypeAgentFramework,
		Patterns:     []string{"strands-agents", "strands_agents"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "AWS Strands Agents SDK",
	},
	{
		Name:         "Dify",
		Type:         TypeAgentFramework,
		Patterns:     []string{"dify"},
		FilePatterns: []string{"requirements.txt", "package.json", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "Dify LLM 应用开发平台",
	},
	{
		Name:         "Flowise",
		Type:         TypeAgentFramework,
		Patterns:     []string{"flowiseai", "flowise"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "Flowise 低代码 LLM 工作流",
	},
	{
		Name:         "n8n",
		Type:         TypeAgentFramework,
		Patterns:     []string{"n8n", "@n8n/"},
		FilePatterns: []string{"package.json"},
		Severity:     SeverityLow,
		Description:  "n8n AI 工作流自动化",
	},
	{
		Name:         "Coze SDK",
		Type:         TypeAgentFramework,
		Patterns:     []string{"cozepy", "coze-sdk"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "字节 Coze 扣子 Bot 平台 SDK",
	},

	// ========== ML 框架 ==========
	{
		Name:         "Hugging Face",
		Type:         TypeMLFramework,
		Patterns:     []string{"huggingface", "transformers", "datasets", "accelerate"},
		FilePatterns: []string{"requirements.txt", "package.json", "Cargo.toml"},
		Severity:     SeverityLow,
		Description:  "Hugging Face Transformers",
	},
	{
		Name:         "TensorFlow",
		Type:         TypeMLFramework,
		Patterns:     []string{"tensorflow", "tf.", "keras"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityLow,
		Description:  "TensorFlow ML framework",
	},
	{
		Name:         "PyTorch",
		Type:         TypeMLFramework,
		Patterns:     []string{"torch", "pytorch", "lightning"},
		FilePatterns: []string{"requirements.txt", "package.json", "Cargo.toml"},
		Severity:     SeverityLow,
		Description:  "PyTorch ML framework",
	},
	{
		Name:         "ONNX Runtime",
		Type:         TypeMLFramework,
		Patterns:     []string{"onnxruntime", "onnx"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityLow,
		Description:  "ONNX Runtime",
	},
	{
		Name:         "Ollama",
		Type:         TypeLLMFramework,
		Patterns:     []string{"ollama"},
		FilePatterns: []string{"go.mod", "requirements.txt"},
		Severity:     SeverityLow,
		Description:  "Ollama 本地 LLM",
	},
	{
		Name:         "Vercel AI SDK",
		Type:         TypeLLMFramework,
		Patterns:     []string{"ai", "@ai-sdk/"},
		FilePatterns: []string{"package.json"},
		Severity:     SeverityLow,
		Description:  "Vercel AI SDK",
	},
	{
		Name:         "LiteLLM",
		Type:         TypeLLMFramework,
		Patterns:     []string{"litellm"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "LiteLLM 统一 LLM API 网关",
	},
	{
		Name:         "OpenClaw",
		Type:         TypeAgentFramework,
		Patterns:     []string{"openclaw", "@openclaw"},
		FilePatterns: []string{"package.json", "go.mod", "config.yaml", "config.yml"},
		Severity:     SeverityLow,
		Description:  "OpenClaw AI Agent Gateway",
	},

	// ========== 向量数据库 ==========
	{
		Name:         "Pinecone",
		Type:         TypeVectorDB,
		Patterns:     []string{"pinecone"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityMedium,
		Description:  "Pinecone 向量数据库",
	},
	{
		Name:         "Chroma",
		Type:         TypeVectorDB,
		Patterns:     []string{"chromadb", "chroma"},
		FilePatterns: []string{"requirements.txt"},
		Severity:     SeverityLow,
		Description:  "Chroma 开源向量数据库",
	},
	{
		Name:         "Milvus",
		Type:         TypeVectorDB,
		Patterns:     []string{"milvus", "pymilvus"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityMedium,
		Description:  "Milvus 分布式向量数据库",
	},
	{
		Name:         "Weaviate",
		Type:         TypeVectorDB,
		Patterns:     []string{"weaviate", "weaviate-client"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityMedium,
		Description:  "Weaviate 向量搜索引擎",
	},
	{
		Name:         "Qdrant",
		Type:         TypeVectorDB,
		Patterns:     []string{"qdrant", "qdrant-client"},
		FilePatterns: []string{"requirements.txt", "package.json", "Cargo.toml"},
		Severity:     SeverityLow,
		Description:  "Qdrant 向量数据库",
	},
	{
		Name:         "FAISS",
		Type:         TypeVectorDB,
		Patterns:     []string{"faiss", "faiss-cpu", "faiss-gpu"},
		FilePatterns: []string{"requirements.txt"},
		Severity:     SeverityLow,
		Description:  "Facebook AI Similarity Search",
	},
	{
		Name:         "pgvector",
		Type:         TypeVectorDB,
		Patterns:     []string{"pgvector"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity:     SeverityLow,
		Description:  "PostgreSQL 向量扩展",
	},
	{
		Name:         "Redis Vector",
		Type:         TypeVectorDB,
		Patterns:     []string{"redisvl", "redis-vector"},
		FilePatterns: []string{"requirements.txt"},
		Severity:     SeverityLow,
		Description:  "Redis Vector Library",
	},
	{
		Name:         "Elasticsearch",
		Type:         TypeVectorDB,
		Patterns:     []string{"elasticsearch", "opensearch"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod", "pom.xml"},
		Severity:     SeverityLow,
		Description:  "Elasticsearch/OpenSearch 向量搜索",
	},
	{
		Name:         "Vespa",
		Type:         TypeVectorDB,
		Patterns:     []string{"vespa"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "Vespa 向量搜索引擎",
	},
	{
		Name:         "LanceDB",
		Type:         TypeVectorDB,
		Patterns:     []string{"lancedb"},
		FilePatterns: []string{"requirements.txt", "package.json", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "LanceDB 嵌入式向量数据库",
	},

	// ========== Embedding 模型 ==========
	{
		Name:         "OpenAI Embeddings",
		Type:         TypeEmbedding,
		Patterns:     []string{"text-embedding", "embedding-3", "embedding-ada"},
		FilePatterns: []string{".py", ".js", ".ts", ".go"},
		Severity:     SeverityLow,
		Description:  "OpenAI Embedding 模型",
	},
	{
		Name:         "Sentence Transformers",
		Type:         TypeEmbedding,
		Patterns:     []string{"sentence-transformers", "sentence_transformers"},
		FilePatterns: []string{"requirements.txt"},
		Severity:     SeverityLow,
		Description:  "Sentence Transformers 嵌入模型",
	},
	{
		Name:         "BGE Embeddings",
		Type:         TypeEmbedding,
		Patterns:     []string{"bge-", "bge_m3", "flagembedding"},
		FilePatterns: []string{"requirements.txt"},
		Severity:     SeverityLow,
		Description:  "BGE 中文 Embedding 模型",
	},

	// ========== 部署/推理框架 ==========
	{
		Name:         "vLLM",
		Type:         TypeDeployment,
		Patterns:     []string{"vllm"},
		FilePatterns: []string{"requirements.txt", "Dockerfile"},
		Severity:     SeverityLow,
		Description:  "vLLM 大模型推理引擎",
	},
	{
		Name:         "Text Generation Inference",
		Type:         TypeDeployment,
		Patterns:     []string{"text-generation-inference", "tgi"},
		FilePatterns: []string{"requirements.txt", "Dockerfile"},
		Severity:     SeverityLow,
		Description:  "HuggingFace TGI 推理服务",
	},
	{
		Name:         "Triton Inference Server",
		Type:         TypeDeployment,
		Patterns:     []string{"triton", "tritonserver"},
		FilePatterns: []string{"requirements.txt", "Dockerfile", "config.pbtxt"},
		Severity:     SeverityLow,
		Description:  "NVIDIA Triton 推理服务器",
	},
	{
		Name:         "BentoML",
		Type:         TypeDeployment,
		Patterns:     []string{"bentoml"},
		FilePatterns: []string{"requirements.txt", "bentofile.yaml"},
		Severity:     SeverityLow,
		Description:  "BentoML 模型服务框架",
	},
	{
		Name:         "SGLang",
		Type:         TypeDeployment,
		Patterns:     []string{"sglang"},
		FilePatterns: []string{"requirements.txt", "Dockerfile", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "SGLang 高性能 LLM 推理引擎",
	},
	{
		Name:         "llama.cpp",
		Type:         TypeDeployment,
		Patterns:     []string{"llama-cpp", "llama_cpp", "llama-server"},
		FilePatterns: []string{"requirements.txt", "Dockerfile", "CMakeLists.txt"},
		Severity:     SeverityLow,
		Description:  "llama.cpp 本地 LLM 推理(C/C++)",
	},
	{
		Name:         "LM Studio",
		Type:         TypeDeployment,
		Patterns:     []string{"lmstudio", "lm-studio"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "LM Studio 本地 LLM 桌面运行时",
	},
	{
		Name:         "MLX",
		Type:         TypeDeployment,
		Patterns:     []string{"mlx-lm", "mlx-vlm", "mlx-whisper"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Apple MLX 本地推理框架",
	},
	{
		Name:         "KTransformers",
		Type:         TypeDeployment,
		Patterns:     []string{"ktransformers"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "KTransformers 异构推理框架",
	},
	{
		Name:         "NVIDIA NIM",
		Type:         TypeDeployment,
		Patterns:     []string{"nvidia-nim", "nvidia_nim", "nvcr.io/nim"},
		FilePatterns: []string{"requirements.txt", "Dockerfile", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "NVIDIA NIM 推理微服务",
	},
	{
		Name:         "TEI",
		Type:         TypeDeployment,
		Patterns:     []string{"text-embeddings-inference"},
		FilePatterns: []string{"requirements.txt", "Dockerfile"},
		Severity:     SeverityLow,
		Description:  "HuggingFace TEI 嵌入推理服务",
	},
	{
		Name:         "LocalAI",
		Type:         TypeDeployment,
		Patterns:     []string{"localai"},
		FilePatterns: []string{"requirements.txt", "Dockerfile", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "LocalAI 本地 OpenAI 兼容推理",
	},
	{
		Name:         "Open WebUI",
		Type:         TypeDeployment,
		Patterns:     []string{"open-webui", "open_webui"},
		FilePatterns: []string{"requirements.txt", "package.json", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "Open WebUI LLM 前端界面",
	},
	{
		Name:         "MCP SDK",
		Type:         TypeAgentFramework,
		Patterns:     []string{"@modelcontextprotocol", "mcp-sdk", "mcp_server"},
		FilePatterns: []string{"package.json", "requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Model Context Protocol SDK",
	},

	// ========== 监控/可观测性 ==========
	{
		Name:         "LangSmith",
		Type:         TypeMonitoring,
		Patterns:     []string{"langsmith"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "LangSmith LLM 应用监控",
	},
	{
		Name:         "Langfuse",
		Type:         TypeMonitoring,
		Patterns:     []string{"langfuse"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "Langfuse LLM 可观测性",
	},
	{
		Name:         "Weights & Biases",
		Type:         TypeMonitoring,
		Patterns:     []string{"wandb"},
		FilePatterns: []string{"requirements.txt"},
		Severity:     SeverityLow,
		Description:  "W&B 实验跟踪",
	},
	{
		Name:         "MLflow",
		Type:         TypeMonitoring,
		Patterns:     []string{"mlflow"},
		FilePatterns: []string{"requirements.txt"},
		Severity:     SeverityLow,
		Description:  "MLflow 模型生命周期管理",
	},
	{
		Name:         "Helicone",
		Type:         TypeMonitoring,
		Patterns:     []string{"helicone"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "Helicone LLM 可观测性网关",
	},
	{
		Name:         "Arize Phoenix",
		Type:         TypeMonitoring,
		Patterns:     []string{"arize-phoenix", "arize_phoenix"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "Arize Phoenix LLM Tracing/评估",
	},
	{
		Name:         "Traceloop OpenLLMetry",
		Type:         TypeMonitoring,
		Patterns:     []string{"traceloop", "openllmetry"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "OpenLLMetry LLM 可观测性",
	},
	{
		Name:         "Portkey",
		Type:         TypeMonitoring,
		Patterns:     []string{"portkey", "portkey-ai"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "Portkey AI 网关/可观测性",
	},
	{
		Name:         "Braintrust",
		Type:         TypeMonitoring,
		Patterns:     []string{"braintrust"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "Braintrust LLM 评估平台",
	},

	// ========== 文档处理 ==========
	{
		Name:         "Unstructured",
		Type:         TypeRAGTool,
		Patterns:     []string{"unstructured"},
		FilePatterns: []string{"requirements.txt"},
		Severity:     SeverityLow,
		Description:  "Unstructured 文档解析",
	},
	{
		Name:         "PyPDF",
		Type:         TypeRAGTool,
		Patterns:     []string{"pypdf", "pypdf2", "pypdf3"},
		FilePatterns: []string{"requirements.txt"},
		Severity:     SeverityLow,
		Description:  "PDF 文档处理",
	},
	{
		Name:         "BeautifulSoup",
		Type:         TypeRAGTool,
		Patterns:     []string{"beautifulsoup", "bs4"},
		FilePatterns: []string{"requirements.txt"},
		Severity:     SeverityLow,
		Description:  "HTML/XML 解析（网页抓取）",
	},
	{
		Name:         "Haystack",
		Type:         TypeRAGTool,
		Patterns:     []string{"farm-haystack", "haystack-ai"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "deepset Haystack RAG 管道框架",
	},
	{
		Name:         "RAGFlow",
		Type:         TypeRAGTool,
		Patterns:     []string{"ragflow"},
		FilePatterns: []string{"requirements.txt", "package.json", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "RAGFlow 深度文档理解 RAG 引擎",
	},
	{
		Name:         "LightRAG",
		Type:         TypeRAGTool,
		Patterns:     []string{"lightrag"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "LightRAG 轻量 RAG 框架",
	},
	{
		Name:         "GraphRAG",
		Type:         TypeRAGTool,
		Patterns:     []string{"graphrag"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "微软 GraphRAG 知识图谱增强 RAG",
	},
	{
		Name:         "txtai",
		Type:         TypeRAGTool,
		Patterns:     []string{"txtai"},
		FilePatterns: []string{"requirements.txt"},
		Severity:     SeverityLow,
		Description:  "txtai 语义搜索/嵌入管道",
	},

	// ========== 2025-2026 新兴: AI 编码代理 ==========
	{
		Name:         "Aider",
		Type:         TypeAgentFramework,
		Patterns:     []string{"aider", "aider-chat"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml", ".aider.conf.yml"},
		Severity:     SeverityLow,
		Description:  "Aider 终端 AI 结对编程",
	},
	{
		Name:         "Cline",
		Type:         TypeAgentFramework,
		Patterns:     []string{"cline"},
		FilePatterns: []string{"package.json"},
		Severity:     SeverityLow,
		Description:  "Cline 自主 AI 编码代理",
	},
	{
		Name:         "OpenAI Codex CLI",
		Type:         TypeAgentFramework,
		Patterns:     []string{"@openai/codex", "codex-cli"},
		FilePatterns: []string{"package.json", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "OpenAI Codex 终端编码代理",
	},
	{
		Name:         "opencode",
		Type:         TypeAgentFramework,
		Patterns:     []string{"opencode"},
		FilePatterns: []string{"package.json", "opencode.json"},
		Severity:     SeverityLow,
		Description:  "opencode 终端 AI 编码代理(sst)",
	},

	// ========== 2025-2026 新兴: Agent 框架 ==========
	{
		Name:         "Langflow",
		Type:         TypeAgentFramework,
		Patterns:     []string{"langflow"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml", "package.json"},
		Severity:     SeverityLow,
		Description:  "Langflow 可视化 AI 代理构建平台",
	},
	{
		Name:         "Mastra",
		Type:         TypeAgentFramework,
		Patterns:     []string{"@mastra/core", "mastra"},
		FilePatterns: []string{"package.json"},
		Severity:     SeverityLow,
		Description:  "Mastra TypeScript AI Agent 框架",
	},
	{
		Name:         "Semantic Kernel",
		Type:         TypeAgentFramework,
		Patterns:     []string{"semantic-kernel", "semantic_kernel"},
		FilePatterns: []string{"requirements.txt", "package.json", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Microsoft Semantic Kernel 编排 SDK",
	},
	{
		Name:         "MS Agent Framework",
		Type:         TypeAgentFramework,
		Patterns:     []string{"agent-framework", "agent_framework"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Microsoft Agent Framework (AutoGen+SK 合并)",
	},
	{
		Name:         "DSPy",
		Type:         TypeAgentFramework,
		Patterns:     []string{"dspy"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "DSPy 声明式 LLM 优化框架(Stanford)",
	},
	{
		Name:         "Instructor",
		Type:         TypeAgentFramework,
		Patterns:     []string{"instructor"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Instructor 结构化 LLM 输出",
	},
	{
		Name:         "Pipecat",
		Type:         TypeAgentFramework,
		Patterns:     []string{"pipecat"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Pipecat 语音/多模态 AI Agent 框架",
	},

	// ========== 2025-2026 新兴: AI 网关 ==========
	{
		Name:         "One API",
		Type:         TypeLLMFramework,
		Patterns:     []string{"one-api"},
		FilePatterns: []string{"docker-compose.yaml", "docker-compose.yml"},
		Severity:     SeverityMedium,
		Description:  "One API LLM API 管理分发网关",
	},
	{
		Name:         "New API",
		Type:         TypeLLMFramework,
		Patterns:     []string{"new-api", "QuantumNous"},
		FilePatterns: []string{"docker-compose.yaml", "docker-compose.yml"},
		Severity:     SeverityMedium,
		Description:  "New API 统一 AI 模型聚合平台",
	},
	{
		Name:         "Bifrost",
		Type:         TypeLLMFramework,
		Patterns:     []string{"bifrost"},
		FilePatterns: []string{"docker-compose.yaml", "go.mod", "requirements.txt"},
		Severity:     SeverityLow,
		Description:  "Bifrost 企业级 AI 网关",
	},

	// ========== 2025-2026 新兴: 推理服务/多模态 ==========
	{
		Name:         "Xinference",
		Type:         TypeDeployment,
		Patterns:     []string{"xinference"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "Xorbits Inference 模型推理平台",
	},
	{
		Name:         "TabbyAPI",
		Type:         TypeDeployment,
		Patterns:     []string{"tabbyapi", "exllamav2"},
		FilePatterns: []string{"requirements.txt", "config.yml", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "TabbyAPI ExLlamaV2 高速推理",
	},
	{
		Name:         "ComfyUI",
		Type:         TypeDeployment,
		Patterns:     []string{"comfyui", "comfy-nodes"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml", "custom_nodes"},
		Severity:     SeverityMedium,
		Description:  "ComfyUI 节点式图像/视频生成",
	},
	{
		Name:         "LiveKit",
		Type:         TypeDeployment,
		Patterns:     []string{"livekit", "@livekit/"},
		FilePatterns: []string{"package.json", "requirements.txt", "go.mod"},
		Severity:     SeverityLow,
		Description:  "LiveKit 实时语音/视频 AI 通信",
	},
	{
		Name:         "Kokoro TTS",
		Type:         TypeDeployment,
		Patterns:     []string{"kokoro"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Kokoro 轻量 TTS 引擎",
	},
	{
		Name:         "Faster Whisper",
		Type:         TypeDeployment,
		Patterns:     []string{"faster-whisper", "faster_whisper", "whisperx", "whisper.cpp"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "faster-whisper/whisper.cpp 语音识别",
	},

	// ========== 2025-2026 新兴: RAG/文档/评估 ==========
	{
		Name:         "Docling",
		Type:         TypeRAGTool,
		Patterns:     []string{"docling"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "IBM Docling 文档解析/转换",
	},
	{
		Name:         "Promptfoo",
		Type:         TypeMonitoring,
		Patterns:     []string{"promptfoo"},
		FilePatterns: []string{"package.json", "promptfooconfig.yaml", "promptfoo.config.yaml"},
		Severity:     SeverityLow,
		Description:  "Promptfoo LLM 评估/红队测试",
	},
	{
		Name:         "DeepEval",
		Type:         TypeMonitoring,
		Patterns:     []string{"deepeval"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "DeepEval LLM 评估框架",
	},
	{
		Name:         "Ragas",
		Type:         TypeMonitoring,
		Patterns:     []string{"ragas"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Ragas RAG 管道评估",
	},
	{
		Name:         "Vanna AI",
		Type:         TypeRAGTool,
		Patterns:     []string{"vanna"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Vanna Text-to-SQL RAG",
	},

	// ========== 2025-2026 新兴: AI 聊天前端 ==========
	{
		Name:         "LobeChat",
		Type:         TypeDeployment,
		Patterns:     []string{"lobe-chat", "lobechat", "@lobehub/"},
		FilePatterns: []string{"package.json", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "LobeChat 开源 AI 聊天框架",
	},
	{
		Name:         "NextChat",
		Type:         TypeDeployment,
		Patterns:     []string{"nextchat", "chatgpt-next-web"},
		FilePatterns: []string{"package.json", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "NextChat (ChatGPT-Next-Web) AI 前端",
	},

	// ========== 2025-2026 新兴: 模型仓库 ==========
	{
		Name:         "ModelScope",
		Type:         TypeMLFramework,
		Patterns:     []string{"modelscope"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "阿里 ModelScope 模型仓库",
	},

	// ========== 2026 新兴: AI 编码代理(第二波) ==========
	{
		Name:         "Claude Code",
		Type:         TypeAgentFramework,
		Patterns:     []string{"@anthropic-ai/claude-code", "claude-code", "claude code"},
		FilePatterns: []string{"package.json", "pyproject.toml", ".claude", "claude.md"},
		Severity:     SeverityLow,
		Description:  "Anthropic Claude Code 终端编码代理",
	},
	{
		Name:         "OpenHands",
		Type:         TypeAgentFramework,
		Patterns:     []string{"openhands", "OpenHands-AI", "opendevin"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml", "package.json", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "OpenHands (原 OpenDevin) 开源自主编码代理",
	},
	{
		Name:         "Roo Code",
		Type:         TypeAgentFramework,
		Patterns:     []string{"roo-code", "roo_code", "@roo-code"},
		FilePatterns: []string{"package.json", ".roomodes", ".roo"},
		Severity:     SeverityLow,
		Description:  "Roo Code 自主 AI 编码代理(Cline 分支)",
	},
	{
		Name:         "Goose",
		Type:         TypeAgentFramework,
		Patterns:     []string{"goose-cli", "block/goose", "@block/goose"},
		FilePatterns: []string{"package.json", "goose.yaml", ".goosehints"},
		Severity:     SeverityLow,
		Description:  "Block Goose 开源 AI 编码/通用代理",
	},
	{
		Name:         "Continue",
		Type:         TypeAgentFramework,
		Patterns:     []string{"continue.dev", "@continuedev"},
		FilePatterns: []string{"package.json", "config.yaml", "config.json"},
		Severity:     SeverityLow,
		Description:  "Continue.dev 开源 AI 编码助手",
	},

	// ========== 2026 新兴: Agent 记忆层 ==========
	{
		Name:         "Letta",
		Type:         TypeAgentFramework,
		Patterns:     []string{"letta", "letta-client", "memgpt"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml", "package.json"},
		Severity:     SeverityLow,
		Description:  "Letta (原 MemGPT) 有状态 AI Agent 框架/记忆服务",
	},
	{
		Name:         "Mem0",
		Type:         TypeAgentFramework,
		Patterns:     []string{"mem0ai", "mem0"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Mem0 智能体持久记忆层",
	},
	{
		Name:         "Zep / Graphiti",
		Type:         TypeAgentFramework,
		Patterns:     []string{"graphiti", "zep-python", "zep_client"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Zep/Graphiti 时序知识图谱 Agent 记忆",
	},

	// ========== 2026 新兴: 浏览器/网页自动化 ==========
	{
		Name:         "Browser Use",
		Type:         TypeAgentFramework,
		Patterns:     []string{"browser-use", "browser_use"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Browser Use AI 浏览器自动化 Agent",
	},
	{
		Name:         "Crawl4AI",
		Type:         TypeRAGTool,
		Patterns:     []string{"crawl4ai"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Crawl4AI 面向 LLM 的异步网页爬虫",
	},
	{
		Name:         "Firecrawl",
		Type:         TypeRAGTool,
		Patterns:     []string{"firecrawl-py", "firecrawl"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml", "package.json"},
		Severity:     SeverityLow,
		Description:  "Firecrawl 网页转 Markdown/结构化数据 API",
	},

	// ========== 2026 新兴: 微调框架 ==========
	{
		Name:         "Unsloth",
		Type:         TypeMLFramework,
		Patterns:     []string{"unsloth"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Unsloth 高效 LLM 微调框架(2-5x 提速)",
	},
	{
		Name:         "LLaMA-Factory",
		Type:         TypeMLFramework,
		Patterns:     []string{"llama-factory", "llamafactory"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "LLaMA-Factory 统一微调框架(100+ 模型)",
	},
	{
		Name:         "Axolotl",
		Type:         TypeMLFramework,
		Patterns:     []string{"axolotl"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml", "axolotl.yaml"},
		Severity:     SeverityLow,
		Description:  "Axolotl YAML 驱动 LLM 微调工具",
	},
	{
		Name:         "TRL",
		Type:         TypeMLFramework,
		Patterns:     []string{"trl", "trl-transformers"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "HuggingFace TRL 强化学习微调库",
	},
	{
		Name:         "ms-swift",
		Type:         TypeMLFramework,
		Patterns:     []string{"ms-swift", "swift-finetune"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "ModelScope Swift 多模态微调/部署框架",
	},

	// ========== 2026 新兴: 图像/视频生成服务 ==========
	{
		Name:         "Stable Diffusion WebUI",
		Type:         TypeDeployment,
		Patterns:     []string{"stable-diffusion-webui", "sd-webui", "auto1111"},
		FilePatterns: []string{"requirements.txt", "webui-user.bat", "webui-user.sh"},
		Severity:     SeverityMedium,
		Description:  "Automatic1111 Stable Diffusion WebUI",
	},
	{
		Name:         "InvokeAI",
		Type:         TypeDeployment,
		Patterns:     []string{"invokeai", "InvokeAI"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "InvokeAI 专业图像生成工作台",
	},
	{
		Name:         "Fooocus",
		Type:         TypeDeployment,
		Patterns:     []string{"fooocus"},
		FilePatterns: []string{"requirements.txt", "Dockerfile"},
		Severity:     SeverityLow,
		Description:  "Fooocus 简洁 Stable Diffusion 界面",
	},
	{
		Name:         "Diffusers",
		Type:         TypeMLFramework,
		Patterns:     []string{"diffusers"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "HuggingFace Diffusers 扩散模型库",
	},

	// ========== 2026 新兴: 语音 AI ==========
	{
		Name:         "Vapi",
		Type:         TypeAgentFramework,
		Patterns:     []string{"@vapi-ai", "vapi-python"},
		FilePatterns: []string{"package.json", "requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Vapi 语音 AI Agent 平台 SDK",
	},
	{
		Name:         "Fish Audio",
		Type:         TypeDeployment,
		Patterns:     []string{"fish-audio", "fishaudio", "openaudio"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Fish Audio/OpenAudio TTS/语音克隆",
	},
	{
		Name:         "Ultravox",
		Type:         TypeAgentFramework,
		Patterns:     []string{"ultravox"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Ultravox 实时语音理解开源模型/框架",
	},

	// ========== 2026 新兴: 自托管 AI 前端/运行时 ==========
	{
		Name:         "LibreChat",
		Type:         TypeDeployment,
		Patterns:     []string{"librechat", "@librechat"},
		FilePatterns: []string{"package.json", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "LibreChat 自托管多模型 AI 聊天平台",
	},
	{
		Name:         "AnythingLLM",
		Type:         TypeDeployment,
		Patterns:     []string{"anythingllm"},
		FilePatterns: []string{"package.json", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "AnythingLLM 一体化自托管 RAG/聊天桌面+服务",
	},
	{
		Name:         "Jan",
		Type:         TypeDeployment,
		Patterns:     []string{"janhq", "@janhq"},
		FilePatterns: []string{"package.json", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "Jan 离线本地 AI 桌面运行时",
	},
	{
		Name:         "Open Interpreter",
		Type:         TypeAgentFramework,
		Patterns:     []string{"open-interpreter", "open_interpreter"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Open Interpreter 本地代码执行 AI 代理",
	},

	// ========== 2026 新兴: AI 搜索 API ==========
	{
		Name:         "Tavily",
		Type:         TypeRAGTool,
		Patterns:     []string{"tavily"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml", "package.json"},
		Severity:     SeverityLow,
		Description:  "Tavily AI 搜索 API(Agent 专用)",
	},
	{
		Name:         "Exa",
		Type:         TypeRAGTool,
		Patterns:     []string{"exa-py", "exa-js"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml", "package.json"},
		Severity:     SeverityLow,
		Description:  "Exa 神经/语义搜索 API",
	},

	// ========== 2026 新兴: 向量数据库补充 ==========
	{
		Name:         "Marqo",
		Type:         TypeVectorDB,
		Patterns:     []string{"marqo"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Marqo 多模态向量搜索引擎",
	},
	{
		Name:         "Vald",
		Type:         TypeVectorDB,
		Patterns:     []string{"vdaas/vald"},
		FilePatterns: []string{"go.mod", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "Vald 云原生分布式向量搜索引擎",
	},

	// ========== 2026 新兴: 评估/可观测性补充 ==========
	{
		Name:         "TruLens",
		Type:         TypeMonitoring,
		Patterns:     []string{"trulens", "trulens-eval"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "TruLens LLM 应用评估/追踪",
	},
	{
		Name:         "Giskard",
		Type:         TypeMonitoring,
		Patterns:     []string{"giskard"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Giskard AI 安全测试/评估",
	},
	{
		Name:         "Opik",
		Type:         TypeMonitoring,
		Patterns:     []string{"opik", "comet-ml/opik"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml", "package.json"},
		Severity:     SeverityLow,
		Description:  "Comet Opik LLM Tracing/评估平台",
	},
	{
		Name:         "LangWatch",
		Type:         TypeMonitoring,
		Patterns:     []string{"langwatch"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity:     SeverityLow,
		Description:  "LangWatch LLM 可观测性/分析",
	},
	{
		Name:         "AgentOps",
		Type:         TypeMonitoring,
		Patterns:     []string{"agentops"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "AgentOps AI Agent 运行时监控",
	},
	{
		Name:         "Weave",
		Type:         TypeMonitoring,
		Patterns:     []string{"weave", "wandb-weave"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "W&B Weave LLM 追踪/评估",
	},

	// ========== 2026 新兴: AI 安全/护栏 ==========
	{
		Name:         "Guardrails AI",
		Type:         TypeAgentFramework,
		Patterns:     []string{"guardrails-ai", "guardrails_ai"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Guardrails AI LLM 输出校验框架",
	},
	{
		Name:         "NeMo Guardrails",
		Type:         TypeAgentFramework,
		Patterns:     []string{"nemoguardrails", "nemo-guardrails"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml", "config.yml"},
		Severity:     SeverityLow,
		Description:  "NVIDIA NeMo Guardrails 对话护栏",
	},

	// ========== 2026-10-12 侦察: AI 编码 CLI(第三波) ==========
	{
		Name:         "Gemini CLI",
		Type:         TypeAgentFramework,
		Patterns:     []string{"@google/gemini-cli", "gemini-cli"},
		FilePatterns: []string{"package.json"},
		Severity:     SeverityLow,
		Description:  "Google Gemini CLI 开源终端 AI 编码代理",
	},
	{
		Name:         "Qwen Code",
		Type:         TypeAgentFramework,
		Patterns:     []string{"@qwen-code/qwen-code", "qwen-code"},
		FilePatterns: []string{"package.json"},
		Severity:     SeverityLow,
		Description:  "阿里 Qwen Code 终端 AI 编码代理",
	},
	{
		Name:         "SWE-agent",
		Type:         TypeAgentFramework,
		Patterns:     []string{"sweagent"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Princeton SWE-agent 软件工程自主代理",
	},

	// ========== 2026-10-12 侦察: 推理引擎/训练框架补充 ==========
	{
		Name:         "TensorRT-LLM",
		Type:         TypeDeployment,
		Patterns:     []string{"tensorrt-llm", "tensorrt_llm"},
		FilePatterns: []string{"requirements.txt", "Dockerfile"},
		Severity:     SeverityLow,
		Description:  "NVIDIA TensorRT-LLM 高性能 LLM 推理引擎",
	},
	{
		Name:         "LMDeploy",
		Type:         TypeDeployment,
		Patterns:     []string{"lmdeploy"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "OpenMMLab LMDeploy/TurboMind LLM 推理服务",
	},
	{
		Name:         "OpenVINO",
		Type:         TypeDeployment,
		Patterns:     []string{"openvino"},
		FilePatterns: []string{"requirements.txt", "Dockerfile"},
		Severity:     SeverityLow,
		Description:  "Intel OpenVINO 模型推理工具包",
	},
	{
		Name:         "DeepSpeed",
		Type:         TypeMLFramework,
		Patterns:     []string{"deepspeed"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Microsoft DeepSpeed 分布式训练框架",
	},
	{
		Name:         "Megatron-LM",
		Type:         TypeMLFramework,
		Patterns:     []string{"megatron-core", "megatron_lm"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "NVIDIA Megatron-LM 大规模模型训练框架",
	},
	{
		Name:         "NeMo",
		Type:         TypeMLFramework,
		Patterns:     []string{"nemo-toolkit", "nemo_toolkit"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "NVIDIA NeMo 对话式 AI 训练框架",
	},
	{
		Name:         "PaddlePaddle",
		Type:         TypeMLFramework,
		Patterns:     []string{"paddlepaddle", "paddlepaddle-gpu"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "百度 PaddlePaddle 深度学习框架",
	},
	{
		Name:         "PaddleOCR",
		Type:         TypeMLFramework,
		Patterns:     []string{"paddleocr"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "百度 PaddleOCR 文字识别工具库",
	},
	{
		Name:         "Ultralytics",
		Type:         TypeMLFramework,
		Patterns:     []string{"ultralytics"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Ultralytics YOLO 目标检测/分割框架",
	},

	// ========== 2026-10-12 侦察: 视觉数据平台 ==========
	{
		Name:         "FiftyOne",
		Type:         TypeMLFramework,
		Patterns:     []string{"fiftyone"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Voxel51 FiftyOne CV 数据集管理平台",
	},
	{
		Name:         "Label Studio",
		Type:         TypeDeployment,
		Patterns:     []string{"label-studio"},
		FilePatterns: []string{"requirements.txt", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "Label Studio 数据标注平台",
	},
	{
		Name:         "CVAT",
		Type:         TypeDeployment,
		Patterns:     []string{"cvat-sdk", "cvat-cli"},
		FilePatterns: []string{"requirements.txt", "docker-compose.yaml"},
		Severity:     SeverityLow,
		Description:  "CVAT 计算机视觉标注平台",
	},

	// ========== 2026-10-12 侦察: RAG/文档/图像工具链 ==========
	{
		Name:         "Kotaemon",
		Type:         TypeRAGTool,
		Patterns:     []string{"kotaemon"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Cinnamon Kotaemon 开源 RAG UI",
	},
	{
		Name:         "GPT Researcher",
		Type:         TypeAgentFramework,
		Patterns:     []string{"gpt-researcher"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "GPT Researcher 自主深度研究 Agent",
	},
	{
		Name:         "DB-GPT",
		Type:         TypeAgentFramework,
		Patterns:     []string{"dbgpt"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "DB-GPT 数据库原生 AI 应用框架",
	},
	{
		Name:         "MinerU",
		Type:         TypeRAGTool,
		Patterns:     []string{"magic-pdf", "mineru"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "OpenDataLab MinerU PDF/文档解析",
	},
	{
		Name:         "Rembg",
		Type:         TypeDeployment,
		Patterns:     []string{"rembg"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Rembg AI 图像背景移除",
	},
	{
		Name:         "IOPaint",
		Type:         TypeDeployment,
		Patterns:     []string{"iopaint", "lama-cleaner"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "IOPaint AI 图像修复/擦除",
	},
	{
		Name:         "Real-ESRGAN",
		Type:         TypeDeployment,
		Patterns:     []string{"realesrgan"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Real-ESRGAN 图像超分辨率",
	},

	// ========== 2026-10-12 侦察: AI 安全/红队 ==========
	{
		Name:         "PyRIT",
		Type:         TypeAgentFramework,
		Patterns:     []string{"pyrit"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Microsoft PyRIT AI 红队测试框架",
	},
	{
		Name:         "LLM Guard",
		Type:         TypeAgentFramework,
		Patterns:     []string{"llm-guard", "llm_guard"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Protect AI LLM Guard 安全扫描",
	},
	{
		Name:         "garak",
		Type:         TypeAgentFramework,
		Patterns:     []string{"garak"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "NVIDIA garak LLM 漏洞扫描",
	},
	{
		Name:         "Hayhooks",
		Type:         TypeRAGTool,
		Patterns:     []string{"hayhooks"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Haystack Hayhooks 管道即服务",
	},
	{
		Name:         "GPT4All",
		Type:         TypeDeployment,
		Patterns:     []string{"gpt4all"},
		FilePatterns: []string{"requirements.txt", "pyproject.toml"},
		Severity:     SeverityLow,
		Description:  "Nomic GPT4All 本地 LLM 运行时",
	},
}

// CommonModelFiles 常见 AI 模型文件
var CommonModelFiles = []ModelFileInfo{
	{Extension: ".gguf", Name: "GGUF Model", Type: TypeModelFile, Severity: SeverityHigh, Description: "GGUF format model (Llama, etc.)"},
	{Extension: ".ggml", Name: "GGML Model", Type: TypeModelFile, Severity: SeverityHigh, Description: "GGML format model"},
	{Extension: ".bin", Name: "Binary Model", Type: TypeModelFile, Severity: SeverityMedium, Description: "Generic binary model file"},
	{Extension: ".pt", Name: "PyTorch Model", Type: TypeModelFile, Severity: SeverityHigh, Description: "PyTorch model file"},
	{Extension: ".pth", Name: "PyTorch Model", Type: TypeModelFile, Severity: SeverityHigh, Description: "PyTorch model file"},
	{Extension: ".onnx", Name: "ONNX Model", Type: TypeModelFile, Severity: SeverityHigh, Description: "ONNX format model"},
	{Extension: ".pb", Name: "TensorFlow Model", Type: TypeModelFile, Severity: SeverityHigh, Description: "TensorFlow protobuf model"},
	{Extension: ".h5", Name: "Keras Model", Type: TypeModelFile, Severity: SeverityHigh, Description: "Keras/TF HDF5 model"},
	{Extension: ".tflite", Name: "TensorFlow Lite", Type: TypeModelFile, Severity: SeverityMedium, Description: "TensorFlow Lite model"},
	{Extension: ".mlmodel", Name: "CoreML Model", Type: TypeModelFile, Severity: SeverityMedium, Description: "Apple CoreML model"},
	{Extension: ".safetensors", Name: "SafeTensors", Type: TypeModelFile, Severity: SeverityHigh, Description: "Hugging Face SafeTensors"},
	{Extension: ".ckpt", Name: "Checkpoint", Type: TypeModelFile, Severity: SeverityHigh, Description: "Model checkpoint file"},
}

// APIKeyPatterns API 密钥检测模式
var APIKeyPatterns = map[string]string{
	"OPENAI_API_KEY":         "OpenAI API Key",
	"ANTHROPIC_API_KEY":      "Anthropic API Key",
	"GOOGLE_API_KEY":         "Google API Key",
	"HUGGINGFACE_TOKEN":      "HuggingFace Token",
	"HUGGINGFACE_API_TOKEN":  "HuggingFace API Token",
	"AZURE_OPENAI_KEY":       "Azure OpenAI Key",
	"AZURE_OPENAI_ENDPOINT":  "Azure OpenAI Endpoint",
	"COHERE_API_KEY":         "Cohere API Key",
	"PINECONE_API_KEY":       "Pinecone API Key",
	"PINECONE_ENVIRONMENT":   "Pinecone Environment",
	"MILVUS_URI":             "Milvus URI",
	"WEAVIATE_API_KEY":       "Weaviate API Key",
	"QDRANT_API_KEY":         "Qdrant API Key",
	"QDRANT_URL":             "Qdrant URL",
	"MISTRAL_API_KEY":        "Mistral API Key",
	"MOONSHOT_API_KEY":       "Moonshot API Key",
	"DASHSCOPE_API_KEY":      "DashScope API Key",
	"DEEPSEEK_API_KEY":       "DeepSeek API Key",
	"XAI_API_KEY":            "xAI API Key",
	"GROQ_API_KEY":           "Groq API Key",
	"TOGETHER_API_KEY":       "Together AI Key",
	"OPENROUTER_API_KEY":     "OpenRouter API Key",
	"SILICONFLOW_API_KEY":    "SiliconFlow API Key",
	"VOLCENGINE_API_KEY":     "Volcengine API Key",
	"ARK_API_KEY":            "Ark (Doubao) API Key",
	"NVIDIA_NIM_API_KEY":     "NVIDIA NIM API Key",
	"ZHIPU_API_KEY":          "Zhipu API Key",
	"BAIDU_API_KEY":          "Baidu API Key",
	"XUNFEI_API_KEY":         "Xunfei API Key",
	"LANGCHAIN_API_KEY":      "LangSmith API Key",
	"LANGFUSE_API_KEY":       "Langfuse API Key",
	"AI_API_KEY":             "Generic AI API Key",
	"OPENCLAW_API_KEY":       "OpenClaw API Key",
	"OPENCLAW_GATEWAY_TOKEN": "OpenClaw Gateway Token",
	"VOYAGE_API_KEY":         "Voyage AI API Key",
	"JINA_API_KEY":           "Jina AI API Key",
	"ELEVENLABS_API_KEY":     "ElevenLabs API Key",
	"REPLICATE_API_TOKEN":    "Replicate API Token",
	"DEEPGRAM_API_KEY":       "Deepgram API Key",
	"PERPLEXITY_API_KEY":     "Perplexity API Key",
	"FIREWORKS_API_KEY":      "Fireworks AI Key",
	"CO_API_KEY":             "Cohere API Key Alt",
	"HUGGING_FACE_HUB_TOKEN": "HuggingFace Hub Token",
	"TAVILY_API_KEY":         "Tavily Search API Key",
	"EXA_API_KEY":            "Exa Search API Key",
	"E2B_API_KEY":            "E2B Code Interpreter API Key",
	"BROWSERBASE_API_KEY":    "Browserbase API Key",
	"HELICONE_API_KEY":       "Helicone API Key",
	"LANGWATCH_API_KEY":      "LangWatch API Key",
	"MEM0_API_KEY":           "Mem0 Memory API Key",
	"VAPI_API_KEY":           "Vapi Voice Agent API Key",
	"VAPI_PRIVATE_KEY":       "Vapi Private Key",
	"CARTESIA_API_KEY":       "Cartesia TTS API Key",
	"ASSEMBLYAI_API_KEY":     "AssemblyAI API Key",
	"MODAL_TOKEN_ID":         "Modal Token ID",
	"AGENTOPS_API_KEY":       "AgentOps API Key",
	"OPIK_API_KEY":           "Comet Opik API Key",
	"GSK_API_KEY":            "Giskard API Key",
	"GEMINI_API_KEY":         "Google Gemini API Key",
	"LLAMA_CLOUD_API_KEY":    "LlamaParse API Key",
	"WANDB_API_KEY":          "Weights & Biases API Key",
	"BRAINTRUST_API_KEY":     "Braintrust API Key",
	"ARIZE_API_KEY":          "Arize Phoenix API Key",
	"FAL_KEY":                "fal.ai API Key",
	"CEREBRAS_API_KEY":       "Cerebras API Key",
	"UNSTRUCTURED_API_KEY":   "Unstructured API Key",
	"BRAVE_API_KEY":          "Brave Search API Key",
}

// GetComponentStats 返回各类组件数量统计
func GetComponentStats(components []AIComponent) map[AIComponentType]int {
	stats := make(map[AIComponentType]int)
	for _, c := range components {
		stats[c.Type]++
	}
	return stats
}

// RuntimeScanResult 运行时扫描结果
type RuntimeScanResult struct {
	ScanTime       string         `json:"scan_time"`
	ScanDuration   string         `json:"scan_duration,omitempty"`
	ProcessCount   int            `json:"process_count"`
	PortCount      int            `json:"port_count"`
	ContainerCount int            `json:"container_count"`
	Components     []AIComponent  `json:"components"`
	Skills         []SkillInfo    `json:"skills,omitempty"`  // 新增：识别的 Skill 列表
	Errors         []string       `json:"errors,omitempty"`
}

// SkillInfo 表示识别到的 Skill 信息
type SkillInfo struct {
	Name        string `json:"name"`
	Source      string `json:"source"`       // openclaw, claude, etc.
	Description string `json:"description,omitempty"`
	Location    string `json:"location"`     // 配置文件路径
	Version     string `json:"version,omitempty"`
	Enabled     bool   `json:"enabled"`      // 是否启用
}

// FullScanResult 完整扫描结果（静态+运行时）
type FullScanResult struct {
	ProjectPath      string          `json:"project_path"`
	StaticScan       *ScanResult     `json:"static_scan,omitempty"`
	RuntimeScan      *RuntimeScanResult `json:"runtime_scan,omitempty"`
	TotalComponents  int             `json:"total_components"`
	ScanTime         string          `json:"scan_time"`
}