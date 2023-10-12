package container

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

// it will clone a new process as parent process and calling /proc/self/exe with init as the first argument
func InitNewParentProcess(tty bool, command string) *exec.Cmd {

	args := []string{"init", command}

	// excute itself by calling /proc/self/exe with init as the first argument
	cmd := exec.Command("/proc/self/exe", args...)

	// create a new namespace as parent process
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
