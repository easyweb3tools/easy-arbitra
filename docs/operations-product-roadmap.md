# 运营化产品路线（Codex 执行文档）

## 1. 目标与范围
目标：把当前“分析工具”升级为“可持续运营内容产品”，核心围绕「高潜力高收益钱包」获取流量与留存。

本期范围（P0-P2）：
1. 榜单化：每日可消费榜单 + 首页运营位。
2. 传播化：钱包分享卡 + 可追踪来源参数。
3. 订阅化：关注钱包与更新通知（站内优先）。
4. 透明化：数据更新时间、规则解释、样本规模。

## 2. 当前基础（已完成）
1. 潜力钱包接口：`GET /api/v1/wallets/potential`
2. AI 批处理：`ai_batch_analyzer`
3. 首页/钱包页：已突出潜力钱包
4. 中英文：核心页面可切换 `en/zh`

## 3. 分期开发计划

### P0（1-2 天）：运营首页 V1
后端新增：
1. `GET /api/v1/ops/highlights`
- 返回字段：
  - `as_of`（UTC）
  - `new_potential_wallets_24h`
  - `top_realized_pnl_24h[]`（wallet_id, address, pnl, trades, has_ai_report）
  - `top_ai_confidence[]`（wallet_id, smart_score, info_edge_level, summary）

前端改造：
1. 首页新增模块：
- 今日新增潜力钱包
- 24h 收益 Top5
- AI 高置信 Top5
2. 每个卡片提供 CTA：`查看详情` / `加入关注`

验收标准：
1. 首页 3 个运营模块可见，加载失败有降级文案。
2. 每个模块至少支持 5 条数据与跳转。

### P1（2-3 天）：传播与增长
后端新增：
1. `GET /api/v1/wallets/:id/share-card`
- 返回钱包分享所需聚合字段（昵称、收益、交易数、策略、AI 标签、更新时间）

前端新增：
1. 钱包详情页增加“分享卡片”区域（复制链接 + 预览图占位）
2. 全站链接统一支持 UTM 参数透传：`utm_source/utm_campaign`

验收标准：
1. 可复制分享链接，带 UTM 参数。
2. 分享卡片信息与详情页数据一致。

### P2（3-5 天）：关注与运营触达
数据表新增：
1. `watchlist`（id, wallet_id, user_fingerprint/session_id, created_at）
2. `wallet_update_event`（wallet_id, event_type, payload, created_at）

后端新增：
1. `POST /api/v1/watchlist`（添加关注）
2. `GET /api/v1/watchlist`（我的关注列表）
3. `GET /api/v1/watchlist/feed`（关注更新流）

前端新增：
1. 钱包卡片和详情页支持“关注/取消关注”
2. 新页面：`/watchlist`

验收标准：
1. 用户可关注钱包并在 watchlist 页看到更新。
2. 更新流至少包含：新 AI 报告、异常告警、收益跃迁。

## 4. 技术实现约束
1. 保持现有 `Next.js App Router + Gin + PostgreSQL` 架构。
2. 所有新接口必须提供分页与默认排序。
3. 所有文案接入 `frontend/src/lib/i18n.ts`（en/zh）。
4. 新增 SQL 使用索引友好写法；避免全表扫描长查询。

## 5. Codex 执行顺序（严格）
1. 先后端接口与 SQL（含测试）
2. 再前端页面与组件
3. 再 i18n 文案补齐
4. 最后运行：
- `cd backend && go test ./...`
- `cd frontend && npm run build`
- `docker compose up -d --build`

## 6. 运营指标（上线后跟踪）
1. 首页点击率：潜力钱包模块 CTR
2. 钱包详情停留时长
3. 分享率：分享按钮点击/复制次数
4. 关注转化率：详情页到 watchlist
5. 回访率：1日/7日留存

## 7. Definition of Done
1. P0/P1/P2 的接口与页面都可在本地 compose 环境跑通。
2. `runbook.md` 更新新增运营接口和排障说明。
3. 至少提供 1 份“运营周报模板”JSON（可由接口直接生成字段）。
