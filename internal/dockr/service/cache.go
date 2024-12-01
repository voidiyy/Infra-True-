package service

import (
	"log"
	"sync"
)

var _ Service = &SrvCache{}

const (
	imageRedisLatest = "redis:latest"
	imageMemcached   = "memcached:latest"
)

// defaultCacheImages is a list of docker cache images
var defaultCacheImages = []string{
	imageRedisLatest,
	imageMemcached,
}

type SrvCache struct {
	logic *SrvLogic
	mu    *sync.RWMutex
}

func DefaultSrvLogicCache() *SrvCache {
	return &SrvCache{
		logic: &SrvLogic{
			IsDefault:     true,
			Type:          Cache,
			Image:         "redis:latest", // наприклад, Redis для кешу
			ContainerName: "cache-container",
			Status:        SrvStatusRun(),
			Network:       "default-network",
			Ports:         []string{"6379:6379"},
			Volumes:       []string{"/data:/data"},
			EnvVars:       map[string]string{"REDIS_PASSWORD": "defaultpassword"},
			RestartPolicy: "always",
		},
		mu: &sync.RWMutex{},
	}
}

func NewSrvLogicCache(conf *SrvConfig) *SrvCache {
	return &SrvCache{
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

//-----------------------------------------------------------------------------

func (sc *SrvCache) GetLogic() *SrvLogic {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.logic
}

func (sc *SrvCache) GetType() SrvType {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.logic.Type
}

func (sc *SrvCache) GetStatus() SrvStatus {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.logic.Status
}

func (sc *SrvCache) SetStatus(status SrvStatus) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.logic.Status = status
}

func (sc *SrvCache) ListImages() []string {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return defaultCacheImages
}

func (sc *SrvCache) Info() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	log.Println("----- Cache Info -----")
	log.Printf("Service Type: %s", sc.logic.Type)
	log.Printf("Current Image: %s", sc.logic.Image)
	log.Printf("Container Name: %s", sc.logic.ContainerName)
	log.Printf("Status: %s", sc.logic.Status)
	log.Printf("Network: %s", sc.logic.Network)
	log.Printf("Ports: %v", sc.logic.Ports)
	log.Printf("Volumes: %v", sc.logic.Volumes)
	log.Printf("Environment Variables: %v", sc.logic.EnvVars)
	log.Printf("Restart Policy: %s", sc.logic.RestartPolicy)
	log.Println("-------------------")
}

func (sc *SrvCache) GetImage() string {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.logic.Image
}
