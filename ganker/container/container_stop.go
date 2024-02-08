package container

import (
	"go_docker_learning/ganker/cgroup"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	syscall "golang.org/x/sys/unix"
)

func StopContainer(containerId string) {
	containerInfo, err := getContainerInfo(containerId)
	if err != nil {
		log.Errorf("Get container %s info error:  %v", containerId, err)
		return
	}

	if containerInfo.Status != RUNNING {
		log.Errorf("Container %s is not running now", containerId)
		return
	}

	pid, err1 := strconv.Atoi(containerInfo.Pid)
	if err1 != nil {
		log.Errorf("Conver pid %s error %v", pid, err)
		return
	}

	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		log.Errorf("Stop container %s error %v", containerId, err)
		return
	}

	cGroupManager := cgroup.NewCgroupManager("GankerCgroup" + "/" + containerId)
	if err := cGroupManager.Delete(); err != nil {
		log.Errorf("Destroy cgroup fail %v", err)
		return
	}

	containerDir := StorageRootPath + containerId + "/"
	deleteWorkSpace(containerDir, containerInfo.Volume)

	quitContainer(containerId)
}

func checkProcessIsAlive(pId int) bool {
	process, err := os.FindProcess(pId)
	if err != nil {
		log.Errorf("Find process %d error %v", pId, err)
		return false
	}

	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false
	} else {
		return true
	}
}
