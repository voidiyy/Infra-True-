package main

import (
	"Infra/internal/dockr/config"
	"Infra/internal/dockr/dockr"
	"context"
	"log"
)

func main() {
	
	doc, err := dockr.NewDockr(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	
	conf, err := config.NewContainersConfig(
		config.ApiGatewayConfig,
		config.HaproxyConfig,
		config.MongoConfig,
	)
	
	if err != nil {
		log.Fatal(err)
	}
	
	err = doc.InitContainers(conf)
	if err != nil {
		log.Fatal(err)
	}
}


