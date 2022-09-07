package service

import (
	"github.com/KirillMironov/kappa/internal/kappa/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	var parser Parser
	var yaml = `
name: test
command: ["/bin/sh", "-c"]
env:
  - name: FOO
    value: BAR
workingDir: "/"
`

	pod, err := parser.Parse([]byte(yaml))
	assert.NoError(t, err)
	assert.Equal(t, core.Pod{
		Name:    "test",
		Command: []string{"/bin/sh", "-c"},
		Environment: []core.Environment{
			{Name: "FOO", Value: "BAR"},
		},
		WorkingDir: "/",
	}, pod)
}
