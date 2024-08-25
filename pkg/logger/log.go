package logger

import (
	"fmt"
	"go.uber.org/zap"
)

func init() {
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

func LogNotImplemented(feature string) {
	zap.L().WithOptions(zap.AddCallerSkip(1)).Warn(fmt.Sprintf("not implemented: %v", feature))
}
