package initialize

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"webServer/inventory/global"
	"webServer/oss/utils"
)

func InitNacos() {
	readNacosConfig()

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
			ContextPath: "/nacos",
			Port:        uint64(global.NacosConfig.Port),
			Scheme:      "http",
		},
	}
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.Dataid,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		zap.S().DPanic(err.Error())
	}

	if err := json.Unmarshal([]byte(content), global.ServerConfig); err != nil {
		zap.S().DPanic(err.Error())
	}

	global.ServerConfig.Port = utils.GetFreePort(global.ServerConfig.Port)
}

// 读取本地配置
func readNacosConfig() {
	path, _ := filepath.Abs(".")
	if getEnv() {
		path = fmt.Sprintf("%s/config/debug.yaml", path)
	} else {
		path = fmt.Sprintf("%s/config/config.yaml", path)
	}

	vi := viper.New()
	vi.SetConfigFile(path)
	if err := vi.ReadInConfig(); err != nil {
		panic(err.Error())
	}

	if err := vi.Unmarshal(&global.NacosConfig); err != nil {
		zap.S().DPanic(err.Error())
	}

	// 动态更新
	vi.WatchConfig()
	vi.OnConfigChange(func(in fsnotify.Event) {
		vi.ReadInConfig()
		vi.Unmarshal(&global.NacosConfig)
	})
}

func getEnv() bool {
	viper.AutomaticEnv()
	return viper.GetBool("GO_WEBSERVER_DEBUG_CONFIG")
}
