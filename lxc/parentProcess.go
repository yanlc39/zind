package lxc

import (
	"os"
	"os/exec"
	"syscall"
)

func CreateParentProcess(interactive bool, tty bool, args []string) *exec.Cmd {
	containerCommand := exec.Command("/proc/self/exe", "child", args[0])

	containerCommand.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	if tty {
		if interactive {
			containerCommand.Stdin = os.Stdin
		}
		containerCommand.Stdout = os.Stdout
		containerCommand.Stderr = os.Stderr
	}

	return containerCommand
}
