package container

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// CommitContainer commit the container to image
func CommitContainer(containerId, image string) {
	containerDir := StorageRootPath + containerId
	if exist, err := checkFileOrDirExist(containerDir); err != nil {
		log.Errorf("check container dir exist failed ", err)
		return
	} else if !exist {
		log.Infof("container %v not found", containerId)
		return
	}

	imageDir := ImageRootPath + image + ".tar"
	if exist, err := checkFileOrDirExist(imageDir); err != nil {
		log.Errorf("check image dir exist failed ", err)
		return
	} else if exist {
		log.Infof("image name already exist")
		return
	}

	if err := exec.Command("tar", "-czf", imageDir, "-C", containerDir+"/"+MergeLayerName, ".").Run(); err != nil {
		log.Errorf("package container dir failed ", err)
		return
	}

	log.Infof("package container %v to image %v success", containerId, image)
}
