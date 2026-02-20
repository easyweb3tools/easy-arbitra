CREATE TABLE IF NOT EXISTS token (
    id BIGSERIAL PRIMARY KEY,
    market_id BIGINT NOT NULL REFERENCES market(id),
    token_id VARCHAR(80) UNIQUE NOT NULL,
    side SMALLINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_token_market_id ON token(market_id);

CREATE TABLE IF NOT EXISTS trade_fill (
    id BIGSERIAL PRIMARY KEY,
    token_id BIGINT NOT NULL REFERENCES token(id),
    maker_wallet_id BIGINT REFERENCES wallet(id),
    taker_wallet_id BIGINT REFERENCES wallet(id),
    side SMALLINT NOT NULL,
    price NUMERIC(18,8) NOT NULL,
    size NUMERIC(36,18) NOT NULL,
    fee_paid NUMERIC(18,8) DEFAULT 0,
    tx_hash BYTEA,
    block_number BIGINT,
    block_time TIMESTAMPTZ NOT NULL,
    source SMALLINT NOT NULL DEFAULT 0,
    uniq_key VARCHAR(130) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_trade_fill_maker_wallet_id ON trade_fill(maker_wallet_id);
CREATE INDEX IF NOT EXISTS idx_trade_fill_taker_wallet_id ON trade_fill(taker_wallet_id);
CREATE INDEX IF NOT EXISTS idx_trade_fill_token_id ON trade_fill(token_id);
CREATE INDEX IF NOT EXISTS idx_trade_fill_block_time ON trade_fill(block_time);

CREATE TABLE IF NOT EXISTS offchain_event (
    id BIGSERIAL PRIMARY KEY,
    market_id BIGINT REFERENCES market(id),
    event_time TIMESTAMPTZ NOT NULL,
    event_type VARCHAR(30) NOT NULL,
    source_name VARCHAR(100),
    title TEXT,
    payload JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_offchain_event_market_id ON offchain_event(market_id);
CREATE INDEX IF NOT EXISTS idx_offchain_event_event_time ON offchain_event(event_time);

CREATE TABLE IF NOT EXISTS wallet_features_daily (
    wallet_id BIGINT NOT NULL REFERENCES wallet(id),
    feature_date DATE NOT NULL,
    pnl_30d NUMERIC(20,2) DEFAULT 0,
    maker_ratio NUMERIC(5,4) DEFAULT 0,
    trade_count_30d INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (wallet_id, feature_date)
);
