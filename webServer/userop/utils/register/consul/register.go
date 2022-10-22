package consul

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

type ConsulClient struct {
	client    *consul.Client
	serviceId string
}

func NewClient(host string, port int, serviceId string) *ConsulClient {
	cfg := consul.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", host, port)
	client, err := consul.NewClient(cfg)
	if err != nil {
		zap.S().DPanicf("consul连接失败:", err.Error())
	}

	return &ConsulClient{
		client:    client,
		serviceId: serviceId,
	}
}

func (c *ConsulClient) Register(name string, host string, port int, tags []string) error {
	ss := &consul.AgentServiceRegistration{
		ID:      c.serviceId,
		Name:    name,
		Tags:    tags,
		Port:    port,
		Address: host,
	}
	ss.Check = &consul.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/v1/health", host, port),
		Timeout:                        "5s",
		Interval:                       "60s",
		DeregisterCriticalServiceAfter: "5s",
	}

	return c.client.Agent().ServiceRegister(ss)
}

func (c *ConsulClient) Deregister() error {
	return c.client.Agent().ServiceDeregister(c.serviceId)
}
