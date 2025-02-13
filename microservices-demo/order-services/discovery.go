package discovery

import (
    "fmt"
    "log"

    consulapi "github.com/hashicorp/consul/api"
)

type ServiceRegistry struct {
    client *consulapi.Client
}

func NewServiceRegistry() (*ServiceRegistry, error) {
    config := consulapi.DefaultConfig()
    client, err := consulapi.NewClient(config)
    if err != nil {
        return nil, err
    }
    return &ServiceRegistry{client: client}, nil
}

func (sr *ServiceRegistry) RegisterService(name string, port int) error {
    reg := &consulapi.AgentServiceRegistration{
        Name: name,
        Port: port,
        Check: &consulapi.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://localhost:%d/health", port),
            Interval: "10s",
            Timeout:  "5s",
        },
    }
    return sr.client.Agent().ServiceRegister(reg)
}