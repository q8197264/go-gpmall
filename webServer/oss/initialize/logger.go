package initialize

import (
	"fmt"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

func InitLogger() {
	path, _ := filepath.Abs(".")
	t := time.Date(2020, 12, 12, 0, 0, 0, 0, time.UTC)
	y, m, d := t.Date()

	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{fmt.Sprintf("%s/logs/%d-%d-%d.log", path, y, m, d)}
	logger, err := cfg.Build()
	if err != nil {
		zap.S().DPanic("日志库初始化失败:", err.Error())
	}

	zap.ReplaceGlobals(logger)
	zap.S().Debug("日志模块初始化成功")
}
