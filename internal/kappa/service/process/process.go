package process

type Process struct {
	Command string
	Args    []string
	pid     int
}

func New(command string, args []string) *Process {
	return &Process{
		Command: command,
		Args:    args,
	}
}

func (p *Process) Getpid() int {
	return p.pid
}
