package service

import (
	"context"
	"errors"
	"github.com/KirillMironov/kappa/internal/kappa/domain"
	"github.com/KirillMironov/kappa/pkg/logger"
	"os/exec"
	"reflect"
	"strings"
	"sync"
	"syscall"
)

type Deployer struct {
	inProgress map[reflect.Type]context.CancelFunc
	mu         sync.Mutex
	logger     logger.Logger
}

func NewDeployer(logger logger.Logger) *Deployer {
	return &Deployer{
		inProgress: make(map[reflect.Type]context.CancelFunc),
		logger:     logger,
	}
}

func (d *Deployer) Deploy(deployment domain.Deployment) error {
	ctx, cancel := context.WithCancel(context.Background())
	d.mu.Lock()
	d.inProgress[reflect.TypeOf(deployment)] = cancel
	d.mu.Unlock()

	for _, service := range deployment.Services {
		pid, err := d.startProcess(ctx, service)
		if err != nil {
			_ = d.Cancel(deployment)
			d.logger.Errorf("failed to deploy service %q: %v", service.Name, err)
			return err
		}
		d.logger.Infof("service %q started with pid %d", service.Name, pid)
	}

	return nil
}

func (d *Deployer) Cancel(deployment domain.Deployment) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	cancel, ok := d.inProgress[reflect.TypeOf(deployment)]
	if !ok {
		return errors.New("deployment not found")
	}

	cancel()
	delete(d.inProgress, reflect.TypeOf(deployment))

	return nil
}

func (d *Deployer) startProcess(ctx context.Context, service domain.Service) (pid int, _ error) {
	splitted := strings.Split(service.Command, " ")
	if len(splitted) < 1 {
		return 0, errors.New("no command specified")
	}

	command := splitted[0]

	var args []string
	if len(splitted) > 1 {
		args = splitted[1:]
	}

	cmd := exec.CommandContext(ctx, command, args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:   true,
		Pdeathsig: syscall.SIGKILL,
	}

	err := cmd.Start()
	if err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}
