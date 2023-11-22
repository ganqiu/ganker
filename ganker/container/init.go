package container

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// CreateRootDir create the root dir of image and container
func CreateRootDir() {
	// image root dir
	// check if the root dir of image exist
	if err := os.MkdirAll(ImageRootPath, 0777); err != nil {
		log.Panic("Fail to create root dir of image: " + err.Error())
	}

	// container root dir
	// check if the root dir of container exist

	if err := os.MkdirAll(ContainerRootPath, 0777); err != nil {
		log.Panic("Fail to create root dir of container: " + err.Error())
	}
}
