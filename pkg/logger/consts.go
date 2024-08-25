//go:build !benchmark

package logger

import "go.uber.org/zap"

var logger, _ = zap.NewDevelopment()
