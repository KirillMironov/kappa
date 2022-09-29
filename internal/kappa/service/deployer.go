package service

import (
	"fmt"
	"github.com/KirillMironov/kappa/internal/kappa/domain"
	"github.com/KirillMironov/kappa/pkg/logger"
	"github.com/KirillMironov/kappa/pkg/process"
	"strings"
	"sync"
)

type Deployer struct {
	inProgress map[string]*process.Process
	mu         sync.Mutex
	logger     logger.Logger
}

func NewDeployer(logger logger.Logger) *Deployer {
	return &Deployer{
		inProgress: make(map[string]*process.Process),
		logger:     logger,
	}
}

func (d *Deployer) Deploy(deployment domain.Deployment) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	for name := range d.inProgress {
		for _, service := range deployment.Services {
			if service.Name == name {
				return fmt.Errorf("service %q is already deployed", service.Name)
			}
		}
	}

	for _, service := range deployment.Services {
		command, args := d.splitCommand(&service)

		proc := process.New(command, args...)

		err := proc.Start()
		if err != nil {
			d.logger.Errorf("failed to deploy service %q: %v", service.Name, err)
			d.Cancel(deployment)
			return err
		}

		d.inProgress[service.Name] = proc

		d.logger.Infof("service %q started with pid %d", service.Name, proc.Getpid())
	}

	return nil
}

func (d *Deployer) Cancel(deployment domain.Deployment) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, service := range deployment.Services {
		proc, ok := d.inProgress[service.Name]
		if !ok {
			continue
		}

		proc.Terminate()

		delete(d.inProgress, service.Name)
	}
}

func (d *Deployer) splitCommand(service *domain.Service) (command string, args []string) {
	splitted := strings.Split(service.Command, " ")
	if len(splitted) < 1 {
		return "", nil
	}

	command = splitted[0]

	if len(splitted) > 1 {
		args = splitted[1:]
	}

	return command, args
}
