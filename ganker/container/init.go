package container

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

// NewParentProcess is used to create a parent process
func InitRunContainerProcess() error {

	// get the command array from pipe
	cmdArray := readParentProcess()
	if cmdArray == nil {
		return fmt.Errorf("get parent command error, cmdArray is nil")
	}

	if len(cmdArray) != 0 {
		logrus.Infof("Run Command %v", cmdArray)
	}

	// set /proc as a private mount namespace's mount point
	if err := syscall.Mount("", "/proc", "proc", syscall.MS_REC|syscall.MS_PRIVATE, ""); err != nil {
		logrus.WithField("method", "syscall.Mount").Error(err)
		return err
	}

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

	// mount /proc to /proc
	if err := syscall.Mount("", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		logrus.Errorf("mount proc error %v", err)
		return err
	}

	//LookPath searches for an executable named file in the directories named by the PATH environment variable.
	// so it will return the absolute path of the command
	// for example, if the command is "ls", it will return "/bin/ls"
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		logrus.Errorf("exec loop path error %v", err)
		return err
	}

	// syscall.Exec will replace the current process with the command
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		logrus.Errorf(err.Error())
		return err
	}
	return nil
}

func readParentProcess() []string {

	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := io.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
