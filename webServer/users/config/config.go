package config

// jwt signKey config
type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

// grpc dato-server config
type GrpcAddr struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

// register center
type Consul struct {
	Host string   `mapstructure:"host" json:"host"`
	Port int      `mapstructure:"port" json:"port"`
	Tags []string `mapstructure:"tags" json:"tags"`
}

type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port"  json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

// web server config
type ServerConfig struct {
	Name      string       `mapstructure:"name" json:"name"`
	Host      string       `mapstructure:"host" json:"host"`
	Port      int          `mapstructure:"port" json:"port"`
	GrpcAddr  GrpcAddr     `mapstructure:"grpc_addr" json:"grpc_addr"`
	Consul    Consul       `mapstructure:"consul" json:"consul"`
	JWTConfig JWTConfig    `mapstructure:"jwt" json:"jwt"`
	Jaeger    JaegerConfig `mapstructure:"jaeger" `
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
