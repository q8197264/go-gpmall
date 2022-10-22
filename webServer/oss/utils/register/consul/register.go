package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

type client struct {
	host      string
	port      int
	serviceID string
}

func NewClient(host string, port int, serviceID string) *client {
	c := &client{
		host:      host,
		port:      port,
		serviceID: serviceID,
	}
	return c
}

// 注册服务
func (c *client) Register(
	host string,
	port int,
	name string,
	tags []string,
) error {
	cfg := &api.AgentServiceRegistration{
		ID:      c.serviceID,
		Name:    name,
		Tags:    tags,
		Port:    c.port,
		Address: c.host,
	}
	cfg.Check = &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", host, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	err := c.getClient().Agent().ServiceRegister(cfg)

	return err
}

// 注销服务
func (c *client) Deregister() error {
	return c.getClient().Agent().ServiceDeregister(c.serviceID)
}

// 连接注册中心
func (c *client) getClient() *api.Client {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", c.host, c.port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().DPanic(err.Error())
	}

	return client
}
