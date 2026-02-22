CREATE TABLE IF NOT EXISTS watchlist (
  id BIGSERIAL PRIMARY KEY,
  wallet_id BIGINT NOT NULL REFERENCES wallet(id) ON DELETE CASCADE,
  user_fingerprint VARCHAR(120) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (wallet_id, user_fingerprint)
);

CREATE INDEX IF NOT EXISTS idx_watchlist_user_fingerprint ON watchlist(user_fingerprint);
CREATE INDEX IF NOT EXISTS idx_watchlist_wallet_id ON watchlist(wallet_id);

CREATE TABLE IF NOT EXISTS wallet_update_event (
  id BIGSERIAL PRIMARY KEY,
  wallet_id BIGINT NOT NULL REFERENCES wallet(id) ON DELETE CASCADE,
  event_type VARCHAR(40) NOT NULL,
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_wallet_update_event_wallet_id ON wallet_update_event(wallet_id);
CREATE INDEX IF NOT EXISTS idx_wallet_update_event_created_at ON wallet_update_event(created_at DESC);
