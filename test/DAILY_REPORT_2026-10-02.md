# AI组件测试日报 — 2026-10-02

## 项目状态
- **停滞天数：203天**（最后代码提交 2026-03-13，commit c33c786）
- 已验证组件：**9/58**（16%）
- 今日无新代码提交、无新测试进展

## 测试日志摘要
- `test/REAL_INSTALL_TEST.md`：Chroma（✅运行中 v1.4.1）、Qdrant（Docker拉取超时失败），最后更新 2026-03-05
- `test/INSTALLATION_TEST.md`：4个模拟服务端口检测全部通过（Qdrant/Milvus/Jupyter/MLflow），最后更新 2026-03-05
- `test/INSTALL_TEST_LOG.md`：Chroma/OpenAI API/LangChain 三重验证成功，最后更新 2026-03-05
- `test/TEST_PROGRESS_2026-03-10.md`：8组件验证通过（Ollama/PyTorch/TensorFlow/LlamaIndex/Gradio/Streamlit/Jupyter/Chroma），最后更新 2026-03-10
- `test/scan_result.txt`：Chroma/LangChain 检测记录（含重复实例问题）

## 待推进事项（无外部依赖，可立即开工）
1. Gradio 进程识别修复
2. pip 进程误报过滤
3. 重复实例去重

## 下一步计划
- 三个无依赖修复任选其一启动，预计每个 1-2 天工作量
- Docker 环境恢复后复测 Qdrant（当前 registry 拉取超时）
- 高优先级待测：vLLM、Weaviate、Transformers
