package runtime

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wz200210/AIGenerateProject/pkg/ai/types"
)

// RuntimeScanner 运行时扫描器
type RuntimeScanner struct {
	processPatterns map[string]*regexp.Regexp
	portPatterns    map[int]string
	containerImages []string
}

// NewRuntimeScanner 创建运行时扫描器
func NewRuntimeScanner() *RuntimeScanner {
	rs := &RuntimeScanner{
		processPatterns: make(map[string]*regexp.Regexp),
		portPatterns:    make(map[int]string),
		containerImages: []string{},
	}
	rs.compilePatterns()
	return rs
}

// compilePatterns 编译检测模式
func (rs *RuntimeScanner) compilePatterns() {
	// 进程名称匹配模式
	patterns := map[string]string{
		"ollama":         `(?i)ollama`,
		"vllm":           `(?i)vllm|vllm\.entrypoints`,
		"text-generation": `(?i)text-generation-inference|tgi`,
		"transformers":   `(?i)transformers|huggingface`,
		"langchain":      `(?i)langchain`,
		"llamaindex":     `(?i)llamaindex`,
		"openai":         `(?i)openai`,
		"torch":          `(?i)torch|pytorch`,
		"tensorflow":     `(?i)tensorflow`,
		"onnx":           `(?i)onnxruntime`,
		"milvus":         `(?i)milvus`,
		"chroma":         `(?i)chromadb|chroma`,
		"weaviate":       `(?i)weaviate`,
		"qdrant":         `(?i)qdrant`,
		"pinecone":       `(?i)pinecone`,
		"redis":          `(?i)redis.*vector|redisvl`,
		"elasticsearch":  `(?i)elasticsearch`,
		"python-llm":     `(?i)python.*llm|python.*gpt`,
		"node-llm":       `(?i)node.*ai|npm.*langchain`,
		"jupyter":        `(?i)jupyter.*notebook|ipython`,
		"mlflow":         `(?i)mlflow`,
		"wandb":          `(?i)wandb`,
		"triton":         `(?i)tritonserver`,
		"bentoml":        `(?i)bentoml`,
	}

	for name, pattern := range patterns {
		if re, err := regexp.Compile(pattern); err == nil {
			rs.processPatterns[name] = re
		}
	}

	// 常见 AI 服务端口
	rs.portPatterns = map[int]string{
		11434: "Ollama API",
		8000:  "vLLM / TGI / FastAPI",
		8080:  "通用 AI 服务",
		5000:  "Flask AI 服务",
		3000:  "Node.js AI 服务",
		6333:  "Qdrant",
		6334:  "Qdrant gRPC",
		19530: "Milvus",
		9091:  "Milvus Proxy",
		6379:  "Redis",
		9200:  "Elasticsearch",
		5601:  "Kibana",
		8888:  "Jupyter",
		5001:  "MLflow",
	}

	// Docker 镜像关键词
	rs.containerImages = []string{
		"ollama", "vllm", "transformers", "langchain", "llamaindex",
		"milvus", "chroma", "weaviate", "qdrant", "pinecone",
		"elasticsearch", "redis", "jupyter", "mlflow", "triton",
		"bentoml", "huggingface", "torch", "tensorflow",
	}
}

// ScanProcesses 扫描运行中的进程
func (rs *RuntimeScanner) ScanProcesses() ([]types.AIComponent, error) {
	var components []types.AIComponent

	// 读取 /proc 目录
	entries, err := os.ReadDir("/proc")
	if err != nil {
		// 可能不是 Linux 系统，尝试使用 ps 命令
		return rs.scanProcessesWithPS()
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// 检查是否是 PID 目录
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		// 读取进程命令行
		cmdline, err := rs.readProcCmdline(pid)
		if err != nil || cmdline == "" {
			continue
		}

		// 匹配 AI 进程
		for name, re := range rs.processPatterns {
			if re.MatchString(cmdline) {
				comp := types.AIComponent{
					Name:        name,
					Type:        types.TypeDeployment,
					FilePath:    fmt.Sprintf("/proc/%d", pid),
					Confidence:  0.9,
					Severity:    types.SeverityMedium,
					Description: fmt.Sprintf("Running process: %s (PID: %d)", cmdline[:min(len(cmdline), 100)], pid),
					RawContent:  cmdline,
				}
				if !rs.componentExists(components, comp) {
					components = append(components, comp)
				}
			}
		}

		// 扫描环境变量中的 API Key
		envComponents := rs.scanProcEnviron(pid)
		components = append(components, envComponents...)
	}

	return components, nil
}

// scanProcessesWithPS 使用 ps 命令扫描进程（备用方案）
func (rs *RuntimeScanner) scanProcessesWithPS() ([]types.AIComponent, error) {
	var components []types.AIComponent

	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return components, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		for name, re := range rs.processPatterns {
			if re.MatchString(line) {
				comp := types.AIComponent{
					Name:        name,
					Type:        types.TypeDeployment,
					Confidence:  0.85,
					Severity:    types.SeverityMedium,
					Description: fmt.Sprintf("Running process: %s", line[:min(len(line), 150)]),
					RawContent:  line,
				}
				if !rs.componentExists(components, comp) {
					components = append(components, comp)
				}
			}
		}
	}

	return components, scanner.Err()
}

// readProcCmdline 读取进程命令行
func (rs *RuntimeScanner) readProcCmdline(pid int) (string, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return "", err
	}
	// cmdline 以 null 分隔
	return strings.ReplaceAll(string(data), "\x00", " "), nil
}

// scanProcEnviron 扫描进程环境变量
func (rs *RuntimeScanner) scanProcEnviron(pid int) []types.AIComponent {
	var components []types.AIComponent

	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/environ", pid))
	if err != nil {
		return components
	}

	env := strings.ReplaceAll(string(data), "\x00", "\n")

	for key, name := range types.APIKeyPatterns {
		if strings.Contains(env, key) {
			comp := types.AIComponent{
				Name:        name,
				Type:        types.TypeAPIKey,
				FilePath:    fmt.Sprintf("/proc/%d/environ", pid),
				Confidence:  0.95,
				Severity:    types.SeverityCritical,
				Description: fmt.Sprintf("API Key found in process environment (PID: %d)", pid),
				RawContent:  key + "=***HIDDEN***",
			}
			components = append(components, comp)
		}
	}

	return components
}

// ScanPorts 扫描网络端口
func (rs *RuntimeScanner) ScanPorts() ([]types.AIComponent, error) {
	var components []types.AIComponent

	// 使用 ss 命令获取监听端口
	cmd := exec.Command("ss", "-tlnp")
	output, err := cmd.Output()
	if err != nil {
		// 尝试使用 netstat
		cmd = exec.Command("netstat", "-tlnp")
		output, err = cmd.Output()
		if err != nil {
			return components, err
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()

		// 解析端口号
		for port, service := range rs.portPatterns {
			portStr := fmt.Sprintf(":%d", port)
			if strings.Contains(line, portStr) {
				comp := types.AIComponent{
					Name:        service,
					Type:        types.TypeDeployment,
					Confidence:  0.8,
					Severity:    types.SeverityMedium,
					Description: fmt.Sprintf("Service listening on port %d", port),
					RawContent:  line,
				}
				if !rs.componentExists(components, comp) {
					components = append(components, comp)
				}
			}
		}
	}

	return components, scanner.Err()
}

// ScanDockerContainers 扫描 Docker 容器
func (rs *RuntimeScanner) ScanDockerContainers() ([]types.AIComponent, error) {
	var components []types.AIComponent

	cmd := exec.Command("docker", "ps", "--format", "{{.Names}}|{{.Image}}|{{.Ports}}|{{.Status}}")
	output, err := cmd.Output()
	if err != nil {
		// Docker 可能未安装或未运行
		return components, nil
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			continue
		}

		name := parts[0]
		image := strings.ToLower(parts[1])

		// 匹配 AI 相关镜像
		for _, keyword := range rs.containerImages {
			if strings.Contains(image, keyword) {
				comp := types.AIComponent{
					Name:        fmt.Sprintf("Docker: %s", keyword),
					Type:        types.TypeDeployment,
					Confidence:  0.9,
					Severity:    types.SeverityMedium,
					Description: fmt.Sprintf("Container: %s (Image: %s)", name, image),
					RawContent:  line,
				}
				if !rs.componentExists(components, comp) {
					components = append(components, comp)
				}
				break
			}
		}
	}

	return components, scanner.Err()
}

// ScanAll 执行所有运行时扫描
func (rs *RuntimeScanner) ScanAll() (*types.RuntimeScanResult, error) {
	result := &types.RuntimeScanResult{
		ScanTime:   time.Now().Format(time.RFC3339),
		Components: []types.AIComponent{},
		Errors:     []string{},
	}

	// 扫描进程
	if comps, err := rs.ScanProcesses(); err == nil {
		result.Components = append(result.Components, comps...)
		result.ProcessCount = len(comps)
	} else {
		result.Errors = append(result.Errors, fmt.Sprintf("Process scan error: %v", err))
	}

	// 扫描端口
	if comps, err := rs.ScanPorts(); err == nil {
		result.Components = append(result.Components, comps...)
		result.PortCount = len(comps)
	} else {
		result.Errors = append(result.Errors, fmt.Sprintf("Port scan error: %v", err))
	}

	// 扫描 Docker
	if comps, err := rs.ScanDockerContainers(); err == nil {
		result.Components = append(result.Components, comps...)
		result.ContainerCount = len(comps)
	} else {
		result.Errors = append(result.Errors, fmt.Sprintf("Docker scan error: %v", err))
	}

	return result, nil
}

// componentExists 检查组件是否已存在
func (rs *RuntimeScanner) componentExists(components []types.AIComponent, new types.AIComponent) bool {
	for _, c := range components {
		if c.Name == new.Name && c.FilePath == new.FilePath {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}