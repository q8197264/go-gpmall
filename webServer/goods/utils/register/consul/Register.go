package consul

import (
	"fmt"
	"webServer/goods/global"

	"github.com/hashicorp/consul/api"
	// uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, serverID string) error
	DeRegister(serverID string) error
}

type Registry struct {
	Host   string
	Port   int
	Client *api.Client
}

func NewRegistryClient(host string, port int) RegistryClient {
	r := &Registry{
		Host: host,
		Port: port,
	}
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Panic(err.Error())
	}
	r.Client = client

	return r
}

func (r *Registry) Register(address string, port int, name string, tags []string, serverID string) error {
	s := &api.AgentServiceRegistration{
		ID:      serverID,
		Name:    name,
		Tags:    tags,
		Port:    port,
		Address: address,
	}
	s.Check = &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/v1/health", address, port),
		Timeout:                        "5s",
		Interval:                       "15s",
		DeregisterCriticalServiceAfter: "5s",
	}

	err := r.Client.Agent().ServiceRegister(s)
	if err != nil {
		zap.S().Warnf("%s 服务注销失败:%s", global.ServerConfig.Name, err.Error())
	}

	return err
}

func (r *Registry) DeRegister(serverID string) error {
	err := r.Client.Agent().ServiceDeregister(serverID)
	if err != nil {
		zap.S().Warnf("%s 服务注销失败:%s", global.ServerConfig.Name, err.Error())
	} else {
		fmt.Printf("%s 服务注销成功\n", global.ServerConfig.Name)
		zap.S().Infof("%s 服务注销成功", global.ServerConfig.Name)
	}
	return err
}
