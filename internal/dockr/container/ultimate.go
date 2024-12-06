package entity

import (
	"Infra/internal/dockr/config"
	"fmt"
	"log"
	"sync"
)

type UltimateContainer struct {
	Containers map[string]ContainerConfiguration

	mu *sync.RWMutex
}

func NewUltimateContainer(configs *config.UltimateConfig) (*UltimateContainer,error) {
	ulti := make(map[string]ContainerConfiguration, len(configs.Containers))
	
	if configs.Containers == nil {
		return nil, fmt.Errorf("conatainer configuration is nil")
	}
	
	for _, v := range configs.Containers {
		cont, err := NewContainer(v.GetFull())
		if err != nil {
			return nil, fmt.Errorf("container creation error: %s", err)
		}
		
		ulti[v.GetService()] = cont
		log.Println("added container: ", v.GetService())
	}
	
	if len(configs.Containers) != len(ulti) {
		return nil, fmt.Errorf("mismatch created containers")
	}
	
	return &UltimateContainer{
		Containers: ulti,
		mu: &sync.RWMutex{},
	}, nil
}

func (uc *UltimateContainer) RemoveContainer(str string) error {
	_, ok := uc.Containers[str]
	if !ok {
		return fmt.Errorf("no conteiner to remove %s", str)
	}	
	
	delete(uc.Containers, str)
	return nil
}

func (uc *UltimateContainer) GetContainerByService(service string) (ContainerConfiguration, error) {
    uc.mu.RLock()
    defer uc.mu.RUnlock()
    container, ok := uc.Containers[service]
    if !ok {
        return nil, fmt.Errorf("container with service %s not found", service)
    }
    return container, nil
}

