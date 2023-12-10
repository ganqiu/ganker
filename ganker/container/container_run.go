package container

import (
	"fmt"
	"go_docker_learning/ganker/cgroup"
	"go_docker_learning/ganker/cgroup/subsystem"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func RunContainer(tty bool, comArray []string, imageName string, volume string, resourceConfig *subsystem.ResourceConfig, containerName string, env []string) {
	id := generateContainerId(15)
	if containerName == "" {
		containerName = imageName + "-" + id[:10]
	}

	parent, writePipe, containerDir, containerId := initNewParentProcess(tty, id, imageName, volume, env)
	if parent == nil {
		logrus.Errorf("fail to init new parent process")
		return
	}

	if err := parent.Start(); err != nil {
		logrus.Error(err)
		return
	}

	if err := recordContainerInfo(strconv.Itoa(parent.Process.Pid), id, containerName, imageName, volume, comArray); err != nil {
		logrus.Error("fail to create container info: %v", err)
		return
	}

	// Initialize cGroup manager
	cGroupManager := cgroup.NewCgroupManager("GankerCgroup" + "/" + containerId)
	// set resource limitation
	if err := cGroupManager.Set(resourceConfig); err != nil {
		logrus.Errorf("%v", err)
		os.Exit(-1)
	}
	// add the process into cGroup
	if err := cGroupManager.Apply(parent.Process.Pid); err != nil {
		logrus.Errorf("%v", err)
		os.Exit(-1)
	}
	// send command to child process
	if err := sendInitCommand(comArray, writePipe); err != nil {
		logrus.Errorf("%v", err)
		os.Exit(-1)
	}
	// if tty,it means that the container is running in foreground
	if tty {
		// wait for child process to exit
		if err := parent.Wait(); err != nil {
			logrus.Errorf("%v", err)
			os.Exit(-1)
		}

		err := cGroupManager.Delete()
		if err != nil {
			logrus.Errorf("%v", err)
			os.Exit(-1)
		}

		deleteWorkSpace(containerDir, volume)
		quitContainer(containerId)

	} else {
		fmt.Println("containerId: ", containerId)
	}
}
