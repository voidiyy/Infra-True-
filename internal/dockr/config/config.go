package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

var _ ContainerConfiguration = &ContainerConfig{}

type ContainerConfiguration interface {
	GetNetworkID() string
	GetRestartPolicy() container.RestartPolicy
	GetVolumes() []string
	GetCMD() strslice.StrSlice
	GetWorkingDir() string
	GetHostname() string
	GetDefault() bool
	GetService() string
	GetPorts() nat.PortMap
	GetEnvVars() []string
	GetImage() string
	GetNetworkMode() container.NetworkMode
	GetLoadLevel() int 
	
	GetFull() *ContainerConfig 
	
	GetHealthTest() []string
	GetHealthInterval() time.Duration
	GetHealthTimeout() time.Duration
	GetHealthRetries() int
	GetHealthStartPeriod() time.Duration
}

// ContainerConfig represents the full configuration for a Docker container.
// It is unified for all services.
type ContainerConfig struct {
	// If IsDefault == true, configuration will be created with default values.
	LoadLevel int `yaml:"load_level" json:"load_level"`
	IsDefault     bool   `yaml:"is_default" json:"is_default"`  // Indicates whether to use default values for this configuration.
	ContainerService string `yaml:"container_service" json:"container_service"` // Type of the container (e.g., "web", "db", etc.).

	// Docker &container.Config{}
	Image         string            `yaml:"image" json:"image"`         // The image to use for the container.
	Hostname     string            `yaml:"hostname" json:"hostname"` // The hostname to use for the container.
	EnvVars       map[string]string `yaml:"env_vars" json:"env_vars"`   // Environment variables to set in the container.
	WorkingDir    string            `yaml:"working_dir" json:"working_dir"` // The working directory for commands to run in.
	Cmd []string `yaml:"cmd" json:"cmd"` // Command to run in the container on startup.
	
	// Docker &container.HostConfig{}
	Volumes       []string          `yaml:"volumes" json:"volumes"`     // List of volumes to mount into the container.
	NetworkMode   string            `yaml:"network_mode" json:"network_mode"` // The network mode for the container.
	Ports         []string          `yaml:"ports" json:"ports"`         // List of ports to expose from the container.
	RestartPolicy string            `yaml:"restart_policy" json:"restart_policy"` // Docker restart policy (e.g., "always", "on-failure").
	
	// Docker &network.NetworkingConfig{}
	NetworkID       string            `yaml:"network" json:"network"`     // The name of the network for the container.

	// HealthCheck configuration for Docker health check (ping of server every 5 minutes, or similar).
	HealthCheck HealthCheckConfig `yaml:"health_check" json:"health_check"`
}

func (c *ContainerConfig) GetLoadLevel() int {
	return c.LoadLevel
}

func (c *ContainerConfig) GetFull() *ContainerConfig {
	return c
}

func (c *ContainerConfig) GetNetworkID() string {
	return c.NetworkID
}

func (c *ContainerConfig) GetRestartPolicy() container.RestartPolicy {
	switch c.RestartPolicy {
		case "no": return container.RestartPolicy{
			Name: container.RestartPolicyDisabled,
		}
		case "always": return container.RestartPolicy{
			Name: container.RestartPolicyAlways,
		}
		case "on-failure": return container.RestartPolicy{
			Name: container.RestartPolicyOnFailure,
		}
		case "unless-stopped": return container.RestartPolicy{
			Name: container.RestartPolicyUnlessStopped,
		}
	}
	return container.RestartPolicy{
		Name: container.RestartPolicyAlways,
	}
}

func (c *ContainerConfig) GetVolumes() []string {
	return c.Volumes
}

func (c *ContainerConfig) GetCMD() strslice.StrSlice{
	return strslice.StrSlice(c.Cmd)
}

func (c *ContainerConfig) GetWorkingDir() string {
	return c.WorkingDir
}

func (c *ContainerConfig) GetHostname() string {
	return c.Hostname
}

func (c *ContainerConfig) GetDefault() bool {
	return c.IsDefault
}

func (c *ContainerConfig) GetService() string {
	return c.ContainerService
}

func (c *ContainerConfig) GetImage() string {
	return c.Image
}

func (c *ContainerConfig) GetEnvVars() []string {
	env := make([]string, len(c.EnvVars))
	for key, value := range c.EnvVars {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}


func (c *ContainerConfig) GetPorts() nat.PortMap {
	portBind := make(nat.PortMap)
	for _, portMapping := range c.Ports {
		parts := strings.Split(portMapping, ":")
		if len(parts) == 2 {
			hostPort := parts[0]
			containerPort := fmt.Sprintf("%s/tcp", parts[1])
			portBind[nat.Port(containerPort)] = []nat.PortBinding{
				{HostPort: hostPort},
				{HostIP: "0.0.0.0"},
			}
		}
	}
	return portBind
}

func (c *ContainerConfig) GetNetworkMode() container.NetworkMode {
	return container.NetworkMode(c.NetworkMode)
}

// HealthCheckConfig defines the configuration for Docker container health checks.
type HealthCheckConfig struct {
	// Test command to perform health check.
	Test        []string `yaml:"test" json:"test"` // Command for checking health.
	Interval    string   `yaml:"interval" json:"interval"` // Time between checks (e.g., "30s").
	Timeout     string   `yaml:"timeout" json:"timeout"` // Timeout for health check (e.g., "5s").
	Retries     int      `yaml:"retries" json:"retries"` // Number of retries before considering the container unhealthy.
	StartPeriod string   `yaml:"start_period" json:"start_period"` // Initial delay before the first health check (e.g., "10s").
}

//------------------- HEALTH CHECK ------------------------

func (c *ContainerConfig) GetHealthTest() []string {
	return c.HealthCheck.Test
}

func (c *ContainerConfig) GetHealthInterval() time.Duration {
	inerv, _ := time.ParseDuration(c.HealthCheck.Interval)
	return inerv
}

func (c *ContainerConfig) GetHealthTimeout() time.Duration {
	timeout, _ := time.ParseDuration(c.HealthCheck.Timeout)
	return timeout
}

func (c *ContainerConfig) GetHealthRetries() int {
	return c.HealthCheck.Retries
}

func (c *ContainerConfig) GetHealthStartPeriod() time.Duration {
	strat, _ := time.ParseDuration(c.HealthCheck.StartPeriod)
	return strat
}

