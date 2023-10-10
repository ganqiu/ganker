package container

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

func InitNewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	// if a tty is needed, it will set for user input and output
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}

func RunContainer(tty bool, command string) {
	parent := InitNewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}
