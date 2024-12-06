package entity

import (
	"Infra/internal/dockr/config"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

var _ ContainerConfiguration = &containerEntity{}

type ContainerStatus string

func ContainerStatusFailed() ContainerStatus  { return "failed" }
func ContainerStatusCreated() ContainerStatus { return "created" }
func ContainerStatusRunning() ContainerStatus { return "running" }
func ContainerStatusStopped() ContainerStatus { return "stopped" }


type ContainerConfiguration interface {
	// in app config
	GetContainerConfig() config.ContainerConfiguration

	// docker internal config
	GetNetworkConfig() *network.NetworkingConfig
	GetHostConfig() *container.HostConfig
	GetConfig() *container.Config
  GetHealthCheckConfig() *container.HealthConfig 
	
	GetID() string
  GetService() string
}

// Container is a representation of running docker container (service.Service + Docker API)
type containerEntity struct {
	id              string
	service 				string
	status          ContainerStatus

	containerConfig config.ContainerConfiguration
	
	// docker internals 
	config          *container.Config
	hostConfig      *container.HostConfig
	networkConfig   *network.NetworkingConfig
	healthCheckConfig *container.HealthConfig

	mu *sync.RWMutex
}

func NewContainer(conf config.ContainerConfiguration) (ContainerConfiguration, error) {

	var res = container.Resources{}

	switch conf.GetLoadLevel() {
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
			Test:        conf.GetHealthTest(),       
			Interval:    conf.GetHealthInterval(),   
			Timeout:     conf.GetHealthTimeout(),    
			Retries:     conf.GetHealthRetries(),    
			StartPeriod: conf.GetHealthStartPeriod(),
	}
	networkConfig := &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				conf.GetNetworkID(): {
					NetworkID: conf.GetNetworkID(),
					Aliases:   []string{conf.GetHostname()},
				},
			},
	}

	return &containerEntity{
		id:              "",
		containerConfig: conf,
		config:          containerConfig,
		hostConfig:      hostConfig,
		networkConfig:   networkConfig,
		healthCheckConfig: healthCheckConfig,
		status:          ContainerStatusCreated(),
		mu:              &sync.RWMutex{},
	}, nil
}

//--------------------------------------

func (c *containerEntity) StatusStart() {
 	c.mu.Lock()
  defer c.mu.Unlock()
	c.status = ContainerStatusRunning()
}

func (c *containerEntity) StatusStop() {
 	c.mu.Lock()
  defer c.mu.Unlock()
	c.status = ContainerStatusStopped()
}

func (c *containerEntity) GetID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.id
}

func (c *containerEntity) GetService() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.service
}

//--------------------------------------

func (c *containerEntity) GetHealthCheckConfig() *container.HealthConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.healthCheckConfig
}

func (c *containerEntity) GetContainerConfig() config.ContainerConfiguration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.containerConfig
}

func (c *containerEntity) GetConfig() *container.Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

func (c *containerEntity) GetHostConfig() *container.HostConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hostConfig
}

func (c *containerEntity) GetNetworkConfig() *network.NetworkingConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.networkConfig
}

//--------------------------------------