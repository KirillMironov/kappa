package service

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestLoader_load(t *testing.T) {
	var (
		tempDir        = t.TempDir()
		numDeployments = 10
		loader         = NewLoader(nil, tempDir, time.Second, Parser{}, logrus.New())
	)

	for i := 0; i < numDeployments; i++ {
		file, err := os.CreateTemp(tempDir, "deployment-*.yaml")
		require.NoError(t, err)
		file.Close()

		file, err = os.CreateTemp(tempDir, "file-*.txt")
		require.NoError(t, err)
		file.Close()

		_, err = os.MkdirTemp(tempDir, "dir-*")
		require.NoError(t, err)
	}

	deployments, err := loader.load()
	assert.NoError(t, err)
	assert.Len(t, deployments, numDeployments)
}
