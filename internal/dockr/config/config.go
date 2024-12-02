package config

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
)

var mu *sync.Mutex

// ContainerConfig is a representation of docker container config
// it UNI for all services
type ContainerConfig struct {
	// config
	IsDefault     bool   `yaml:"is_default" json:"is_default"`
	ContainerType string `yaml:"container_type" json:"container_type"`

	//docker config
	Image         string            `yaml:"image" json:"image"`
	ContainerName string            `yaml:"container_name" json:"container_name"`
	Network       string            `yaml:"network" json:"network"`
	Ports         []string          `yaml:"ports" json:"ports"`
	Volumes       []string          `yaml:"volumes" json:"volumes"`
	EnvVars       map[string]string `yaml:"env_vars" json:"env_vars"`
	RestartPolicy string            `yaml:"restart_policy" json:"restart_policy"`
}

func (c *ContainerConfig) Validate() error {
	// Check for required fields
	if c.ContainerType == "" {
		return errors.New("container type is required")
	}

	if c.Image == "" {
		return errors.New("image is required")
	}

	if c.ContainerName == "" {
		return errors.New("container name is required")
	}

	if c.Network == "" {
		return errors.New("network is required")
	}

	// Validate Ports
	for _, portMapping := range c.Ports {
		if err := validatePortMapping(portMapping); err != nil {
			return fmt.Errorf("invalid port mapping '%s': %v", portMapping, err)
		}
	}

	// Validate Volumes
	for _, volume := range c.Volumes {
		if !strings.Contains(volume, ":") {
			return fmt.Errorf("invalid volume format '%s', must be 'source:target'", volume)
		}
	}

	// Restart policy validation
	validPolicies := map[string]bool{
		"always":     true,
		"on-failure": true,
		"no":         true,
	}
	if _, ok := validPolicies[c.RestartPolicy]; !ok {
		return fmt.Errorf("invalid restart policy '%s'", c.RestartPolicy)
	}

	return nil
}

func (c *ContainerConfig) GetDefault() bool {
	mu.Lock()
	defer mu.Unlock()
	return c.IsDefault
}

func (c *ContainerConfig) GetType() string {
	mu.Lock()
	defer mu.Unlock()
	return c.ContainerType
}

func (c *ContainerConfig) GetImage() string {
	mu.Lock()
	defer mu.Unlock()
	return c.Image
}

func (c *ContainerConfig) GetContainerName() string {
	mu.Lock()
	defer mu.Unlock()
	return c.ContainerName
}

func (c *ContainerConfig) GetNetwork() string {
	mu.Lock()
	defer mu.Unlock()
	return c.Network
}

func (c *ContainerConfig) GetPorts() []string {
	mu.Lock()
	defer mu.Unlock()
	return c.Ports
}

func (c *ContainerConfig) GetVolumes() []string {
	mu.Lock()
	defer mu.Unlock()
	return c.Volumes
}

func (c *ContainerConfig) GetEnvVars() []string {
	env := make([]string, 0)
	for key, value := range c.EnvVars {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	mu.Lock()
	defer mu.Unlock()
	return env
}

func (c *ContainerConfig) GetRestartPolicy() string {
	mu.Lock()
	defer mu.Unlock()
	return c.RestartPolicy
}

// validatePortMapping checks if a port mapping string is valid (e.g., "80:80").
func validatePortMapping(mapping string) error {
	parts := strings.Split(mapping, ":")
	if len(parts) != 2 {
		return errors.New("must have format 'hostPort:containerPort'")
	}

	hostPort, containerPort := parts[0], parts[1]
	if err := validatePort(hostPort); err != nil {
		return fmt.Errorf("invalid host port '%s': %v", hostPort, err)
	}
	if err := validatePort(containerPort); err != nil {
		return fmt.Errorf("invalid container port '%s': %v", containerPort, err)
	}

	return nil
}

// validatePort checks if a port is a valid numeric value within the allowed range.
func validatePort(port string) error {
	p, err := net.LookupPort("tcp", port)
	if err != nil || p < 1 || p > 65535 {
		return errors.New("must be a number between 1 and 65535")
	}
	return nil
}
