package container

import (
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

// NewParentProcess is used to create a parent process
func InitRunContainerProcess(command string, args []string) error {

	logrus.Infof("Run Command %s", command)

	// set /proc as a private mount namespace's mount point
	if err := syscall.Mount("", "/proc", "proc", syscall.MS_REC|syscall.MS_PRIVATE, ""); err != nil {
		logrus.WithField("method", "syscall.Mount").Error(err)
		return err
	}

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// mount /proc to /proc
	if err := syscall.Mount("", "/proc", "proc", uintptr(defaultMountFlags), ""); err != nil {
		logrus.Errorf("mount proc error %v", err)
	}

	argv := []string{command}

	// syscall.Exec will replace the current process with the command
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
		return err
	}
	return nil
}
