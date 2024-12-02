package entity

import "sync"

type UltimateContainer struct {
	Containers map[string]*Container

	mu *sync.RWMutex
}

func NewUltimateContainer() *UltimateContainer {
	return &UltimateContainer{
		Containers: make(map[string]*Container),
		mu:         &sync.RWMutex{},
	}
}

func (uc *UltimateContainer) AddContainer(id string, container *Container) {
	uc.mu.Lock()

}
func (uc *UltimateContainer) GetContainers() map[string]*Container {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	return uc.Containers
}

func (uc *UltimateContainer) GetContainer(id string) *Container {
	uc.mu.Lock()
	defer uc.mu.Unlock()
	return uc.Containers[id]
}
