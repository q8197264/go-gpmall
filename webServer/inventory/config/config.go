package config

type NacosConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      int    `mapstructure:"port" json:"port"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
	Dataid    string `mapstructure:"dataid" json:"dataid"`
	Group     string `mapstructure:"group" json:"group"`
}

type ConsulConfig struct {
	Host string   `mapstructure:"host" json:"host"`
	Port int      `mapstructure:"port" json:"port"`
	Tags []string `mapstructure:"tags" json:"tags"`
}

type GrpcClientConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type JwtConfig struct {
	Key string `mapstructure:"key" json:"key"`
}

type JaegerConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name       string           `mapstructure:"name" json:"name"`
	Host       string           `mapstructure:"host" json:"host"`
	Port       int              `mapstructure:"port" json:"port"`
	Consul     ConsulConfig     `mapstructure:"consul" json:"consul"`
	GrpcClient GrpcClientConfig `mapstructure:"grpc_addr" json:"grpc_addr"`
	Jwt        JwtConfig        `mapstructure:"jwt" json:"jwt"`
	Jaeger     JaegerConfig     `mapstructure:"jaeger"`
}
