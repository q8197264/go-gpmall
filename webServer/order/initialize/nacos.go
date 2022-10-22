package initialize

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"webServer/order/global"
)

func InitNacosConfig() {
	readLocalNacosConf()

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
	// Another way of create config client for dynamic configuration (recommend)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
	}
	err = json.Unmarshal([]byte(content), global.ServerConfig)
	if err != nil {
	}
}

func readLocalNacosConf() {
	realPath, _ := filepath.Abs(".")
	filename := "develop.yaml"
	if os.Getenv("GO") == "false" {
		filename = "produce.yaml"
	}
	path := filepath.Join(realPath, "config", filename)
	vp := viper.New()
	// vp.AddConfigPath(path)
	vp.SetConfigFile(path)
	if err := vp.ReadInConfig(); err != nil {
		zap.S().DPanic(err.Error())
	}
	vp.Unmarshal(&global.NacosConfig)

	vp.WatchConfig()
	vp.OnConfigChange(func(in fsnotify.Event) {
		vp.ReadInConfig()
		vp.Unmarshal(&global.NacosConfig)
	})
}
