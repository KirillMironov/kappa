package service

import (
	"context"
	"github.com/KirillMironov/kappa/internal/kappa/core"
	"github.com/KirillMironov/kappa/pkg/log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Loader struct {
	pods         chan<- []core.Pod
	PodsDir      string
	LoadInterval time.Duration
	parser       parser
	logger       log.Logger
}

type parser interface {
	Parse(data []byte) (core.Pod, error)
}

func NewLoader(pods chan<- []core.Pod, podsDir string, loadInterval time.Duration, parser parser,
	logger log.Logger) *Loader {
	return &Loader{
		pods:         pods,
		PodsDir:      podsDir,
		LoadInterval: loadInterval,
		parser:       parser,
		logger:       logger,
	}
}

func (l Loader) Start(ctx context.Context) {
	timer := time.NewTicker(l.LoadInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			pods, err := l.load()
			if err != nil {
				l.logger.Errorf("failed to read deployments dir: %v", err)
				continue
			}

			l.pods <- pods
		}
	}
}

func (l Loader) load() (pods []core.Pod, err error) {
	files, err := os.ReadDir(l.PodsDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yaml") && !file.IsDir() {
			path := filepath.Join(l.PodsDir, file.Name())

			data, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}

			pod, err := l.parser.Parse(data)
			if err != nil {
				return nil, err
			}

			pods = append(pods, pod)
		}
	}

	return pods, nil
}
