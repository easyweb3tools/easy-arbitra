# Easy Arbitra Implementation Plan (Codex-Executable)

## 0. Plan Contract
- Goal: 在 Hackathon 截止前交付一个可演示、可部署、可扩展的 Smart Wallet Analyzer。
- Scope: 同时覆盖 `backend`、`frontend`、`docker/ops`。
- Execution rule for Codex:
1. 严格按阶段顺序执行。
2. 每阶段先完成 `Backend`，再联调 `Frontend`，最后补 `Docker/Ops`。
3. 每阶段结束必须满足对应验收标准 (Definition of Done)。

## 1. Current Baseline (As-Is)
- Backend skeleton 已存在：`Go + Gin + GORM + Viper + Zap`。
- 已有核心表模型：`wallet`, `market`, `wallet_score`, `ai_analysis_report`。
- 已有 API:
1. `GET /healthz`
2. `GET /api/v1/wallets`
3. `GET /api/v1/wallets/:id`
4. `GET /api/v1/markets`
5. `GET /api/v1/markets/:id`
6. `POST /api/v1/ai/analyze/:wallet_id`
7. `GET /api/v1/ai/report/:wallet_id`
- AI 分析层当前为可替换 mock analyzer。

## 2. Target Architecture (To-Be)
- Backend: 事实层采集 + 三层归因计算 + AI 报告服务。
- Frontend: Dashboard + Wallet Profile + AI Analysis + Methodology。
- Docker/Ops: 本地一键启动、环境隔离、可观测性基础、CI 可复用构建。

## 3. Workstream A - Backend Plan

### Phase A1 (Foundation Hardening, 1-2 days)
- Tasks:
1. 新增 migrations: `token`, `trade_fill`, `offchain_event`, `wallet_features_daily`。
2. 引入中间件: request_id, structured_error, rate_limit, CORS。
3. 扩展 API 查询参数: pagination/filter/sort。
4. 增加 seed 命令: 钱包、市场、评分样例数据。
- Deliverables:
1. `backend/migrations/00x_*.sql`
2. `backend/internal/api/middleware/*`
3. `backend/cmd/seed/main.go`
- DoD:
1. `go test ./...` 通过。
2. 本地可通过 seed 数据看到非空列表响应。

### Phase A2 (Data Ingestion + Layer1, 2-4 days)
- Tasks:
1. 实现 `Gamma/DataAPI/CLOB/Subgraph` 客户端基础版。
2. 实现 worker: `MarketSyncer`, `TradeSyncer`, `OffchainEventSyncer`。
3. 计算 Layer1: realized/unrealized/trading/maker/fee 分解。
4. 提供 wallet profile 聚合接口。
- Deliverables:
1. `backend/internal/client/*`
2. `backend/internal/worker/*`
3. `backend/internal/service/pnl_service.go`
- DoD:
1. 同一交易重复拉取不产生重复记录 (upsert idempotent)。
2. Wallet profile 返回 Layer1 指标。

### Phase A3 (Layer2 + AI Report, 3-5 days)
- Tasks:
1. 构建 `wallet_features_daily` (至少 10 个核心特征)。
2. 实现策略分类 v1: `market_maker/event_trader/quant/arb_hedge/noise/lucky`。
3. AI 报告缓存策略 (24h) + 历史报告列表。
4. 用 Bedrock 实现真实 analyzer，替换 mock。
- Deliverables:
1. `backend/internal/service/feature_service.go`
2. `backend/internal/service/classification_service.go`
3. `backend/internal/ai/bedrock_client.go` (real implementation)
- DoD:
1. `/api/v1/ai/report/:wallet_id` 返回结构化三层报告。
2. 重复请求 24h 内命中缓存。

### Phase A4 (Layer3 + Alerts, 3-5 days)
- Tasks:
1. 实现 Δt timing 计算与事件研究统计。
2. 实现 anomaly engine 与 evidence 输出。
3. 增加 explanations/disclosures API。
- Deliverables:
1. `backend/internal/service/info_edge_service.go`
2. `backend/internal/service/anomaly_service.go`
3. `backend/internal/api/handler/explanation_handler.go`
- DoD:
1. 可生成至少 3 类异常告警。
2. 每条告警都可回放证据。

## 4. Workstream B - Frontend Plan

### Phase B1 (UI Foundation, parallel with A1)
- Tasks:
1. 初始化 `frontend` (Next.js App Router + TailwindCSS + TypeScript)。
2. 搭建全局布局: sidebar/header/footer。
3. 建立 API client 与 typed models。
- Deliverables:
1. `frontend/src/app/layout.tsx`
2. `frontend/src/lib/api.ts`
3. `frontend/src/lib/types.ts`
- DoD:
1. 首页与 wallets/markets 页面可访问。
2. 通过 mock 或 seed API 完成基础列表展示。

### Phase B2 (Core Pages, parallel with A2)
- Tasks:
1. Dashboard: 平台概览 + strategy distribution。
2. Wallet Profile: PnL 分解、交易列表、持仓列表。
3. Leaderboard: 过滤、排序、分页。
- Deliverables:
1. `frontend/src/app/page.tsx`
2. `frontend/src/app/wallets/[address]/page.tsx`
3. `frontend/src/app/leaderboard/page.tsx`
- DoD:
1. Wallet 页面展示 Layer1 指标。
2. 列表页支持 query 参数联动。

### Phase B3 (AI + Explainability, parallel with A3/A4)
- Tasks:
1. AIAnalysisPanel + AIInsightCard。
2. ThreeLayerPanel + InfoEdgeTimeline + EvidencePanel。
3. Methodology + DisclosureBanner 全站接入。
- Deliverables:
1. `frontend/src/components/ai/*`
2. `frontend/src/components/wallet/ThreeLayerPanel.tsx`
3. `frontend/src/app/methodology/page.tsx`
- DoD:
1. 每个评分页显示 disclosure。
2. 钱包页支持查看 AI 报告历史。

## 5. Workstream C - Docker/Ops Opportunities

### C1. Local DevOps (必须)
- Opportunities:
1. `docker-compose.yml` 一键启动 `postgres + backend + frontend`。
2. 开发环境热更新：backend 使用 `air`，frontend 使用 `next dev`。
3. `.env.example` 统一配置注入。
- Expected output:
1. `docker-compose.yml`
2. `backend/Dockerfile` (dev/prod multi-stage)
3. `frontend/Dockerfile` (dev/prod multi-stage)

### C2. Build & Release (高价值)
- Opportunities:
1. GitHub Actions: lint/test/build 镜像。
2. 镜像标签策略：`sha-<short>` + `latest`。
3. SBOM/镜像扫描（如 Trivy）用于安全加分。
- Expected output:
1. `.github/workflows/ci.yml`
2. `.github/workflows/release.yml`

### C3. Runtime Reliability (加分项)
- Opportunities:
1. 健康检查和启动顺序控制（readiness/liveness）。
2. 基础监控: Prometheus metrics + structured logs。
3. 数据库备份脚本与恢复演练说明。
- Expected output:
1. `ops/monitoring/*`
2. `ops/backup/*`
3. `docs/runbook.md`

## 6. Codex Task Queue (Machine-Friendly)

```yaml
execution_order:
  - A1
  - B1
  - C1
  - A2
  - B2
  - C2
  - A3
  - B3
  - A4
  - C3

milestones:
  A1:
    status: completed
    blockers: []
    next_actions: []
  B1:
    status: completed
    blockers: []
  C1:
    status: completed
    blockers: []
  A2:
    status: completed
    blockers: []
  B2:
    status: completed
    blockers: []
  C2:
    status: completed
    blockers: []
  A3:
    status: completed
    blockers: []
  B3:
    status: completed
    blockers: []
  C3:
    status: completed
    blockers: []
  A4:
    status: completed
    blockers: []

acceptance_gates:
  - name: compile
    command: "cd backend && go test ./..."
  - name: api_health
    command: "curl -sf http://localhost:8080/healthz"
  - name: e2e_stack
    command: "docker compose up -d && docker compose ps"
```

## 7. Risks and Mitigation
- Risk: Polymarket API schema 变动导致 ingestion 失败。
- Mitigation: 客户端加 contract tests + fallback parser。
- Risk: AI 成本超预算。
- Mitigation: 24h 缓存 + 分层分析（普通请求走轻量模板）。
- Risk: 告警误报导致产品风险。
- Mitigation: 默认使用 anomaly wording + 证据可回放 + 强制免责声明。

## 8. Immediate Next Sprint (建议本周执行)
1. 已完成 API 错误语义统一与 request_id 透出。
2. 已完成多格式数据源适配与 HTTP 重试测试。
3. 已完成 UI 自动化 smoke（Playwright）接入 CI。
