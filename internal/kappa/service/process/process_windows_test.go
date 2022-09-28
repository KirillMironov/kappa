//go:build windows

package process

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcess(t *testing.T) {
	process := New("notepad.exe")

	err := process.Start()
	require.NoError(t, err)

	process.Terminate()

	processState, _ := process.process.Wait()
	assert.True(t, processState.Exited())
}
