package detector

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wz200210/AIGenerateProject/pkg/ai/types"
)

// Detector 负责检测 AI 组件
type Detector struct {
	frameworkPatterns map[string]*regexp.Regexp
	apiKeyPatterns    map[string]*regexp.Regexp
}

// NewDetector 创建新的检测器
func NewDetector() *Detector {
	d := &Detector{
		frameworkPatterns: make(map[string]*regexp.Regexp),
		apiKeyPatterns:    make(map[string]*regexp.Regexp),
	}
	d.compilePatterns()
	return d
}

// compilePatterns 编译所有检测正则
func (d *Detector) compilePatterns() {
	// 编译框架检测模式
	for _, fw := range types.CommonAIFrameworks {
		for _, pattern := range fw.Patterns {
			if _, exists := d.frameworkPatterns[pattern]; !exists {
				if re, err := regexp.Compile(`(?i)` + regexp.QuoteMeta(pattern)); err == nil {
					d.frameworkPatterns[pattern] = re
				}
			}
		}
	}

	// 编译 API Key 检测模式
	for key, name := range types.APIKeyPatterns {
		// 匹配 KEY=xxx 或 "KEY": "xxx" 等格式
		pattern := `(?i)(` + regexp.QuoteMeta(key) + `)\s*[=:]\s*["']?([a-zA-Z0-9_\-\.]+)["']?`
		if re, err := regexp.Compile(pattern); err == nil {
			d.apiKeyPatterns[name] = re
		}
	}
}

// DetectInFile 检测单个文件中的 AI 组件
func (d *Detector) DetectInFile(filePath string) ([]types.AIComponent, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var components []types.AIComponent
	ext := strings.ToLower(filepath.Ext(filePath))
	baseName := filepath.Base(filePath)

	// 检测模型文件
	if modelFile := d.detectModelFile(filePath, ext); modelFile != nil {
		components = append(components, *modelFile)
	}

	// 按行扫描代码文件
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		lowerLine := strings.ToLower(line)

		// 检测框架引用
		for _, fw := range types.CommonAIFrameworks {
			if d.shouldScanFile(baseName, fw.FilePatterns) {
				for _, pattern := range fw.Patterns {
					if re, ok := d.frameworkPatterns[pattern]; ok && re.MatchString(lowerLine) {
						comp := types.AIComponent{
							Name:        fw.Name,
							Type:        fw.Type,
							FilePath:    filePath,
							LineNumber:  lineNum,
							Confidence:  0.85,
							Severity:    fw.Severity,
							Description: fw.Description,
							RawContent:  strings.TrimSpace(line),
						}
						if !d.componentExists(components, comp) {
							components = append(components, comp)
						}
					}
				}
			}
		}

		// 检测 API Key
		for name, re := range d.apiKeyPatterns {
			if matches := re.FindStringSubmatch(line); len(matches) > 0 {
				// 隐藏实际的 key 值
				masked := re.ReplaceAllString(line, matches[1]+"=***HIDDEN***")
				comp := types.AIComponent{
					Name:        name,
					Type:        types.TypeAPIKey,
					FilePath:    filePath,
					LineNumber:  lineNum,
					Confidence:  0.95,
					Severity:    types.SeverityCritical,
					Description: "Potential hardcoded API key detected",
					RawContent:  masked,
				}
				if !d.componentExists(components, comp) {
					components = append(components, comp)
				}
			}
		}
	}

	return components, scanner.Err()
}

// detectModelFile 检测模型文件
func (d *Detector) detectModelFile(filePath, ext string) *types.AIComponent {
	for _, model := range types.CommonModelFiles {
		if ext == model.Extension {
			return &types.AIComponent{
				Name:        model.Name,
				Type:        model.Type,
				FilePath:    filePath,
				Confidence:  1.0,
				Severity:    model.Severity,
				Description: model.Description,
			}
		}
	}
	return nil
}

// shouldScanFile 检查是否应该扫描该文件
func (d *Detector) shouldScanFile(fileName string, patterns []string) bool {
	lowerName := strings.ToLower(fileName)
	for _, pattern := range patterns {
		if strings.ToLower(pattern) == lowerName {
			return true
		}
	}
	// 代码文件也扫描
	codeExts := []string{".go", ".py", ".js", ".ts", ".java", ".rs", ".cpp", ".c", ".h"}
	for _, ext := range codeExts {
		if strings.HasSuffix(lowerName, ext) {
			return true
		}
	}
	return false
}

// componentExists 检查组件是否已存在（去重）
func (d *Detector) componentExists(components []types.AIComponent, new types.AIComponent) bool {
	for _, c := range components {
		if c.Name == new.Name && c.FilePath == new.FilePath && c.LineNumber == new.LineNumber {
			return true
		}
	}
	return false
}