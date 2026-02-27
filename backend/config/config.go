package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Nova       NovaConfig       `mapstructure:"nova"`
	Polymarket PolymarketConfig `mapstructure:"polymarket"`
	Worker     WorkerConfig     `mapstructure:"worker"`
}

type AuthConfig struct {
	JWTSecret         string        `mapstructure:"jwt_secret"`
	JWTExpiry         time.Duration `mapstructure:"jwt_expiry"`
	GoogleClientID    string        `mapstructure:"google_client_id"`
	GoogleSecret      string        `mapstructure:"google_client_secret"`
	GoogleRedirectURL string        `mapstructure:"google_redirect_url"`
	FrontendURL       string        `mapstructure:"frontend_url"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	SSLMode      string `mapstructure:"sslmode"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	AutoMigrate  bool   `mapstructure:"auto_migrate"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type NovaConfig struct {
	Enabled            bool          `mapstructure:"enabled"`
	Provider           string        `mapstructure:"provider"`
	Region             string        `mapstructure:"region"`
	APIBaseURL         string        `mapstructure:"api_base_url"`
	APIKey             string        `mapstructure:"api_key"`
	AnalysisModel      string        `mapstructure:"analysis_model"`
	MaxTokens          int           `mapstructure:"max_tokens"`
	Temperature        float32       `mapstructure:"temperature"`
	AnalysisCacheHours time.Duration `mapstructure:"analysis_cache_hours"`
}

type PolymarketConfig struct {
	GammaAPIURL       string        `mapstructure:"gamma_api_url"`
	DataAPIURL        string        `mapstructure:"data_api_url"`
	OffchainEventsURL string        `mapstructure:"offchain_events_url"`
	RequestTO         time.Duration `mapstructure:"request_timeout"`
}

type WorkerConfig struct {
	Enabled                      bool          `mapstructure:"enabled"`
	MarketSyncerInterval         time.Duration `mapstructure:"market_syncer_interval"`
	TradeSyncerInterval          time.Duration `mapstructure:"trade_syncer_interval"`
	TradeSyncerMaxPages          int           `mapstructure:"trade_syncer_max_pages"`
	TradeSyncerCursorLookback    time.Duration `mapstructure:"trade_syncer_cursor_lookback"`
	TradeBackfillSyncerInterval  time.Duration `mapstructure:"trade_backfill_syncer_interval"`
	OffchainSyncerInterval       time.Duration `mapstructure:"offchain_event_syncer_interval"`
	OffchainSyncerMaxPages       int           `mapstructure:"offchain_event_syncer_max_pages"`
	OffchainSyncerCursorLookback time.Duration `mapstructure:"offchain_event_syncer_cursor_lookback"`
	FeatureBuilderInterval       time.Duration `mapstructure:"feature_builder_interval"`
	ScoreCalculatorInterval      time.Duration `mapstructure:"score_calculator_interval"`
	AnomalyDetectorInterval      time.Duration `mapstructure:"anomaly_detector_interval"`
	RunOnStartup                 bool          `mapstructure:"run_on_startup"`
	MaxTradesPerSync             int           `mapstructure:"max_trades_per_sync"`
	BackfillWalletsPerSync       int           `mapstructure:"backfill_wallets_per_sync"`
	BackfillPagesPerWallet       int           `mapstructure:"backfill_pages_per_wallet"`
	BackfillPageSize             int           `mapstructure:"backfill_page_size"`
	BackfillConcurrency          int           `mapstructure:"backfill_concurrency"`
	BackfillTargetMinTrades      int64         `mapstructure:"backfill_target_min_trades"`
	AIBatchEnabled               bool          `mapstructure:"ai_batch_enabled"`
	AIBatchAnalyzerInterval      time.Duration `mapstructure:"ai_batch_analyzer_interval"`
	AIBatchSize                  int           `mapstructure:"ai_batch_size"`
	AIBatchCooldown              time.Duration `mapstructure:"ai_batch_cooldown"`
	AIBatchRequestSpacing        time.Duration `mapstructure:"ai_batch_request_spacing"`
	AIBatchMinTrades             int64         `mapstructure:"ai_batch_min_trades"`
	AIBatchMinRealizedPnL        float64       `mapstructure:"ai_batch_min_realized_pnl"`
	MaxMarketsPerSync            int           `mapstructure:"max_markets_per_sync"`
	MaxOffchainEventsPerSync     int           `mapstructure:"max_offchain_events_per_sync"`
	CopyTradeSyncerInterval      time.Duration `mapstructure:"copy_trade_syncer_interval"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("server.port", "8080")
	v.SetDefault("server.mode", "debug")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.auto_migrate", true)
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "json")
	v.SetDefault("auth.jwt_secret", "change-me-in-production")
	v.SetDefault("auth.jwt_expiry", "168h")
	v.SetDefault("auth.google_client_id", "")
	v.SetDefault("auth.google_client_secret", "")
	v.SetDefault("auth.google_redirect_url", "http://localhost:8080/api/v1/auth/google/callback")
	v.SetDefault("auth.frontend_url", "http://localhost:3000")
	v.SetDefault("nova.enabled", false)
	v.SetDefault("nova.provider", "devapi")
	v.SetDefault("nova.region", "us-east-1")
	v.SetDefault("nova.api_base_url", "https://api.nova.amazon.com/v1")
	v.SetDefault("nova.api_key", "")
	v.SetDefault("nova.analysis_model", "nova-pro-v1")
	v.SetDefault("nova.max_tokens", 2048)
	v.SetDefault("nova.temperature", 0.0)
	v.SetDefault("nova.analysis_cache_hours", "24h")
	v.SetDefault("polymarket.gamma_api_url", "https://gamma-api.polymarket.com")
	v.SetDefault("polymarket.data_api_url", "https://data-api.polymarket.com")
	v.SetDefault("polymarket.offchain_events_url", "https://gamma-api.polymarket.com")
	v.SetDefault("polymarket.request_timeout", "20s")
	v.SetDefault("worker.enabled", false)
	v.SetDefault("worker.market_syncer_interval", "10m")
	v.SetDefault("worker.trade_syncer_interval", "5m")
	v.SetDefault("worker.trade_syncer_max_pages", 20)
	v.SetDefault("worker.trade_syncer_cursor_lookback", "2m")
	v.SetDefault("worker.trade_backfill_syncer_interval", "5m")
	v.SetDefault("worker.offchain_event_syncer_interval", "15m")
	v.SetDefault("worker.offchain_event_syncer_max_pages", 5)
	v.SetDefault("worker.offchain_event_syncer_cursor_lookback", "2h")
	v.SetDefault("worker.feature_builder_interval", "30m")
	v.SetDefault("worker.score_calculator_interval", "1h")
	v.SetDefault("worker.anomaly_detector_interval", "10m")
	v.SetDefault("worker.run_on_startup", true)
	v.SetDefault("worker.max_trades_per_sync", 200)
	v.SetDefault("worker.backfill_wallets_per_sync", 100)
	v.SetDefault("worker.backfill_pages_per_wallet", 5)
	v.SetDefault("worker.backfill_page_size", 200)
	v.SetDefault("worker.backfill_concurrency", 16)
	v.SetDefault("worker.backfill_target_min_trades", 100)
	v.SetDefault("worker.ai_batch_enabled", true)
	v.SetDefault("worker.ai_batch_analyzer_interval", "10m")
	v.SetDefault("worker.ai_batch_size", 3)
	v.SetDefault("worker.ai_batch_cooldown", "24h")
	v.SetDefault("worker.ai_batch_request_spacing", "25s")
	v.SetDefault("worker.ai_batch_min_trades", 100)
	v.SetDefault("worker.ai_batch_min_realized_pnl", 0.0)
	v.SetDefault("worker.max_markets_per_sync", 100)
	v.SetDefault("worker.max_offchain_events_per_sync", 50)
	v.SetDefault("worker.copy_trade_syncer_interval", "5m")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	// Ensure env overrides are honored for nested fields.
	cfg.Server.Port = v.GetString("server.port")
	cfg.Server.Mode = v.GetString("server.mode")

	cfg.Database.Host = v.GetString("database.host")
	cfg.Database.Port = v.GetInt("database.port")
	cfg.Database.User = v.GetString("database.user")
	cfg.Database.Password = v.GetString("database.password")
	cfg.Database.DBName = v.GetString("database.dbname")
	cfg.Database.SSLMode = v.GetString("database.sslmode")
	cfg.Database.MaxIdleConns = v.GetInt("database.max_idle_conns")
	cfg.Database.MaxOpenConns = v.GetInt("database.max_open_conns")
	cfg.Database.AutoMigrate = v.GetBool("database.auto_migrate")

	cfg.Logger.Level = v.GetString("logger.level")
	cfg.Logger.Format = v.GetString("logger.format")

	cfg.Auth.JWTSecret = v.GetString("auth.jwt_secret")
	cfg.Auth.JWTExpiry = v.GetDuration("auth.jwt_expiry")
	cfg.Auth.GoogleClientID = v.GetString("auth.google_client_id")
	cfg.Auth.GoogleSecret = v.GetString("auth.google_client_secret")
	cfg.Auth.GoogleRedirectURL = v.GetString("auth.google_redirect_url")
	cfg.Auth.FrontendURL = v.GetString("auth.frontend_url")

	cfg.Nova.Enabled = v.GetBool("nova.enabled")
	cfg.Nova.Provider = v.GetString("nova.provider")
	cfg.Nova.Region = v.GetString("nova.region")
	cfg.Nova.APIBaseURL = v.GetString("nova.api_base_url")
	cfg.Nova.APIKey = v.GetString("nova.api_key")
	cfg.Nova.AnalysisModel = v.GetString("nova.analysis_model")
	cfg.Nova.MaxTokens = v.GetInt("nova.max_tokens")
	cfg.Nova.Temperature = float32(v.GetFloat64("nova.temperature"))
	cfg.Nova.AnalysisCacheHours = v.GetDuration("nova.analysis_cache_hours")

	cfg.Polymarket.GammaAPIURL = v.GetString("polymarket.gamma_api_url")
	cfg.Polymarket.DataAPIURL = v.GetString("polymarket.data_api_url")
	cfg.Polymarket.OffchainEventsURL = v.GetString("polymarket.offchain_events_url")
	cfg.Polymarket.RequestTO = v.GetDuration("polymarket.request_timeout")

	cfg.Worker.Enabled = v.GetBool("worker.enabled")
	cfg.Worker.MarketSyncerInterval = v.GetDuration("worker.market_syncer_interval")
	cfg.Worker.TradeSyncerInterval = v.GetDuration("worker.trade_syncer_interval")
	cfg.Worker.TradeSyncerMaxPages = v.GetInt("worker.trade_syncer_max_pages")
	cfg.Worker.TradeSyncerCursorLookback = v.GetDuration("worker.trade_syncer_cursor_lookback")
	cfg.Worker.TradeBackfillSyncerInterval = v.GetDuration("worker.trade_backfill_syncer_interval")
	cfg.Worker.OffchainSyncerInterval = v.GetDuration("worker.offchain_event_syncer_interval")
	cfg.Worker.OffchainSyncerMaxPages = v.GetInt("worker.offchain_event_syncer_max_pages")
	cfg.Worker.OffchainSyncerCursorLookback = v.GetDuration("worker.offchain_event_syncer_cursor_lookback")
	cfg.Worker.FeatureBuilderInterval = v.GetDuration("worker.feature_builder_interval")
	cfg.Worker.ScoreCalculatorInterval = v.GetDuration("worker.score_calculator_interval")
	cfg.Worker.AnomalyDetectorInterval = v.GetDuration("worker.anomaly_detector_interval")
	cfg.Worker.RunOnStartup = v.GetBool("worker.run_on_startup")
	cfg.Worker.MaxTradesPerSync = v.GetInt("worker.max_trades_per_sync")
	cfg.Worker.BackfillWalletsPerSync = v.GetInt("worker.backfill_wallets_per_sync")
	cfg.Worker.BackfillPagesPerWallet = v.GetInt("worker.backfill_pages_per_wallet")
	cfg.Worker.BackfillPageSize = v.GetInt("worker.backfill_page_size")
	cfg.Worker.BackfillConcurrency = v.GetInt("worker.backfill_concurrency")
	cfg.Worker.BackfillTargetMinTrades = int64(v.GetInt("worker.backfill_target_min_trades"))
	cfg.Worker.AIBatchEnabled = v.GetBool("worker.ai_batch_enabled")
	cfg.Worker.AIBatchAnalyzerInterval = v.GetDuration("worker.ai_batch_analyzer_interval")
	cfg.Worker.AIBatchSize = v.GetInt("worker.ai_batch_size")
	cfg.Worker.AIBatchCooldown = v.GetDuration("worker.ai_batch_cooldown")
	cfg.Worker.AIBatchRequestSpacing = v.GetDuration("worker.ai_batch_request_spacing")
	cfg.Worker.AIBatchMinTrades = int64(v.GetInt("worker.ai_batch_min_trades"))
	cfg.Worker.AIBatchMinRealizedPnL = v.GetFloat64("worker.ai_batch_min_realized_pnl")
	cfg.Worker.MaxMarketsPerSync = v.GetInt("worker.max_markets_per_sync")
	cfg.Worker.MaxOffchainEventsPerSync = v.GetInt("worker.max_offchain_events_per_sync")
	cfg.Worker.CopyTradeSyncerInterval = v.GetDuration("worker.copy_trade_syncer_interval")

	return &cfg, nil
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}
