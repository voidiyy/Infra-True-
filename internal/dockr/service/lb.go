package service

import (
	"log"
	"sync"
)

var (
	_ Service = &SrvLB{}
)

const (
	imageNginxLatest   = "nginx:latest"
	imageTraefikLatest = "traefik:latest"
)

// defaultLBImages is a list of docker load balancer images
var defaultLBImages = []string{
	imageNginxLatest,
	imageTraefikLatest,
}

type SrvLB struct {
	logic *SrvLogic
	mu    *sync.RWMutex
}

func DefaultSrvLogicLB() *SrvLB {
	return &SrvLB{
		logic: &SrvLogic{
			IsDefault:     true,
			Type:          LB,
			Image:         "haproxy:latest",
			ContainerName: "load-balancer",
			Network:       "default-network",
			Ports:         []string{"80:80", "443:443"},
			Volumes:       []string{"/etc/haproxy:/usr/local/etc/haproxy"},
			EnvVars:       map[string]string{},
			RestartPolicy: "always",
		},
		mu: &sync.RWMutex{},
	}
}

func NewSrvLogicLB(conf *SrvConfig) *SrvLB {
	return &SrvLB{
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

//--------------------------------------------------------

func (sl *SrvLB) Info() {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	log.Println("----- LB Info -----")
	log.Printf("Service Type: %s", sl.logic.Type)
	log.Printf("Current Image: %s", sl.logic.Image)
	log.Printf("Container Name: %s", sl.logic.ContainerName)
	log.Printf("Status: %s", sl.logic.Status)
	log.Printf("Network: %s", sl.logic.Network)
	log.Printf("Ports: %v", sl.logic.Ports)
	log.Printf("Volumes: %v", sl.logic.Volumes)
	log.Printf("Environment Variables: %v", sl.logic.EnvVars)
	log.Printf("Restart Policy: %s", sl.logic.RestartPolicy)
	log.Println("-------------------")
}

func (sl *SrvLB) GetLogic() *SrvLogic {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	return sl.logic
}

func (sl *SrvLB) GetImage() string {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	return sl.logic.Image
}

func (sl *SrvLB) GetType() SrvType {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	return sl.logic.Type
}

func (sl *SrvLB) GetStatus() SrvStatus {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	return sl.logic.Status
}

func (sl *SrvLB) SetStatus(status SrvStatus) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.logic.Status = status
}

func (sl *SrvLB) ListImages() []string {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	return defaultLBImages
}
