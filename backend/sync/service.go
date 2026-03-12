package profilesync

import (
	"context"
	"log"
	"time"

	"github.com/brucexwang/easy-arbitra/backend/discovery"
	"github.com/brucexwang/easy-arbitra/backend/leaderboard"
	"github.com/brucexwang/easy-arbitra/backend/polymarket"
	"github.com/brucexwang/easy-arbitra/backend/profileai"
	"github.com/brucexwang/easy-arbitra/backend/storage"
)

type Service struct {
	client      *polymarket.Client
	store       *storage.Store
	ai          *profileai.Client
	interval    time.Duration
	topLimit    int
	walletLimit int
}

func NewService(client *polymarket.Client, store *storage.Store, ai *profileai.Client, interval time.Duration, topLimit, walletLimit int) *Service {
	if interval <= 0 {
		interval = 4 * time.Hour
	}
	if topLimit <= 0 {
		topLimit = 100
	}
	if walletLimit <= 0 {
		walletLimit = 3000
	}
	return &Service{
		client:      client,
		store:       store,
		ai:          ai,
		interval:    interval,
		topLimit:    topLimit,
		walletLimit: walletLimit,
	}
}

func (s *Service) Start(ctx context.Context) {
	go func() {
		s.runOnceWithLogging(ctx)

		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.runOnceWithLogging(ctx)
			}
		}
	}()
}

func (s *Service) RunOnce(ctx context.Context) error {
	entries, err := leaderboard.FetchNBALeaderboard(ctx, s.topLimit)
	if err != nil {
		return err
	}

	tracked := make([]storage.TrackedWallet, 0, len(entries))
	wallets := make([]string, 0, len(entries))
	metaByWallet := map[string]leaderboard.Entry{}
	for _, entry := range entries {
		tracked = append(tracked, storage.TrackedWallet{
			WalletAddress:     entry.WalletAddress,
			DisplayName:       entry.DisplayName,
			Source:            "polymarketanalytics_nba",
			SourceRank:        entry.Rank,
			SourcePredictions: entry.Predictions,
			SourceWins:        entry.Wins,
			VolumeUSD:         entry.VolumeUSD,
			LossUSD:           entry.LossUSD,
			WinRate:           entry.WinRate,
			OpenPositionsUSD:  entry.OpenPositionsUSD,
			PnlUSD:            entry.PnlUSD,
			LastSeenAt:        entry.FetchedAt,
		})
		wallets = append(wallets, entry.WalletAddress)
		metaByWallet[entry.WalletAddress] = entry
	}

	if err := s.store.UpsertTrackedWallets(ctx, tracked); err != nil {
		return err
	}

	candidates, err := discovery.ScoreWallets(ctx, s.client, wallets, discovery.Options{
		Sport:       "nba",
		OutputLimit: len(wallets),
		WalletLimit: s.walletLimit,
	})
	if err != nil {
		return err
	}

	candidateByWallet := map[string]discovery.Candidate{}
	for _, candidate := range candidates {
		candidateByWallet[candidate.Wallet] = candidate
	}

	for _, wallet := range wallets {
		entry := metaByWallet[wallet]
		candidate, ok := candidateByWallet[wallet]
		if !ok {
			continue
		}

		styleResult, err := s.ai.Classify(ctx, profileai.Input{
			Wallet:                  wallet,
			DisplayName:             candidate.DisplayName,
			SourceRank:              entry.Rank,
			WinRate:                 entry.WinRate,
			PnlUSD:                  entry.PnlUSD,
			NbaTrades:               candidate.NbaTrades,
			RecentMarkets:           candidate.RecentMarkets,
			EntryTimingHours:        candidate.EntryTimingHours,
			SizeRatioPct:            candidate.SizeRatioPct,
			Conviction:              candidate.Conviction,
			DeterministicStyleLabel: candidate.StyleLabel,
			PresentationScore:       candidate.PresentationScore,
		})
		if err != nil {
			styleResult = profileai.Result{
				StyleLabel: candidate.StyleLabel,
				Summary:    "AI tagging failed; using deterministic style label.",
				Source:     "fallback",
				Model:      "",
			}
		}

		if err := s.store.UpsertWalletProfile(ctx, storage.WalletProfile{
			WalletAddress:           wallet,
			DisplayName:             candidate.DisplayName,
			SourceRank:              entry.Rank,
			WinRate:                 entry.WinRate,
			PnlUSD:                  entry.PnlUSD,
			NbaTrades:               candidate.NbaTrades,
			RecentMarkets:           candidate.RecentMarkets,
			EntryTimingHours:        candidate.EntryTimingHours,
			SizeRatioPct:            candidate.SizeRatioPct,
			Conviction:              candidate.Conviction,
			DeterministicStyleLabel: candidate.StyleLabel,
			AIStyleLabel:            styleResult.StyleLabel,
			AIStyleSummary:          styleResult.Summary,
			ExplanationSource:       styleResult.Source,
			Model:                   styleResult.Model,
			PresentationScore:       candidate.PresentationScore,
			AnalyzedAt:              time.Now().UTC(),
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) runOnceWithLogging(ctx context.Context) {
	runCtx, cancel := context.WithTimeout(ctx, 2*time.Hour)
	defer cancel()

	start := time.Now()
	if err := s.RunOnce(runCtx); err != nil {
		log.Printf("wallet sync failed: %v", err)
		return
	}
	log.Printf("wallet sync completed in %s", time.Since(start).Round(time.Second))
}
