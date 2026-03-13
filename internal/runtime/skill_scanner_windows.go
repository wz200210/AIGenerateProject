//go:build windows

package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wz200210/AIGenerateProject/internal/config"
	"github.com/wz200210/AIGenerateProject/pkg/ai/types"
	"gopkg.in/yaml.v3"
)

// SkillScanner Skill 配置文件扫描器
type SkillScanner struct {
	configs []config.SkillScanConfig
}

// NewSkillScanner 创建 Skill 扫描器
func NewSkillScanner(configs []config.SkillScanConfig) *SkillScanner {
	var enabled []config.SkillScanConfig
	for _, cfg := range configs {
		if cfg.Enabled {
			enabled = append(enabled, cfg)
		}
	}
	return &SkillScanner{configs: enabled}
}

// ScanAll 扫描所有配置的 Skill 源
func (ss *SkillScanner) ScanAll() ([]types.SkillInfo, error) {
	var allSkills []types.SkillInfo

	for _, cfg := range ss.configs {
		skills, err := ss.scanSource(cfg)
		if err != nil {
			continue
		}
		allSkills = append(allSkills, skills...)
	}

	return allSkills, nil
}

// scanSource 扫描单个 Skill 源
func (ss *SkillScanner) scanSource(cfg config.SkillScanConfig) ([]types.SkillInfo, error) {
	var skills []types.SkillInfo

	// 1. 从配置文件解析 Skill 列表
	for _, path := range cfg.ConfigPaths {
		expandedPath := ss.expandPath(path)
		if _, err := os.Stat(expandedPath); err == nil {
			parsedSkills, err := ss.parseConfigFile(expandedPath, cfg)
			if err == nil {
				skills = append(skills, parsedSkills...)
			}
		}
	}

	// 2. 从 Skill 目录扫描
	for _, dir := range cfg.SkillDirs {
		expandedDir := ss.expandPath(dir)
		dirSkills, err := ss.scanSkillDirectory(expandedDir, cfg)
		if err == nil {
			skills = append(skills, dirSkills...)
		}
	}

	return ss.deduplicateSkills(skills), nil
}

// expandPath 扩展路径（处理 Windows 用户目录和环境变量）
func (ss *SkillScanner) expandPath(path string) string {
	// 处理 ~ 扩展到用户主目录
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(homeDir, path[2:])
		}
	}
	
	// 处理环境变量
	path = os.ExpandEnv(path)
	
	return path
}

// parseConfigFile 解析配置文件
func (ss *SkillScanner) parseConfigFile(path string, cfg config.SkillScanConfig) ([]types.SkillInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	switch cfg.ParseRules.ConfigFormat {
	case "yaml", "yml":
		return ss.parseYAMLConfig(data, path, cfg)
	case "json":
		return ss.parseJSONConfig(data, path, cfg)
	default:
		if json.Valid(data) {
			return ss.parseJSONConfig(data, path, cfg)
		}
		return ss.parseYAMLConfig(data, path, cfg)
	}
}

// parseYAMLConfig 解析 YAML 配置
func (ss *SkillScanner) parseYAMLConfig(data []byte, path string, cfg config.SkillScanConfig) ([]types.SkillInfo, error) {
	var root interface{}
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	return ss.extractSkillsFromInterface(root, path, cfg)
}

// parseJSONConfig 解析 JSON 配置
func (ss *SkillScanner) parseJSONConfig(data []byte, path string, cfg config.SkillScanConfig) ([]types.SkillInfo, error) {
	var root interface{}
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	return ss.extractSkillsFromInterface(root, path, cfg)
}

// extractSkillsFromInterface 从解析后的数据中提取 Skill 列表
func (ss *SkillScanner) extractSkillsFromInterface(root interface{}, path string, cfg config.SkillScanConfig) ([]types.SkillInfo, error) {
	var skills []types.SkillInfo

	rootMap, ok := root.(map[string]interface{})
	if !ok {
		return skills, nil
	}

	skillsKey := "skills"
	if cfg.ParseRules.SkillNamePath != "" {
		parts := strings.Split(cfg.ParseRules.SkillNamePath, ".")
		if len(parts) > 0 {
			skillsKey = parts[0]
		}
	}

	skillsData, ok := rootMap[skillsKey]
	if !ok {
		for _, key := range []string{"skills", "agents", "tools", "extensions", "mcpServers"} {
			if data, exists := rootMap[key]; exists {
				skillsData = data
				break
			}
		}
	}

	if skillsData == nil {
		return skills, nil
	}

	if skillsArray, ok := skillsData.([]interface{}); ok {
		for _, item := range skillsArray {
			if skill := ss.parseSkillItem(item, path, cfg); skill != nil {
				skills = append(skills, *skill)
			}
		}
	}

	if skillsMap, ok := skillsData.(map[string]interface{}); ok {
		for name, item := range skillsMap {
			if skill := ss.parseSkillItem(item, path, cfg); skill != nil {
				skill.Name = name
				skills = append(skills, *skill)
			}
		}
	}

	return skills, nil
}

// parseSkillItem 解析单个 Skill 项
func (ss *SkillScanner) parseSkillItem(item interface{}, path string, cfg config.SkillScanConfig) *types.SkillInfo {
	skill := &types.SkillInfo{
		Source:   cfg.Name,
		Location: path,
		Enabled:  true,
	}

	itemMap, ok := item.(map[string]interface{})
	if !ok {
		if name, ok := item.(string); ok {
			skill.Name = name
			return skill
		}
		return nil
	}

	if cfg.ParseRules.SkillNamePath != "" {
		skill.Name = ss.getValueByPath(itemMap, cfg.ParseRules.SkillNamePath)
	} else {
		for _, key := range []string{"name", "id", "key", "title"} {
			if val, ok := itemMap[key]; ok {
				skill.Name = fmt.Sprintf("%v", val)
				break
			}
		}
	}

	if skill.Name == "" {
		return nil
	}

	if cfg.ParseRules.SkillDescPath != "" {
		skill.Description = ss.getValueByPath(itemMap, cfg.ParseRules.SkillDescPath)
	} else {
		for _, key := range []string{"description", "desc", "summary", "purpose"} {
			if val, ok := itemMap[key]; ok {
				skill.Description = fmt.Sprintf("%v", val)
				break
			}
		}
	}

	if cfg.ParseRules.SkillEnablePath != "" {
		enabledStr := ss.getValueByPath(itemMap, cfg.ParseRules.SkillEnablePath)
		skill.Enabled = enabledStr == "true" || enabledStr == "1"
	}

	return skill
}

// getValueByPath 通过路径获取值
func (ss *SkillScanner) getValueByPath(data map[string]interface{}, path string) string {
	parts := strings.Split(path, ".")
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			if val, ok := current[part]; ok {
				return fmt.Sprintf("%v", val)
			}
			return ""
		}

		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			return ""
		}
	}

	return ""
}

// scanSkillDirectory 扫描 Skill 目录
func (ss *SkillScanner) scanSkillDirectory(dir string, cfg config.SkillScanConfig) ([]types.SkillInfo, error) {
	var skills []types.SkillInfo

	entries, err := os.ReadDir(dir)
	if err != nil {
		return skills, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillName := entry.Name()
		skillDir := filepath.Join(dir, skillName)

		description := ""
		skillFile := filepath.Join(skillDir, cfg.ParseRules.SkillFilePattern)
		if cfg.ParseRules.SkillFilePattern == "" {
			skillFile = filepath.Join(skillDir, "SKILL.md")
		}

		if data, err := os.ReadFile(skillFile); err == nil {
			description = ss.extractDescriptionFromSkillMD(string(data))
		}

		skills = append(skills, types.SkillInfo{
			Name:        skillName,
			Source:      cfg.Name,
			Description: description,
			Location:    skillDir,
			Enabled:     true,
		})
	}

	return skills, nil
}

// extractDescriptionFromSkillMD 从 SKILL.md 提取描述
func (ss *SkillScanner) extractDescriptionFromSkillMD(content string) string {
	lines := strings.Split(content, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if len(line) > 100 {
			return line[:100] + "..."
		}
		return line
	}
	
	return ""
}

// deduplicateSkills 去重
func (ss *SkillScanner) deduplicateSkills(skills []types.SkillInfo) []types.SkillInfo {
	seen := make(map[string]bool)
	var unique []types.SkillInfo

	for _, skill := range skills {
		key := fmt.Sprintf("%s|%s|%s", skill.Source, skill.Name, skill.Location)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, skill)
		}
	}

	return unique
}