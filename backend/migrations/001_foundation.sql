CREATE TABLE IF NOT EXISTS wallet (
    id BIGSERIAL PRIMARY KEY,
    address BYTEA NOT NULL,
    chain_id INT NOT NULL DEFAULT 137,
    pseudonym VARCHAR(100),
    is_tracked BOOLEAN DEFAULT FALSE,
    first_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(address, chain_id)
);

CREATE TABLE IF NOT EXISTS market (
    id BIGSERIAL PRIMARY KEY,
    condition_id VARCHAR(66) UNIQUE NOT NULL,
    slug VARCHAR(255),
    title TEXT NOT NULL,
    category VARCHAR(50),
    status SMALLINT NOT NULL DEFAULT 0,
    has_fee BOOLEAN DEFAULT FALSE,
    resolution_time TIMESTAMPTZ,
    resolved_outcome SMALLINT,
    volume DECIMAL(20,2) DEFAULT 0,
    liquidity DECIMAL(20,2) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wallet_score (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallet(id),
    scored_at TIMESTAMPTZ NOT NULL,
    strategy_type VARCHAR(30),
    strategy_confidence DECIMAL(5,4) DEFAULT 0,
    info_edge_level VARCHAR(20),
    info_edge_confidence DECIMAL(5,4) DEFAULT 0,
    smart_score INT DEFAULT 0,
    scoring_detail JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ai_analysis_report (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallet(id),
    model_id VARCHAR(50) NOT NULL,
    report JSONB NOT NULL,
    nl_summary TEXT,
    risk_warnings JSONB,
    input_tokens INT,
    output_tokens INT,
    latency_ms INT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_wallet_score_wallet_id ON wallet_score(wallet_id);
CREATE INDEX IF NOT EXISTS idx_ai_report_wallet_id ON ai_analysis_report(wallet_id);
