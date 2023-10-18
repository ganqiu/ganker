package container

import (
	"fmt"
	"go_docker_learning/ganker/cgroup"
	"go_docker_learning/ganker/cgroup/subsystem"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

// InitNewParentProcess it will clone a new process as parent process and calling /proc/self/exe with init as the first argument
func InitNewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	// create a pipe,it will be used to send command to child process,in another word, it can be used to send command to init process
	readPipe, writePipe, err := newPipe()
	if err != nil {
		logrus.Errorf("New pipe create error %v", err)
		return nil, nil
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

	// ExtraFiles specifies additional open files to be inherited by the new process,
	// it will deliver the pipe file to child process
	// as file descriptor 0,1,2 are used for stdin,stdout,stderr,
	// we use 3 to describe the first pipe file,and 4 to the second pipe file, and so on
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}

func RunContainer(tty bool, comArray []string, resourceConfig *subsystem.ResourceConfig) {
	parent, writePipe := InitNewParentProcess(tty)

	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}

	// Initialize cgroup manager
	cgroupManager := cgroup.NewCgroupManager("GankerCgroup")

	defer cgroupManager.Delete()

	// set resource limitation
	if err := cgroupManager.Set(resourceConfig); err != nil {
		logrus.Errorf("%v", err)
		os.Exit(-1)
	}
	// add the process into cgroup
	if err := cgroupManager.Apply(parent.Process.Pid); err != nil {
		logrus.Errorf("%v", err)
		os.Exit(-1)
	}
	// send command to child process
	if err := sendInitCommand(comArray, writePipe); err != nil {
		logrus.Errorf("%v", err)
		os.Exit(-1)
	}
	// wait for child process to exit
	if err := parent.Wait(); err != nil {
		logrus.Errorf("%v", err)
		os.Exit(-1)
	}
	os.Exit(-1)
}

func newPipe() (*os.File, *os.File, error) {
	// create a pipe
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return readPipe, writePipe, nil
}

func sendInitCommand(comArray []string, writePipe *os.File) error {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)
	if _, err := writePipe.WriteString(command); err != nil {
		return fmt.Errorf("write pipe error %v", err)
	}
	writePipe.Close()
	return nil
}
