package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
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
	for _, config := range configs {
				ult.Containers[config.ContainerService] = append(ult.Containers[config.ContainerService], &config)
	}

	log.Printf("loaded %v configs\n", len(ult.Containers))
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
	
		ulti[c.ContainerService] = append(ulti[c.ContainerService], c)
	}

	log.Printf("loaded %v configs\n", len(ulti))
	
	return &UltimateConfig{Containers: ulti}, nil
}
