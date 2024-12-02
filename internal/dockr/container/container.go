package entity

import (
	"Infra/internal/dockr/config"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"strings"
	"sync"
)

type ContainerStatus string

func ContainerStatusFailed() ContainerStatus  { return "failed" }
func ContainerStatusCreated() ContainerStatus { return "created" }
func ContainerStatusRunning() ContainerStatus { return "running" }
func ContainerStatusStopped() ContainerStatus { return "stopped" }

// Container is a representation of running docker container (service.Service + Docker API)
type Container struct {
	ID              string
	ContainerConfig *config.ContainerConfig
	Config          *container.Config
	HostConfig      *container.HostConfig
	NetworkConfig   *network.NetworkingConfig
	Status          ContainerStatus

	mu *sync.RWMutex
}

func NewContainer(config *config.ContainerConfig) (*Container, error) {

	containerConfig := &container.Config{
		Image: config.GetImage(),
		Env:   config.GetEnvVars(),
	}

	hostConfig := &container.HostConfig{
		PortBindings: parsePortBindings(config.GetPorts()),
		Binds:        config.GetVolumes(),
		RestartPolicy: container.RestartPolicy{
			Name: container.RestartPolicyAlways,
		},
	}

	networkConfig := &network.NetworkingConfig{}

	return &Container{
		ID:              "",
		ContainerConfig: config,
		Config:          containerConfig,
		HostConfig:      hostConfig,
		NetworkConfig:   networkConfig,
		Status:          ContainerStatusCreated(),
		mu:              &sync.RWMutex{},
	}, nil
}

func parsePortBindings(ports []string) nat.PortMap {
	portBindings := make(nat.PortMap)
	for _, mapping := range ports {
		parts := strings.Split(mapping, ":")
		if len(parts) == 2 {
			hostPort := parts[0]
			containerPort := fmt.Sprintf("%s/tcp", parts[1])
			portBindings[nat.Port(containerPort)] = []nat.PortBinding{
				{HostPort: hostPort},
			}
		}
	}
	return portBindings
}

func (c *Container) StatusStart() {
	c.Status = ContainerStatusRunning()
}

func (c *Container) StatusStop() {
	c.Status = ContainerStatusStopped()
}

func (c *Container) GetID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ID
}

func (c *Container) GetContainerConfig() *config.ContainerConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ContainerConfig
}
