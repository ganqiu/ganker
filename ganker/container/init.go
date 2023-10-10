package container

import (
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

// NewParentProcess is used to create a parent process
func InitRunContainerProcess(command string, args []string) error {
	logrus.Infof("command %s", command)
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
		return err
	}
	return nil

}
