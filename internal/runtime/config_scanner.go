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

	"github.com/wz200210/AIGenerateProject/internal/config"
	"github.com/wz200210/AIGenerateProject/pkg/ai/types"
)

// ConfigBasedScanner 基于配置的扫描器
type ConfigBasedScanner struct {
	configLoader    *config.Loader
	serviceConfigs  []config.ServiceConfig
	apiKeyPatterns  []config.APIKeyPattern
	globalConfig    config.GlobalConfig
	
	// 编译后的正则表达式缓存
	processPatterns map[string]*regexp.Regexp
}

// NewConfigBasedScanner 创建基于配置的扫描器
func NewConfigBasedScanner(configPath string) (*ConfigBasedScanner, error) {
	loader := config.NewLoader(configPath)
	if err := loader.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	cs := &ConfigBasedScanner{
		configLoader:    loader,
		serviceConfigs:  loader.GetAllServices(),
		apiKeyPatterns:  loader.GetAPIKeyPatterns(),
		globalConfig:    loader.GetGlobalConfig(),
		processPatterns: make(map[string]*regexp.Regexp),
	}

	// 预编译正则表达式
	cs.compilePatterns()

	return cs, nil
}

// compilePatterns 编译所有正则表达式
func (cs *ConfigBasedScanner) compilePatterns() {
	for _, svc := range cs.serviceConfigs {
		for _, pattern := range svc.ProcessPatterns {
			if _, exists := cs.processPatterns[pattern]; !exists {
				if re, err := regexp.Compile(`(?i)` + pattern); err == nil {
					cs.processPatterns[pattern] = re
				}
			}
		}
	}
}

// ScanAll 执行完整扫描
func (cs *ConfigBasedScanner) ScanAll() (*types.RuntimeScanResult, error) {
	result := &types.RuntimeScanResult{
		ScanTime:   time.Now().Format(time.RFC3339),
		Components: []types.AIComponent{},
		Errors:     []string{},
	}

	fmt.Println("🔍 Scanning process tree...")
	processes, err := cs.scanProcessTree()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Process scan error: %v", err))
	}

	// 分析每个进程
	for _, proc := range processes {
		if component := cs.analyzeProcess(proc); component != nil {
			result.Components = append(result.Components, *component)
			result.ProcessCount++
		}
	}

	fmt.Println("🔍 Scanning network services...")
	portComponents, err := cs.scanNetworkServices(processes)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Network scan error: %v", err))
	}
	result.Components = append(result.Components, portComponents...)
	result.PortCount = len(portComponents)

	fmt.Println("🔍 Scanning Docker containers...")
	containerComponents, err := cs.scanDockerContainers()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Docker scan error: %v", err))
	}
	result.Components = append(result.Components, containerComponents...)
	result.ContainerCount = len(containerComponents)

	// 智能去重
	result.Components = cs.deduplicateComponents(result.Components)

	return result, nil
}

// scanProcessTree 扫描进程树
func (cs *ConfigBasedScanner) scanProcessTree() (map[int]*ProcessInfo, error) {
	processes := make(map[int]*ProcessInfo)

	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	// 收集所有进程
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		proc := cs.readProcessInfo(pid)
		if proc != nil {
			processes[pid] = proc
		}
	}

	// 建立父子关系
	for _, proc := range processes {
		if parent, ok := processes[proc.PPID]; ok {
			proc.Parent = parent
			parent.Children = append(parent.Children, proc)
		}
	}

	return processes, nil
}

// readProcessInfo 读取进程信息
func (cs *ConfigBasedScanner) readProcessInfo(pid int) *ProcessInfo {
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

	// 读取状态
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

	// 读取可执行文件
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

	proc.StartTime = time.Now()
	return proc
}

// analyzeProcess 分析进程
func (cs *ConfigBasedScanner) analyzeProcess(proc *ProcessInfo) *types.AIComponent {
	cmdlineLower := strings.ToLower(proc.Cmdline)

	for _, svc := range cs.serviceConfigs {
		for _, pattern := range svc.ProcessPatterns {
			if re, ok := cs.processPatterns[pattern]; ok && re.MatchString(cmdlineLower) {
				return cs.createComponentFromService(proc, svc)
			}
		}
	}

	return nil
}

// createComponentFromService 从服务配置创建组件
func (cs *ConfigBasedScanner) createComponentFromService(proc *ProcessInfo, svc config.ServiceConfig) *types.AIComponent {
	comp := &types.AIComponent{
		Name:        svc.Name,
		Type:        types.AIComponentType(svc.Type),
		FilePath:    fmt.Sprintf("/proc/%d", proc.PID),
		Confidence:  cs.calculateConfidence(proc, svc),
		Severity:    types.Severity(svc.Severity),
		Description: cs.buildDescription(proc, svc),
	}

	// 提取版本
	version := cs.extractVersion(proc, svc)
	comp.Version = version

	// 检测 API Key
	cs.detectAPIKeys(proc, comp)

	return comp
}

// extractVersion 提取版本号
func (cs *ConfigBasedScanner) extractVersion(proc *ProcessInfo, svc config.ServiceConfig) string {
	if svc.VersionProbe == nil {
		return ""
	}

	for _, method := range svc.VersionProbe.Methods {
		switch method.Type {
		case "cli_arg":
			for _, pattern := range method.Patterns {
				if re, err := regexp.Compile(pattern); err == nil {
					if matches := re.FindStringSubmatch(proc.Cmdline); len(matches) > 1 {
						return matches[1]
					}
				}
			}

		case "exec":
			if proc.Executable != "" {
				version := cs.executeVersionCommand(proc.Executable, method.Command, method.Parser)
				if version != "" {
					return version
				}
			}

		case "http_api":
			for _, port := range svc.DefaultPorts {
				if version := cs.probeHTTPVersion(port, method.Endpoint, method.JSONPath); version != "" {
					return version
				}
			}
		}
	}

	// 从环境变量检查
	for _, envVar := range []string{"VERSION", "APP_VERSION", "SERVICE_VERSION"} {
		if version, ok := proc.Environment[envVar]; ok && version != "" {
			return version
		}
	}

	return ""
}

// executeVersionCommand 执行版本命令
func (cs *ConfigBasedScanner) executeVersionCommand(executable, command, parser string) string {
	args := []string{command}
	cmd := exec.Command(executable, args...)
	cmd.Env = []string{}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	if parser != "" {
		if re, err := regexp.Compile(parser); err == nil {
			if matches := re.FindStringSubmatch(string(output)); len(matches) > 1 {
				return matches[1]
			}
		}
	}

	// 默认解析
	return cs.parseVersionOutput(string(output))
}

// parseVersionOutput 解析版本输出
func (cs *ConfigBasedScanner) parseVersionOutput(output string) string {
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

// probeHTTPVersion HTTP 探测版本
func (cs *ConfigBasedScanner) probeHTTPVersion(port int, endpoint, jsonPath string) string {
	url := fmt.Sprintf("http://localhost:%d%s", port, endpoint)
	
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if jsonPath != "" {
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err == nil {
			// 简单支持一级路径
			if version, ok := result[jsonPath].(string); ok {
				return version
			}
		}
	}

	return cs.parseVersionOutput(string(body))
}

// calculateConfidence 计算置信度
func (cs *ConfigBasedScanner) calculateConfidence(proc *ProcessInfo, svc config.ServiceConfig) float64 {
	weights := cs.globalConfig.ConfidenceWeights
	confidence := 0.5

	// 进程名匹配
	for _, pattern := range svc.ProcessPatterns {
		if strings.Contains(strings.ToLower(proc.Name), pattern) {
			confidence += weights["process_name_match"]
			break
		}
	}

	// 有监听端口
	if len(svc.DefaultPorts) > 0 {
		confidence += weights["has_listening_port"]
	}

	// 有 API 环境变量
	for _, envVar := range svc.EnvIndicators {
		if _, ok := proc.Environment[envVar]; ok {
			confidence += weights["has_api_env_var"]
			break
		}
	}

	// 存活时间
	if time.Since(proc.StartTime) > time.Minute {
		confidence += weights["long_uptime"]
	}

	// 父进程检查
	if proc.Parent != nil {
		parentLower := strings.ToLower(proc.Parent.Cmdline)
		for _, pattern := range svc.ProcessPatterns {
			if strings.Contains(parentLower, pattern) {
				confidence += weights["parent_is_ai"]
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
func (cs *ConfigBasedScanner) buildDescription(proc *ProcessInfo, svc config.ServiceConfig) string {
	parts := []string{fmt.Sprintf("PID: %d", proc.PID)}

	if len(svc.DefaultPorts) > 0 {
		portsStr := make([]string, len(svc.DefaultPorts))
		for i, port := range svc.DefaultPorts {
			portsStr[i] = strconv.Itoa(port)
		}
		parts = append(parts, fmt.Sprintf("Default Ports: %s", strings.Join(portsStr, ", ")))
	}

	if proc.Executable != "" {
		parts = append(parts, fmt.Sprintf("Exe: %s", filepath.Base(proc.Executable)))
	}

	return strings.Join(parts, " | ")
}

// detectAPIKeys 检测 API Key
func (cs *ConfigBasedScanner) detectAPIKeys(proc *ProcessInfo, comp *types.AIComponent) {
	envStr := ""
	for k, v := range proc.Environment {
		envStr += k + "=" + v + " "
	}

	for _, pattern := range cs.apiKeyPatterns {
		if strings.Contains(envStr, pattern.Key) {
			comp.Severity = types.Severity(pattern.Severity)
			if comp.Description != "" {
				comp.Description += fmt.Sprintf(" | ⚠️ %s detected", pattern.Name)
			}
			break
		}
	}
}

// scanNetworkServices 扫描网络服务 - 优化版：端口+进程名双重验证
func (cs *ConfigBasedScanner) scanNetworkServices(processes map[int]*ProcessInfo) ([]types.AIComponent, error) {
	var components []types.AIComponent
	
	// 获取端口到 PID 的映射
	portToPID := cs.getPortToPIDMapping()
	
	// 获取所有监听端口的列表（用于日志输出）
	listeningPorts := cs.getListeningPorts()
	
	// 遍历所有服务配置
	for _, svc := range cs.serviceConfigs {
		// 跳过没有配置端口的服务
		if len(svc.DefaultPorts) == 0 {
			continue
		}
		
		// 检查该服务的端口是否有进程在监听
		for _, port := range svc.DefaultPorts {
			pid, ok := portToPID[port]
			if !ok {
				// 端口没有在监听，跳过
				continue
			}
			
			// 获取监听该端口的进程
			proc, ok := processes[pid]
			if !ok {
				// 进程信息获取失败，记录错误但继续
				continue
			}
			
			// 关键：进程名必须匹配服务的 ProcessPatterns 才算
			if !cs.matchProcessPatterns(proc, svc) {
				// 进程名不匹配，可能是其他服务占用了这个端口，跳过
				continue
			}
			
			// 进程名匹配成功，创建组件
			comp := cs.createComponentFromNetworkService(proc, svc, port)
			if comp != nil {
				components = append(components, *comp)
			}
		}
	}
	
	// 记录发现的监听端口信息（用于调试）
	_ = listeningPorts
	
	return components, nil
}

// matchProcessPatterns 检查进程是否匹配服务的进程模式
func (cs *ConfigBasedScanner) matchProcessPatterns(proc *ProcessInfo, svc config.ServiceConfig) bool {
	// 检查进程名
	procNameLower := strings.ToLower(proc.Name)
	cmdlineLower := strings.ToLower(proc.Cmdline)
	exeLower := strings.ToLower(filepath.Base(proc.Executable))
	
	for _, pattern := range svc.ProcessPatterns {
		if re, ok := cs.processPatterns[pattern]; ok {
			// 匹配进程名
			if re.MatchString(procNameLower) {
				return true
			}
			// 匹配命令行
			if re.MatchString(cmdlineLower) {
				return true
			}
			// 匹配可执行文件名
			if re.MatchString(exeLower) {
				return true
			}
		}
	}
	
	return false
}

// createComponentFromNetworkService 从网络服务创建组件（带版本验证）
// 只有当版本号能获取到时才返回组件，否则返回 nil
func (cs *ConfigBasedScanner) createComponentFromNetworkService(proc *ProcessInfo, svc config.ServiceConfig, port int) *types.AIComponent {
	// 提取版本号 - 关键：版本号必须能获取到才算匹配
	version := cs.extractVersion(proc, svc)
	
	// 如果无法获取版本号，则视为未匹配成功，避免误报
	if version == "" {
		return nil
	}
	
	// 计算置信度
	confidence := cs.calculateConfidence(proc, svc)
	// 网络服务检测有端口+进程名+版本号三重验证，置信度更高
	confidence += 0.2
	if confidence > 1.0 {
		confidence = 1.0
	}
	
	// 构建描述
	description := fmt.Sprintf("Network service detected | PID: %d | Port: %d", proc.PID, port)
	if proc.Executable != "" {
		description += fmt.Sprintf(" | Exe: %s", filepath.Base(proc.Executable))
	}
	if version != "" {
		description += fmt.Sprintf(" | Version: %s", version)
	}
	
	comp := &types.AIComponent{
		Name:        svc.Name,
		Type:        types.AIComponentType(svc.Type),
		Version:     version,
		FilePath:    fmt.Sprintf("/proc/%d (port %d)", proc.PID, port),
		Confidence:  confidence,
		Severity:    types.Severity(svc.Severity),
		Description: description,
		RawContent:  fmt.Sprintf("Process: %s | Cmdline: %s", proc.Name, proc.Cmdline),
	}
	
	// 检测 API Key
	cs.detectAPIKeys(proc, comp)
	
	return comp
}

// getListeningPorts 获取所有正在监听的端口列表
func (cs *ConfigBasedScanner) getListeningPorts() []int {
	var ports []int
	
	for _, file := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		
		lines := strings.Split(string(data), "\n")
		for _, line := range lines[1:] {
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}
			
			// 检查是否为监听状态 (0x0A = TCP_LISTEN)
			state := fields[3]
			if state != "0A" {
				continue
			}
			
			localAddr := fields[1]
			if idx := strings.LastIndex(localAddr, ":"); idx > 0 {
				portHex := localAddr[idx+1:]
				port, _ := strconv.ParseInt(portHex, 16, 32)
				if port > 0 {
					ports = append(ports, int(port))
				}
			}
		}
	}
	
	return ports
}

// getPortToPIDMapping 获取端口到 PID 映射
func (cs *ConfigBasedScanner) getPortToPIDMapping() map[int]int {
	mapping := make(map[int]int)

	for _, file := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines[1:] {
			fields := strings.Fields(line)
			if len(fields) < 10 {
				continue
			}

			localAddr := fields[1]
			if idx := strings.LastIndex(localAddr, ":"); idx > 0 {
				portHex := localAddr[idx+1:]
				port, _ := strconv.ParseInt(portHex, 16, 32)
				inode := fields[9]

				if inode != "0" {
					if pid := cs.findPIDByInode(inode); pid > 0 {
						mapping[int(port)] = pid
					}
				}
			}
		}
	}

	return mapping
}

// findPIDByInode 通过 inode 查找 PID
func (cs *ConfigBasedScanner) findPIDByInode(inode string) int {
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

// scanDockerContainers 扫描 Docker 容器
func (cs *ConfigBasedScanner) scanDockerContainers() ([]types.AIComponent, error) {
	var components []types.AIComponent

	cmd := exec.Command("docker", "ps", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.Ports}}|{{.Status}}")
	output, err := cmd.Output()
	if err != nil {
		return components, nil
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			continue
		}

		containerName := parts[1]
		image := strings.ToLower(parts[2])
		_ = parts[3] // ports 未使用
		status := parts[4]

		// 提取版本
		version := ""
		if re := regexp.MustCompile(`:([\d.]+[\w.-]*)`); re != nil {
			if matches := re.FindStringSubmatch(image); len(matches) > 1 {
				version = matches[1]
			}
		}

		// 匹配服务
		for _, svc := range cs.serviceConfigs {
			for _, pattern := range svc.ProcessPatterns {
				if re, ok := cs.processPatterns[pattern]; ok && re.MatchString(image) {
					comp := types.AIComponent{
						Name:        svc.Name,
						Type:        types.AIComponentType(svc.Type),
						Version:     version,
						Confidence:  0.9,
						Severity:    types.Severity(svc.Severity),
						Description: fmt.Sprintf("Docker: %s (Image: %s, Status: %s)", containerName, image, status),
						RawContent:  line,
					}
					components = append(components, comp)
					goto nextContainer
				}
			}
		}
	nextContainer:
	}

	return components, scanner.Err()
}

// deduplicateComponents 去重
func (cs *ConfigBasedScanner) deduplicateComponents(components []types.AIComponent) []types.AIComponent {
	seen := make(map[string]bool)
	var unique []types.AIComponent

	for _, comp := range components {
		key := fmt.Sprintf("%s|%s|%s", comp.Name, comp.Type, comp.FilePath)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, comp)
		}
	}

	return unique
}

// ProcessInfo 进程信息（与之前兼容）
type ProcessInfo struct {
	PID         int
	PPID        int
	Name        string
	Cmdline     string
	Executable  string
	Environment map[string]string
	Ports       []int
	Parent      *ProcessInfo
	Children    []*ProcessInfo
	StartTime   time.Time
}
