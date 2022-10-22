package consul

import (
	"fmt"
	"webServer/users/global"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, serverID string) error
	DeRegister(serviceID string) error
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
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Warnf("%s", err.Error())
	}
	r.Client = client

	return r
}

func (r *Registry) Register(address string, port int, name string, tags []string, serverID string) error {
	s := &api.AgentServiceRegistration{
		ID:      serverID,
		Name:    name,
		Address: address,
		Port:    port,
		Tags:    tags,
	}
	s.Check = &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/v1/health", address, port),
		Timeout:                        "5s",
		Interval:                       "15s",
		DeregisterCriticalServiceAfter: "5s",
	}
	err := r.Client.Agent().ServiceRegister(s)

	return err
}

func (r *Registry) DeRegister(serviceID string) error {
	err := r.Client.Agent().ServiceDeregister(serviceID)

	return err
}
