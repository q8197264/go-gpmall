package config

type NacosConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      int    `mapstructure:"port" json:"port"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	Group     string `mapstructure:"group" json:"group"`
	DataId    string `mapstructure:"dataid" json:"dataid"`
}

type ConsulConfig struct {
	Tags []string `mapstructure:"tags" json:"tags"`
	Host string   `mapstructure:"host" json:"host"`
	Port int      `mapstructure:"port" json:"port"`
}

type JwtConfig struct {
	Key string `mapstructure:"key"`
}

type GrpcConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type JaegerConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type AlipayConfig struct {
	AppId        string `mapstructure:"appid" json:""`
	AliPublicKey string `mapstructure:"ali_public_key" json:"ali_public_key"`
	PrivateKey   string `mapstructure:"private_key" json:"private_key"`
	CallbackUrl  string `mapstructure:"callback_url" json:"callback_url"`
	NotifyUrl    string `mapstructure:"notify_url" json:"notify_url"`
}

type ServerConfig struct {
	Name   string       `mapstructure:"name" json:"name"`
	Host   string       `mapstructure:"host" json:"host"`
	Port   int          `mapstructure:"port" json:"port"`
	Jwt    JwtConfig    `mapstructure:"jwt" json:"jwt"`
	Consul ConsulConfig `mapstructure:"consul"`
	Grpc   GrpcConfig   `mapstructure:"grpc"`
	Alipay AlipayConfig `mapstructure:"alipay"`
	Jaeger JaegerConfig `mapstructure:"jaeger"`
}
