## SportStyle AI Explainer

### 1. 项目目标

SportStyle AI Explainer 是一个面向 Amazon Nova Hackathon 的可执行 Demo 项目，目标是在 3 分钟内清楚展示以下能力：

1. Amazon Nova 作为分析与编排大脑，能够理解用户提供的钱包身份信息。
2. Nova 通过 MCP 调用 Go 工具链，拉取真实体育交易数据并完成结构化分析。
3. 前端将分析过程和结果可视化，输出一个可解释、可复现的研究报告。

这个项目不做“预测未来收益”，只做“解释历史交易风格”。

---

### 2. 一句话定位

用户输入一个钱包地址，或直接粘贴一个 Polymarket 用户详情页链接；Nova Agent 负责解析目标用户、规划调用链路、输出解释日志；Go MCP Tools 负责从 Polymarket 拉取该钱包的体育市场交易数据，并用确定性规则计算风格指标；前端最终展示钱包风格卡片、雷达图和可解释报告。

---

### 3. Demo 场景

本项目只支持一个主场景：分析单个 Polymarket 钱包的体育交易风格。

系统只接受以下两种输入：

1. 钱包地址
2. Polymarket 用户详情页链接

完整 Demo 流程固定如下：

1. 用户在首页输入钱包地址或 Polymarket 用户详情页链接。
2. Nova Agent 解析钱包身份并启动分析流程。
3. Nova 调用 `resolve_wallet_target` 标准化目标钱包信息。
4. Nova 调用 `fetch_sports_trades` 获取该钱包的 NBA 相关交易样本。
5. Nova 调用 `calculate_style_metrics` 计算钱包风格指标。
6. Nova 调用 `build_report_payload` 生成结构化报告数据。
7. 前端展示决策日志、钱包风格卡片、雷达图和自然语言解释。

---

### 4. 核心设计原则

#### 4.1 Nova 是大脑

Nova 负责：

1. 解析用户输入的钱包地址或 Polymarket 用户链接。
2. 决定工具调用顺序。
3. 汇总每一步结果。
4. 输出自然语言解释和决策日志。

Nova 不负责：

1. 改写指标定义。
2. 动态改变计算口径。
3. 直接充当数据源。

#### 4.2 MCP Tools 是事实引擎

所有核心指标由 Go MCP Tools 按固定规则计算，保证：

1. 结果可复现。
2. 结果可验证。
3. 前后演示一致。

#### 4.3 前端是证据层

前端只做三件事：

1. 展示分析过程。
2. 展示结构化结果。
3. 展示 Nova 的解释内容。

---

### 5. 系统架构

```text
用户 (Next.js)
  ↓
前端调用 Nova / Bedrock
  ↓
Nova Agent
  ├── 理解用户意图
  ├── 调用 MCP Tools
  └── 输出解释和日志
  ↓
Go MCP Server
  ├── resolve_wallet_target
  ├── fetch_sports_trades
  ├── calculate_style_metrics
  └── build_report_payload
  ↓
Polymarket Data API / CLOB WebSocket
  ↓
Next.js Dashboard 展示结果
```

调度策略：

1. 默认由用户主动触发分析。
2. Go 服务内部维护体育市场 WebSocket 长连接。
3. 如需刷新缓存，只使用分钟级调度，不依赖秒级任务。

---

### 6. 技术选型

#### 6.1 AI 层

1. Amazon Nova 2 Lite：用于推理、工具编排、生成解释。
2. Amazon Nova Act：用于 Agent 工作流与后续 UI Automation 扩展。
3. MCP：用于标准化暴露后端工具。

#### 6.2 后端

1. Go 1.23
2. `mark3labs/mcp-go`
3. Polymarket Data API / CLOB WebSocket

#### 6.3 前端

1. Next.js 15 App Router
2. Tailwind CSS
3. shadcn/ui
4. Recharts

---

### 7. 必须实现的 4 个 MCP Tool

#### 7.1 `resolve_wallet_target`

用途：

将用户输入标准化为可查询的钱包目标。

输入：

1. 钱包地址
2. 或 Polymarket 用户详情页链接

输出：

1. 标准化钱包地址
2. 用户名或展示名称（如果能解析到）
3. 原始输入类型

要求：

1. 支持识别钱包地址。
2. 支持识别 Polymarket 用户详情页链接。
3. 统一返回结构化 JSON。

#### 7.2 `fetch_sports_trades`

用途：

拉取指定钱包在 Polymarket 体育市场中的 NBA 交易样本。

输入：

1. 标准化钱包地址
2. `sport`: 固定为 `nba`
3. `limit`: 返回的交易条数

输出：

1. 钱包地址
2. 市场标识
3. 交易时间
4. 交易金额
5. 市场成交量
6. 市场开始时间或可用于近似计算入场时间的时间字段

要求：

1. 只返回体育市场样本。
2. 优先返回 NBA 市场。
3. 结果必须与指定钱包绑定。
4. 返回结构化 JSON。

#### 7.3 `calculate_style_metrics`

用途：

对单个钱包计算固定风格指标。

输入：

1. 钱包地址
2. 该钱包关联交易样本

输出：

1. `entry_timing_hours`
2. `size_ratio_pct`
3. `win_rate` 或 `roi`

固定口径：

1. `entry_timing_hours`
   - 定义：钱包首次入场时间相对市场开始时间的小时差

2. `size_ratio_pct`
   - 定义：钱包交易金额占对应市场成交量的比例

3. `win_rate` 或 `roi`
   - 二选一，开发时只保留一个口径
   - 若结算数据不足，优先保留 `roi`

要求：

1. 计算逻辑固定写在 Go 工具中。
2. 同一输入必须得到稳定输出。
3. 返回结构化 JSON，供前端直接绘图。

#### 7.4 `build_report_payload`

用途：

将计算结果转换成前端和 Nova 可直接使用的结构化报告数据。

输入：

1. 钱包基础信息
2. 风格指标

输出：

1. 钱包卡片数据
2. 雷达图数据
3. 报告摘要字段
4. 供 Nova 解释的结构化上下文

要求：

1. 不在该 Tool 中重新计算指标。
2. 只负责组织结构化结果。
3. 输出必须适合前端直接消费。

---

### 8. Nova Agent 调用规则

Nova Agent 必须按固定顺序执行：

1. 调用 `resolve_wallet_target`
2. 调用 `fetch_sports_trades`
3. 调用 `calculate_style_metrics`
4. 调用 `build_report_payload`

Nova Agent 输出内容必须包含：

1. 当前步骤日志
2. 调用的 Tool 名称
3. 每一步的核心结果
4. 最终自然语言解释

Nova Agent 的解释目标：

用自然语言回答“这个钱包在体育市场中表现出什么交易风格”。

示例解释：

`该钱包在体育市场中通常会在市场开始后 2.3 小时内完成首次入场，单笔仓位占市场成交量 4.1%，整体 ROI 较高，因此呈现出明显的早入大额型风格。`

---

### 9. 前端页面范围

#### 9.1 首页 `/`

功能：

1. 输入钱包地址或 Polymarket 用户详情页链接
2. 点击按钮启动分析
3. 跳转到 `/dashboard`

页面必须有：

1. 标题
2. 钱包地址 / 链接输入框
3. 启动按钮

#### 9.2 仪表盘 `/dashboard`

功能：

1. 显示 Nova 决策日志
2. 显示钱包分析结果
3. 显示雷达图
4. 显示自然语言解释

页面必须有：

1. 日志面板
2. 钱包卡片
3. 雷达图
4. 报告摘要区

如果时间不足，可以只实现这两个页面，不再实现其他页面。

---

### 10. 数据与口径约束

为了保证 Demo 稳定，必须遵守以下规则：

1. 只分析体育市场，不扩展到政治或加密市场。
2. 只优先支持 NBA，不同时覆盖过多联赛。
3. 只保留 3 个核心指标，不增加 embedding 或黑箱特征。
4. 所有指标口径固定写入 Go 代码，不交给 Agent 动态推断。
5. 结果以“解释历史风格”为目标，不输出投资建议。
6. 用户输入必须先标准化为钱包地址，再进入交易分析流程。

---

### 11. 交付物

本项目最终必须交付以下内容：

1. 一个本地可运行的 Go MCP Server
2. 一个本地可运行的 Next.js 前端
3. 一条完整可运行的 Nova 调用链
4. 一个稳定可复现的 3 分钟 Demo 流程
5. 一张可用于 Devpost 的系统架构图

---

### 12. 开发计划

#### Day 1

1. 初始化 Go MCP Server
2. 注册 4 个 Tool 骨架
3. 接入 Polymarket 数据源

#### Day 2

1. 完成 `resolve_wallet_target`
2. 完成 `fetch_sports_trades`
3. 固化 3 个指标口径

#### Day 3

1. 完成 `calculate_style_metrics`
2. 完成 `build_report_payload`
3. 跑通 Nova Agent -> MCP Tool 调用链并输出结构化日志

#### Day 4

1. 完成首页 `/`
2. 完成仪表盘 `/dashboard`
3. 接入雷达图和钱包卡片

#### Day 5

1. 打磨文案和解释语句
2. 固定 Demo 演示脚本
3. 录制 3 分钟视频

#### Day 6

1. 联调
2. 修复演示问题
3. 准备提交材料

---

### 13. 开发指令

基于本方案实现：

1. 用 Go + `mark3labs/mcp-go` 实现 MCP Server
2. 实现 `resolve_wallet_target`、`fetch_sports_trades`、`calculate_style_metrics`、`build_report_payload`
3. 用 Next.js 15 实现 `/` 和 `/dashboard`
4. 用 Recharts 渲染 3 维度雷达图
5. 保留 Nova 作为编排与解释大脑的调用示例
