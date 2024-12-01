package entity

import (
	"Infra/internal/dockr/service"
	"fmt"
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
	ID            string
	Service       service.Service
	Config        *container.Config
	HostConfig    *container.HostConfig
	NetworkConfig *network.NetworkingConfig
	Status        ContainerStatus
}

func ConvertEnvVars(envVars map[string]string) []string {
	envList := make([]string, 0, len(envVars))
	for key, value := range envVars {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}
	return envList
}
