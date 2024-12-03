package dockr

import (
	"Infra/internal/dockr/config"
	entity "Infra/internal/dockr/container"
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

type Dockr struct {
	cli        *client.Client
	ctx        context.Context
	config     *config.UltimateConfig
	containers *entity.UltimateContainer
	logger     *zap.SugaredLogger
}

func NewDockr(ctx context.Context, ultimateConfig *config.UltimateConfig, logger *zap.SugaredLogger) (*Dockr, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	if logger == nil {
		logger = zap.NewNop().Sugar()
	}

	return &Dockr{
		cli:        cli,
		ctx:        ctx,
		config:     ultimateConfig,
		containers: entity.NewUltimateContainer(),
		logger:     logger,
	}, nil
}

func (d *Dockr) InitContainers() {
	
}

// Close closes the docker client session
func (d *Dockr) Close() {
	err := d.cli.Close()
	if err != nil {
		return
	}
}
