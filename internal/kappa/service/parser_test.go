package service

import (
	"github.com/KirillMironov/kappa/internal/kappa/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	var parser Parser
	var yaml = `
name: example
command: go
args:
  - version
env:
  - name: FOO
    value: BAR
workingDir: "/"
`

	pod, err := parser.Parse([]byte(yaml))
	assert.NoError(t, err)
	assert.Equal(t, domain.Pod{
		Name:    "example",
		Command: "go",
		Args:    []string{"version"},
		Environment: []domain.Environment{
			{Name: "FOO", Value: "BAR"},
		},
		WorkingDir: "/",
	}, pod)
}
