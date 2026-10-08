package detector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wz200210/AIGenerateProject/pkg/ai/types"
)

// writeFixture 写入测试固件文件
func writeFixture(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write fixture %s: %v", name, err)
	}
	return path
}

// detectNames 扫描文件并返回检测到的组件名集合
func detectNames(t *testing.T, d *Detector, path string) map[string]bool {
	t.Helper()
	comps, err := d.DetectInFile(path)
	if err != nil {
		t.Fatalf("detect %s: %v", path, err)
	}
	names := make(map[string]bool)
	for _, c := range comps {
		names[c.Name] = true
	}
	return names
}

// TestDetectNewLLMServices 验证 2026-10 新增 LLM 服务静态识别
func TestDetectNewLLMServices(t *testing.T) {
	dir := t.TempDir()
	d := NewDetector()

	reqs := writeFixture(t, dir, "requirements.txt", `
openai>=1.0.0
deepseek>=1.0.0
groq>=0.9.0
openrouter>=0.2.0
siliconflow>=0.1.0
volcengine>=1.0.0
xai-sdk>=1.0.0
together>=1.3.0
`)
	names := detectNames(t, d, reqs)

	for _, want := range []string{"DeepSeek", "xAI Grok", "OpenRouter", "SiliconFlow", "Volcengine 豆包", "Groq", "Together AI", "OpenAI"} {
		if !names[want] {
			t.Errorf("requirements.txt 未检测到 %s", want)
		}
	}
}

// TestDetectNewAgentFrameworks 验证新增 Agent 框架识别
func TestDetectNewAgentFrameworks(t *testing.T) {
	dir := t.TempDir()
	d := NewDetector()

	reqs := writeFixture(t, dir, "requirements.txt", `
openai-agents>=0.1.0
google-adk>=0.1.0
pydantic-ai>=0.0.7
agno>=1.0.0
smolagents>=1.0.0
strands-agents>=0.1.0
camel-ai>=0.2.0
dify>=0.15.0
`)
	names := detectNames(t, d, reqs)

	for _, want := range []string{"OpenAI Agents SDK", "Google ADK", "Pydantic AI", "Agno", "smolagents", "Strands Agents", "Camel AI", "Dify"} {
		if !names[want] {
			t.Errorf("requirements.txt 未检测到 %s", want)
		}
	}
}

// TestDetectNewRAGAndVector 验证新增 RAG/向量数据库识别
func TestDetectNewRAGAndVector(t *testing.T) {
	dir := t.TempDir()
	d := NewDetector()

	reqs := writeFixture(t, dir, "requirements.txt", `
lancedb>=0.6.0
farm-haystack>=2.0.0
ragflow>=0.9.0
lightrag>=1.0.0
graphrag>=0.3.0
txtai>=7.0.0
`)
	names := detectNames(t, d, reqs)

	for _, want := range []string{"LanceDB", "Haystack", "RAGFlow", "LightRAG", "GraphRAG", "txtai"} {
		if !names[want] {
			t.Errorf("requirements.txt 未检测到 %s", want)
		}
	}
}

// TestDetectNewInferenceFrameworks 验证新增推理框架识别
func TestDetectNewInferenceFrameworks(t *testing.T) {
	dir := t.TempDir()
	d := NewDetector()

	reqs := writeFixture(t, dir, "requirements.txt", `
sglang>=0.4.0
llama-cpp-python>=0.2.0
ktransformers>=0.2.0
localai>=0.1.0
text-embeddings-inference>=1.5.0
mlx-lm>=0.20.0
`)
	names := detectNames(t, d, reqs)

	for _, want := range []string{"SGLang", "llama.cpp", "KTransformers", "LocalAI", "TEI", "MLX"} {
		if !names[want] {
			t.Errorf("requirements.txt 未检测到 %s", want)
		}
	}
}

// TestDetectNewMonitoring 验证新增可观测性组件识别
func TestDetectNewMonitoring(t *testing.T) {
	dir := t.TempDir()
	d := NewDetector()

	reqs := writeFixture(t, dir, "requirements.txt", `
helicone>=1.0.0
arize-phoenix>=4.0.0
traceloop-sdk>=0.15.0
portkey>=0.1.0
braintrust>=0.1.0
`)
	names := detectNames(t, d, reqs)

	for _, want := range []string{"Helicone", "Arize Phoenix", "Traceloop OpenLLMetry", "Portkey", "Braintrust"} {
		if !names[want] {
			t.Errorf("requirements.txt 未检测到 %s", want)
		}
	}
}

// TestDetectPackageJSON 验证 package.json 中新组件识别（n8n / Flowise / MCP / Open WebUI）
func TestDetectPackageJSON(t *testing.T) {
	dir := t.TempDir()
	d := NewDetector()

	pkg := writeFixture(t, dir, "package.json", `{
  "name": "ai-app",
  "dependencies": {
    "n8n": "^1.0.0",
    "flowiseai": "^2.0.0",
    "@modelcontextprotocol/sdk": "^1.0.0",
    "open-webui": "^0.1.0"
  }
}`)
	names := detectNames(t, d, pkg)

	for _, want := range []string{"n8n", "Flowise", "MCP SDK", "Open WebUI"} {
		if !names[want] {
			t.Errorf("package.json 未检测到 %s", want)
		}
	}
}

// TestDetectNewAPIKeys 验证新增 API Key 环境变量识别
func TestDetectNewAPIKeys(t *testing.T) {
	dir := t.TempDir()
	d := NewDetector()

	env := writeFixture(t, dir, ".env", `
DEEPSEEK_API_KEY=sk-test123
GROQ_API_KEY=gsk-test456
OPENROUTER_API_KEY=sk-or-test789
SILICONFLOW_API_KEY=sk-sf-test
VOLCENGINE_API_KEY=ak-test
ARK_API_KEY=ak-ark-test
XAI_API_KEY=xai-test
TOGETHER_API_KEY=tg-test
`)
	names := detectNames(t, d, env)

	for _, want := range []string{"DeepSeek API Key", "Groq API Key", "OpenRouter API Key", "SiliconFlow API Key", "Volcengine API Key", "Ark (Doubao) API Key", "xAI API Key", "Together AI Key"} {
		if !names[want] {
			t.Errorf(".env 未检测到 %s", want)
		}
	}
}

// TestComponentCount 记录当前组件库规模
func TestComponentCount(t *testing.T) {
	fwCount := len(types.CommonAIFrameworks)
	keyCount := len(types.APIKeyPatterns)

	if fwCount < 60 {
		t.Errorf("框架数量异常偏少: %d", fwCount)
	}
	if keyCount < 25 {
		t.Errorf("API Key 模式数量异常偏少: %d", keyCount)
	}
	t.Logf("框架: %d, API Key 模式: %d", fwCount, keyCount)
}
