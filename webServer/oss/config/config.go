package config

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Dataid    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

type ConsulConfig struct {
	Host string   `mapstructure:"host" json:"host"`
	Port int      `mapstructure:"prot" json:"port"`
	Tags []string `mapstructure:"tags" json:"tags"`
}

type GrpcConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type KodoConfig struct {
	AccessKey   string `mapstructure:"accesskey" json:"accesskey"`
	SecretKey   string `mapstructure:"secretKey" json:"secretKey"`
	Bucket      string `mapstructure:"bucket" json:"bucket"`
	CallbackURL string `mapstructure:"callback_url" json:"callback_url"`
}

type JWTConfig struct {
	Key string `mapstructure:"key" json:"key"`
}

type JaegerConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name     string       `mapstructure:"name" json:"name"`
	Host     string       `mapstructure:"host" json:"host"`
	Port     int          `mapstructure:"port" json:"port"`
	Consul   ConsulConfig `mapstructure:"consul" json:"consul"`
	GrpcAddr GrpcConfig   `mapstructure:"grpc_addr" json:"grpc_addr"`
	Kodo     KodoConfig   `mapstructure:"kodo" json:"kodo"`
	Jwt      JWTConfig    `mapstructure:"jwt" json:"jwt"`
	Jaeger   JaegerConfig `mapstructure:"jaeger"`
}
