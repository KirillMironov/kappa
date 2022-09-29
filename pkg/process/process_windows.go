//go:build windows

package process

import (
	"golang.org/x/sys/windows"
	"os/exec"
	"strconv"
	"syscall"
	"unsafe"
)

func (p *Process) Start() error {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return err
	}

	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}

	_, err = windows.SetInformationJobObject(
		job,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info)),
	)
	if err != nil {
		return err
	}

	cmd := exec.Command(p.Command, p.Args...)

	err = cmd.Start()
	if err != nil {
		return err
	}

	p.process = cmd.Process

	type process struct {
		pid    int
		handle syscall.Handle
	}

	handle := (*process)(unsafe.Pointer(cmd.Process)).handle

	return windows.AssignProcessToJobObject(job, windows.Handle(handle))
}

func (p *Process) Terminate() {
	taskkill := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(p.process.Pid))

	_ = taskkill.Run()
	_ = p.process.Kill()
}
