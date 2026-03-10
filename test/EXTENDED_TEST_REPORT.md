# AI 组件扩展测试报告

## 测试时间
2026-03-09

## 测试环境
- 系统: Ubuntu 24.04 LTS
- Python: 3.12.3
- 磁盘: 28GB 剩余
- 扫描器版本: v0.4.0

---

## 已测试组件

### 1. Jupyter Notebook ✅
- **类型**: monitoring
- **安装**: `pip3 install jupyter`
- **启动**: `jupyter notebook --ip=0.0.0.0 --port=8888`
- **端口**: 8888
- **PID**: 121535
- **检测结果**: ✅ 成功识别
- **扫描器输出**:
```
[monitoring]
  • Jupyter [low]
    Source: /proc/121535
    PID: 121535 | Default Ports: 8888 | Exe: python3.12
```

### 2. Streamlit ✅
- **类型**: deployment
- **安装**: `pip3 install streamlit`
- **启动**: `streamlit run app.py --server.port=8501`
- **端口**: 8501
- **PID**: 121609
- **检测结果**: ✅ 成功识别
- **扫描器输出**:
```
[deployment]
  • Streamlit [low]
    Source: /proc/121609
    PID: 121609 | Default Ports: 8501 | Exe: python3.12
```

### 3. Gradio ⚠️
- **类型**: deployment
- **安装**: `pip3 install gradio`
- **启动**: `python3 gradio_app.py` (端口 7860)
- **端口**: 7860 ✅ 服务运行正常
- **PID**: 121620
- **检测结果**: ❌ 扫描器未识别
- **原因分析**: Gradio 进程名显示为 `python3`，扫描器的进程模式匹配 `\bgradio\b` 未能命中
- **进程信息**: `python3 /tmp/gradio_app.py`

### 4. LangChain ⚠️
- **类型**: llm_framework
- **安装**: `pip3 install langchain langchain-community`
- **检测结果**: N/A
- **说明**: LangChain 是开发库，非独立服务进程，扫描器设计用于检测运行中的服务

### 5. Chroma (历史测试) ✅
- **类型**: vector_database
- **状态**: 已在前期测试中验证成功
- **端口**: 8000

---

## 扫描结果汇总

| 组件 | 类型 | 安装方式 | 检测状态 | 备注 |
|------|------|---------|---------|------|
| Jupyter | monitoring | pip | ✅ | 进程名匹配成功 |
| Streamlit | deployment | pip | ✅ | 进程名匹配成功 |
| Gradio | deployment | pip | ⚠️ | 服务运行但进程名未匹配 |
| Chroma | vector_database | pip | ✅ | 历史测试 |
| LangChain | llm_framework | pip | N/A | 非服务进程 |

---

## 发现的问题

### 1. Gradio 进程识别问题
Gradio 应用运行时进程名为 `python3` 而非 `gradio`，导致扫描器的正则 `\bgradio\b` 无法匹配。

**建议**: 增强语义分析器，检测 Python 进程中是否导入 gradio 模块。

### 2. 端口扫描显示异常
扫描器输出显示 `Ports scanned: 0`，但实际服务在运行。

**待排查**: 端口扫描逻辑是否正常。

---

## 下一步测试计划

- [ ] Ollama - LLM 推理服务 (安装中)
- [ ] 使用模拟端口测试更多组件类型
- [ ] 验证 Docker 容器检测功能

---

*报告生成时间: 2026-03-09*
