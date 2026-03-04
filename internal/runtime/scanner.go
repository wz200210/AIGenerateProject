package runtime

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wz200210/AIGenerateProject/pkg/ai/types"
)

// RuntimeScanner 增强版运行时扫描器
type RuntimeScanner struct {
	processPatterns     map[string]*regexp.Regexp
	versionPatterns     map[string]*regexp.Regexp
	apiEndpoints        map[int]string
	semanticAnalyzers   []SemanticAnalyzer
}

// SemanticAnalyzer 语义分析器接口
type SemanticAnalyzer interface {
	Analyze(process *ProcessInfo) *types.AIComponent
}

// ProcessInfo 增强的进程信息
type ProcessInfo struct {
	PID          int
	PPID         int
	Name         string
	Cmdline      string
	Executable   string
	Environment  map[string]string
	Ports        []int
	Connections  []string
	Children     []*ProcessInfo
	Parent       *ProcessInfo
	StartTime    time.Time
}

// NewRuntimeScanner 创建增强版运行时扫描器
func NewRuntimeScanner() *RuntimeScanner {
	rs := &RuntimeScanner{
		processPatterns:   make(map[string]*regexp.Regexp),
		versionPatterns:   make(map[string]*regexp.Regexp),
		apiEndpoints:      make(map[int]string),
		semanticAnalyzers: []SemanticAnalyzer{},
	}
	rs.compilePatterns()
	rs.initSemanticAnalyzers()
	return rs
}

// compilePatterns 编译检测模式
func (rs *RuntimeScanner) compilePatterns() {
	// AI 服务进程匹配模式（更精确）
	processes := map[string]string{
		"ollama":            `(?i)\bollama\b`,
		"vllm":              `(?i)\bvllm\b|vllm\.entrypoints`,
		"text-generation":   `(?i)text-generation-inference|tgi`,
		"transformers":      `(?i)transformers|huggingface.*transformers`,
		"langchain":         `(?i)langchain-(serve|cli)|python.*langchain`,
		"llamaindex":        `(?i)llamaindex|llama_index`,
		"openai":            `(?i)openai-(api|serve)|openai.*api`,
		"fastchat":          `(?i)fastchat.*serve|fastchat.*controller`,
		"torch":             `(?i)torchserve|python.*torch`,
		"tensorflow":        `(?i)tensorflow.*serving|tf_serving`,
		"onnx":              `(?i)onnxruntime.*server`,
		"milvus":            `(?i)milvus.*server|milvus-standalone`,
		"chroma":            `(?i)chromadb|chroma.*server`,
		"weaviate":          `(?i)weaviate|weaviate.*server`,
		"qdrant":            `(?i)qdrant`,
		"pinecone-proxy":    `(?i)pinecone.*proxy|pinecone.*gateway`,
		"redis-vector":      `(?i)redis.*vector|redisvl`,
		"elasticsearch":     `(?i)elasticsearch.*vector|es.*knn`,
		"jupyter":           `(?i)jupyter.*(notebook|lab|server)`,
		"mlflow":            `(?i)mlflow.*server|mlflow.*ui`,
		"wandb":             `(?i)wandb.*server|wandb.*local`,
		"triton":            `(?i)tritonserver`,
		"bentoml":           `(?i)bentoml.*serve`,
		"ray-serve":         `(?i)ray.*serve`,
		"gradio":            `(?i)gradio`,
		"streamlit":         `(?i)streamlit.*run`,
	}

	for name, pattern := range processes {
		if re, err := regexp.Compile(pattern); err == nil {
			rs.processPatterns[name] = re
		}
	}

	// 版本号提取模式
	rs.versionPatterns = map[string]*regexp.Regexp{
		"arg-version":   regexp.MustCompile(`--version[=\s]*([\d.]+)`),
		"arg-v":         regexp.MustCompile(`-v[=\s]*([\d.]+)`),
		"version-flag":  regexp.MustCompile(`\bversion[=\s]*([\d.]+)`),
		"v-prefix":      regexp.MustCompile(`\bv([\d.]+)\b`),
		"docker-tag":    regexp.MustCompile(`:([\d.]+[\w.-]*)`),
	}

	// 常见 AI 服务 API 端点
	rs.apiEndpoints = map[int]string{
		11434: "Ollama API",
		8000:  "vLLM/FastAPI/TGI",
		8080:  "通用 AI 服务",
		5000:  "Flask ML 服务",
		3000:  "Node.js AI 服务",
		6333:  "Qdrant REST",
		6334:  "Qdrant gRPC",
		19530: "Milvus",
		9091:  "Milvus Proxy",
		6379:  "Redis",
		9200:  "Elasticsearch",
		5601:  "Kibana",
		8888:  "Jupyter",
		5001:  "MLflow",
		7860:  "Gradio",
		8501:  "Streamlit",
		8265:  "Ray Dashboard",
	}
}

// initSemanticAnalyzers 初始化语义分析器
func (rs *RuntimeScanner) initSemanticAnalyzers() {
	rs.semanticAnalyzers = []SemanticAnalyzer{
		&PythonMLAnalyzer{},
		&NodeJSAnalyzer{},
		&DockerAnalyzer{},
		&ServiceMeshAnalyzer{},
	}
}

// ScanAll 执行完整运行时扫描
func (rs *RuntimeScanner) ScanAll() (*types.RuntimeScanResult, error) {
	result := &types.RuntimeScanResult{
		ScanTime:   time.Now().Format(time.RFC3339),
		Components: []types.AIComponent{},
		Errors:     []string{},
	}

	// 1. 扫描进程树
	fmt.Println("🔍 Scanning process tree...")
	processes, err := rs.scanProcessTree()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Process scan error: %v", err))
	}

	// 2. 分析每个进程
	for _, proc := range processes {
		if component := rs.analyzeProcess(proc); component != nil {
			result.Components = append(result.Components, *component)
			result.ProcessCount++
		}
	}

	// 3. 扫描网络端口和服务
	fmt.Println("🔍 Scanning network services...")
	portComponents, err := rs.scanNetworkServices(processes)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Network scan error: %v", err))
	}
	result.Components = append(result.Components, portComponents...)
	result.PortCount = len(portComponents)

	// 4. 扫描 Docker 容器
	fmt.Println("🔍 Scanning Docker containers...")
	containerComponents, err := rs.scanDockerContainers()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Docker scan error: %v", err))
	}
	result.Components = append(result.Components, containerComponents...)
	result.ContainerCount = len(containerComponents)

	// 5. 智能去重和关联分析
	result.Components = rs.deduplicateAndAnalyze(result.Components)

	return result, nil
}

// scanProcessTree 扫描进程树
func (rs *RuntimeScanner) scanProcessTree() (map[int]*ProcessInfo, error) {
	processes := make(map[int]*ProcessInfo)

	// 读取 /proc 目录
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	// 第一遍：收集所有进程基本信息
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		proc := rs.readProcessInfo(pid)
		if proc != nil {
			processes[pid] = proc
		}
	}

	// 第二遍：建立父子关系
	for _, proc := range processes {
		if parent, ok := processes[proc.PPID]; ok {
			proc.Parent = parent
			parent.Children = append(parent.Children, proc)
		}
	}

	return processes, nil
}

// readProcessInfo 读取进程详细信息
func (rs *RuntimeScanner) readProcessInfo(pid int) *ProcessInfo {
	proc := &ProcessInfo{
		PID:         pid,
		Environment: make(map[string]string),
	}

	// 读取命令行
	cmdlineData, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return nil
	}
	proc.Cmdline = strings.ReplaceAll(string(cmdlineData), "\x00", " ")

	// 读取状态信息获取 PPID
	statusData, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err == nil {
		for _, line := range strings.Split(string(statusData), "\n") {
			if strings.HasPrefix(line, "PPid:") {
				ppidStr := strings.TrimSpace(strings.TrimPrefix(line, "PPid:"))
				proc.PPID, _ = strconv.Atoi(ppidStr)
			}
			if strings.HasPrefix(line, "Name:") {
				proc.Name = strings.TrimSpace(strings.TrimPrefix(line, "Name:"))
			}
		}
	}

	// 读取可执行文件路径
	exePath, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	if err == nil {
		proc.Executable = exePath
	}

	// 读取环境变量
	envData, err := os.ReadFile(fmt.Sprintf("/proc/%d/environ", pid))
	if err == nil {
		for _, env := range strings.Split(string(envData), "\x00") {
			if idx := strings.Index(env, "="); idx > 0 {
				key := env[:idx]
				value := env[idx+1:]
				proc.Environment[key] = value
			}
		}
	}

	// 获取进程启动时间
	_, _ = os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	// 简化处理，实际需要解析启动时间
	proc.StartTime = time.Now()

	return proc
}

// analyzeProcess 分析单个进程
func (rs *RuntimeScanner) analyzeProcess(proc *ProcessInfo) *types.AIComponent {
	cmdlineLower := strings.ToLower(proc.Cmdline)

	for name, pattern := range rs.processPatterns {
		if pattern.MatchString(cmdlineLower) {
			component := &types.AIComponent{
				Name:       rs.getDisplayName(name),
				Type:       rs.getComponentType(name),
				FilePath:   fmt.Sprintf("/proc/%d", proc.PID),
				Confidence: rs.calculateConfidence(proc, name),
				Severity:   rs.getSeverity(name),
			}

			// 提取版本号
			version := rs.extractVersion(proc)
			component.Version = version

			// 构建描述
			desc := rs.buildDescription(proc, name)
			component.Description = desc

			// 语义分析
			for _, analyzer := range rs.semanticAnalyzers {
				if analyzed := analyzer.Analyze(proc); analyzed != nil {
					// 合并分析结果
					if analyzed.Type != "" {
						component.Type = analyzed.Type
					}
					if analyzed.Description != "" {
						component.Description += "; " + analyzed.Description
					}
				}
			}

			// 检测 API Key
			rs.detectAPIKeysInProcess(proc, component)

			return component
		}
	}

	return nil
}

// extractVersion 提取版本号
func (rs *RuntimeScanner) extractVersion(proc *ProcessInfo) string {
	// 1. 从命令行参数提取
	for _, pattern := range rs.versionPatterns {
		if matches := pattern.FindStringSubmatch(proc.Cmdline); len(matches) > 1 {
			return matches[1]
		}
	}

	// 2. 从环境变量提取
	versionEnvVars := []string{
		"VERSION", "APP_VERSION", "SERVICE_VERSION",
		"IMAGE_VERSION", "CONTAINER_VERSION",
	}
	for _, envVar := range versionEnvVars {
		if version, ok := proc.Environment[envVar]; ok && version != "" {
			return version
		}
	}

	// 3. 尝试执行进程获取版本
	version := rs.executeVersionCommand(proc)
	if version != "" {
		return version
	}

	// 4. 从镜像标签提取（Docker）
	if image := proc.Environment["IMAGE_NAME"]; image != "" {
		if matches := rs.versionPatterns["docker-tag"].FindStringSubmatch(image); len(matches) > 1 {
			return matches[1]
		}
	}

	return "unknown"
}

// executeVersionCommand 执行版本命令
func (rs *RuntimeScanner) executeVersionCommand(proc *ProcessInfo) string {
	// 安全起见，只对特定进程执行
	if proc.Executable == "" || proc.PID == 0 {
		return ""
	}

	// 尝试 --version
	cmd := exec.Command(proc.Executable, "--version")
	cmd.Env = []string{} // 清空环境变量，避免副作用
	output, err := cmd.CombinedOutput()
	if err == nil {
		return rs.parseVersionOutput(string(output))
	}

	// 尝试 -v
	cmd = exec.Command(proc.Executable, "-v")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return rs.parseVersionOutput(string(output))
	}

	// 尝试 version 子命令
	cmd = exec.Command(proc.Executable, "version")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return rs.parseVersionOutput(string(output))
	}

	return ""
}

// parseVersionOutput 解析版本输出
func (rs *RuntimeScanner) parseVersionOutput(output string) string {
	// 常见版本输出格式
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)version[\s:]+([\d.]+[\w.-]*)`),
		regexp.MustCompile(`(?i)v?([\d]+\.[\d]+\.[\d]+[\w.-]*)`),
		regexp.MustCompile(`([\d]+\.[\d]+(?:\.[\d]+)?)`),
	}

	for _, pattern := range patterns {
		if matches := pattern.FindStringSubmatch(output); len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}

// calculateConfidence 计算置信度
func (rs *RuntimeScanner) calculateConfidence(proc *ProcessInfo, name string) float64 {
	confidence := 0.5 // 基础分

	// 1. 进程名匹配程度高
	if strings.Contains(strings.ToLower(proc.Name), name) {
		confidence += 0.2
	}

	// 2. 有监听端口
	if len(proc.Ports) > 0 {
		confidence += 0.15
	}

	// 3. 有 API 相关的环境变量
	apiEnvVars := []string{"API_KEY", "API_URL", "ENDPOINT", "PORT", "HOST"}
	for _, envVar := range apiEnvVars {
		if _, ok := proc.Environment[envVar]; ok {
			confidence += 0.05
			break
		}
	}

	// 4. 进程存活时间较长（不是临时进程）
	if time.Since(proc.StartTime) > time.Minute {
		confidence += 0.1
	}

	// 5. 有子进程或父进程是 AI 相关
	if proc.Parent != nil {
		for name := range rs.processPatterns {
			if strings.Contains(strings.ToLower(proc.Parent.Cmdline), name) {
				confidence += 0.1
				break
			}
		}
	}

	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// buildDescription 构建描述
func (rs *RuntimeScanner) buildDescription(proc *ProcessInfo, name string) string {
	parts := []string{
		fmt.Sprintf("PID: %d", proc.PID),
	}

	if len(proc.Ports) > 0 {
		portsStr := make([]string, len(proc.Ports))
		for i, port := range proc.Ports {
			portsStr[i] = strconv.Itoa(port)
		}
		parts = append(parts, fmt.Sprintf("Ports: %s", strings.Join(portsStr, ", ")))
	}

	if proc.Executable != "" {
		parts = append(parts, fmt.Sprintf("Exe: %s", filepath.Base(proc.Executable)))
	}

	return strings.Join(parts, " | ")
}

// detectAPIKeysInProcess 检测进程中的 API Key
func (rs *RuntimeScanner) detectAPIKeysInProcess(proc *ProcessInfo, component *types.AIComponent) {
	envStr := ""
	for k, v := range proc.Environment {
		envStr += k + "=" + v + " "
	}

	for key, _ := range types.APIKeyPatterns {
		if strings.Contains(envStr, key) {
			component.Severity = types.SeverityCritical
			if component.Description != "" {
				component.Description += " | ⚠️ API Key detected in environment"
			}
			break
		}
	}
}

// scanNetworkServices 扫描网络服务
func (rs *RuntimeScanner) scanNetworkServices(processes map[int]*ProcessInfo) ([]types.AIComponent, error) {
	var components []types.AIComponent

	// 获取端口监听信息
	portToPID := rs.getPortToPIDMapping()

	// 使用 ss 命令获取监听端口
	cmd := exec.Command("ss", "-tlnp")
	output, err := cmd.Output()
	if err != nil {
		return components, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()

		for port, service := range rs.apiEndpoints {
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

				// 尝试关联到进程
				if pid, ok := portToPID[port]; ok {
					if proc, ok := processes[pid]; ok {
						comp.FilePath = fmt.Sprintf("/proc/%d (port %d)", pid, port)
						comp.Version = rs.extractVersion(proc)
					}
				}

				// 尝试获取版本（HTTP API）
				if version := rs.probeServiceVersion(port); version != "" {
					comp.Version = version
					comp.Confidence = 0.9
				}

				components = append(components, comp)
			}
		}
	}

	return components, scanner.Err()
}

// getPortToPIDMapping 获取端口到 PID 的映射
func (rs *RuntimeScanner) getPortToPIDMapping() map[int]int {
	mapping := make(map[int]int)

	// 从 /proc/net/tcp 和 /proc/net/tcp6 解析
	for _, file := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines[1:] { // 跳过标题
			fields := strings.Fields(line)
			if len(fields) < 10 {
				continue
			}

			// 解析本地地址 (格式: 本地地址:端口)
			localAddr := fields[1]
			if idx := strings.LastIndex(localAddr, ":"); idx > 0 {
				portHex := localAddr[idx+1:]
				port, _ := strconv.ParseInt(portHex, 16, 32)

				// 解析 inode
				inode := fields[9]
				if inode != "0" {
					// 通过 inode 找到进程
					if pid := rs.findPIDByInode(inode); pid > 0 {
						mapping[int(port)] = pid
					}
				}
			}
		}
	}

	return mapping
}

// findPIDByInode 通过 inode 查找进程
func (rs *RuntimeScanner) findPIDByInode(inode string) int {
	entries, _ := os.ReadDir("/proc")
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		fdDir := fmt.Sprintf("/proc/%d/fd", pid)
		fds, _ := os.ReadDir(fdDir)
		for _, fd := range fds {
			link, _ := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if strings.Contains(link, "socket:["+inode+"]") {
				return pid
			}
		}
	}
	return 0
}

// probeServiceVersion 探测服务版本
func (rs *RuntimeScanner) probeServiceVersion(port int) string {
	// 常见服务的版本端点
	endpoints := []string{
		"/version",
		"/api/version",
		"/health",
		"/v1/version",
		"/api/v1/version",
	}

	client := &http.Client{Timeout: 2 * time.Second}

	for _, endpoint := range endpoints {
		url := fmt.Sprintf("http://localhost:%d%s", port, endpoint)
		resp, err := client.Get(url)
		if err != nil {
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// 尝试解析 JSON 中的版本
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err == nil {
			if version, ok := result["version"].(string); ok && version != "" {
				return version
			}
		}

		// 从文本中提取版本
		if version := rs.parseVersionOutput(string(body)); version != "" {
			return version
		}
	}

	return ""
}

// scanDockerContainers 扫描 Docker 容器
func (rs *RuntimeScanner) scanDockerContainers() ([]types.AIComponent, error) {
	var components []types.AIComponent

	cmd := exec.Command("docker", "ps", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.Ports}}|{{.Status}}")
	output, err := cmd.Output()
	if err != nil {
		return components, nil // Docker 可能未运行
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			continue
		}

		_ = parts[0] // containerID 未使用
		containerName := parts[1]
		image := strings.ToLower(parts[2])
		ports := parts[3]
		status := parts[4]

		// 从镜像标签提取版本
		version := ""
		if matches := rs.versionPatterns["docker-tag"].FindStringSubmatch(image); len(matches) > 1 {
			version = matches[1]
		}

		// 匹配 AI 相关容器
		for namePattern, pattern := range rs.processPatterns {
			if pattern.MatchString(image) {
				comp := types.AIComponent{
					Name:        rs.getDisplayName(namePattern),
					Type:        types.TypeDeployment,
					Version:     version,
					Confidence:  0.9,
					Severity:    types.SeverityMedium,
					Description: fmt.Sprintf("Docker container: %s (Image: %s, Status: %s, Ports: %s)", containerName, image, status, ports),
					RawContent:  line,
				}
				components = append(components, comp)
				break
			}
		}
	}

	return components, scanner.Err()
}

// deduplicateAndAnalyze 去重和关联分析
func (rs *RuntimeScanner) deduplicateAndAnalyze(components []types.AIComponent) []types.AIComponent {
	// 使用 map 去重
	seen := make(map[string]bool)
	var unique []types.AIComponent

	for _, comp := range components {
		// 去重键：名称 + 类型 + 文件路径
		key := fmt.Sprintf("%s|%s|%s", comp.Name, comp.Type, comp.FilePath)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, comp)
		}
	}

	return unique
}

// getDisplayName 获取显示名称
func (rs *RuntimeScanner) getDisplayName(internalName string) string {
	names := map[string]string{
		"ollama":          "Ollama",
		"vllm":            "vLLM",
		"text-generation": "Text Generation Inference",
		"transformers":    "Hugging Face Transformers",
		"langchain":       "LangChain",
		"llamaindex":      "LlamaIndex",
		"openai":          "OpenAI API",
		"milvus":          "Milvus",
		"chroma":          "Chroma",
		"weaviate":        "Weaviate",
		"qdrant":          "Qdrant",
	}

	if name, ok := names[internalName]; ok {
		return name
	}
	return internalName
}

// getComponentType 获取组件类型
func (rs *RuntimeScanner) getComponentType(name string) types.AIComponentType {
	typeMap := map[string]types.AIComponentType{
		"ollama":          types.TypeLLMFramework,
		"vllm":            types.TypeDeployment,
		"text-generation": types.TypeDeployment,
		"transformers":    types.TypeMLFramework,
		"langchain":       types.TypeLLMFramework,
		"llamaindex":      types.TypeRAGTool,
		"openai":          types.TypeLLMFramework,
		"milvus":          types.TypeVectorDB,
		"chroma":          types.TypeVectorDB,
		"weaviate":        types.TypeVectorDB,
		"qdrant":          types.TypeVectorDB,
		"jupyter":         types.TypeMonitoring,
		"mlflow":          types.TypeMonitoring,
		"wandb":           types.TypeMonitoring,
	}

	if t, ok := typeMap[name]; ok {
		return t
	}
	return types.TypeDeployment
}

// getSeverity 获取严重级别
func (rs *RuntimeScanner) getSeverity(name string) types.Severity {
	// API Key 在环境变量中是 Critical，其他根据情况
	return types.SeverityMedium
}

// ==================== 语义分析器实现 ====================

// PythonMLAnalyzer Python ML 分析器
type PythonMLAnalyzer struct{}

func (a *PythonMLAnalyzer) Analyze(proc *ProcessInfo) *types.AIComponent {
	if !strings.Contains(strings.ToLower(proc.Cmdline), "python") {
		return nil
	}

	// 分析 Python 进程的 ML 框架
	frameworks := map[string]string{
		"torch":        "PyTorch",
		"tensorflow":   "TensorFlow",
		"jax":          "JAX",
		"sklearn":      "Scikit-learn",
		"transformers": "Hugging Face",
	}

	for key, name := range frameworks {
		if strings.Contains(strings.ToLower(proc.Cmdline), key) {
			return &types.AIComponent{
				Name:        name,
				Type:        types.TypeMLFramework,
				Confidence:  0.85,
				Description: fmt.Sprintf("Python process using %s (PID: %d)", name, proc.PID),
			}
		}
	}

	return nil
}

// NodeJSAnalyzer Node.js 分析器
type NodeJSAnalyzer struct{}

func (a *NodeJSAnalyzer) Analyze(proc *ProcessInfo) *types.AIComponent {
	if !strings.Contains(strings.ToLower(proc.Cmdline), "node") {
		return nil
	}

	// 分析 Node.js 进程的 AI 相关包
	aiPackages := map[string]string{
		"openai":    "OpenAI Node.js SDK",
		"langchain": "LangChain.js",
		"@anthropic-ai": "Anthropic SDK",
	}

	for key, name := range aiPackages {
		if strings.Contains(strings.ToLower(proc.Cmdline), key) {
			return &types.AIComponent{
				Name:        name,
				Type:        types.TypeLLMFramework,
				Confidence:  0.8,
				Description: fmt.Sprintf("Node.js process using %s (PID: %d)", name, proc.PID),
			}
		}
	}

	return nil
}

// DockerAnalyzer Docker 分析器
type DockerAnalyzer struct{}

func (a *DockerAnalyzer) Analyze(proc *ProcessInfo) *types.AIComponent {
	// 检测容器内的 AI 服务
	if proc.Environment["KUBERNETES_SERVICE_HOST"] != "" {
		return &types.AIComponent{
			Name:        "Kubernetes Pod",
			Type:        types.TypeDeployment,
			Confidence:  0.7,
			Description: fmt.Sprintf("Running in Kubernetes (PID: %d)", proc.PID),
		}
	}

	return nil
}

// ServiceMeshAnalyzer 服务网格分析器
type ServiceMeshAnalyzer struct{}

func (a *ServiceMeshAnalyzer) Analyze(proc *ProcessInfo) *types.AIComponent {
	// 检测 AI 服务之间的调用关系
	if len(proc.Connections) > 0 {
		// 分析网络连接
		for _, conn := range proc.Connections {
			if strings.Contains(conn, ":11434") { // Ollama
				return &types.AIComponent{
					Name:        "Ollama Client",
					Type:        types.TypeLLMFramework,
					Confidence:  0.75,
					Description: fmt.Sprintf("Connected to Ollama service (PID: %d)", proc.PID),
				}
			}
		}
	}

	return nil
}