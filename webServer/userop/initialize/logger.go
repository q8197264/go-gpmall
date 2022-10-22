package initialize

import (
	"fmt"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

func InitLogger() {
	yy, mm, dd := time.Now().Date()

	root, _ := filepath.Abs(".")
	path := filepath.Join(root, "logs", fmt.Sprintf("%d-%d-%d.log", yy, mm, dd))

	cfg := zap.NewDevelopmentConfig()
	// cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.OutputPaths = append(cfg.OutputPaths, path)
	logger, err := cfg.Build()
	if err != nil {
		panic(err.Error())
	}
	zap.ReplaceGlobals(logger)
}
