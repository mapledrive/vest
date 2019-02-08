package log

import (
	"fmt"

	"go.uber.org/zap"
)

var (
	sugar *zap.SugaredLogger
)

func init() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"vest.log",
	}
	logger, err := cfg.Build()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	defer logger.Sync()

	sugar = logger.Sugar()
}

func Error(args ...interface{}) {
	sugar.Error(args...)
}

func Info(args ...interface{}) {
	sugar.Info(args...)
}
