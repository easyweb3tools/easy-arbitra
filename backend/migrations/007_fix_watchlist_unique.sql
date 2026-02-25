-- Fix Bug 1: watchlist missing UNIQUE (wallet_id, user_fingerprint)
-- Required by ON CONFLICT in watchlist add/batch operations.
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint
    WHERE conrelid = 'watchlist'::regclass
      AND contype = 'u'
      AND conkey @> ARRAY[
        (SELECT attnum FROM pg_attribute WHERE attrelid = 'watchlist'::regclass AND attname = 'wallet_id'),
        (SELECT attnum FROM pg_attribute WHERE attrelid = 'watchlist'::regclass AND attname = 'user_fingerprint')
      ]
  ) THEN
    ALTER TABLE watchlist ADD CONSTRAINT uq_watchlist_wallet_user UNIQUE (wallet_id, user_fingerprint);
  END IF;
END $$;

-- Fix Bug 2: offchain_event partial unique index cannot be matched by ON CONFLICT.
-- Drop the partial index (WHERE source_event_id <> '') and create a full unique index.
DROP INDEX IF EXISTS idx_offchain_event_source_event;
CREATE UNIQUE INDEX idx_offchain_event_source_event
    ON offchain_event (source_name, source_event_id);
