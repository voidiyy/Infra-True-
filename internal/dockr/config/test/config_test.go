package config

import (
	"Infra/internal/dockr/config"
	"fmt"
	"testing"
)


func TestConfig(t *testing.T) {

	var (
	ulti *config.UltimateConfig
	err error
	
)
	
	t.Run(".json", func(t *testing.T){
		ulti, err = config.LoadContainersConfig("conf.json")
		if err != nil {
			t.Error(err)
		}
		
		if len(ulti.Containers) == 0 {
			t.Errorf("zero configs loaded")
		}
		
		db := ulti.Containers["db"]
		lb := ulti.Containers["lb"]
		
		fmt.Println(db[0])
		fmt.Println(lb[0])
	})
	
	t.Run(".yaml", func(t *testing.T){
		ulti, err = config.LoadContainersConfig("conf.yaml")
		if err != nil {
			t.Error(err)
		}
		
		if len(ulti.Containers) == 0 {
			t.Errorf("zero configs loaded")
		}
		
		db := ulti.Containers["db"]
		lb := ulti.Containers["lb"]
		
		fmt.Println(db[0])
		fmt.Println(lb[0])
	})
	
	t.Run(".yml", func(t *testing.T){
		ulti, err = config.LoadContainersConfig("conf.yml")
		if err != nil {
			t.Error(err)
		}
		
		if len(ulti.Containers) == 0 {
			t.Errorf("zero configs loaded")
		}
		
		db := ulti.Containers["db"]
		lb := ulti.Containers["lb"]
		
		fmt.Println(db[0])
		fmt.Println(lb[0])
	})
	
	t.Run("slice", func(t *testing.T){
		ulti, err = config.NewContainersConfig(configs...)
		if err != nil {
			t.Error(err)
		}
		
		if len(ulti.Containers) == 0 {
			t.Errorf("zero configs loaded")
		}
		
		db := ulti.Containers["db"]
		lb := ulti.Containers["lb"]
		
		fmt.Println(db[0])
		fmt.Println(lb[0])
	})
}



var configs = []config.ContainerConfig{
		{
			LoadLevel:        1,
			IsDefault:        true,
			ContainerService: "lb",
			Image:            "nginx:latest",
			Hostname:         "web-service",
			EnvVars: map[string]string{
				"ENV": "production",
			},
			WorkingDir:    "/var/www",
			Cmd:           []string{"nginx", "-g", "daemon off;"},
			Volumes:       []string{"/host/web:/var/www"},
			NetworkMode:   "bridge",
			Ports:         []string{"80:80"},
			RestartPolicy: "always",
			NetworkID:     "web-network",
			HealthCheck: config.HealthCheckConfig{
				Interval:    "30s",
				Timeout:     "5s",
				Retries:     3,
				Test: []string{"CMD", "curl", "-f", "http://localhost"},
			},
		},
		{
			LoadLevel:        2,
			IsDefault:        false,
			ContainerService: "db",
			Image:            "postgres:latest",
			Hostname:         "db-service",
			EnvVars: map[string]string{
				"POSTGRES_USER":     "admin",
				"POSTGRES_PASSWORD": "password",
				"POSTGRES_DB":       "mydb",
			},
			WorkingDir:    "/var/lib/postgresql/data",
			Cmd:           []string{"postgres"},
			Volumes:       []string{"/host/db:/var/lib/postgresql/data"},
			NetworkMode:   "bridge",
			Ports:         []string{"5432:5432"},
			RestartPolicy: "on-failure",
			NetworkID:     "db-network",
			HealthCheck: config.HealthCheckConfig{
				Interval:    "1m",
				Timeout:     "10s",
				Retries:     5,
				Test: []string{"CMD-SHELL", "pg_isready -U admin"},
			},
		},
}	