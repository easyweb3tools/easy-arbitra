ALTER TABLE wallet_score
  ADD COLUMN IF NOT EXISTS pool_tier VARCHAR(20) DEFAULT 'observation',
  ADD COLUMN IF NOT EXISTS pool_tier_updated_at TIMESTAMP,
  ADD COLUMN IF NOT EXISTS suitable_for VARCHAR(50),
  ADD COLUMN IF NOT EXISTS risk_level VARCHAR(10),
  ADD COLUMN IF NOT EXISTS suggested_position VARCHAR(20),
  ADD COLUMN IF NOT EXISTS momentum VARCHAR(20);

ALTER TABLE wallet_update_event
  ADD COLUMN IF NOT EXISTS action_required BOOLEAN DEFAULT FALSE,
  ADD COLUMN IF NOT EXISTS suggestion TEXT,
  ADD COLUMN IF NOT EXISTS suggestion_zh TEXT;

CREATE TABLE IF NOT EXISTS portfolio (
  id           BIGSERIAL PRIMARY KEY,
  name         VARCHAR(100) NOT NULL,
  name_zh      VARCHAR(100),
  description  TEXT,
  risk_level   VARCHAR(10) NOT NULL DEFAULT 'medium',
  wallet_ids   JSONB NOT NULL DEFAULT '[]'::jsonb,
  is_active    BOOLEAN NOT NULL DEFAULT TRUE,
  created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_wallet_score_pool_tier ON wallet_score(pool_tier);
CREATE INDEX IF NOT EXISTS idx_wallet_update_event_action_required ON wallet_update_event(action_required);
CREATE INDEX IF NOT EXISTS idx_portfolio_active ON portfolio(is_active);
