package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"easy-arbitra/backend/config"
	"easy-arbitra/backend/internal/ai"
	"easy-arbitra/backend/internal/api"
	"easy-arbitra/backend/internal/api/handler"
	"easy-arbitra/backend/internal/client"
	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/internal/repository"
	"easy-arbitra/backend/internal/service"
	"easy-arbitra/backend/internal/worker"
	"easy-arbitra/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	lg, err := logger.New(cfg.Logger.Level, cfg.Logger.Format)
	if err != nil {
		log.Fatalf("build logger: %v", err)
	}
	defer func() { _ = lg.Sync() }()

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("db handle: %v", err)
	}
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)

	if cfg.Database.AutoMigrate {
		if err := db.AutoMigrate(
			&model.Wallet{},
			&model.Market{},
			&model.Token{},
			&model.TradeFill{},
			&model.OffchainEvent{},
			&model.WalletFeaturesDaily{},
			&model.WalletScore{},
			&model.AIAnalysisReport{},
			&model.AnomalyAlert{},
		); err != nil {
			log.Fatalf("auto migrate: %v", err)
		}
	}

	walletRepo := repository.NewWalletRepository(db)
	marketRepo := repository.NewMarketRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	tradeRepo := repository.NewTradeRepository(db)
	featureRepo := repository.NewFeatureRepository(db)
	scoreRepo := repository.NewScoreRepository(db)
	anomalyRepo := repository.NewAnomalyRepository(db)
	aiReportRepo := repository.NewAIReportRepository(db)

	infoEdgeService := service.NewInfoEdgeService(tradeRepo)
	walletService := service.NewWalletService(walletRepo, scoreRepo, tradeRepo, infoEdgeService)
	marketService := service.NewMarketService(marketRepo)
	statsService := service.NewStatsService(walletRepo, marketRepo, scoreRepo)
	anomalyService := service.NewAnomalyService(anomalyRepo, walletRepo, tradeRepo, infoEdgeService)
	explainService := service.NewExplanationService(walletRepo, featureRepo, scoreRepo, tradeRepo, infoEdgeService)
	classifier := service.NewClassificationService(featureRepo, scoreRepo)
	analyzer := ai.NewAnalyzer(cfg.Nova, lg)
	aiService := service.NewAIService(walletRepo, scoreRepo, tradeRepo, aiReportRepo, analyzer, cfg.Nova.AnalysisCacheHours)

	h := handler.New(
		walletService, marketService, statsService, anomalyService, explainService, infoEdgeService, aiService,
		func(c *gin.Context) error {
			return sqlDB.PingContext(c.Request.Context())
		},
	)
	r := api.NewRouter(h)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if cfg.Worker.Enabled {
		gammaClient := client.NewGammaClient(cfg.Polymarket.GammaAPIURL, cfg.Polymarket.RequestTO)
		dataClient := client.NewDataAPIClient(cfg.Polymarket.DataAPIURL, cfg.Polymarket.RequestTO)
		offchainClient := client.NewOffchainClient()

		mgr := worker.NewManager(lg,
			worker.ScheduledSyncer{Syncer: worker.NewMarketSyncer(gammaClient, marketRepo, cfg.Worker.MaxMarketsPerSync), Interval: cfg.Worker.MarketSyncerInterval},
			worker.ScheduledSyncer{Syncer: worker.NewTradeSyncer(dataClient, walletRepo, marketRepo, tokenRepo, tradeRepo, cfg.Worker.MaxTradesPerSync), Interval: cfg.Worker.TradeSyncerInterval},
			worker.ScheduledSyncer{Syncer: worker.NewOffchainEventSyncer(offchainClient, cfg.Worker.MaxOffchainEventsPerSync), Interval: cfg.Worker.OffchainSyncerInterval},
			worker.ScheduledSyncer{Syncer: worker.NewFeatureBuilder(featureRepo), Interval: cfg.Worker.FeatureBuilderInterval},
			worker.ScheduledSyncer{Syncer: worker.NewScoreCalculator(walletRepo, classifier), Interval: cfg.Worker.ScoreCalculatorInterval},
			worker.ScheduledSyncer{Syncer: worker.NewAnomalyDetector(anomalyService), Interval: cfg.Worker.AnomalyDetectorInterval},
		)
		mgr.Start(ctx, cfg.Worker.RunOnStartup)
	}

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
