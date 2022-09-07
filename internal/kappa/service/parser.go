package service

import (
	"github.com/KirillMironov/kappa/internal/kappa/core"
	"gopkg.in/yaml.v3"
)

type Parser struct{}

func (Parser) Parse(data []byte) (pod core.Pod, _ error) {
	return pod, yaml.Unmarshal(data, &pod)
}
