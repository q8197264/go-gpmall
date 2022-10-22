package main

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

func main2() {
	// deregister("user-web_1")
	// register("192.168.1.122", 5353, "user-web", "user-web_1", true)
	// getServices("192.168.1.122", 8500)
	// getServicefilter("192.168.1.122", 8500, "user-web")
}

func register(ip string, port int, name string, serverID string, check bool) {
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.1.122:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err.Error())
	}

	c := &api.AgentServiceRegistration{
		ID:      serverID,
		Name:    name,
		Tags:    []string{"gpmall"},
		Port:    port,
		Address: ip,
	}
	if check {
		check := &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/v1/health", ip, port),
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "5s",
		}
		c.Check = check
	}
	err = client.Agent().ServiceRegister(c)
	if err != nil {
		println("fail", err.Error())
	} else {
		println("success")
	}
}

func deregister(serverID string) {
	cfg := api.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err.Error())
	}

	err = client.Agent().ServiceDeregister(serverID)
	if err != nil {
		println("fail:", err.Error())
	} else {
		println("success")
	}
}

func getServices(ip string, port int) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", ip, port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err.Error())
	}
	list, err := client.Agent().Services()
	if err != nil {
		panic(err.Error())
	}
	for v, k := range list {
		println(v, k)
	}

}

func getServicefilter(ip string, port int, serverID string) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", ip, port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err.Error())
	}

	list, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, serverID))
	if err != nil {
		panic(err.Error())
	}
	for v, k := range list {
		fmt.Printf("%s <= %v\n", v, k)
	}
}
