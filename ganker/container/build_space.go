package container

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// LowerName  the prefix name of lower layer
const LowerName = "lower"

// UpperName the prefix name of upper layer
const UpperName = "diff"

// WorkSpaceName the prefix name of container layer
const WorkSpaceName = "work"

// MergeLayerName the prefix name of merge layer
const MergeLayerName = "merged"

// ImageRootPath is the root path of image compressed file
const ImageRootPath = "./images/"

// StorageRootPath is the root path of container
const StorageRootPath = "./storage/"

// VolumeRootPath is the root path of volume
const VolumeRootPath = "./volumes/"

// ContainerRootPath is the root path of container
const ContainerRootPath = "./containers/"

// NewWorkSpace create the work space for the container
func newWorkSpace(id, image, volume string) (string, string) {
	var volumeArray []string
	var err error
	if volume != "" {
		volumeArray, err = extractVolume(volume)
		if volumeArray == nil {
			log.Errorf("Fail to extract the volume: " + err.Error())
			os.Exit(-1)
		}
	}
	imagePath := ImageRootPath + image + ".tar"
	if exists, err := checkFileOrDirExist(imagePath); err != nil {
		log.Errorf("Fail to judge if the root dir of image exist: " + err.Error())
		os.Exit(-1)
	} else if !exists {
		log.Infof("Image not found")
		os.Exit(-1)
	}
	containerDir := StorageRootPath + id + "/"
	if err := os.MkdirAll(containerDir, 0777); err != nil {
		log.Errorf("Fail to create the container dir:" + err.Error())
		os.Exit(-1)
	}

	// create the layers
	newUpperLayer(containerDir)
	newLowerLayer(containerDir, imagePath)
	newMergeLayer(containerDir)
	newWorkLayer(containerDir)

	// mount the overlay file system
	execMountFS(containerDir)

	if volume != "" {
		mountVolume(volumeArray, containerDir)
	}
	fmt.Printf("container: %v is created \n", id)
	return containerDir, id
}

// NewReadWriteLayer create the read-write layer
func newUpperLayer(containerDir string) string {
	// create the upper layer
	upperLayerPath := containerDir + UpperName
	if err := os.MkdirAll(upperLayerPath, 0777); err != nil {
		log.Errorf("Fail to create the root dir: " + err.Error())
		os.Exit(-1)
	}
	return upperLayerPath
}

// newWorkLayer create the work layer that is used to store the container's data (the data will be cleared after the container is mounted)
func newWorkLayer(containerDir string) string {
	// create the work layer
	workLayerPath := containerDir + WorkSpaceName
	if err := os.MkdirAll(workLayerPath, 0777); err != nil {
		log.Errorf("Fail to create the work layer: " + err.Error())
		os.Exit(-1)
	}
	return workLayerPath
}

// newMergeLayer create the merge layer that is used to merge the upper layer and lower layer to a new layer
func newMergeLayer(containerDir string) string {
	// create the merge layer
	mergeLayerPath := containerDir + MergeLayerName
	if err := os.MkdirAll(mergeLayerPath, 0777); err != nil {
		log.Errorf("Fail to create the merge layer: " + err.Error())
		os.Exit(-1)
	}
	return mergeLayerPath
}

// newLowerLayer create the read only layer
func newLowerLayer(containerDir, imagePath string) string {

	// create the read only layer
	lowerDir := containerDir + LowerName
	if err := os.MkdirAll(lowerDir, 0777); err != nil {
		log.Errorf("Fail to create lower dir: " + err.Error())
		os.Exit(-1)
	}

	// decompression the image to the lower dir
	if err := exec.Command("tar", "-xvf", imagePath, "-C", lowerDir).Run(); err != nil {
		log.Errorf("Fail to extract the image to the lower dir: " + err.Error())
		os.Exit(-1)
	}

	return lowerDir
}

// execMountFS mount the overlay file system
func execMountFS(containerDir string) {

	data := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", containerDir+LowerName, containerDir+UpperName, containerDir+WorkSpaceName)
	if err := syscall.Mount("overlay", containerDir+MergeLayerName, "overlay", 0, data); err != nil {
		log.Errorf("Fail to mount the overlay file system: " + err.Error())
		os.Exit(-1)
	}

}

func deleteMountPoint(containerDir string) {
	mountPoint := containerDir + MergeLayerName
	if err := syscall.Unmount(mountPoint, syscall.MNT_DETACH); err != nil {
		log.Errorf("Fail to unmount the overlay file system: " + err.Error())
	}

	if err := os.RemoveAll(mountPoint); err != nil {
		log.Errorf("Fail to delete the mount point: " + err.Error())
	} else {
		log.Println("Unmount the overlay file system successfully")
	}
}

func deleteWorkSpace(containerDir, volume string) {
	// check if the volume is empty
	if volume != "" {
		volumeArray, err := extractVolume(volume)
		if volumeArray == nil {
			log.Errorf("Fail to extract the volume when delete: " + err.Error())
		}
		deleteVolume(containerDir, volumeArray[1])
	}
	deleteMountPoint(containerDir)
}

func checkFileOrDirExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		// file or dir exist
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		// file or dir not exist
		return false, nil
	}
	return false, err
}

// generateContainerId returns a random string with a fixed length
func generateContainerId(n int) string {
	rand.New(rand.NewSource(time.Now().Unix()))
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
