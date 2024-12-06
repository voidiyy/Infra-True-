package dockr

import (
	"Infra/internal/dockr/config"
	entity "Infra/internal/dockr/container"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types/image"
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

func NewDockr(ctx context.Context, logger *zap.SugaredLogger) (*Dockr,error) {
	if logger == nil {
		logger = zap.NewNop().Sugar()
	}
	
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("error init Docker API client: %s", err)
	}
	
	ping, err := cli.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ping Docker API client: %s", err) 
	}
	
	log.Printf("\nAPI client initialized with version: %s\nOS version: %s", ping.APIVersion, ping.OSType)
	
	return &Dockr{
		cli: cli,
		ctx: ctx,
		config: &config.UltimateConfig{},
		containers: &entity.UltimateContainer{},
		logger: logger,
	}, nil
}

func (d *Dockr) InitContainers(configs *config.UltimateConfig) error {
	if configs == nil {
		return errors.New("ultimate config id nil")
	}
	
	ultiContainers, err := entity.NewUltimateContainer(configs)
	if err != nil {
		return fmt.Errorf("error create ultimate containers %s", err)
	}
	
	var failedImages []string
	for _, v := range ultiContainers.Containers {
    out, err := d.cli.ImagePull(d.ctx, v.GetContainerConfig().GetImage(), image.PullOptions{})
    if err != nil {
        failedImages = append(failedImages, v.GetContainerConfig().GetImage())
        d.logger.Errorf("error pulling image %s: %v", v.GetContainerConfig().GetImage(), err)
        continue
    }
    io.Copy(os.Stdout, out)
    out.Close()
	}
	if len(failedImages) > 0 {
    return fmt.Errorf("failed to pull images: %v", failedImages)
	}
	
	for k, v := range ultiContainers.Containers {
		d.logger.Infof("image %s successfully pulled for %s ", v.GetContainerConfig().GetImage(), k)
	} 
	return nil
}

// Close closes the docker client session
func (d *Dockr) Close() {
	err := d.cli.Close()
	if err != nil {
		return
	}
}
