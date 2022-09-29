//go:build linux

package process

import (
	"os/exec"
	"syscall"
)

func (p *Process) start() error {
	cmd := exec.Command(p.Command, p.Args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:   true,
		Pdeathsig: syscall.SIGKILL,
	}

	err := cmd.Start()
	if err != nil {
		return err
	}

	p.process = cmd.Process

	return nil
}

func (p *Process) terminate() {
	_ = p.process.Kill()
	_, _ = p.process.Wait()
}
