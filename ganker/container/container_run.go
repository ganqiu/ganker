package container

import (
	"go_docker_learning/ganker/cgroup"
	"go_docker_learning/ganker/cgroup/subsystem"
	"os"

	"github.com/sirupsen/logrus"
)

func RunContainer(tty bool, comArray []string, resourceConfig *subsystem.ResourceConfig) {
	imageName := "busybox"
	parent, writePipe, containerDir, id := initNewParentProcess(tty, imageName)

	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}

	// Initialize cgroup manager
	cgroupManager := cgroup.NewCgroupManager("GankerCgroup" + "/" + id)
	defer cgroupManager.Delete()
	defer deleteWorkSpace(containerDir)
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

}
