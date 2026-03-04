package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/wz200210/AIGenerateProject/internal/runtime"
	"github.com/wz200210/AIGenerateProject/internal/scanner"
	"github.com/wz200210/AIGenerateProject/pkg/ai/types"
)

var (
	projectPath    string
	outputFormat   string
	scanMode       string
	includeRuntime bool
	version        = "0.2.0"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "scanner",
		Short: "AI Component Scanner - 扫描项目中的 AI 组件",
		Long: `scanner 是一个用于识别项目中 AI/ML 组件的工具。

扫描模式:
  • 静态扫描 (static)  - 扫描代码文件、依赖、配置
  • 运行时扫描 (runtime) - 扫描运行中的进程、端口、容器
  • 完整扫描 (full)    - 静态扫描 + 运行时扫描

它可以检测:
  • LLM 框架 (OpenAI, Claude, Gemini, LangChain, etc.)
  • ML 框架 (TensorFlow, PyTorch, Hugging Face, etc.)
  • 向量数据库 (Pinecone, Milvus, Chroma, etc.)
  • AI 模型文件 (.gguf, .pt, .onnx, etc.)
  • 运行中的 AI 服务进程
  • API Keys 和配置`,
		Version: version,
	}

	// 静态扫描命令
	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "扫描项目目录",
		Long:  "扫描指定目录，识别其中的 AI 组件和依赖",
		RunE:  runScan,
	}
	scanCmd.Flags().StringVarP(&projectPath, "path", "p", ".", "要扫描的项目路径")
	scanCmd.Flags().StringVarP(&outputFormat, "output", "o", "console", "输出格式 (console|json)")
	scanCmd.Flags().BoolVarP(&includeRuntime, "runtime", "r", false, "同时扫描运行时进程和端口")

	// 运行时扫描命令
	runtimeCmd := &cobra.Command{
		Use:   "runtime",
		Short: "扫描运行时 AI 组件",
		Long: `扫描当前系统中运行的 AI 相关进程、服务和容器。

检测内容:
  • 运行中的 AI/ML 进程 (ollama, vllm, python, etc.)
  • 监听的 AI 服务端口 (11434, 8000, 6333, etc.)
  • Docker 容器中的 AI 服务
  • 进程环境变量中的 API Key`,
		RunE: runRuntimeScan,
	}
	runtimeCmd.Flags().StringVarP(&outputFormat, "output", "o", "console", "输出格式 (console|json)")

	// 完整扫描命令
	fullCmd := &cobra.Command{
		Use:   "full",
		Short: "执行完整扫描（静态+运行时）",
		Long:  "同时扫描项目文件和运行时进程，获取最全面的 AI 组件信息",
		RunE:  runFullScan,
	}
	fullCmd.Flags().StringVarP(&projectPath, "path", "p", ".", "要扫描的项目路径")
	fullCmd.Flags().StringVarP(&outputFormat, "output", "o", "console", "输出格式 (console|json)")

	rootCmd.AddCommand(scanCmd, runtimeCmd, fullCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// runScan 执行静态扫描
func runScan(cmd *cobra.Command, args []string) error {
	s := scanner.NewScanner()

	fmt.Printf("🔍 Scanning: %s\n", projectPath)

	result, err := s.Scan(projectPath)
	if err != nil {
		return fmt.Errorf("扫描失败: %w", err)
	}

	// 如果需要，同时扫描运行时
	if includeRuntime {
		rt := runtime.NewRuntimeScanner()
		rtResult, _ := rt.ScanAll()
		// 合并结果
		result.Components = append(result.Components, rtResult.Components...)
	}

	switch outputFormat {
	case "json":
		return scanner.PrintJSONReport(result)
	default:
		scanner.PrintConsoleReport(result)
	}

	return nil
}

// runRuntimeScan 执行运行时扫描
func runRuntimeScan(cmd *cobra.Command, args []string) error {
	rt := runtime.NewRuntimeScanner()

	fmt.Println("🔍 Runtime Scanning: 正在扫描运行中的 AI 组件...")
	fmt.Println()

	result, err := rt.ScanAll()
	if err != nil {
		return fmt.Errorf("运行时扫描失败: %w", err)
	}

	switch outputFormat {
	case "json":
		return scanner.PrintRuntimeJSONReport(result)
	default:
		scanner.PrintRuntimeConsoleReport(result)
	}

	return nil
}

// runFullScan 执行完整扫描
func runFullScan(cmd *cobra.Command, args []string) error {
	fullResult := &types.FullScanResult{
		ProjectPath: projectPath,
		ScanTime:    time.Now().Format(time.RFC3339),
	}

	// 1. 静态扫描
	fmt.Printf("🔍 Step 1/2: Static Scanning %s...\n", projectPath)
	s := scanner.NewScanner()
	staticResult, err := s.Scan(projectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Static scan warning: %v\n", err)
	} else {
		fullResult.StaticScan = staticResult
		fullResult.TotalComponents += len(staticResult.Components)
	}

	// 2. 运行时扫描
	fmt.Println("🔍 Step 2/2: Runtime Scanning...")
	rt := runtime.NewRuntimeScanner()
	runtimeResult, err := rt.ScanAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Runtime scan warning: %v\n", err)
	} else {
		fullResult.RuntimeScan = runtimeResult
		fullResult.TotalComponents += len(runtimeResult.Components)
	}

	switch outputFormat {
	case "json":
		return scanner.PrintFullJSONReport(fullResult)
	default:
		scanner.PrintFullConsoleReport(fullResult)
	}

	return nil
}