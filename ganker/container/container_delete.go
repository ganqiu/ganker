package container

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func DeleteContainer(containerId string) {
	containerInfo, err := getContainerInfo(containerId)
	if err != nil {
		log.Errorf("get container %s info error %v", containerId, err)
	}

	if containerInfo.Status != EXIT {
		log.Errorf("couldn't remove running container, please stop it first")
		return
	}

	deleteContainerInfo(containerId)

	storagePath := StorageRootPath + containerId
	if err := os.RemoveAll(storagePath); err != nil {
		log.Errorf("remove container %s error %v", containerId, err)
	}

}
