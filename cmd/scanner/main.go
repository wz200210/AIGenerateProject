package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/wz200210/AIGenerateProject/internal/config"
	"github.com/wz200210/AIGenerateProject/internal/runtime"
	"github.com/wz200210/AIGenerateProject/internal/scanner"
	"github.com/wz200210/AIGenerateProject/pkg/ai/types"
)

var (
	outputFormat string
	configPath   string
	version      = "0.4.0"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "scanner",
		Short: "AI Component Scanner - 配置化运行时检测工具",
		Long: `scanner v0.4.0 - 基于配置文件的运行时 AI 组件检测

🚀 全新架构：
  • 组件特征完全外置到 YAML 配置文件
  • 支持热更新，无需重新编译
  • 灵活的版本探测策略配置
  • 可扩展的语义分析器

检测能力：
  • LLM 推理服务 (Ollama, vLLM, TGI, etc.)
  • 向量数据库 (Milvus, Chroma, Weaviate, etc.)
  • ML 框架服务 (Transformers, PyTorch, etc.)
  • Agent/RAG 框架 (LangChain, LlamaIndex, etc.)
  • Docker/K8s 容器中的 AI 服务
  • API Key 泄露检测`,
		Version: version,
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "配置文件路径 (默认: ./config/rules.yaml)")

	// 扫描命令
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "扫描运行时 AI 组件",
		Long:  "扫描当前系统中运行的 AI 相关进程、服务和容器",
		RunE:  runScan,
	}
	scanCmd.Flags().StringVarP(&outputFormat, "output", "o", "console", "输出格式 (console|json)")

	// 验证配置命令
	validateCmd := &cobra.Command{
		Use:   "validate-config",
		Short: "验证配置文件",
		Long:  "检查配置文件格式和规则有效性",
		RunE:  runValidateConfig,
	}

	// 列出规则命令
	listRulesCmd := &cobra.Command{
		Use:   "list-rules",
		Short: "列出所有检测规则",
		Long:  "显示配置文件中定义的所有 AI 组件检测规则",
		RunE:  runListRules,
	}

	rootCmd.AddCommand(scanCmd, validateCmd, listRulesCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// getConfigPath 获取配置文件路径
func getConfigPath() string {
	if configPath != "" {
		return configPath
	}
	return config.DefaultConfigPath()
}

// runScan 执行扫描
func runScan(cmd *cobra.Command, args []string) error {
	cfgPath := getConfigPath()

	fmt.Printf("🔍 AI Component Runtime Scanner v%s\n", version)
	fmt.Printf("📄 Config: %s\n", cfgPath)
	fmt.Println("═══════════════════════════════════════════════════════")
	fmt.Println()

	// 创建基于配置的扫描器
	rs, err := runtime.NewConfigBasedScanner(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to initialize scanner: %w", err)
	}

	start := time.Now()
	result, err := rs.ScanAll()
	elapsed := time.Since(start)

	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Scan warning: %v\n", err)
	}

	result.ScanDuration = elapsed.String()

	switch outputFormat {
	case "json":
		return scanner.PrintRuntimeJSONReport(result)
	default:
		scanner.PrintRuntimeConsoleReport(result)
	}

	return nil
}

// runValidateConfig 验证配置
func runValidateConfig(cmd *cobra.Command, args []string) error {
	cfgPath := getConfigPath()

	fmt.Printf("📄 Validating config: %s\n", cfgPath)
	fmt.Println()

	loader := config.NewLoader(cfgPath)
	if err := loader.Load(); err != nil {
		fmt.Printf("❌ Config validation failed: %v\n", err)
		return err
	}

	cfg := loader.GetConfig()
	services := loader.GetAllServices()
	apiKeys := loader.GetAPIKeyPatterns()

	fmt.Println("✅ Config file is valid")
	fmt.Println()
	fmt.Printf("📊 Summary:\n")
	fmt.Printf("  • LLM Services: %d\n", len(cfg.LLMServices))
	fmt.Printf("  • Vector Databases: %d\n", len(cfg.VectorDatabases))
	fmt.Printf("  • ML Frameworks: %d\n", len(cfg.MLFrameworks))
	fmt.Printf("  • Agent Frameworks: %d\n", len(cfg.AgentFrameworks))
	fmt.Printf("  • Deployment Tools: %d\n", len(cfg.DeploymentTools))
	fmt.Printf("  • Monitoring Tools: %d\n", len(cfg.MonitoringTools))
	fmt.Printf("  • Total Services: %d\n", len(services))
	fmt.Printf("  • API Key Patterns: %d\n", len(apiKeys))

	return nil
}

// runListRules 列出规则
func runListRules(cmd *cobra.Command, args []string) error {
	cfgPath := getConfigPath()

	loader := config.NewLoader(cfgPath)
	if err := loader.Load(); err != nil {
		return err
	}

	services := loader.GetAllServices()

	fmt.Println("🔍 AI Component Detection Rules")
	fmt.Println("═══════════════════════════════════════════════════════")
	fmt.Println()

	// 按类型分组
	byType := make(map[types.AIComponentType][]config.ServiceConfig)
	for _, svc := range services {
		byType[types.AIComponentType(svc.Type)] = append(byType[types.AIComponentType(svc.Type)], svc)
	}

	for typeName, svcs := range byType {
		fmt.Printf("\n[%s]\n", typeName)
		for _, svc := range svcs {
			fmt.Printf("  • %s (id: %s)\n", svc.Name, svc.ID)
			fmt.Printf("    Ports: %v\n", svc.DefaultPorts)
			if len(svc.ProcessPatterns) > 0 {
				fmt.Printf("    Patterns: %d regex rules\n", len(svc.ProcessPatterns))
			}
		}
	}

	return nil
}