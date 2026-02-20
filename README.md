# Easy Arbitra — Polymarket Smart Wallet Analyzer

> **An AI-powered prediction market intelligence platform that tracks profitable wallets on Polymarket and explains *how* they make money — powered by Amazon Nova AI.**

[![Amazon Nova AI Hackathon](https://img.shields.io/badge/Hackathon-Amazon%20Nova%20AI-orange)](https://amazon-nova.devpost.com/)
[![Track](https://img.shields.io/badge/Track-Agentic%20AI-blue)]()
[![Built with](https://img.shields.io/badge/Built%20with-Go%20%7C%20Next.js%20%7C%20PostgreSQL-green)]()

---

## Inspiration

Polymarket processes over **$1 billion in monthly trading volume**, yet its participants range from casual bettors to sophisticated algorithmic traders — and possibly those with non-public information advantages. While tools like Nansen and Arkham exist for DeFi, **no purpose-built analytics platform** deeply explains *why* a Polymarket wallet is profitable.

We were inspired by a simple but powerful question: **Is a profitable wallet genuinely skilled, or just lucky?** Answering this requires more than surface-level metrics like win rate or ROI. It demands a rigorous statistical framework that decomposes performance into facts, strategy patterns, and information timing advantages.

The Amazon Nova AI Hackathon gave us the perfect opportunity to build this. Nova 2 Lite's **agentic capabilities** — tool calling, extended thinking, code interpreter, and structured output — are exactly what we needed to create an autonomous AI analyst that can reason through complex, multi-step wallet attribution workflows.

---

## What It Does

Easy Arbitra is a **Three-Layer Attribution Framework** powered by Amazon Nova AI that progressively explains how any Polymarket wallet makes money:

### Layer 1: Facts (PnL Accounting)
- Decomposes wallet profits into **Trading PnL, Maker Rebates, and Fee Attribution**
- Calculates performance metrics: win rate, ROI, max drawdown, Sharpe ratio
- Identifies concentration risk: are profits spread across markets or dominated by one lucky bet?

### Layer 2: Strategy Classification
- Classifies wallets into **6 behavioral archetypes** using clustering + rule engines:
  - **Market Maker** — Earns spreads and rebates via two-sided quoting
  - **Event Trader** — Concentrates bets around specific news events
  - **Quant/Analytical** — Diversified, model-driven pricing with consistent edge
  - **Arbitrage/Hedge** — Cross-market or cross-platform hedging strategies
  - **Noise/Retail** — Chases momentum with no persistent advantage
  - **Lucky** — Small sample, one-time windfall, statistically insignificant

### Layer 3: Information Edge Attribution
- Measures **information timing advantage (Δt)** — how early a wallet trades relative to public news events
- Applies statistical tests (Wilcoxon signed-rank, bootstrap counterfactual) to distinguish:
  - **Luck-driven** — Insignificant sample, concentrated in 1-2 events
  - **Quant/Analytical** — Diversified, consistent edge across markets
  - **Processing Edge** — Systematically trades before news, but using public information faster
  - **Insider-Suspected** — Multiple converging evidence of non-public information (presented with disclaimers, never as accusations)

### Key Features
- **AI Chat Interface** — Ask questions about any wallet in natural language ("Why is this wallet profitable?")
- **Voice Interaction** — Speech-to-speech analysis powered by Nova 2 Sonic
- **Similar Wallet Discovery** — Find wallets with matching trading behavior via embedding similarity
- **Automated Data Collection** — Browser agent scrapes Polymarket pages for data not available via API
- **Compliance-First Design** — Three-tier labeling (neutral/attention/warning), evidence-replayable, mandatory disclaimers

---

## How We Built It

### Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                    Frontend (Next.js 14+)                     │
│  App Router · TailwindCSS · Recharts · AI Chat · Voice UI    │
└──────────────────────┬───────────────────────────────────────┘
                       │ REST API
┌──────────────────────▼───────────────────────────────────────┐
│                    Backend (Go + Gin)                         │
│  GORM · Viper · Zap · robfig/cron · aws-sdk-go-v2           │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐  │
│  │           Amazon Nova AI Analysis Layer                 │  │
│  │                                                        │  │
│  │  Nova 2 Lite ── Agentic AI Core                        │  │
│  │    ├── Tool Calling → Polymarket API orchestration      │  │
│  │    ├── Extended Thinking → Deep attribution reasoning   │  │
│  │    ├── Code Interpreter → Statistical tests (Python)    │  │
│  │    └── Structured Output → JSON wallet profiles         │  │
│  │                                                        │  │
│  │  Nova Multimodal Embeddings ── Behavior Vectorization   │  │
│  │    └── pgvector similarity → Wallet clustering          │  │
│  │                                                        │  │
│  │  Nova 2 Sonic ── Voice Interaction                      │  │
│  │    └── Speech-to-speech real-time analysis              │  │
│  │                                                        │  │
│  │  Nova Act ── Browser Automation                         │  │
│  │    └── Scrape Polymarket pages for non-API data         │  │
│  └────────────────────────────────────────────────────────┘  │
└──────────────────────┬───────────────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────────────┐
│              PostgreSQL + pgvector                            │
│  Time-partitioned fact tables · Wallet embeddings             │
└──────────────────────────────────────────────────────────────┘
```

### Amazon Nova AI Integration (All 4 Models)

| Nova Model | Role | How It's Used |
|------------|------|---------------|
| **Nova 2 Lite** | Agentic AI Core | The brain of the system. Uses **tool calling** to autonomously orchestrate 7+ Polymarket data tools, **extended thinking** (intensity: high) for multi-step attribution reasoning across 1M token context, **code interpreter** to run statistical tests (bootstrap, Wilcoxon, Herfindahl) in a Python sandbox, and **constrained decoding** to output structured JSON wallet profiles with >95% schema compliance. |
| **Nova Multimodal Embeddings** | Behavior Vectorization | Converts wallet trading behavior (48-dimensional feature vectors from 5 feature families) into dense embeddings. Stored in PostgreSQL with **pgvector** for cosine similarity search, enabling "find similar wallets" and anomaly detection (distance from known cluster centroids). |
| **Nova 2 Sonic** | Voice Interface | Enables **speech-to-speech** interaction via WebSocket streaming. Users can ask "Tell me about wallet 0xABC..." and receive spoken analysis reports. Supports real-time bidirectional audio with automatic voice activity detection. |
| **Nova Act** | Browser Automation | Headless browser agent that navigates Polymarket pages to collect data unavailable through APIs — market resolution details, comment sentiment, UI-only statistics. Also captures screenshots for multi-modal verification. |

### Tech Stack

- **Frontend:** Next.js 14+ (App Router), TailwindCSS, Recharts/ECharts, WebSocket audio streaming
- **Backend:** Go, Gin, GORM, Viper (config), Zap (logging), robfig/cron (scheduling)
- **Database:** PostgreSQL with pgvector extension, time-partitioned fact tables
- **AI:** Amazon Bedrock (Nova 2 Lite, Nova 2 Sonic, Nova Multimodal Embeddings), Nova Act SDK
- **Data Sources:** Polymarket Gamma API, CLOB API, Data API, Goldsky Subgraphs (Positions/Orders/Activity/OI/PNL), Polygon RPC
- **Infrastructure:** Docker, AWS (Bedrock, S3)

### Agentic AI Workflow (Nova 2 Lite)

The core innovation is the **agentic analysis loop**:

1. User submits a wallet address or asks a question
2. Nova 2 Lite receives the query with the Three-Layer Attribution system prompt
3. The agent **autonomously decides** which tools to call (trades, positions, market info, price history, offchain events, pre-computed features, similar wallets)
4. After collecting data, it uses **extended thinking** to reason through all three attribution layers
5. The **code interpreter** executes statistical tests (bootstrap counterfactual, Wilcoxon signed-rank for Δt distribution, rolling window stability)
6. Results are output as **structured JSON** via constrained decoding, then rendered as interactive dashboards
7. A natural language summary is generated for non-technical users

---

## Challenges We Ran Into

- **PnL Accounting Precision:** Polymarket's CLOB has nuances — maker rebates, fee-enabled markets, split/merge operations for CTF tokens — that make accurate PnL decomposition significantly harder than typical DEX analytics. We had to carefully design the accounting pipeline to handle all edge cases.

- **Information Edge Attribution is Inherently Noisy:** Measuring whether a wallet traded *before* a news event requires reliable off-chain event timestamps, which are imprecise. We mitigated this with statistical testing (requiring p < 0.05) rather than hard thresholds.

- **Balancing Transparency with Responsibility:** Labeling wallets as "insider-suspected" carries serious implications. We designed a three-tier compliance system where the highest-risk labels are never public — only shown as "anomaly + confidence + disclaimer" with replayable evidence.

- **Nova Agent Tool Orchestration:** Getting the agentic loop right — where Nova 2 Lite autonomously decides the sequence of 7+ tool calls — required careful prompt engineering and tool description design to avoid infinite loops or redundant API calls.

- **Embedding Heterogeneity:** Wallet behaviors are high-dimensional (48 features across 5 families). Designing meaningful embeddings that capture behavioral similarity while remaining interpretable was a core technical challenge.

---

## Accomplishments That We're Proud Of

- **Three-Layer Attribution Framework** — A rigorous, academically-grounded approach to wallet analysis that goes far beyond simple "smart money" labels. Each layer builds on the previous with increasing analytical depth.

- **Fully Agentic Analysis** — Nova 2 Lite doesn't just answer questions — it autonomously orchestrates data collection, statistical testing, and report generation through multi-step tool calling with extended thinking. This is a genuine Agentic AI application, not a simple prompt-response wrapper.

- **Statistical Rigor** — We use proper statistical tests (bootstrap counterfactual, Wilcoxon signed-rank, Herfindahl concentration index, rolling window stability) instead of arbitrary thresholds. Every classification comes with confidence intervals and p-values.

- **All 4 Nova Models Working Together** — Each Nova model serves a distinct, complementary purpose: Lite for reasoning, Embeddings for similarity, Sonic for voice, Act for data collection. This demonstrates the breadth of the Nova AI ecosystem.

- **Compliance-First Design** — Built ethical guardrails from day one: no PII, no accusatory language, evidence-replayable, mandatory disclaimers. The product outputs "probabilistic judgments + replayable evidence," never definitive accusations.

---

## What We Learned

- **Agentic AI requires careful tool design.** The quality of an agentic system depends heavily on how tools are described and scoped. Vague tool descriptions lead to poor orchestration; overly narrow tools lead to incomplete analysis.

- **Statistical rigor matters in crypto analytics.** The crypto space is full of survivorship bias and post-hoc narratives. Applying proper hypothesis testing (with pre-registered significance levels) dramatically reduces false positives in "smart money" detection.

- **Nova 2 Lite's extended thinking is a game-changer for complex reasoning.** Multi-step attribution — where you need to reason from raw trades → PnL decomposition → strategy classification → information edge — benefits enormously from structured thinking before responding.

- **Voice interaction creates a new UX paradigm for analytics.** Being able to ask "Why is this wallet profitable?" by voice and get a spoken analysis summary makes complex data accessible to non-technical users.

- **Prediction markets are underserved by analytics tools.** Despite Polymarket's massive volume, the analytics ecosystem is nascent compared to DeFi. There's significant opportunity for purpose-built intelligence platforms.

---

## What's Next

- **Real-time Streaming Pipeline** — Replace batch processing with a streaming architecture that updates wallet profiles within minutes of new trades
- **Multi-Platform Expansion** — Extend analysis to other prediction markets (Kalshi, Metaculus) and cross-platform arbitrage detection
- **Wallet Entity Resolution** — Use on-chain fund flow analysis and embedding clustering to link related wallets controlled by the same entity
- **Advanced Anomaly Detection** — Time-series anomaly detection on the information edge timeline, with automated alert generation
- **Public API & Embeddable Widgets** — Allow third-party platforms to integrate Easy Arbitra's wallet intelligence
- **Community-Driven Labeling** — Enable verified analysts to contribute strategy labels, creating a supervised training dataset for future model fine-tuning

---

## Built With

`amazon-nova` `amazon-bedrock` `nova-2-lite` `nova-2-sonic` `nova-multimodal-embeddings` `nova-act` `go` `gin` `gorm` `postgresql` `pgvector` `next.js` `tailwindcss` `typescript` `docker` `polymarket` `prediction-markets` `blockchain-analytics` `agentic-ai`
