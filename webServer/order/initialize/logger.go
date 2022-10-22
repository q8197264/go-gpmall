package initialize

import (
	"fmt"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// 设置日志存储目录
func InitLogger() {
	rootDir, _ := filepath.Abs(".")
	y, m, d := time.Now().Date()
	path := filepath.Join(rootDir, "logs", fmt.Sprintf("%d-%d-%d.log", y, m, d))
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = append(cfg.OutputPaths, path)
	logger, err := cfg.Build()
	if err != nil {
		zap.S().DPanic(err.Error())
	}
	zap.ReplaceGlobals(logger)
}
