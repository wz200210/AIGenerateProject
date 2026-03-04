package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/wz200210/AIGenerateProject/internal/runtime"
	"github.com/wz200210/AIGenerateProject/internal/scanner"
)

var (
	outputFormat string
	version      = "0.3.0"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "scanner",
		Short: "AI Component Scanner - 运行时 AI 组件检测工具",
		Long: `scanner v0.3.0 - 纯运行时 AI 组件检测工具 (已废弃文件扫描)

🚀 全新架构：
  • 完全基于运行时进程检测，不再扫描代码文件
  • 智能语义分析，准确识别 AI 服务类型
  • 自动版本探测，支持多种版本获取方式
  • 进程关系分析，理解服务依赖拓扑

检测能力：
  • LLM 推理服务 (Ollama, vLLM, TGI, OpenAI API, etc.)
  • 向量数据库 (Milvus, Chroma, Weaviate, Qdrant, etc.)
  • ML 框架服务 (Transformers, TorchServe, Triton, etc.)
  • RAG/Agent 框架 (LangChain, LlamaIndex, etc.)
  • 监控工具 (MLflow, W&B, etc.)
  • Docker/K8s 容器中的 AI 服务
  • API Key 泄露检测

⚠️  注意：v0.3.0 起已废弃文件扫描功能，专注于运行时检测。`,
		Version: version,
	}

	// 运行时扫描命令（默认）
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "扫描运行时 AI 组件（默认）",
		Long: `扫描当前系统中运行的 AI 相关进程、服务和容器。

检测内容:
  • 运行中的 AI/ML 进程 (ollama, vllm, python, etc.)
  • 监听的 AI 服务端口 (11434, 8000, 6333, etc.)
  • Docker 容器中的 AI 服务
  • 进程环境变量中的 API Key
  • 进程版本信息自动探测
  • 服务依赖关系分析`,
		RunE: runScan,
	}
	scanCmd.Flags().StringVarP(&outputFormat, "output", "o", "console", "输出格式 (console|json)")

	// 详细扫描命令
	detailCmd := &cobra.Command{
		Use:   "detail",
		Short: "详细扫描（包含进程树和网络连接）",
		Long:  "执行更详细的扫描，包含完整的进程关系和网络连接分析",
		RunE:  runDetailScan,
	}
	detailCmd.Flags().StringVarP(&outputFormat, "output", "o", "console", "输出格式 (console|json)")

	// 版本探测命令
	versionCmd := &cobra.Command{
		Use:   "version-check",
		Short: "检查 AI 服务版本",
		Long:  "主动探测运行的 AI 服务的版本信息",
		RunE:  runVersionCheck,
	}

	// 遗留命令（提示已废弃）
	legacyScanCmd := &cobra.Command{
		Use:    "static",
		Short:  "[已废弃] 静态文件扫描",
		Long:   "⚠️  此功能已在 v0.3.0 中废弃，请使用 'scan' 命令进行运行时检测",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("⚠️  静态文件扫描功能已在 v0.3.0 中废弃")
			fmt.Println("📝 请使用 'scanner scan' 进行运行时 AI 组件检测")
			fmt.Println("")
			fmt.Println("原因：")
			fmt.Println("  1. 文件扫描只能发现代码中存在，无法确认是否在运行")
			fmt.Println("  2. 运行时检测能发现容器化、远程服务等实际运行的 AI")
			fmt.Println("  3. 版本号只能从运行中的进程准确获取")
			return nil
		},
	}

	rootCmd.AddCommand(scanCmd, detailCmd, versionCmd, legacyScanCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// runScan 执行运行时扫描
func runScan(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 AI Component Runtime Scanner v" + version)
	fmt.Println("═══════════════════════════════════════════════════════")
	fmt.Println()

	rs := runtime.NewRuntimeScanner()

	start := time.Now()
	result, err := rs.ScanAll()
	elapsed := time.Since(start)

	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Scan warning: %v\n", err)
	}

	// 添加扫描统计
	result.ScanDuration = elapsed.String()

	switch outputFormat {
	case "json":
		return scanner.PrintRuntimeJSONReport(result)
	default:
		scanner.PrintRuntimeConsoleReport(result)
	}

	return nil
}

// runDetailScan 执行详细扫描
func runDetailScan(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 AI Component Detailed Scanner v" + version)
	fmt.Println("═══════════════════════════════════════════════════════")
	fmt.Println()

	rs := runtime.NewRuntimeScanner()

	start := time.Now()
	result, err := rs.ScanAll()
	elapsed := time.Since(start)

	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Scan warning: %v\n", err)
	}

	result.ScanDuration = elapsed.String()

	// TODO: 实现更详细的报告，包含进程树
	fmt.Printf("⏱️  Scan completed in %s\n", elapsed)
	fmt.Println()

	switch outputFormat {
	case "json":
		return scanner.PrintRuntimeJSONReport(result)
	default:
		scanner.PrintRuntimeConsoleReport(result)
	}

	return nil
}

// runVersionCheck 执行版本检查
func runVersionCheck(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Checking AI Service Versions...")
	fmt.Println()

	// TODO: 实现专门的版本探测功能
	fmt.Println("此功能正在开发中...")

	return nil
}