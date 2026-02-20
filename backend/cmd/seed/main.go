package main

import (
	"log"
	"time"

	"easy-arbitra/backend/config"
	"easy-arbitra/backend/internal/model"
	"easy-arbitra/backend/pkg/polyaddr"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	if cfg.Database.AutoMigrate {
		if err := db.AutoMigrate(&model.Wallet{}, &model.Market{}, &model.WalletScore{}, &model.AIAnalysisReport{}); err != nil {
			log.Fatalf("auto migrate: %v", err)
		}
	}

	seedWallets(db)
	seedMarkets(db)
	seedScores(db)

	log.Println("seed completed")
}

func seedWallets(db *gorm.DB) {
	rows := []struct {
		Address   string
		Pseudonym string
		Tracked   bool
	}{
		{"0x1111111111111111111111111111111111111111", "alpha_quant", true},
		{"0x2222222222222222222222222222222222222222", "event_hunter", true},
		{"0x3333333333333333333333333333333333333333", "maker_bot", false},
	}

	for _, r := range rows {
		addr, err := polyaddr.HexToBytes(r.Address)
		if err != nil {
			log.Fatalf("invalid seed wallet address: %v", err)
		}
		pseudonym := r.Pseudonym
		wallet := model.Wallet{
			Address:     addr,
			ChainID:     137,
			Pseudonym:   &pseudonym,
			IsTracked:   r.Tracked,
			FirstSeenAt: time.Now().UTC().Add(-48 * time.Hour),
		}
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "address"}, {Name: "chain_id"}},
			DoUpdates: clause.Assignments(map[string]any{"pseudonym": wallet.Pseudonym, "is_tracked": wallet.IsTracked, "updated_at": time.Now().UTC()}),
		}).Create(&wallet).Error; err != nil {
			log.Fatalf("seed wallet: %v", err)
		}
	}
}

func seedMarkets(db *gorm.DB) {
	rows := []model.Market{
		{ConditionID: "0xaaa", Slug: "us-election-2028", Title: "Who wins US election 2028?", Category: "Politics", Status: 0, HasFee: true, Volume: 120000, Liquidity: 45000},
		{ConditionID: "0xbbb", Slug: "btc-200k-2026", Title: "BTC above 200k by end of 2026?", Category: "Crypto", Status: 0, HasFee: false, Volume: 88000, Liquidity: 30000},
		{ConditionID: "0xccc", Slug: "nba-finals-2026", Title: "Who wins NBA Finals 2026?", Category: "Sports", Status: 0, HasFee: false, Volume: 64000, Liquidity: 27000},
	}

	for _, row := range rows {
		m := row
		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "condition_id"}},
			DoUpdates: clause.Assignments(map[string]any{
				"title":      m.Title,
				"category":   m.Category,
				"status":     m.Status,
				"volume":     m.Volume,
				"liquidity":  m.Liquidity,
				"updated_at": time.Now().UTC(),
			}),
		}).Create(&m).Error; err != nil {
			log.Fatalf("seed market: %v", err)
		}
	}
}

func seedScores(db *gorm.DB) {
	var wallets []model.Wallet
	if err := db.Find(&wallets).Error; err != nil {
		log.Fatalf("load wallets for scores: %v", err)
	}

	for i, w := range wallets {
		score := model.WalletScore{
			WalletID:           w.ID,
			ScoredAt:           time.Now().UTC(),
			StrategyType:       []string{"quant", "event_trader", "market_maker"}[i%3],
			StrategyConfidence: 0.72,
			InfoEdgeLevel:      []string{"processing_edge", "luck", "quant"}[i%3],
			InfoEdgeConfidence: 0.64,
			SmartScore:         []int{78, 61, 69}[i%3],
			ScoringDetail:      datatypes.JSON([]byte(`{"consistency":0.74,"risk_control":0.66}`)),
		}
		if err := db.Create(&score).Error; err != nil {
			log.Fatalf("seed score: %v", err)
		}
	}
}
