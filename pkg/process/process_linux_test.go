//go:build linux

package process

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestProcess(t *testing.T) {
	process := New("sleep", "5")

	err := process.Start()
	require.NoError(t, err)

	process.Terminate()

	_, err = os.Stat(fmt.Sprintf("/proc/%d", process.Getpid()))
	assert.True(t, os.IsNotExist(err))
}
