package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

type TrackedWallet struct {
	WalletAddress     string
	DisplayName       string
	Source            string
	SourceRank        int
	SourcePredictions int
	SourceWins        int
	VolumeUSD         float64
	LossUSD           float64
	WinRate           float64
	OpenPositionsUSD  float64
	PnlUSD            float64
	LastSeenAt        time.Time
}

type WalletProfile struct {
	WalletAddress           string
	DisplayName             string
	SourceRank              int
	WinRate                 float64
	PnlUSD                  float64
	NbaTrades               int
	RecentMarkets           int
	EntryTimingHours        float64
	SizeRatioPct            float64
	Conviction              float64
	DeterministicStyleLabel string
	AIStyleLabel            string
	AIStyleSummary          string
	ExplanationSource       string
	Model                   string
	PresentationScore       float64
	AnalyzedAt              time.Time
}

type StyleWallet struct {
	WalletAddress     string  `json:"wallet_address"`
	DisplayName       string  `json:"display_name"`
	SourceRank        int     `json:"source_rank"`
	WinRate           float64 `json:"win_rate"`
	PnlUSD            float64 `json:"pnl_usd"`
	NbaTrades         int     `json:"nba_trades"`
	EntryTimingHours  float64 `json:"entry_timing_hours"`
	SizeRatioPct      float64 `json:"size_ratio_pct"`
	Conviction        float64 `json:"conviction"`
	StyleLabel        string  `json:"style_label"`
	StyleSummary      string  `json:"style_summary"`
	ExplanationSource string  `json:"explanation_source"`
}

type StyleGroup struct {
	Label   string        `json:"label"`
	Wallets []StyleWallet `json:"wallets"`
}

func Open(ctx context.Context, databaseURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	store := &Store{pool: pool}
	if err := store.initSchema(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() {
	if s != nil && s.pool != nil {
		s.pool.Close()
	}
}

func (s *Store) initSchema(ctx context.Context) error {
	const schema = `
CREATE TABLE IF NOT EXISTS tracked_wallets (
  wallet_address TEXT PRIMARY KEY,
  display_name TEXT NOT NULL,
  source TEXT NOT NULL,
  source_rank INTEGER NOT NULL DEFAULT 0,
  source_predictions INTEGER NOT NULL DEFAULT 0,
  source_wins INTEGER NOT NULL DEFAULT 0,
  volume_usd DOUBLE PRECISION NOT NULL DEFAULT 0,
  loss_usd DOUBLE PRECISION NOT NULL DEFAULT 0,
  win_rate DOUBLE PRECISION NOT NULL DEFAULT 0,
  open_positions_usd DOUBLE PRECISION NOT NULL DEFAULT 0,
  pnl_usd DOUBLE PRECISION NOT NULL DEFAULT 0,
  first_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wallet_profiles (
  wallet_address TEXT PRIMARY KEY REFERENCES tracked_wallets(wallet_address) ON DELETE CASCADE,
  nba_trades INTEGER NOT NULL DEFAULT 0,
  recent_markets INTEGER NOT NULL DEFAULT 0,
  entry_timing_hours DOUBLE PRECISION NOT NULL DEFAULT 0,
  size_ratio_pct DOUBLE PRECISION NOT NULL DEFAULT 0,
  conviction DOUBLE PRECISION NOT NULL DEFAULT 0,
  deterministic_style_label TEXT NOT NULL DEFAULT '',
  ai_style_label TEXT NOT NULL DEFAULT '',
  ai_style_summary TEXT NOT NULL DEFAULT '',
  explanation_source TEXT NOT NULL DEFAULT 'fallback',
  model TEXT NOT NULL DEFAULT '',
  presentation_score DOUBLE PRECISION NOT NULL DEFAULT 0,
  analyzed_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_wallet_profiles_ai_style_label
  ON wallet_profiles (ai_style_label, analyzed_at DESC);

CREATE INDEX IF NOT EXISTS idx_tracked_wallets_source_rank
  ON tracked_wallets (source_rank ASC);`

	_, err := s.pool.Exec(ctx, schema)
	if err != nil {
		return fmt.Errorf("init schema: %w", err)
	}
	return nil
}

func (s *Store) UpsertTrackedWallets(ctx context.Context, wallets []TrackedWallet) error {
	if len(wallets) == 0 {
		return nil
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tracked wallet upsert: %w", err)
	}
	defer tx.Rollback(ctx)

	const query = `
INSERT INTO tracked_wallets (
  wallet_address, display_name, source, source_rank, source_predictions, source_wins,
  volume_usd, loss_usd, win_rate, open_positions_usd, pnl_usd, last_seen_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6,
  $7, $8, $9, $10, $11, $12, NOW()
)
ON CONFLICT (wallet_address) DO UPDATE SET
  display_name = EXCLUDED.display_name,
  source = EXCLUDED.source,
  source_rank = EXCLUDED.source_rank,
  source_predictions = EXCLUDED.source_predictions,
  source_wins = EXCLUDED.source_wins,
  volume_usd = EXCLUDED.volume_usd,
  loss_usd = EXCLUDED.loss_usd,
  win_rate = EXCLUDED.win_rate,
  open_positions_usd = EXCLUDED.open_positions_usd,
  pnl_usd = EXCLUDED.pnl_usd,
  last_seen_at = EXCLUDED.last_seen_at,
  updated_at = NOW()`

	for _, wallet := range wallets {
		_, err := tx.Exec(ctx, query,
			wallet.WalletAddress,
			wallet.DisplayName,
			wallet.Source,
			wallet.SourceRank,
			wallet.SourcePredictions,
			wallet.SourceWins,
			wallet.VolumeUSD,
			wallet.LossUSD,
			wallet.WinRate,
			wallet.OpenPositionsUSD,
			wallet.PnlUSD,
			wallet.LastSeenAt,
		)
		if err != nil {
			return fmt.Errorf("upsert tracked wallet %s: %w", wallet.WalletAddress, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tracked wallet upsert: %w", err)
	}
	return nil
}

func (s *Store) UpsertWalletProfile(ctx context.Context, profile WalletProfile) error {
	const query = `
INSERT INTO wallet_profiles (
  wallet_address, nba_trades, recent_markets, entry_timing_hours, size_ratio_pct, conviction,
  deterministic_style_label, ai_style_label, ai_style_summary, explanation_source, model,
  presentation_score, analyzed_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6,
  $7, $8, $9, $10, $11,
  $12, $13, NOW()
)
ON CONFLICT (wallet_address) DO UPDATE SET
  nba_trades = EXCLUDED.nba_trades,
  recent_markets = EXCLUDED.recent_markets,
  entry_timing_hours = EXCLUDED.entry_timing_hours,
  size_ratio_pct = EXCLUDED.size_ratio_pct,
  conviction = EXCLUDED.conviction,
  deterministic_style_label = EXCLUDED.deterministic_style_label,
  ai_style_label = EXCLUDED.ai_style_label,
  ai_style_summary = EXCLUDED.ai_style_summary,
  explanation_source = EXCLUDED.explanation_source,
  model = EXCLUDED.model,
  presentation_score = EXCLUDED.presentation_score,
  analyzed_at = EXCLUDED.analyzed_at,
  updated_at = NOW()`

	_, err := s.pool.Exec(ctx, query,
		profile.WalletAddress,
		profile.NbaTrades,
		profile.RecentMarkets,
		profile.EntryTimingHours,
		profile.SizeRatioPct,
		profile.Conviction,
		profile.DeterministicStyleLabel,
		profile.AIStyleLabel,
		profile.AIStyleSummary,
		profile.ExplanationSource,
		profile.Model,
		profile.PresentationScore,
		profile.AnalyzedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert wallet profile %s: %w", profile.WalletAddress, err)
	}
	return nil
}

func (s *Store) ListStyleGroups(ctx context.Context, limitPerGroup int) ([]StyleGroup, error) {
	if limitPerGroup <= 0 {
		limitPerGroup = 6
	}

	const query = `
WITH ranked AS (
  SELECT
    wp.ai_style_label,
    tw.wallet_address,
    tw.display_name,
    tw.source_rank,
    tw.win_rate,
    tw.pnl_usd,
    wp.nba_trades,
    wp.entry_timing_hours,
    wp.size_ratio_pct,
    wp.conviction,
    wp.ai_style_summary,
    wp.explanation_source,
    ROW_NUMBER() OVER (
      PARTITION BY wp.ai_style_label
      ORDER BY tw.source_rank ASC, wp.presentation_score DESC, wp.analyzed_at DESC
    ) AS row_num
  FROM wallet_profiles wp
  JOIN tracked_wallets tw ON tw.wallet_address = wp.wallet_address
  WHERE wp.ai_style_label <> ''
)
SELECT
  ai_style_label,
  wallet_address,
  display_name,
  source_rank,
  win_rate,
  pnl_usd,
  nba_trades,
  entry_timing_hours,
  size_ratio_pct,
  conviction,
  ai_style_summary,
  explanation_source
FROM ranked
WHERE row_num <= $1
ORDER BY ai_style_label ASC, source_rank ASC`

	rows, err := s.pool.Query(ctx, query, limitPerGroup)
	if err != nil {
		return nil, fmt.Errorf("list style groups: %w", err)
	}
	defer rows.Close()

	groupMap := map[string][]StyleWallet{}
	order := []string{}
	for rows.Next() {
		var label string
		var wallet StyleWallet
		if err := rows.Scan(
			&label,
			&wallet.WalletAddress,
			&wallet.DisplayName,
			&wallet.SourceRank,
			&wallet.WinRate,
			&wallet.PnlUSD,
			&wallet.NbaTrades,
			&wallet.EntryTimingHours,
			&wallet.SizeRatioPct,
			&wallet.Conviction,
			&wallet.StyleSummary,
			&wallet.ExplanationSource,
		); err != nil {
			return nil, fmt.Errorf("scan style group row: %w", err)
		}
		wallet.StyleLabel = label
		if _, ok := groupMap[label]; !ok {
			order = append(order, label)
		}
		groupMap[label] = append(groupMap[label], wallet)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate style groups: %w", err)
	}

	groups := make([]StyleGroup, 0, len(order))
	for _, label := range order {
		groups = append(groups, StyleGroup{
			Label:   label,
			Wallets: groupMap[label],
		})
	}
	return groups, nil
}
