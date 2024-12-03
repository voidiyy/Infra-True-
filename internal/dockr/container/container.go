package entity

import (
	"Infra/internal/dockr/config"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type ContainerStatus string

func ContainerStatusFailed() ContainerStatus  { return "failed" }
func ContainerStatusCreated() ContainerStatus { return "created" }
func ContainerStatusRunning() ContainerStatus { return "running" }
func ContainerStatusStopped() ContainerStatus { return "stopped" }

// Container is a representation of running docker container (service.Service + Docker API)
type Container struct {
	ID              string
	Service 				string
	Status          ContainerStatus

	ContainerConfig *config.ContainerConfig
	Config          *container.Config
	
	HostConfig      *container.HostConfig
	NetworkConfig   *network.NetworkingConfig
	HealthCheckConfig *container.HealthConfig

	mu *sync.RWMutex
}

func NewContainer(conf *config.ContainerConfig) (*Container, error) {

	var res = container.Resources{}

	switch conf.LoadLevel {
		case 0: res = config.LowLoadConfig
		case 1: res = config.MediumLoadConfig
		case 2: res = config.HighLoadConfig
		default:
		res = config.MediumLoadConfig
	}

	containerConfig := &container.Config{
		Image: conf.GetImage(),
		Env:   conf.GetEnvVars(),
		Hostname: conf.GetHostname(),
		WorkingDir: conf.GetWorkingDir(),
		Cmd: conf.GetCMD(),
	}

	hostConfig := &container.HostConfig{
		Binds: conf.GetVolumes(),
		NetworkMode: conf.GetNetworkMode(),
		PortBindings: conf.GetPorts(),
		RestartPolicy: conf.GetRestartPolicy(),
		Resources: res,
	}
	
	healthCheckConfig := &container.HealthConfig{
			Test:        conf.GetHealthTest(),       // Health check command
			Interval:    conf.GetHealthInterval(),   // Health check interval
			Timeout:     conf.GetHealthTimeout(),    // Health check timeout
			Retries:     conf.GetHealthRetries(),    // Number of retries before marking unhealthy
			StartPeriod: conf.GetHealthStartPeriod(), // Initial delay before starting health checks
	}
	networkConfig := &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				conf.GetNetworkID(): {
					NetworkID: conf.GetNetworkID(),
					Aliases:   []string{conf.GetHostname()},
				},
			},
	}

	return &Container{
		ID:              "",
		ContainerConfig: conf,
		Config:          containerConfig,
		HostConfig:      hostConfig,
		NetworkConfig:   networkConfig,
		HealthCheckConfig: healthCheckConfig,
		Status:          ContainerStatusCreated(),
		mu:              &sync.RWMutex{},
	}, nil
}

//--------------------------------------

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

func (c *Container) GetService() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Service
}

//--------------------------------------

func (c *Container) GetContainerConfig() *config.ContainerConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ContainerConfig
}

func (c *Container) GetConfig() *container.Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Config
}

func (c *Container) GetHostConfig() *container.HostConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.HostConfig
}

func (c *Container) GetNetworkConfig() *network.NetworkingConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.NetworkConfig
}

//--------------------------------------