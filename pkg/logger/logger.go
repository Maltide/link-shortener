// Package logger provides a configured zap logger for application logging.
package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger creates and configures a zap.SugaredLogger with the specified log level.
func Logger(loglevel string) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()

	ourline := strings.ToLower(strings.TrimSpace(loglevel))

	lvl, err := zapcore.ParseLevel(ourline)
	if err != nil {
		lvl = zapcore.DebugLevel
	}

	cfg.Level = zap.NewAtomicLevelAt(lvl)

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	logger.Info("logger succesfully created")

	return logger.Sugar(), nil
}
