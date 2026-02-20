ALTER TABLE wallet_features_daily
    ADD COLUMN IF NOT EXISTS pnl_7d NUMERIC(20,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS pnl_90d NUMERIC(20,2) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS avg_edge NUMERIC(10,6) DEFAULT 0,
    ADD COLUMN IF NOT EXISTS active_days_30d INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS tx_frequency_per_day NUMERIC(10,4) DEFAULT 0;

CREATE TABLE IF NOT EXISTS anomaly_alert (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallet(id),
    market_id BIGINT REFERENCES market(id),
    alert_type VARCHAR(30) NOT NULL,
    severity SMALLINT NOT NULL DEFAULT 0,
    evidence JSONB NOT NULL,
    description TEXT,
    acknowledged BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_anomaly_alert_wallet_id ON anomaly_alert(wallet_id);
CREATE INDEX IF NOT EXISTS idx_anomaly_alert_severity ON anomaly_alert(severity);
CREATE INDEX IF NOT EXISTS idx_anomaly_alert_alert_type ON anomaly_alert(alert_type);
