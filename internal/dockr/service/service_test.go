package service_test

import (
	"Infra/internal/dockr"
	"Infra/internal/dockr/service"
	"testing"
)

func TestServiceGlobal(t *testing.T) {

	t.Parallel()

	t.Run("DB", func(t *testing.T) {
		t.Parallel()
		srv := service.NewSrvLogicDB(configs[0])

		srv.Info()
	})

	t.Run("Init", func(t *testing.T) {
		srv := service.InitServices(configs...)

		_, err := dockr.NewDockr(srv)
		if err != nil {
			t.Error(err)
		}

	})
}

var configs = []*service.SrvConfig{
	&service.SrvConfig{IsDefault: true, Type: service.DB},
	&service.SrvConfig{IsDefault: true, Type: service.LB},
	&service.SrvConfig{IsDefault: true, Type: service.Cache},
	&service.SrvConfig{IsDefault: true, Type: service.ServerMain},
	&service.SrvConfig{IsDefault: true, Type: service.ServerAdd},
}
