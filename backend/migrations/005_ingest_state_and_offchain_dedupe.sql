CREATE TABLE IF NOT EXISTS ingest_cursor (
    source TEXT NOT NULL,
    stream TEXT NOT NULL,
    cursor_value TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (source, stream)
);

CREATE TABLE IF NOT EXISTS ingest_run (
    id BIGSERIAL PRIMARY KEY,
    job_name TEXT NOT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'running',
    stats JSONB NOT NULL DEFAULT '{}'::jsonb,
    error_text TEXT
);

ALTER TABLE offchain_event
    ADD COLUMN IF NOT EXISTS source_event_id VARCHAR(120);

UPDATE offchain_event
SET source_event_id = ''
WHERE source_event_id IS NULL;

ALTER TABLE offchain_event
    ALTER COLUMN source_event_id SET DEFAULT '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_offchain_event_source_event
    ON offchain_event (source_name, source_event_id)
    WHERE source_event_id <> '';
