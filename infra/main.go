package main

import (
	"Infra/internal/dockr"
	"Infra/internal/dockr/service"
	"log"
)

var configs = []*service.SrvConfig{
	&service.SrvConfig{IsDefault: true, Type: service.DB},
	&service.SrvConfig{IsDefault: true, Type: service.LB},
	&service.SrvConfig{IsDefault: true, Type: service.Cache},
	&service.SrvConfig{IsDefault: true, Type: service.ServerMain},
	&service.SrvConfig{IsDefault: true, Type: service.ServerAdd},
}

func main() {
	// Initialize services
	services := service.InitServices(configs...)

	dck, err := dockr.NewDockr(services)
	if err != nil {
		log.Fatal(err)
		return
	}

	dck.ListContainers()
}
