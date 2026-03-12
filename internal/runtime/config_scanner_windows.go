//go:build windows

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

// scanProcessTree 扫描进程树 (Windows 实现)
func (cs *ConfigBasedScanner) scanProcessTree() (map[int]*ProcessInfo, error) {
	processes := make(map[int]*ProcessInfo)

	// 使用 PowerShell 获取进程信息
	cmd := exec.Command("powershell", "-Command",
		`Get-Process | Select-Object Id, ParentId, ProcessName, Path, @{Name="CommandLine";Expression={$_.CommandLine}} | ConvertTo-Json`)
	
	output, err := cmd.Output()
	if err != nil {
		// 降级方案：使用 tasklist
		return cs.scanProcessTreeTasklist()
	}

	// 解析 JSON 输出
	var procList []struct {
		ID           int    `json:"Id"`
		ParentId     int    `json:"ParentId"`
		ProcessName  string `json:"ProcessName"`
		Path         string `json:"Path"`
		CommandLine  string `json:"CommandLine"`
	}

	if err := json.Unmarshal(output, &procList); err != nil {
		// 可能是单个对象，包装成数组
		var singleProc struct {
			ID           int    `json:"Id"`
			ParentId     int    `json:"ParentId"`
			ProcessName  string `json:"ProcessName"`
			Path         string `json:"Path"`
			CommandLine  string `json:"CommandLine"`
		}
		if err := json.Unmarshal(output, &singleProc); err == nil {
			procList = append(procList, singleProc)
		} else {
			// 降级到 tasklist
			return cs.scanProcessTreeTasklist()
		}
	}

	for _, p := range procList {
		proc := &ProcessInfo{
			PID:         p.ID,
			PPID:        p.ParentId,
			Name:        p.ProcessName,
			Cmdline:     p.CommandLine,
			Executable:  p.Path,
			Environment: make(map[string]string),
			StartTime:   time.Now(),
		}
		processes[p.ID] = proc
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

// scanProcessTreeTasklist 使用 tasklist 作为降级方案
func (cs *ConfigBasedScanner) scanProcessTreeTasklist() (map[int]*ProcessInfo, error) {
	processes := make(map[int]*ProcessInfo)

	// 使用 tasklist 获取进程列表
	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run tasklist: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := parseCSVLine(line)
		if len(parts) < 2 {
			continue
		}

		pid, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		if pid == 0 {
			continue
		}

		proc := &ProcessInfo{
			PID:         pid,
			Name:        strings.Trim(strings.TrimSpace(parts[0]), `"`),
			Environment: make(map[string]string),
			StartTime:   time.Now(),
		}
		processes[pid] = proc
	}

	// 使用 wmic 获取命令行和父进程信息
	cmd = exec.Command("wmic", "process", "get", "ProcessId,ParentProcessId,CommandLine", "/FORMAT:CSV")
	output, err = cmd.Output()
	if err == nil {
		scanner = bufio.NewScanner(strings.NewReader(string(output)))
		firstLine := true
		for scanner.Scan() {
			line := scanner.Text()
			if firstLine {
				firstLine = false
				continue // 跳过 CSV 头
			}
			
			parts := strings.Split(line, ",")
			if len(parts) < 4 {
				continue
			}

			pid, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
			ppid, _ := strconv.Atoi(strings.TrimSpace(parts[3]))
			cmdline := strings.Trim(parts[1], `"`)

			if proc, ok := processes[pid]; ok {
				proc.PPID = ppid
				proc.Cmdline = cmdline
				// 从命令行提取可执行文件路径
				if cmdline != "" {
					parts := strings.Fields(cmdline)
					if len(parts) > 0 {
						proc.Executable = strings.Trim(parts[0], `"`)
					}
				}
			}
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

// parseCSVLine 简单解析 CSV 行
func parseCSVLine(line string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false

	for _, r := range line {
		switch r {
		case '"':
			inQuotes = !inQuotes
		case ',':
			if !inQuotes {
				parts = append(parts, current.String())
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}
	parts = append(parts, current.String())
	return parts
}

// analyzeProcess 分析进程
func (cs *ConfigBasedScanner) analyzeProcess(proc *ProcessInfo) *types.AIComponent {
	cmdlineLower := strings.ToLower(proc.Cmdline)
	exeLower := strings.ToLower(filepath.Base(proc.Executable))
	nameLower := strings.ToLower(proc.Name)

	for _, svc := range cs.serviceConfigs {
		for _, pattern := range svc.ProcessPatterns {
			if re, ok := cs.processPatterns[pattern]; ok {
				// 匹配命令行
				if re.MatchString(cmdlineLower) {
					return cs.createComponentFromService(proc, svc)
				}
				// 匹配可执行文件名
				if re.MatchString(exeLower) {
					return cs.createComponentFromService(proc, svc)
				}
				// 匹配进程名
				if re.MatchString(nameLower) {
					return cs.createComponentFromService(proc, svc)
				}
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
		FilePath:    fmt.Sprintf("PID:%d", proc.PID),
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

// scanNetworkServices 扫描网络服务
func (cs *ConfigBasedScanner) scanNetworkServices(processes map[int]*ProcessInfo) ([]types.AIComponent, error) {
	var components []types.AIComponent
	
	// 获取端口到 PID 的映射 (Windows 使用 netstat)
	portToPID := cs.getPortToPIDMapping()
	
	// 遍历所有服务配置
	for _, svc := range cs.serviceConfigs {
		if len(svc.DefaultPorts) == 0 {
			continue
		}
		
		for _, port := range svc.DefaultPorts {
			pid, ok := portToPID[port]
			if !ok {
				continue
			}
			
			proc, ok := processes[pid]
			if !ok {
				continue
			}
			
			if !cs.matchProcessPatterns(proc, svc) {
				continue
			}
			
			comp := cs.createComponentFromNetworkService(proc, svc, port)
			if comp != nil {
				components = append(components, *comp)
			}
		}
	}
	
	return components, nil
}

// matchProcessPatterns 检查进程是否匹配服务的进程模式
func (cs *ConfigBasedScanner) matchProcessPatterns(proc *ProcessInfo, svc config.ServiceConfig) bool {
	procNameLower := strings.ToLower(proc.Name)
	cmdlineLower := strings.ToLower(proc.Cmdline)
	exeLower := strings.ToLower(filepath.Base(proc.Executable))
	
	for _, pattern := range svc.ProcessPatterns {
		if re, ok := cs.processPatterns[pattern]; ok {
			if re.MatchString(procNameLower) {
				return true
			}
			if re.MatchString(cmdlineLower) {
				return true
			}
			if re.MatchString(exeLower) {
				return true
			}
		}
	}
	
	return false
}

// createComponentFromNetworkService 从网络服务创建组件
func (cs *ConfigBasedScanner) createComponentFromNetworkService(proc *ProcessInfo, svc config.ServiceConfig, port int) *types.AIComponent {
	version := cs.extractVersion(proc, svc)
	
	if version == "" {
		return nil
	}
	
	confidence := cs.calculateConfidence(proc, svc)
	confidence += 0.2
	if confidence > 1.0 {
		confidence = 1.0
	}
	
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
		FilePath:    fmt.Sprintf("PID:%d (port %d)", proc.PID, port),
		Confidence:  confidence,
		Severity:    types.Severity(svc.Severity),
		Description: description,
		RawContent:  fmt.Sprintf("Process: %s | Cmdline: %s", proc.Name, proc.Cmdline),
	}
	
	cs.detectAPIKeys(proc, comp)
	
	return comp
}

// getPortToPIDMapping 获取端口到 PID 映射 (Windows 实现)
func (cs *ConfigBasedScanner) getPortToPIDMapping() map[int]int {
	mapping := make(map[int]int)

	// 使用 netstat -ano 获取端口和 PID 映射
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return mapping
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		// 检查是否为 LISTENING 状态
		if !strings.Contains(line, "LISTENING") {
			continue
		}

		// 解析本地地址和端口
		localAddr := fields[1]
		if idx := strings.LastIndex(localAddr, ":"); idx > 0 {
			portStr := localAddr[idx+1:]
			port, err := strconv.Atoi(portStr)
			if err != nil {
				continue
			}

			// 获取 PID
			pidStr := fields[len(fields)-1]
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				continue
			}

			mapping[port] = pid
		}
	}

	return mapping
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

// ProcessInfo 进程信息
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
