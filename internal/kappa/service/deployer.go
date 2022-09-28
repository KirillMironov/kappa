package service

import (
	"github.com/KirillMironov/kappa/internal/kappa/domain"
	"github.com/KirillMironov/kappa/internal/kappa/service/process"
	"github.com/KirillMironov/kappa/pkg/logger"
	"reflect"
	"strings"
	"sync"
)

type Deployer struct {
	inProgress map[reflect.Type]*process.Process
	mu         sync.Mutex
	logger     logger.Logger
}

func NewDeployer(logger logger.Logger) *Deployer {
	return &Deployer{
		inProgress: make(map[reflect.Type]*process.Process),
		logger:     logger,
	}
}

func (d *Deployer) Deploy(deployment domain.Deployment) error {
	for _, service := range deployment.Services {
		command, args := d.splitCommand(&service)

		proc := process.New(command, args...)

		err := proc.Start()
		if err != nil {
			d.logger.Errorf("failed to deploy service %q: %v", service.Name, err)
			d.Cancel(deployment)
			return err
		}

		d.mu.Lock()
		d.inProgress[reflect.TypeOf(service)] = proc
		d.mu.Unlock()

		d.logger.Infof("service %q started with pid %d", service.Name, proc.Getpid())
	}

	return nil
}

func (d *Deployer) Cancel(deployment domain.Deployment) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, service := range deployment.Services {
		proc, ok := d.inProgress[reflect.TypeOf(service)]
		if !ok {
			continue
		}

		err := proc.Terminate()
		if err != nil {
			d.logger.Errorf("failed to terminate service %q: %v", service.Name, err)
			continue
		}

		delete(d.inProgress, reflect.TypeOf(service))
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
