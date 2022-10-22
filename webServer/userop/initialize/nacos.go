package initialize

import (
	"encoding/json"
	"path/filepath"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"webServer/userop/global"
)

func InitConfig() {
	readLocalConfig()
	clientConfig := constant.ClientConfig{
		TimeoutMs:   5000,
		NamespaceId: global.NacosConf.Namespace,
		// UpdateThreadNum:     5,
		NotLoadCacheAtStart: true,
		CacheDir:            "/tmp/nacos/cache",
		LogDir:              "/tmp/nacos/log",
		LogLevel:            "debug",
	}
	serverConfig := []constant.ServerConfig{
		{
			Scheme:      "http",
			ContextPath: "/nacos",
			IpAddr:      global.NacosConf.Host,
			Port:        uint64(global.NacosConf.Port),
		},
	}
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfig,
		},
	)
	if err != nil {
		zap.S().DPanicf("连接nacos失败：", err.Error())
	}

	content, err := client.GetConfig(vo.ConfigParam{
		DataId: global.NacosConf.DataId,
		Group:  global.NacosConf.Group,
	})
	if err != nil {
		zap.S().DPanicf("获取nacos配置失败：", err.Error())
	}
	if err := json.Unmarshal([]byte(content), &global.ServerConfig); err != nil {
		zap.S().DPanic("绑定nacos配置失败", err.Error())
	}
}

func readLocalConfig() {
	root, _ := filepath.Abs(".")
	path := filepath.Join(root, "config", "develop.yaml")
	vp := viper.New()
	vp.SetConfigFile(path)
	if err := vp.ReadInConfig(); err != nil {
		zap.S().DPanicf("读取nacos配置失败:", err.Error())
	}
	if err := vp.Unmarshal(&global.NacosConf); err != nil {
		zap.S().DPanicf("配置赋值失败：", err.Error())
	}
}
