package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// RulesConfig 检测规则配置
type RulesConfig struct {
	LLMServices       []ServiceConfig   `yaml:"llm_services"`
	VectorDatabases   []ServiceConfig   `yaml:"vector_databases"`
	MLFrameworks      []ServiceConfig   `yaml:"ml_frameworks"`
	AgentFrameworks   []ServiceConfig   `yaml:"agent_frameworks"`
	DeploymentTools   []ServiceConfig   `yaml:"deployment_tools"`
	MonitoringTools   []ServiceConfig   `yaml:"monitoring_tools"`
	APIKeyPatterns    []APIKeyPattern   `yaml:"api_key_patterns"`
	Global            GlobalConfig      `yaml:"global"`
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	ID                string               `yaml:"id"`
	Name              string               `yaml:"name"`
	Type              string               `yaml:"type"`
	Severity          string               `yaml:"severity"`
	Description       string               `yaml:"description"`
	ProcessPatterns   []string             `yaml:"process_patterns"`
	DefaultPorts      []int                `yaml:"default_ports"`
	VersionProbe      *VersionProbeConfig  `yaml:"version_probe"`
	HTTPEndpoints     []HTTPEndpointConfig `yaml:"http_endpoints"`
	EnvIndicators     []string             `yaml:"env_indicators"`
	SemanticAnalyzers []SemanticAnalyzerConfig `yaml:"semantic_analyzers"`
}

// VersionProbeConfig 版本探测配置
type VersionProbeConfig struct {
	Methods []VersionMethodConfig `yaml:"methods"`
}

// VersionMethodConfig 版本探测方法
type VersionMethodConfig struct {
	Type     string `yaml:"type"`     // cli_arg, exec, http_api
	Patterns []string `yaml:"patterns,omitempty"`
	Command  string `yaml:"command,omitempty"`
	Endpoint string `yaml:"endpoint,omitempty"`
	JSONPath string `yaml:"json_path,omitempty"`
	Parser   string `yaml:"parser,omitempty"`
}

// HTTPEndpointConfig HTTP 端点配置
type HTTPEndpointConfig struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

// SemanticAnalyzerConfig 语义分析器配置
type SemanticAnalyzerConfig struct {
	Type     string   `yaml:"type"` // python_import, node_import
	Patterns []string `yaml:"patterns"`
}

// APIKeyPattern API Key 检测模式
type APIKeyPattern struct {
	Name     string `yaml:"name"`
	Key      string `yaml:"key"`
	Severity string `yaml:"severity"`
}

// GlobalConfig 全局配置
type GlobalConfig struct {
	Scan struct {
		Timeout       string `yaml:"timeout"`
		MaxProcesses  int    `yaml:"max_processes"`
	} `yaml:"scan"`
	VersionProbe struct {
		Timeout    string `yaml:"timeout"`
		RetryCount int    `yaml:"retry_count"`
	} `yaml:"version_probe"`
	HTTPProbe struct {
		Timeout   string `yaml:"timeout"`
		UserAgent string `yaml:"user_agent"`
	} `yaml:"http_probe"`
	ConfidenceWeights map[string]float64 `yaml:"confidence_weights"`
}

// Loader 配置加载器
type Loader struct {
	configPath string
	config     *RulesConfig
}

// NewLoader 创建配置加载器
func NewLoader(configPath string) *Loader {
	return &Loader{
		configPath: configPath,
	}
}

// Load 加载配置文件
func (l *Loader) Load() error {
	data, err := os.ReadFile(l.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	config := &RulesConfig{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	l.config = config
	return nil
}

// GetConfig 获取配置
func (l *Loader) GetConfig() *RulesConfig {
	return l.config
}

// GetAllServices 获取所有服务配置
func (l *Loader) GetAllServices() []ServiceConfig {
	if l.config == nil {
		return nil
	}

	var services []ServiceConfig
	services = append(services, l.config.LLMServices...)
	services = append(services, l.config.VectorDatabases...)
	services = append(services, l.config.MLFrameworks...)
	services = append(services, l.config.AgentFrameworks...)
	services = append(services, l.config.DeploymentTools...)
	services = append(services, l.config.MonitoringTools...)

	return services
}

// GetServiceByID 根据 ID 获取服务配置
func (l *Loader) GetServiceByID(id string) *ServiceConfig {
	for _, svc := range l.GetAllServices() {
		if svc.ID == id {
			return &svc
		}
	}
	return nil
}

// GetAPIKeyPatterns 获取 API Key 检测模式
func (l *Loader) GetAPIKeyPatterns() []APIKeyPattern {
	if l.config == nil {
		return nil
	}
	return l.config.APIKeyPatterns
}

// GetGlobalConfig 获取全局配置
func (l *Loader) GetGlobalConfig() GlobalConfig {
	if l.config == nil {
		return GlobalConfig{}
	}
	return l.config.Global
}

// Reload 重新加载配置（热更新支持）
func (l *Loader) Reload() error {
	return l.Load()
}

// DefaultConfigPath 返回默认配置文件路径
func DefaultConfigPath() string {
	// 尝试多个位置
	paths := []string{
		"./config/rules.yaml",
		"/etc/ai-scanner/rules.yaml",
		"~/.config/ai-scanner/rules.yaml",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return paths[0] // 返回第一个作为默认
}