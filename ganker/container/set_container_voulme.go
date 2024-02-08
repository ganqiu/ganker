package container

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	syscall "golang.org/x/sys/unix"

	log "github.com/sirupsen/logrus"
)

// extractVolume extract the volume from the volume string
func extractVolume(volume string) ([]string, error) {
	regexpPrefix, _ := regexp.Compile(`^(\./|~/|\.\./|/)`) // ./ ~/ ../ /
	regexpVolume, _ := regexp.Compile(`[a-zA-Z0-9][a-zA-Z0-9_.-]`)
	volumeArray := strings.Split(volume, ":")
	if len(volumeArray) != 2 && volumeArray[0] != "" && volumeArray[1] != "" {
		return nil, fmt.Errorf("invalid volume specification")
	}
	if !path.IsAbs(volumeArray[1]) {
		return nil, fmt.Errorf(" invalid mount config for type \"bind\",mount path must be absolute")
	}
	if regexpPrefix.MatchString(volumeArray[0]) {
		return volumeArray, nil
	}
	if regexpVolume.MatchString(volumeArray[0]) {
		volumeArray[0] = VolumeRootPath + volumeArray[0]
		return volumeArray, nil
	}
	return nil, fmt.Errorf("path includes invalid characters for a local volume name, only '[a-zA-Z0-9][a-zA-Z0-9_.-]' are allowed")
}

// mountVolume mount the volume to the container
func mountVolume(volumeArray []string, containerDir string) {
	if err := os.MkdirAll(volumeArray[0], 0777); err != nil {
		log.Errorf("Fail to create the volume dir: " + err.Error())
		os.Exit(-1)
	}

	containerMountPoint := containerDir + MergeLayerName + volumeArray[1]
	if err := os.MkdirAll(containerMountPoint, 0777); err != nil {
		log.Errorf("Fail to create the container mount point: " + err.Error())
		os.Exit(-1)
	}

	if err := syscall.Mount(volumeArray[0], containerMountPoint, "", syscall.MS_BIND, ""); err != nil {
		log.Errorf("Fail to mount the volume: " + err.Error())
		os.Exit(-1)
	}
}

// deleteVolume delete the volume from the container
func deleteVolume(containerDir, volume string) {
	containerVolumePath := containerDir + MergeLayerName + volume
	if err := syscall.Unmount(containerVolumePath, syscall.MNT_DETACH); err != nil {
		log.Errorf("Fail to unmount the volume: " + err.Error())
	}
	os.Exit(-1)
}
