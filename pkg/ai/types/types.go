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
	Errors         []string       `json:"errors,omitempty"`
}

// FullScanResult 完整扫描结果（静态+运行时）
type FullScanResult struct {
	ProjectPath      string          `json:"project_path"`
	StaticScan       *ScanResult     `json:"static_scan,omitempty"`
	RuntimeScan      *RuntimeScanResult `json:"runtime_scan,omitempty"`
	TotalComponents  int             `json:"total_components"`
	ScanTime         string          `json:"scan_time"`
}