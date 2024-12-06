package config_test

import (
	"Infra/internal/dockr/config"
	"fmt"
	"testing"
)

func TestDefault(t *testing.T) {
	
	t.Run("Default", func(t *testing.T){
		
		ulti,err := config.NewContainersConfig(config.MongoConfig, config.ApiGatewayConfig, config.HaproxyConfig, config.MonitoringConfig)
		if err != nil {
			t.Error(err)
		}
		
		if ulti == nil {
			t.Errorf("ulti is NIL")
		}
		
		fmt.Println("LEN: ", len(ulti.Containers))
		
		for _, v := range ulti.Containers {
			fmt.Println("SERVICE: ",v.GetService())
			fmt.Println(v)
		}
	})
}
