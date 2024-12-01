package service

import (
	"log"
	"sync"
)

var _ Service = &SrvServerAdd{}

const (
	imageNodeJsLatest = "node:latest"
	imageGoLatest     = "golang:latest"
)

var defaultAddServerImages = []string{
	imageNodeJsLatest,
	imageGoLatest,
}

type SrvServerAdd struct {
	logic *SrvLogic
	mu    *sync.RWMutex
}

func DefaultSrvLogicServerAdd() *SrvServerAdd {
	return &SrvServerAdd{
		logic: &SrvLogic{
			IsDefault:     true,
			Type:          ServerAdd,
			Image:         "nginx:latest",
			ContainerName: "additional-server",
			Network:       "default-network",
			Ports:         []string{"8080:80"},
			Volumes:       []string{"/etc/nginx/conf.d:/etc/nginx/conf.d"},
			EnvVars:       map[string]string{},
			RestartPolicy: "always",
		},
		mu: &sync.RWMutex{},
	}
}

func NewSrvLogicServerAdd(conf *SrvConfig) *SrvServerAdd {
	return &SrvServerAdd{
		logic: &SrvLogic{
			IsDefault:     false,
			Type:          conf.Type,
			Image:         conf.Image,
			ContainerName: conf.ContainerName,
			Status:        SrvStatusRun(),
			Network:       conf.Network,
			Ports:         conf.Ports,
			Volumes:       conf.Volumes,
			EnvVars:       conf.EnvVars,
			RestartPolicy: conf.RestartPolicy,
		},
		mu: &sync.RWMutex{},
	}
}

// ----------------------------------------------------------------

func (sa *SrvServerAdd) Info() {
	sa.mu.Lock()
	defer sa.mu.Unlock()

	log.Println("----- Server Add Info -----")
	log.Printf("Service Type: %s", sa.logic.Type)
	log.Printf("Current Image: %s", sa.logic.Image)
	log.Printf("Container Name: %s", sa.logic.ContainerName)
	log.Printf("Status: %s", sa.logic.Status)
	log.Printf("Network: %s", sa.logic.Network)
	log.Printf("Ports: %v", sa.logic.Ports)
	log.Printf("Volumes: %v", sa.logic.Volumes)
	log.Printf("Environment Variables: %v", sa.logic.EnvVars)
	log.Printf("Restart Policy: %s", sa.logic.RestartPolicy)
	log.Println("-------------------")
}

func (sa *SrvServerAdd) GetLogic() *SrvLogic {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	return sa.logic
}

func (sa *SrvServerAdd) GetImage() string {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	return sa.logic.Image
}

func (sa *SrvServerAdd) GetType() SrvType {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	return sa.logic.Type
}

func (sa *SrvServerAdd) GetStatus() SrvStatus {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	return sa.logic.Status
}

func (sa *SrvServerAdd) SetStatus(status SrvStatus) {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	sa.logic.Status = status
}

func (sa *SrvServerAdd) ListImages() []string {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	return sa.logic.Ports
}
