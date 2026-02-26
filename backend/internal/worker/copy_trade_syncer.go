package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"easy-arbitra/backend/internal/copytrade"
	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/internal/repository"

	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CopyTradeSyncer struct {
	copyTradeService *copytrade.Service
	copyTradeRepo    *copytrade.Repository
	tradeRepo        *repository.TradeRepository
	marketRepo       *repository.MarketRepository
	logger           *zap.Logger
	db               *gorm.DB
}

func NewCopyTradeSyncer(
	copyTradeService *copytrade.Service,
	copyTradeRepo *copytrade.Repository,
	tradeRepo *repository.TradeRepository,
	marketRepo *repository.MarketRepository,
	logger *zap.Logger,
	db *gorm.DB,
) *CopyTradeSyncer {
	return &CopyTradeSyncer{
		copyTradeService: copyTradeService,
		copyTradeRepo:    copyTradeRepo,
		tradeRepo:        tradeRepo,
		marketRepo:       marketRepo,
		logger:           logger,
		db:               db,
	}
}

func (s *CopyTradeSyncer) Name() string { return "copy_trade_syncer" }

type copyTradeStats struct {
	WalletsChecked int `json:"wallets_checked"`
	TradesFound    int `json:"trades_found"`
	NewTrades      int `json:"new_trades"`
	DecisionsCopy  int `json:"decisions_copy"`
	DecisionsSkip  int `json:"decisions_skip"`
	Errors         int `json:"errors"`
}

func (s *CopyTradeSyncer) RunOnce(ctx context.Context) error {
	now := time.Now().UTC()

	// Create IngestRun record
	run := model.IngestRun{
		JobName:   s.Name(),
		StartedAt: now,
		Status:    "running",
		Stats:     datatypes.JSON([]byte("{}")),
	}
	if err := s.db.WithContext(ctx).Create(&run).Error; err != nil {
		s.logger.Warn("copy_trade_syncer: failed to create ingest_run", zap.Error(err))
	}

	stats := copyTradeStats{}
	var runErr error

	defer func() {
		endedAt := time.Now().UTC()
		statsJSON, _ := json.Marshal(stats)
		updates := map[string]any{
			"ended_at": endedAt,
			"stats":    datatypes.JSON(statsJSON),
		}
		if runErr != nil {
			updates["status"] = "error"
			errText := runErr.Error()
			updates["error_text"] = errText
		} else {
			updates["status"] = "done"
		}
		if err := s.db.WithContext(ctx).Model(&model.IngestRun{}).Where("id = ?", run.ID).Updates(updates).Error; err != nil {
			s.logger.Warn("copy_trade_syncer: failed to update ingest_run", zap.Error(err))
		}
		s.logger.Info("copy_trade_syncer: run complete",
			zap.Int("wallets_checked", stats.WalletsChecked),
			zap.Int("trades_found", stats.TradesFound),
			zap.Int("new_trades", stats.NewTrades),
			zap.Int("decisions_copy", stats.DecisionsCopy),
			zap.Int("decisions_skip", stats.DecisionsSkip),
			zap.Int("errors", stats.Errors),
			zap.Duration("duration", endedAt.Sub(now)),
		)
	}()

	configs, err := s.copyTradeRepo.ListEnabledConfigs(ctx)
	if err != nil {
		s.logger.Warn("copy_trade_syncer: ListEnabledConfigs failed", zap.Error(err))
		runErr = fmt.Errorf("ListEnabledConfigs: %w", err)
		return err
	}

	s.logger.Info("copy_trade_syncer: starting", zap.Int("enabled_configs", len(configs)))

	for _, cfg := range configs {
		stats.WalletsChecked++
		since := cfg.LastCheckedAt
		if since.IsZero() {
			since = now.Add(-24 * time.Hour)
		}

		trades, _, err := s.tradeRepo.ListByWalletID(ctx, cfg.WalletID, 20, 0)
		if err != nil {
			s.logger.Warn("copy_trade_syncer: ListByWalletID failed",
				zap.Int64("wallet_id", cfg.WalletID),
				zap.Int64("config_id", cfg.ID),
				zap.Error(err),
			)
			stats.Errors++
			continue
		}
		stats.TradesFound += len(trades)

		for _, trade := range trades {
			if !trade.BlockTime.After(since) {
				continue
			}
			stats.NewTrades++

			market, err := s.marketRepo.GetBySlug(ctx, trade.MarketSlug)
			if err != nil {
				s.logger.Warn("copy_trade_syncer: GetBySlug failed",
					zap.String("market_slug", trade.MarketSlug),
					zap.Error(err),
				)
			}

			dec, err := s.copyTradeService.ProcessNewTrade(ctx, cfg.ID, trade, market)
			if err != nil {
				s.logger.Warn("copy_trade_syncer: ProcessNewTrade failed",
					zap.Int64("config_id", cfg.ID),
					zap.Int64("trade_id", trade.TradeID),
					zap.Error(err),
				)
				stats.Errors++
				continue
			}
			if dec != nil {
				if dec.Decision == "copy" {
					stats.DecisionsCopy++
				} else {
					stats.DecisionsSkip++
				}
			}
		}

		if err := s.copyTradeRepo.UpdateLastChecked(ctx, cfg.ID, now); err != nil {
			s.logger.Warn("copy_trade_syncer: UpdateLastChecked failed",
				zap.Int64("config_id", cfg.ID),
				zap.Error(err),
			)
		}
	}

	return nil
}
