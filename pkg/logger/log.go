package logger

import (
	"go.uber.org/zap"
)

func init() {
	logger := zap.NewNop() //zap.NewDevelopment()
	Set(logger)
}

func Log() *zap.Logger {
	return zap.L()
}

func Set(logger *zap.Logger) {
	zap.L().Sync()
	zap.ReplaceGlobals(logger)
}

func Sync() error {
	return zap.L().Sync()
}
