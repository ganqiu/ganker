package container

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// CommitContainer commit the container to image
func CommitContainer(containerId, image string) {
	containerDir := ContainerRootPath + containerId
	if exist, err := checkFileOrDirExist(containerDir); err != nil {
		log.Panic("check container dir exist failed ", err)
	} else if !exist {
		log.Panic("container %v not found", containerId)
	}

	imageDir := ImageRootPath + image
	if exist, err := checkFileOrDirExist(imageDir); err != nil {
		log.Panic("check image dir exist failed ", err)
	} else if exist {
		log.Panic("image name already exist")
	}

	if err := exec.Command("tar", "-czf", imageDir, "-C", containerDir, ".").Run(); err != nil {
		log.Panic("package container dir failed ", err)
	}

	log.Infof("package container %v to image %v success", containerId, image)
}
