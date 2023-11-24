package container

import (
	"fmt"
	"os"
	"os/exec"
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

	// set Mount
	if err := MountSet(); err != nil {
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

// initNewParentProcess it will clone a new process as parent process and calling /proc/self/exe with init as the first argument
func initNewParentProcess(tty bool, imageName, volume string) (*exec.Cmd, *os.File, string, string) {
	// create a pipe,it will be used to send command to child process,in another word, it can be used to send command to init process
	readPipe, writePipe, err := newPipe()
	if err != nil {
		logrus.Errorf("New pipe create error %v", err)
		return nil, nil, "", ""
	}

	// excute itself by calling /proc/self/exe with init as the first argument
	cmd := exec.Command("/proc/self/exe", "init")

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
	containerDir, containerId := newWorkSpace(imageName, volume)
	// ExtraFiles specifies additional open files to be inherited by the new process,
	// it will deliver the pipe file to child process
	// as file descriptor 0,1,2 are used for stdin,stdout,stderr,
	// we use 3 to describe the first pipe file,and 4 to the second pipe file, and so on
	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Dir = containerDir + MergeLayerName
	return cmd, writePipe, containerDir, containerId
}
