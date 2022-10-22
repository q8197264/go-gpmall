package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
)

type ConsulClient struct {
	name      string
	host      string
	port      int
	tags      []string
	serviceId string
	client    *api.Client
}

func New(name string, host string, port int, tags []string, serviceId string) *ConsulClient {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", host, port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().DPanic(err.Error())
	}
	return &ConsulClient{
		name:      name,
		host:      host,
		port:      port,
		tags:      tags,
		serviceId: serviceId,
		client:    client,
	}
}

func (c *ConsulClient) ServiceRegister(host string, port int) error {
	cfg := &api.AgentServiceRegistration{
		ID:      c.serviceId,
		Name:    c.name,
		Tags:    c.tags,
		Port:    port,
		Address: host,
	}
	cfg.Check = &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/v1/health", host, port),
		Timeout:                        "5s",
		Interval:                       "60s",
		DeregisterCriticalServiceAfter: "5s",
	}

	return c.client.Agent().ServiceRegister(cfg)
}

func (c *ConsulClient) ServiceDeregister() error {
	return c.client.Agent().ServiceDeregister(c.serviceId)
}
