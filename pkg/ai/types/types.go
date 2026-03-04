package types

// AIComponentType 表示 AI 组件的类型
type AIComponentType string

const (
	TypeLLMFramework   AIComponentType = "llm_framework"
	TypeMLFramework    AIComponentType = "ml_framework"
	TypeModelFile      AIComponentType = "model_file"
	TypeAPIKey         AIComponentType = "api_key"
	TypeConfig         AIComponentType = "config"
	TypeDependency     AIComponentType = "dependency"
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
	Name        string
	Type        AIComponentType
	Patterns    []string
	FilePatterns []string
	Severity    Severity
	Description string
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
	{
		Name:     "OpenAI",
		Type:     TypeLLMFramework,
		Patterns: []string{"openai", "gpt-", "chatgpt"},
		FilePatterns: []string{"go.mod", "requirements.txt", "package.json", "Cargo.toml"},
		Severity: SeverityMedium,
		Description: "OpenAI GPT API",
	},
	{
		Name:     "Anthropic Claude",
		Type:     TypeLLMFramework,
		Patterns: []string{"anthropic", "claude"},
		FilePatterns: []string{"go.mod", "requirements.txt", "package.json"},
		Severity: SeverityMedium,
		Description: "Anthropic Claude API",
	},
	{
		Name:     "Google Gemini",
		Type:     TypeLLMFramework,
		Patterns: []string{"google.golang.org/genai", "gemini", "palm"},
		FilePatterns: []string{"go.mod", "requirements.txt"},
		Severity: SeverityMedium,
		Description: "Google Gemini API",
	},
	{
		Name:     "LangChain",
		Type:     TypeLLMFramework,
		Patterns: []string{"langchain"},
		FilePatterns: []string{"go.mod", "requirements.txt", "package.json"},
		Severity: SeverityLow,
		Description: "LangChain framework",
	},
	{
		Name:     "LlamaIndex",
		Type:     TypeLLMFramework,
		Patterns: []string{"llamaindex", "llama_index"},
		FilePatterns: []string{"requirements.txt", "package.json"},
		Severity: SeverityLow,
		Description: "LlamaIndex framework",
	},
	{
		Name:     "Hugging Face",
		Type:     TypeMLFramework,
		Patterns: []string{"huggingface", "transformers"},
		FilePatterns: []string{"requirements.txt", "package.json", "Cargo.toml"},
		Severity: SeverityLow,
		Description: "Hugging Face Transformers",
	},
	{
		Name:     "TensorFlow",
		Type:     TypeMLFramework,
		Patterns: []string{"tensorflow", "tf."},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity: SeverityLow,
		Description: "TensorFlow ML framework",
	},
	{
		Name:     "PyTorch",
		Type:     TypeMLFramework,
		Patterns: []string{"torch", "pytorch"},
		FilePatterns: []string{"requirements.txt", "package.json", "Cargo.toml"},
		Severity: SeverityLow,
		Description: "PyTorch ML framework",
	},
	{
		Name:     "ONNX Runtime",
		Type:     TypeMLFramework,
		Patterns: []string{"onnxruntime", "onnx"},
		FilePatterns: []string{"requirements.txt", "package.json", "go.mod"},
		Severity: SeverityLow,
		Description: "ONNX Runtime",
	},
	{
		Name:     "Ollama",
		Type:     TypeLLMFramework,
		Patterns: []string{"ollama"},
		FilePatterns: []string{"go.mod", "requirements.txt"},
		Severity: SeverityLow,
		Description: "Ollama local LLM",
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
}

// APIKeyPatterns API 密钥检测模式
var APIKeyPatterns = map[string]string{
	"OPENAI_API_KEY":     "OpenAI API Key",
	"ANTHROPIC_API_KEY":  "Anthropic API Key",
	"GOOGLE_API_KEY":     "Google API Key",
	"HUGGINGFACE_TOKEN":  "HuggingFace Token",
	"AZURE_OPENAI_KEY":   "Azure OpenAI Key",
	"COHERE_API_KEY":     "Cohere API Key",
	"AI_API_KEY":         "Generic AI API Key",
}