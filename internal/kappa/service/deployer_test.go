package service

import (
	"github.com/KirillMironov/kappa/internal/kappa/domain"
	"github.com/KirillMironov/kappa/pkg/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeployer_Deploy(t *testing.T) {
	var deployer = NewDeployer(mock.Logger{})

	testCases := map[string]struct {
		deployment  domain.Deployment
		expectError bool
	}{
		"success": {
			deployment: domain.Deployment{
				Services: []domain.Service{
					{Name: "echo", Command: "echo $FOO"},
				},
			},
			expectError: false,
		},
		"no command": {
			deployment: domain.Deployment{
				Services: []domain.Service{
					{Name: "empty", Command: ""},
				},
			},
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := deployer.Deploy(tc.deployment)
			assert.True(t, tc.expectError == (err != nil))
		})
	}
}

func TestDeployer_Cancel(t *testing.T) {
	var deployer = NewDeployer(mock.Logger{})

	var deployment = domain.Deployment{
		Services: []domain.Service{
			{Name: "sleep", Command: "sleep 10"},
		},
	}

	testCases := map[string]struct {
		deployment  domain.Deployment
		preparation func(domain.Deployment)
		expectError bool
	}{
		"success": {
			deployment: deployment,
			preparation: func(deployment domain.Deployment) {
				_ = deployer.Deploy(deployment)
			},
			expectError: false,
		},
		"deployment not found": {
			deployment:  deployment,
			expectError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.preparation != nil {
				tc.preparation(tc.deployment)
			}
			err := deployer.Cancel(tc.deployment)
			assert.True(t, tc.expectError == (err != nil))
		})
	}
}
