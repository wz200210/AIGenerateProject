package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/wz200210/AIGenerateProject/internal/detector"
	"github.com/wz200210/AIGenerateProject/pkg/ai/types"
)

// Scanner 项目扫描器
type Scanner struct {
	detector     *detector.Detector
	ignoreDirs   []string
	ignoreExts   []string
	totalFiles   int
	errors       []string
}

// NewScanner 创建新扫描器
func NewScanner() *Scanner {
	return &Scanner{
		detector:   detector.NewDetector(),
		ignoreDirs: []string{".git", "node_modules", "vendor", ".env", "__pycache__", ".venv", "dist", "build"},
		ignoreExts: []string{".exe", ".dll", ".so", ".dylib", ".zip", ".tar", ".gz", ".rar", ".7z"},
	}
}

// Scan 扫描指定路径
func (s *Scanner) Scan(projectPath string) (*types.ScanResult, error) {
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("无法解析路径: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("无法访问路径: %w", err)
	}

	result := &types.ScanResult{
		ProjectPath: absPath,
		ScanTime:    time.Now().Format(time.RFC3339),
		Components:  []types.AIComponent{},
	}

	if !info.IsDir() {
		// 单文件扫描
		s.totalFiles = 1
		comps, err := s.scanFile(absPath)
		if err != nil {
			s.errors = append(s.errors, fmt.Sprintf("%s: %v", absPath, err))
		} else {
			result.Components = append(result.Components, comps...)
		}
	} else {
		// 目录扫描
		err := filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				s.errors = append(s.errors, fmt.Sprintf("walk error at %s: %v", path, err))
				return nil
			}

			// 跳过忽略的目录
			if info.IsDir() {
				if s.shouldIgnoreDir(path) {
					return filepath.SkipDir
				}
				return nil
			}

			s.totalFiles++

			// 跳过忽略的文件类型
			if s.shouldIgnoreFile(path) {
				return nil
			}

			comps, err := s.scanFile(path)
			if err != nil {
				s.errors = append(s.errors, fmt.Sprintf("%s: %v", path, err))
			} else {
				result.Components = append(result.Components, comps...)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	result.TotalFiles = s.totalFiles
	result.Errors = s.errors

	return result, nil
}

// scanFile 扫描单个文件
func (s *Scanner) scanFile(filePath string) ([]types.AIComponent, error) {
	return s.detector.DetectInFile(filePath)
}

// shouldIgnoreDir 检查是否应该忽略目录
func (s *Scanner) shouldIgnoreDir(path string) bool {
	base := filepath.Base(path)
	for _, dir := range s.ignoreDirs {
		if base == dir {
			return true
		}
	}
	return false
}

// shouldIgnoreFile 检查是否应该忽略文件
func (s *Scanner) shouldIgnoreFile(path string) bool {
	fileExt := strings.ToLower(filepath.Ext(path))
	for _, ignoreExt := range s.ignoreExts {
		if fileExt == ignoreExt {
			return true
		}
	}
	return false
}

// PrintConsoleReport 打印控制台报告
func PrintConsoleReport(result *types.ScanResult) {
	bold := color.New(color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	cyan := color.New(color.FgCyan)

	fmt.Println()
	bold.Println("╔════════════════════════════════════════╗")
	bold.Println("║     AI Component Scan Report           ║")
	bold.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	cyan.Printf("📁 Project: ")
	fmt.Println(result.ProjectPath)
	cyan.Printf("📊 Files scanned: ")
	fmt.Println(result.TotalFiles)
	cyan.Printf("🤖 AI components found: ")
	
	if len(result.Components) == 0 {
		green.Println("0 (No AI components detected)")
	} else {
		red.Println(len(result.Components))
	}
	fmt.Println()

	if len(result.Components) > 0 {
		bold.Println("Detected Components:")
		fmt.Println(strings.Repeat("─", 60))

		// 按类型分组
		byType := make(map[types.AIComponentType][]types.AIComponent)
		for _, c := range result.Components {
			byType[c.Type] = append(byType[c.Type], c)
		}

		for typeName, components := range byType {
			cyan.Printf("\n[%s]\n", typeName)
			for _, c := range components {
				severityColor := getSeverityColor(c.Severity)
				fmt.Printf("  • %s ", c.Name)
				severityColor.Printf("[%s]\n", c.Severity)
				fmt.Printf("    File: %s", c.FilePath)
				if c.LineNumber > 0 {
					fmt.Printf(":%d", c.LineNumber)
				}
				fmt.Println()
				if c.Version != "" {
					fmt.Printf("    Version: %s\n", c.Version)
				}
				if c.Description != "" {
					fmt.Printf("    %s\n", c.Description)
				}
				fmt.Println()
			}
		}
	}

	if len(result.Errors) > 0 {
		yellow.Printf("\n⚠️  Warnings (%d):\n", len(result.Errors))
		for _, err := range result.Errors {
			fmt.Printf("  • %s\n", err)
		}
	}

	fmt.Println()
	fmt.Printf("Scan completed at: %s\n", result.ScanTime)
}

// PrintJSONReport 打印 JSON 格式报告
func PrintJSONReport(result *types.ScanResult) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// getSeverityColor 根据严重级别返回颜色
func getSeverityColor(s types.Severity) *color.Color {
	switch s {
	case types.SeverityCritical:
		return color.New(color.FgRed, color.Bold)
	case types.SeverityHigh:
		return color.New(color.FgRed)
	case types.SeverityMedium:
		return color.New(color.FgYellow)
	default:
		return color.New(color.FgGreen)
	}
}

// PrintRuntimeConsoleReport 打印运行时扫描控制台报告
func PrintRuntimeConsoleReport(result *types.RuntimeScanResult) {
	bold := color.New(color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	cyan := color.New(color.FgCyan)

	fmt.Println()
	bold.Println("╔══════════════════════════════════════════════════════╗")
	bold.Println("║     Runtime AI Component Scan Report                 ║")
	bold.Println("╚══════════════════════════════════════════════════════╝")
	fmt.Println()

	cyan.Printf("⏱️  Scan Time: ")
	fmt.Println(result.ScanTime)
	fmt.Println()

	// 统计信息
	bold.Println("📊 Scan Summary:")
	fmt.Printf("  • Processes scanned: %d\n", result.ProcessCount)
	fmt.Printf("  • Ports scanned: %d\n", result.PortCount)
	fmt.Printf("  • Containers scanned: %d\n", result.ContainerCount)
	fmt.Printf("  • Total components found: ")
	if len(result.Components) == 0 {
		green.Println("0 (No running AI components detected)")
	} else {
		red.Println(len(result.Components))
	}
	fmt.Println()

	if len(result.Components) > 0 {
		bold.Println("Running AI Components:")
		fmt.Println(strings.Repeat("─", 60))

		// 按类型分组
		byType := make(map[types.AIComponentType][]types.AIComponent)
		for _, c := range result.Components {
			byType[c.Type] = append(byType[c.Type], c)
		}

		for typeName, components := range byType {
			cyan.Printf("\n[%s]\n", typeName)
			for _, c := range components {
				severityColor := getSeverityColor(c.Severity)
				fmt.Printf("  • %s ", c.Name)
				severityColor.Printf("[%s]", c.Severity)
				if c.Version != "" {
					fmt.Printf(" v%s", c.Version)
				}
				fmt.Println()
				fmt.Printf("    Source: %s\n", c.FilePath)
				if c.Description != "" {
					fmt.Printf("    %s\n", c.Description)
				}
				fmt.Println()
			}
		}
	}

	if len(result.Errors) > 0 {
		yellow.Printf("\n⚠️  Warnings (%d):\n", len(result.Errors))
		for _, err := range result.Errors {
			fmt.Printf("  • %s\n", err)
		}
	}

	fmt.Println()
}

// PrintRuntimeJSONReport 打印运行时扫描 JSON 报告
func PrintRuntimeJSONReport(result *types.RuntimeScanResult) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// PrintFullConsoleReport 打印完整扫描控制台报告
func PrintFullConsoleReport(result *types.FullScanResult) {
	bold := color.New(color.Bold)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	cyan := color.New(color.FgCyan)

	fmt.Println()
	bold.Println("╔══════════════════════════════════════════════════════╗")
	bold.Println("║     Full AI Component Scan Report                    ║")
	bold.Println("║     (Static + Runtime)                               ║")
	bold.Println("╚══════════════════════════════════════════════════════╝")
	fmt.Println()

	cyan.Printf("📁 Project: ")
	fmt.Println(result.ProjectPath)
	cyan.Printf("⏱️  Scan Time: ")
	fmt.Println(result.ScanTime)
	fmt.Println()

	// 静态扫描结果
	if result.StaticScan != nil {
		bold.Println("📁 Static Scan Results:")
		fmt.Printf("  Files scanned: %d\n", result.StaticScan.TotalFiles)
		fmt.Printf("  Components found: ")
		if len(result.StaticScan.Components) == 0 {
			green.Println("0")
		} else {
			red.Println(len(result.StaticScan.Components))
		}
		fmt.Println()
	}

	// 运行时扫描结果
	if result.RuntimeScan != nil {
		bold.Println("⚙️  Runtime Scan Results:")
		fmt.Printf("  Processes found: %d\n", result.RuntimeScan.ProcessCount)
		fmt.Printf("  Ports found: %d\n", result.RuntimeScan.PortCount)
		fmt.Printf("  Containers found: %d\n", result.RuntimeScan.ContainerCount)
		fmt.Printf("  Components found: ")
		if len(result.RuntimeScan.Components) == 0 {
			green.Println("0")
		} else {
			red.Println(len(result.RuntimeScan.Components))
		}
		fmt.Println()
	}

	// 总计
	bold.Printf("📊 Total Components Found: ")
	if result.TotalComponents == 0 {
		green.Println("0")
	} else {
		red.Println(result.TotalComponents)
	}

	fmt.Println()
}

// PrintFullJSONReport 打印完整扫描 JSON 报告
func PrintFullJSONReport(result *types.FullScanResult) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}