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