package initialize

import (
	"fmt"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

func InitLogger() {
	// fmt.Println(os.Getwd())
	path, _ := filepath.Abs(".")
	day := time.Now().Day()

	// 日志配置
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "console"
	cfg.OutputPaths = []string{
		fmt.Sprintf("%s/logs/%d.log", path, day),
	}
	// logger, err := zap.NewDevelopment()
	logger, err := cfg.Build()
	if err != nil {
		zap.S().Fatalf("can't initialize zap logger: %v", err)
	}

	// defer logger.Sync() // flushes buffer, if any
	zap.ReplaceGlobals(logger)

	// zap.S() 可以获取一个全局的Sugar, 可以让我们自己设置一个全局的logger
	// zap.S().Debug("启动日志模块...")
}
