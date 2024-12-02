package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type UltimateConfig struct {
	Containers map[string][]*ContainerConfig
}

func (c *UltimateConfig) GetContainers() map[string][]*ContainerConfig {
	return c.Containers
}

func (c *UltimateConfig) GetContainerType(containerType string) []*ContainerConfig {
	return c.Containers[containerType]
}

func NewContainersConfig(configs ...ContainerConfig) (*UltimateConfig, error) {
	ult := &UltimateConfig{
		Containers: make(map[string][]*ContainerConfig, len(configs)),
	}
	for i, config := range configs {
		if err := config.Validate(); err != nil {
			return nil, fmt.Errorf("config #%d is invalid: %w", i, err)
		}

		ult.Containers[config.ContainerType] = append(ult.Containers[config.ContainerType], &config)
	}

	return ult, nil
}

func LoadContainersConfig(path string) (*UltimateConfig, error) {
	conf := make([]*ContainerConfig, 0)

	if path == "" {
		return nil, fmt.Errorf("config file path is empty")
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	ext := filepath.Ext(path)

	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(file, &conf)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
		}
	case ".json":
		err = json.Unmarshal(file, &conf)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config file extension: %s", ext)
	}

	ulti := make(map[string][]*ContainerConfig, len(conf))
	for _, c := range conf {
		if err = c.Validate(); err != nil {
			return nil, fmt.Errorf("config is invalid: %w", err)
		}
		ulti[c.ContainerType] = append(ulti[c.ContainerType], c)
	}

	return &UltimateConfig{Containers: ulti}, nil
}
