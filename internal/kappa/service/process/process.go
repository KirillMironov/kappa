package process

import "os"

type Process struct {
	Command string
	Args    []string
	process *os.Process
}

func New(command string, args ...string) *Process {
	return &Process{
		Command: command,
		Args:    args,
	}
}

func (p *Process) Getpid() int {
	return p.process.Pid
}
