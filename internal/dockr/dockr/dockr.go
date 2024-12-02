package dockr

import (
	entity "Infra/internal/dockr/container"
	"Infra/internal/dockr/service"
	"context"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

type Dockr struct {
	cli        *client.Client
	ctx        context.Context
	config     map[service.SrvType][]service.Service
	containers map[service.SrvType][]*entity.Container
	logger     *zap.SugaredLogger
}
