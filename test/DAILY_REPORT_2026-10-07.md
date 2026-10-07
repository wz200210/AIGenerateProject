# AI 组件测试日报 - 2026-10-07

## 项目状态总览

- **最后代码提交**: 2026-03-13（67e1873, Skill 扫描功能文档）
- **停滞天数**: 208 天
- **组件验证进度**: 9/58（16%）
- **本日变化**: 无新增代码、无新增测试记录

## 已验证组件（9 个）

基于 test/INSTALLATION_TEST.md、REAL_INSTALL_TEST.md 及 TEST_PROGRESS_2026-03-10.md 记录：

1. **Ollama** — LLM推理（v0.17.7，进程+端口+版本API）✅
2. **Chroma** — 向量数据库（pip 安装，端口 8000）✅
3. **Qdrant** — 向量数据库（client 1.17.0，端口 6333）✅
4. **PyTorch** — ML框架（v2.10.0，语义分析器）✅
5. **TensorFlow** — ML框架（v2.20.0，语义分析器）✅
6. **LlamaIndex** — RAG工具（语义分析器）✅
7. **Gradio** — 部署工具（端口 7860，语义分析器修复后识别）✅
8. **Streamlit** — 部署工具（端口 8501，进程匹配）✅
9. **Jupyter** — 监控工具（端口 8888，进程匹配）✅

另：OpenAI 代理（mock）、LangChain（mock）、OpenClaw/Claude Skill 扫描有记录但非完整验证。

注：Docker 环境不可用（registry 拉取超时），Qdrant 等组件 Docker 复测仍被阻塞。

## 扫描结果摘要

- 最近有效扫描时间：2026-03-10（9 进程 / 1 端口 / 10 组件检出）
- 扫描日志、结果文件无更新

## 待推进事项（无外部依赖，可立即执行）

1. Gradio 进程识别修复
2. pip 进程误报过滤
3. 重复实例去重

## 待测组件（需 Docker 环境）

- vLLM、Weaviate、Qdrant（Docker 复测）

## 下一步计划

项目持续停滞第 208 天。如需恢复推进，建议优先完成上述 3 个无依赖修复，可一次性提升扫描准确率。
