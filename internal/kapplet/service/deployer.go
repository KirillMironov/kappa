package service

import (
	"context"
	"github.com/KirillMironov/kappa/internal/kappa/domain"
	"github.com/KirillMironov/kappa/pkg/logger"
	"os/exec"
	"reflect"
	"sync"
)

type Deployer struct {
	podsCh     <-chan []domain.Pod
	inProgress map[string]state
	logger     logger.Logger
	mu         sync.Mutex
}

type state struct {
	pod    domain.Pod
	cancel context.CancelFunc
}

func NewDeployer(podsCh <-chan []domain.Pod, logger logger.Logger) *Deployer {
	return &Deployer{
		podsCh:     podsCh,
		inProgress: make(map[string]state),
		logger:     logger,
	}
}

func (d *Deployer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case pods := <-d.podsCh:
			for _, pod := range pods {
				func() {
					d.mu.Lock()
					defer d.mu.Unlock()

					deployed, ok := d.inProgress[pod.Name]
					if ok {
						if reflect.DeepEqual(pod, deployed.pod) {
							return
						}
						deployed.cancel()
					}

					podCtx, cancel := context.WithCancel(ctx)
					d.inProgress[pod.Name] = state{pod: pod, cancel: cancel}

					d.deploy(podCtx, pod)
				}()
			}
		}
	}
}

func (d *Deployer) deploy(ctx context.Context, pod domain.Pod) {
	go func() {
		defer func() {
			d.mu.Lock()
			d.inProgress[pod.Name].cancel()
			delete(d.inProgress, pod.Name)
			d.mu.Unlock()
		}()

		cmd := exec.CommandContext(ctx, pod.Command, pod.Args...)

		err := cmd.Start()
		if err != nil {
			d.logger.Errorf("failed to start pod %s: %v", pod.Name, err)
			return
		}

		err = cmd.Wait()
		if err != nil {
			d.logger.Errorf("failed to execute command: %v", err)
		}
	}()
}
