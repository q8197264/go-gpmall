package initialize

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"webServer/users/global"
	"webServer/users/utils"
)

// 配置初始化
func InitConfig() {
	path, _ := os.Getwd()
	if debug := getEnv("GO_WEBSERVER_DEBUG_CONFIG"); debug {
		path = fmt.Sprintf("%s/config/config-debug.yaml", path)
	} else {
		path = fmt.Sprintf("%s/config/config.yaml", path)
	}
	cfg := viper.New()
	cfg.SetConfigFile(path)
	if err := cfg.ReadInConfig(); err != nil {
		zap.S().Errorf(err.Error())
	}

	// bind struct
	if err := cfg.Unmarshal(global.NacosConfig); err != nil {
		zap.S().Infof("bind struct err:", err.Error())
	}

	//load dynamic
	cfg.WatchConfig()
	cfg.OnConfigChange(func(e fsnotify.Event) {
		cfg.ReadInConfig()
		cfg.Unmarshal(&global.NacosConfig)
	})

	if b := getEnv("GO_WEBSERVER_DEBUG_CONFIG"); !b {
		setFreePort()
	}

	LoadConfig()
}

// 获取系统环境变量
func getEnv(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
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

// 载入配置中心的配置
func LoadConfig() {
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
			ContextPath: "/nacos",
			Scheme:      "http",
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
}
