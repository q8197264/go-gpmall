package initialize

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"webServer/goods/global"
	"webServer/goods/utils"
)

func LoadConfig() {
	path, _ := os.Getwd()
	path = filepath.Join(path, "config")
	if getEnv("GO_WEBSERVER_DEBUG_CONFIG") {
		//debug
		path = fmt.Sprintf("%s/debug.yaml", path)
	} else {
		//product
		path = fmt.Sprintf("%s/config.yaml", path)
	}

	cfg := viper.New()
	cfg.SetConfigFile(path)
	if err := cfg.ReadInConfig(); err != nil {
		zap.S().DPanic(err.Error())
	}

	if err := cfg.Unmarshal(global.NacosConfig); err != nil {
		zap.S().Warn(err.Error())
	}

	cfg.WatchConfig()
	cfg.OnConfigChange(func(in fsnotify.Event) {
		cfg.ReadInConfig()
		cfg.Unmarshal(global.NacosConfig)
	})

	// nacos config
	getNacosConfig()
}

func getEnv(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

// 获取配置中心所在配置
func getNacosConfig() {
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      global.NacosConfig.Host,
			Port:        global.NacosConfig.Port,
			
		},
	}

	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err.Error())
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		panic(err.Error())
	}

	// 要将一个json字符串转化为struct，需要去设置struct的tag为 `json:"field_name"`
	err = json.Unmarshal([]byte(content), global.ServerConfig)
	if err != nil {
		panic(err.Error())
	}

	if !getEnv("GO_WEBSERVER_DEBUG_CONFIG") {
		setFreePort()
	}
}

// 设置空闲端口
func setFreePort() {
	//online
	port, err := utils.GetFreePort()
	if err != nil {
		panic(err.Error())
	}
	global.ServerConfig.Port = port
}
