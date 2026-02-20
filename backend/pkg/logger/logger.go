package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string, format string) (*zap.Logger, error) {
	lvl := zapcore.InfoLevel
	if err := lvl.Set(level); err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	cfg := zap.NewProductionConfig()
	if format == "console" {
		cfg.Encoding = "console"
	}
	cfg.Level = zap.NewAtomicLevelAt(lvl)
	return cfg.Build()
}
