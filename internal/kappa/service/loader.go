package service

import (
	"context"
	"github.com/KirillMironov/kappa/internal/kappa/domain"
	"github.com/KirillMironov/kappa/pkg/logger"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Loader struct {
	pods         chan<- []domain.Pod
	podsDir      string
	loadInterval time.Duration
	parser       parser
	logger       logger.Logger
}

type parser interface {
	Parse(data []byte) (domain.Pod, error)
}

func NewLoader(pods chan<- []domain.Pod, podsDir string, loadInterval time.Duration, parser parser,
	logger logger.Logger) *Loader {
	return &Loader{
		pods:         pods,
		podsDir:      podsDir,
		loadInterval: loadInterval,
		parser:       parser,
		logger:       logger,
	}
}

func (l Loader) Start(ctx context.Context) {
	timer := time.NewTicker(l.loadInterval)

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

func (l Loader) load() (pods []domain.Pod, err error) {
	files, err := os.ReadDir(l.podsDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yaml") && !file.IsDir() {
			path := filepath.Join(l.podsDir, file.Name())

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
