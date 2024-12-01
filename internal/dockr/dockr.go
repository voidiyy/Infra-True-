package dockr

import (
	entity "Infra/internal/dockr/container"
	"Infra/internal/dockr/service"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"go.uber.org/zap"
	"io"
)

type Dockr struct {
	cli        *client.Client
	ctx        context.Context
	config     map[service.SrvType]service.Service
	containers map[service.SrvType]*entity.Container
	logger     *zap.SugaredLogger
}

// NewDockr створює новий інстанс Dockr
func NewDockr(conf map[service.SrvType]service.Service) (*Dockr, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	if _, err = cli.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping Docker API: %w", err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	dockrInstance := &Dockr{
		cli:        cli,
		ctx:        context.Background(),
		config:     conf,
		containers: make(map[service.SrvType]*entity.Container),
		logger:     sugar,
	}

	// Створення контейнерів після ініціалізації Dockr
	if err = dockrInstance.createContainers(); err != nil {
		return nil, fmt.Errorf("failed to create containers: %w", err)
	}

	return dockrInstance, nil
}

// createContainer створює контейнер для кожної служби
func (d *Dockr) createContainer(srv service.Service) (*entity.Container, error) {
	// Перевірка наявності образу
	if err := d.imageExist(srv.GetImage()); err != nil {
		d.logger.Errorf("Image check failed for %s: %v", srv.GetImage(), err)
		return nil, err
	}

	// Формування конфігурації контейнера
	confC, confH, confN := d.configContainer(srv)

	// Створення контейнера
	resp, err := d.cli.ContainerCreate(d.ctx, confC, confH, confN, nil, srv.GetLogic().ContainerName)
	if err != nil {
		d.logger.Errorf("Failed to create container for %s: %v", srv.GetLogic().ContainerName, err)
		return nil, err
	}

	// Створення контейнера і додавання його в мапу
	container := &entity.Container{
		ID:            resp.ID,
		Service:       srv,
		Config:        confC,
		HostConfig:    confH,
		NetworkConfig: confN,
		Status:        entity.ContainerStatusStopped(),
	}

	d.containers[srv.GetType()] = container
	d.logger.Infof("Container %s created successfully", srv.GetLogic().ContainerName)
	return container, nil
}

// createContainers створює контейнери для всіх служб
func (d *Dockr) createContainers() error {
	for srvType, srv := range d.config {
		_, err := d.createContainer(srv)
		if err != nil {
			return fmt.Errorf("failed to create container for service %s: %w", srvType, err)
		}
		d.logger.Infof("Created container for service: %s", srvType)
	}
	return nil
}

// StartContainer запускає контейнер для певної служби
func (d *Dockr) StartContainer(srvType service.SrvType) error {
	cont, exists := d.containers[srvType]
	if !exists {
		return fmt.Errorf("container for service %s not found", srvType)
	}

	if err := d.cli.ContainerStart(d.ctx, cont.ID, container.StartOptions{}); err != nil {
		d.logger.Errorf("Failed to start container %s: %v", cont.ID, err)
		return err
	}

	cont.Status = entity.ContainerStatusRunning()
	d.logger.Infof("Container %s started successfully", cont.ID)
	return nil
}

// StopContainer зупиняє контейнер для певної служби
func (d *Dockr) StopContainer(srvType service.SrvType) error {
	cont, exists := d.containers[srvType]
	if !exists {
		return fmt.Errorf("container for service %s not found", srvType)
	}

	if err := d.cli.ContainerStop(d.ctx, cont.ID, container.StopOptions{}); err != nil {
		d.logger.Errorf("Failed to stop container %s: %v", cont.ID, err)
		return err
	}

	cont.Status = entity.ContainerStatusStopped()
	d.logger.Infof("Container %s stopped successfully", cont.ID)
	return nil
}

// ListContainers виводить список контейнерів
func (d *Dockr) ListContainers() {
	fmt.Println("len of container: ", len(d.config))
	fmt.Println("len of container: ", len(d.containers))

	for srvType, container := range d.containers {
		d.logger.Infof("Service: %s, Container ID: %s, Status: %s", srvType, container.ID, container.Status)
	}
}

// CheckContainerHealth перевіряє здоров'я контейнера
func (d *Dockr) CheckContainerHealth(srvType service.SrvType) (string, error) {
	container, exists := d.containers[srvType]
	if !exists {
		return "", fmt.Errorf("container for service %s not found", srvType)
	}

	inspect, err := d.cli.ContainerInspect(d.ctx, container.ID)
	if err != nil {
		d.logger.Errorf("Failed to inspect container %s: %v", container.ID, err)
		return "", err
	}

	if inspect.State.Health != nil {
		return inspect.State.Health.Status, nil
	}

	return "unknown", nil
}

// CheckAllContainersHealth перевіряє здоров'я всіх контейнерів
func (d *Dockr) CheckAllContainersHealth() {
	for srvType, cont := range d.containers {
		status, err := d.CheckContainerHealth(srvType)
		if err != nil {
			d.logger.Errorf("Error checking health for container %s: %v", cont.ID, err)
		} else {
			d.logger.Infof("Container %s health status: %s", cont.ID, status)
		}
	}
}

// configContainer повертає конфігурацію контейнера
func (d *Dockr) configContainer(s service.Service) (*container.Config, *container.HostConfig, *network.NetworkingConfig) {
	conf := s.GetLogic()

	return &container.Config{
			Image: conf.Image,
			Env:   entity.ConvertEnvVars(conf.EnvVars),
		}, &container.HostConfig{
			PortBindings: nat.PortMap{
				"5432/tcp": []nat.PortBinding{{HostPort: "5432"}},
			},
			Binds:         conf.Volumes,
			RestartPolicy: container.RestartPolicy{Name: container.RestartPolicyAlways},
		}, &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				conf.Network: {},
			},
		}
}

// imageExist перевіряє, чи існує образ Docker і завантажує його, якщо потрібно
func (d *Dockr) imageExist(img string) error {
	_, _, err := d.cli.ImageInspectWithRaw(d.ctx, img)
	if err != nil {
		d.logger.Warnf("Image %s not found, pulling...", img)
		reader, pullErr := d.cli.ImagePull(d.ctx, img, image.PullOptions{})
		if pullErr != nil {
			return fmt.Errorf("failed to pull image %s: %w", img, pullErr)
		}
		defer reader.Close()
		_, _ = io.Copy(io.Discard, reader)
		d.logger.Infof("Image %s pulled successfully", img)
	}
	return nil
}
