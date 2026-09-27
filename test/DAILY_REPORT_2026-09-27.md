# AI组件测试日报 - 2026-09-27

**项目名称**: AIGenerateProject（AI组件扫描器）  
**扫描器版本**: v0.4.0  
**报告时间**: 2026-09-27 18:14 (Asia/Shanghai)  
**cron任务ID**: b099e650-1e3c-4f6b-8f74-c1d959e38b4a

---

## 一、状态摘要

**项目持续停滞，无任何新增进展。**

- 距离上次实质性代码提交已 **198天**（3月13日，`67e1873`）
- 期间仅有 cron 日报提交，无代码变更
- 无新组件验证、无Bug修复、无功能迭代

---

## 二、已验证组件（9个 / 58个配置，约16%）

| 类别 | 组件 | 检测方式 | 状态 |
|------|------|---------|------|
| 向量数据库 | Chroma, Qdrant, Milvus | 端口扫描/模拟服务 | ✅ |
| 监控/部署 | Jupyter, MLflow, Streamlit | 进程匹配/端口扫描 | ✅ |
| LLM框架 | OpenAI API代理, LangChain, Ollama | 进程+端口/语义分析 | ✅ |
| ML框架 | PyTorch, TensorFlow | 语义分析器 | ✅ |
| RAG工具 | LlamaIndex | 语义分析器 | ✅ |
| 部署工具 | Gradio | 语义分析器（进程识别仍失败） | ⚠️ |

---

## 三、测试日志摘要

### 最后一次实质性测试（2026年3月）
- **Chroma**: 模拟服务+真实安装（chroma 1.4.1），端口8000，三重验证通过
- **Ollama**: 真实安装v0.17.7，CLI+HTTP双验证
- **PyTorch/TensorFlow/LlamaIndex**: 语义分析器验证通过
- **Streamlit/Jupyter/MLflow**: 进程匹配或端口扫描正常
- **Qdrant/Milvus**: 模拟服务验证（INSTALLATION_TEST.md 记录）；Docker真实部署因 docker.io 网络超时失败（REAL_INSTALL_TEST.md 记录）
- **Gradio**: ⚠️ 语义分析器可检测，但进程名识别仍失败

### 环境限制（长期未解决）
- Docker Hub访问受限（registry-1.docker.io 连接超时），无法真实部署向量数据库
- 无GPU环境，vLLM等推理服务未测试

---

## 四、已知问题（全部搁置198天）

1. **Gradio 进程名识别失败** — 语义分析器可匹配，进程级检测不到
2. **pip 进程误报** — 安装进程被误识别为运行实例
3. **重复实例去重缺失** — 同一组件多个进程未合并计数

---

## 五、下一步计划（不变）

1. **可立即推进**（无外部依赖）：Gradio识别修复、pip误报过滤、重复实例去重
2. **待环境就绪后推进**：vLLM、Weaviate、Transformers、Qdrant（Docker复测）

---

## 六、结论

第198天日报。项目状态与过去190余天完全一致。三个无外部依赖的修复（Gradio识别、pip误报过滤、重复实例去重）随时可以动手，只差一个决定。
