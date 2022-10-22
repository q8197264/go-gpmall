package initialize

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"webServer/oss/global"
	"webServer/oss/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func getEnv(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

// 本地配置
func localConfig() {
	vi := viper.New()

	path, _ := filepath.Abs(".")
	if getEnv("GO_WEBSERVER_DEBUG_CONFIG") {
		path = fmt.Sprintf("%s/config/debug.yaml", path)
	} else {
		path = fmt.Sprintf("%s/config/config.yaml", path)
	}

	vi.SetConfigFile(path)
	if err := vi.ReadInConfig(); err != nil {
		zap.S().DPanic(err.Error())
	}
	if err := vi.Unmarshal(global.NacosConfig); err != nil {
		zap.S().DPanic(err.Error())
	}

	// 动态更新
	vi.WatchConfig()
	vi.OnConfigChange(func(in fsnotify.Event) {
		vi.ReadInConfig()
		vi.Unmarshal(global.NacosConfig)
	})
}

// 远程配置
func InitNacos() {
	// 载入本地文件配置
	localConfig()

	//create clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	// At least one ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      global.NacosConfig.Host,
			Port:        uint64(global.NacosConfig.Port),
			ContextPath: "/nacos",
			Scheme:      "http",
		},
	}
	// Another way of create config client for dynamic configuration (recommend)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		zap.S().Warnf("nacos 加载失败:", err.Error())
	}

	content, err := configClient.GetConfig(
		vo.ConfigParam{
			DataId: global.NacosConfig.Dataid,
			Group:  global.NacosConfig.Group,
		},
	)
	err = json.Unmarshal([]byte(content), global.ServerConfig)
	if err != nil {
		zap.S().DPanic("获取配置失败：", err.Error())
	}
	global.ServerConfig.Port = utils.GetFreePort(global.ServerConfig.Port)
}
