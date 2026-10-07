# AI组件测试日报 - 2026-09-16

**项目名称**: AIGenerateProject（AI组件扫描器）  
**扫描器版本**: v0.4.0  
**报告时间**: 2026-09-16 18:14 (Asia/Shanghai)  
**cron任务ID**: b099e650-1e3c-4f6b-8f74-c1d959e38b4a

---

## 一、状态摘要

**项目持续停滞，无任何新增进展。**

- 距离上次代码提交已 **187天**（6个月+，3月13日，`67e1873`）
- 距离上次实质性测试进展已 **6个月**（3月初的模拟服务验证）
- 期间无任何代码变更、无新组件验证、无Bug修复

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

### 实质性测试记录（2026年3月）
- **Chroma**: 模拟服务+真实安装（v1.4.1，端口8000），三重验证通过
- **Qdrant/Milvus**: Python http.server 模拟服务验证通过；Docker 真实部署因网络超时失败
- **Ollama**: 真实安装 v0.17.7，CLI+HTTP 双验证
- **PyTorch/TensorFlow/LlamaIndex**: 语义分析器验证通过
- **Streamlit/Jupyter/MLflow**: 进程匹配或端口扫描正常
- **Gradio**: ⚠️ 语义分析器可检测，但进程名识别失败（`python3` vs `gradio`）

### 环境限制（长期未解决）
- Docker Hub 访问受限，无法真实部署向量数据库
- 无 GPU 环境，vLLM 等推理服务未测试

---

## 四、已知问题（全部搁置187天）

1. **Gradio 进程识别失败** — Python进程名显示为`python3`，正则`\bgradio\b`未命中命令行
2. **扫描器重复检测** — 同一组件被进程扫描+网络扫描各发现一次，无去重逻辑
3. **pip进程误报** — Python包管理进程被误识别为AI组件
4. **版本号探测不完整** — 部分Mock服务显示`unknown`
5. **Docker不可用** — 真实容器部署受阻

---

## 五、可立即推进的无依赖项

以下修复不需要外部环境，仅需代码改动：

1. **Gradio识别修复** — 扩展匹配逻辑：除进程名外，增加对命令行参数、脚本文件名的检查
2. **pip误报过滤** — 在进程扫描规则中排除 `pip`/`pip3` 相关进程
3. **重复实例去重** — 当进程扫描和网络扫描发现同一组件时，合并结果并标注多信号源

---

## 六、未测试组件概览（49个）

- **向量数据库**: Weaviate, Pinecone, Vald, Vespa 等
- **LLM框架**: vLLM, HuggingFace TGI, LMDeploy 等
- **ML框架**: JAX, MXNet, PaddlePaddle 等
- **MLOps**: Kubeflow, Airflow, Prefect 等
- **Agent框架**: AutoGen, CrewAI, LangGraph 等

---

## 七、下一步计划

**如果项目重新启动，优先行动：**
1. 修复3个无依赖Bug（预计各需1-2小时）
2. 覆盖率从16%提升至20%+（新增3-4个组件验证）
3. 解决Docker Hub访问问题，恢复向量数据库真实部署测试
4. 补充 vLLM、Weaviate、Transformers、Qdrant 的测试

---

*本报告由 cron 任务自动生成。*
