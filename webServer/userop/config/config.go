package config

type NacosConf struct {
	Name      string `mapstructure:"name" json:"name"`
	Host      string `mapstructure:"host" json:"host"`
	Port      int    `mapstructure:"port" json:"port"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	DataId    string `mapstructure:"dataid" json:"dataid"`
	Group     string `mapstructure:"group" json:"group"`
}

type consulConf struct {
	Host string   `mapstructure:"host" json:"host"`
	Port int      `mapstructure:"port" json:"port"`
	Tags []string `mapstructure:"tags" json:"tags"`
}

type grpcConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type jwt struct {
	Key string `mapstructure:"key" json:"key"`
}

type JaegerConfig struct {
	Name string `mapstructure:"name" json:"name"`
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConf struct {
	Name       string       `mapstructure:"name" json:"name"`
	Host       string       `mapstructure:"host" json:"host"`
	Port       int          `mapstructure:"port" json:"port"`
	Consul     consulConf   `mapstructure:"consul" json:"consul"`
	GrpcClient grpcConfig   `mapstructure:"grpc_addr" json:"grpc_addr"`
	Jwt        jwt          `mapstructure:"jwt" json:"jwt"`
	Jaeger     JaegerConfig `mapstructure:"jaeger"`
}
