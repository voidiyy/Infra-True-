package service

import (
	"fmt"
	"log"
)

type Service interface {
	Info()
	// internal logic

	GetLogic() *SrvLogic
	GetImage() string
	GetType() SrvType
	GetStatus() SrvStatus
	SetStatus(status SrvStatus)
	//ui logic

	ListImages() []string
}

type SrvLogic struct {
	IsDefault     bool              // default or manual config
	Type          SrvType           // type of service
	Image         string            // docker Dim
	Status        SrvStatus         // service status
	ContainerName string            // docker container name
	Network       string            // docker network
	Ports         []string          // docker ports
	Volumes       []string          // docker volume
	EnvVars       map[string]string // docker env variables
	RestartPolicy string
}

func InitServices(configs ...*SrvConfig) map[SrvType]Service {
	srvs := make(map[SrvType]Service)

	for i, config := range configs {
		if config == nil {
			log.Printf("Config #%d is nil, skipping", i)
			continue
		}

		service, err := createService(config)
		if err != nil {
			log.Printf("Failed to initialize service for type %s: %v", config.Type, err)
			continue
		}

		srvs[config.Type] = service
	}

	log.Printf("Initialized %d services", len(srvs))
	return srvs
}

func createService(config *SrvConfig) (Service, error) {
	if config.IsDefault {
		switch config.Type {
		case DB:
			return DefaultSrvLogicDB(), nil
		case LB:
			return DefaultSrvLogicLB(), nil
		case Cache:
			return DefaultSrvLogicCache(), nil
		case ServerMain:
			return DefaultSrvLogicServerMain(), nil
		case ServerAdd:
			return DefaultSrvLogicServerAdd(), nil
		default:
			return nil, fmt.Errorf("unsupported service type: %s", config.Type)
		}
	}

	switch config.Type {
	case DB:
		return NewSrvLogicDB(config), nil
	case LB:
		return NewSrvLogicLB(config), nil
	case Cache:
		return NewSrvLogicCache(config), nil
	case ServerMain:
		return NewSrvLogicServerMain(config), nil
	case ServerAdd:
		return NewSrvLogicServerAdd(config), nil
	default:
		return nil, fmt.Errorf("unsupported service type: %s", config.Type)
	}
}
