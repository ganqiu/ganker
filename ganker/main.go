package main

import (
	"go_docker_learning/ganker/cmd"
	"go_docker_learning/ganker/container"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	cmd.Execute()

}

func init() {
	// image root dir
	// check if the root dir of image exist
	if err := os.MkdirAll(container.ImageRootPath, 0777); err != nil {
		log.Panic("Fail to create root dir of image: " + err.Error())
	}

	// storage root dir, which is used to store the container's data
	// check if the root dir of container exist

	if err := os.MkdirAll(container.StorageRootPath, 0777); err != nil {
		log.Panic("Fail to create root dir of container storage: " + err.Error())
	}

	// volume root dir
	// check if the root dir of volume exist
	if err := os.MkdirAll(container.VolumeRootPath, 0777); err != nil {
		log.Panic("Fail to create root dir of volume: " + err.Error())
	}

	// container root dir, which is used to store the container's info
	// check if the root dir of container exist
	if err := os.MkdirAll(container.ContainerRootPath, 0777); err != nil {
		log.Panic("Fail to create root dir of container: " + err.Error())
	}

	// network root dir, which is used to store the network's info
}
