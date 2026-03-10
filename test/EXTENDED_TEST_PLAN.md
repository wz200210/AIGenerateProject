# AI 组件扩展测试计划

## 测试目标
继续推进 AIGenerateProject 扫描器的真实组件测试，覆盖更多 AI 服务类型。

## 测试组件清单

### 阶段 1: pip 可安装组件
- [x] Jupyter Notebook - 监控工具类 ✅ (端口 8888, PID 121535)
- [x] Streamlit - 部署工具类 ✅ (端口 8501, PID 121609)
- [x] Gradio - 部署工具类 ⚠️ (端口 7860 运行中，扫描器未识别进程名)
- [ ] LangChain - Agent/RAG 框架类

### 阶段 2: 其他安装方式
- [ ] Ollama - LLM 推理服务 (curl 安装)
- [ ] Weaviate - 向量数据库 (Docker 或二进制)

### 阶段 3: 模拟服务测试
- [ ] vLLM - 需要 GPU，使用端口模拟
- [ ] Milvus - Docker 网络问题，使用端口模拟

## 当前状态
- Chroma: ✅ 真实测试通过
- 其他: ⏳ 进行中

---

