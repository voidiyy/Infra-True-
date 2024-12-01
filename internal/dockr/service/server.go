package service

import (
	"log"
	"sync"
)

var (
	_ Service = &SrvServerMain{}
)

const (
	imageHttpdLatest = "httpd:latest" // Apache HTTP Server
	imageCaddyLatest = "caddy:latest"
)

// defaultMainServerImages is a list of docker HTTP server images
var defaultMainServerImages = []string{
	imageHttpdLatest,
	imageCaddyLatest,
}

type SrvServerMain struct {
	logic *SrvLogic
	mu    *sync.RWMutex
}

func DefaultSrvLogicServerMain() *SrvServerMain {
	return &SrvServerMain{
		logic: &SrvLogic{
			IsDefault:     true,
			Type:          ServerMain,
			Image:         "nginx:latest",
			ContainerName: "main-server",
			Network:       "default-network",
			Ports:         []string{"80:80"},
			Volumes:       []string{"/etc/nginx/conf.d:/etc/nginx/conf.d"},
			EnvVars:       map[string]string{},
			RestartPolicy: "always",
		},
		mu: &sync.RWMutex{},
	}
}

func NewSrvLogicServerMain(conf *SrvConfig) *SrvServerMain {
	return &SrvServerMain{
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

// --------------------------------------------------------

func (ss *SrvServerMain) Info() {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	log.Println("----- Server Main Info -----")
	log.Printf("Service Type: %s", ss.logic.Type)
	log.Printf("Current Image: %s", ss.logic.Image)
	log.Printf("Container Name: %s", ss.logic.ContainerName)
	log.Printf("Status: %s", ss.logic.Status)
	log.Printf("Network: %s", ss.logic.Network)
	log.Printf("Ports: %v", ss.logic.Ports)
	log.Printf("Volumes: %v", ss.logic.Volumes)
	log.Printf("Environment Variables: %v", ss.logic.EnvVars)
	log.Printf("Restart Policy: %s", ss.logic.RestartPolicy)
	log.Println("-------------------")
}

func (ss *SrvServerMain) GetLogic() *SrvLogic {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.logic
}

func (ss *SrvServerMain) GetImage() string {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.logic.Image
}

func (ss *SrvServerMain) GetType() SrvType {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.logic.Type
}

func (ss *SrvServerMain) GetStatus() SrvStatus {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.logic.Status
}

func (ss *SrvServerMain) SetStatus(status SrvStatus) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.logic.Status = status
}

func (ss *SrvServerMain) ListImages() []string {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return defaultMainServerImages
}
