package initialize

import (
	"fmt"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

func InitLogger() {
	path, _ := filepath.Abs(".")
	t := time.Now()
	path = fmt.Sprintf("%s/logs/%d-%d-%d.log", path, t.Year(), t.Month(), t.Day())
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{path}
	logger, err := cfg.Build()
	if err != nil {
		panic(err.Error())
	}

	zap.ReplaceGlobals(logger)
}
