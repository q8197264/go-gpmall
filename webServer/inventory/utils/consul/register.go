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

func (c *client) Register(host string, port int, name string, tags []string) error {
	sr := &api.AgentServiceRegistration{
		ID:      c.serviceID,
		Name:    name,
		Tags:    tags,
		Port:    c.port,
		Address: c.host,
	}
	sr.Check = &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", host, port),
		Timeout:                        "5s",
		Interval:                       "15s",
		DeregisterCriticalServiceAfter: "5s",
	}
	return c.getClient().Agent().ServiceRegister(sr)
}

func (c *client) Deregister(name string, serviceID string) error {
	return c.getClient().Agent().ServiceDeregister(serviceID)
}

func (c *client) getClient() *api.Client {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", c.host, c.port)
	cli, err := api.NewClient(cfg)
	if err != nil {
		zap.S().DPanic(err.Error())
	}
	return cli
}
