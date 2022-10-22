package config

// 配置映射struct

// nacos 配置中心
type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

// dao server config
type GrpcAddr struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

// consul center config
type Consul struct {
	Host string   `mapstructure:"host" json:"host"`
	Port int      `mapstructure:"port" json:"port"`
	Tags []string `mapstructure:"tags" json:"tags"`
}

// jwt signKey config
type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type JaegerConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port"  json:"port"`
}

type ServerConfig struct {
	Name      string       `mapstructure:"name" json:"name"`
	Host      string       `mapstructure:"host" json:"host"`
	Port      int          `mapstructure:"port" json:"port"`
	GrpcAddr  GrpcAddr     `mapstructure:"grpc_addr" json:"grpc_addr"`
	Consul    Consul       `mapstructure:"consul" json:"consul"`
	JWTConfig JWTConfig    `mapstructure:"jwt" json:"jwt"`
	Jaeger    JaegerConfig `mapstructure:"jaeger"`
}
