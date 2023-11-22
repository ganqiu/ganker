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
const LowerName = "lower_layer"

// UpperName the prefix name of upper layer
const UpperName = "upper_layer"

// WorkSpaceName the prefix name of container layer
const WorkSpaceName = "work_base_layer"

// MergeLayerName the prefix name of merge layer
const MergeLayerName = "merge_layer"

// ImageRootPath is the root path of image compressed file
const ImageRootPath = "./images/"

// ContainerRootPath is the root path of container
const ContainerRootPath = "./containers/"

func newWorkSpace(image string) (string, string) {
	id := GenerateContainerId(24)
	containerDir := ContainerRootPath + id + "/"
	if err := os.MkdirAll(containerDir, 0777); err != nil {
		log.Panic("Fail to create the container dir: " + err.Error())
	}

	// create the layers
	newUpperLayer(containerDir)
	newLowerLayer(containerDir, image)
	newMergeLayer(containerDir)
	newWorkLayer(containerDir)

	// mount the overlay file system
	execMountFS(containerDir)
	log.Infof("Container id :%s Created", id)
	return containerDir, id
}

// NewReadWriteLayerr create the read-write layer
func newUpperLayer(containerDir string) string {
	// create the upper layer
	upperLayerPath := containerDir + UpperName
	if err := os.MkdirAll(upperLayerPath, 0777); err != nil {
		log.Panic("Fail to create the root dir: " + err.Error())
	}
	return upperLayerPath
}

// newWorkLayer create the work layer that is used to store the container's data (the data will be cleared after the container is mounted)
func newWorkLayer(containerDir string) string {
	// create the work layer
	workLayerPath := containerDir + WorkSpaceName
	if err := os.MkdirAll(workLayerPath, 0777); err != nil {
		log.Panic("Fail to create the work layer: " + err.Error())
	}
	return workLayerPath
}

// newMergeLayer create the merge layer that is used to merge the upper layer and lower layer to a new layer
func newMergeLayer(containerDir string) string {
	// create the merge layer
	mergeLayerPath := containerDir + MergeLayerName
	if err := os.MkdirAll(mergeLayerPath, 0777); err != nil {
		log.Panic("Fail to create the merge layer: " + err.Error())
	}
	return mergeLayerPath
}

// newLowerLayer create the read only layer
func newLowerLayer(containerDir, image string) string {

	iamgePath := ImageRootPath + image + ".tar"
	if exists, err := checkFileOrDirExist(iamgePath); err != nil {
		log.Panic("Fail to judge if the root dir of image exist: " + err.Error())
	} else if !exists {
		log.Panic("Image not found")
	}

	// create the read only layer
	lowerDir := containerDir + LowerName
	if err := os.MkdirAll(lowerDir, 0777); err != nil {
		log.Panic("Fail to create lower dir: " + err.Error())
	}

	// decompression the image to the lower dir
	if err := exec.Command("tar", "-xvf", iamgePath, "-C", lowerDir).Run(); err != nil {
		log.Panic("Fail to extract the image to the lower dir: " + err.Error())
	}

	log.Println("Decompression the image to the lower dir successfully")
	return lowerDir
}

// execMountFS mount the overlay file system
func execMountFS(containerDir string) {

	data := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", containerDir+LowerName, containerDir+UpperName, containerDir+WorkSpaceName)
	if err := syscall.Mount("overlay", containerDir+MergeLayerName, "overlay", 0, data); err != nil {
		log.Panic("Fail to mount the overlay file system: " + err.Error())
	}

	log.Println("Mount the overlay file system successfully")
}

func deleteWriteLayer(containerDir string) {
	writeLayerPath := containerDir + UpperName
	if err := os.RemoveAll(writeLayerPath); err != nil {
		log.Panic("Fail to delete the write layer: " + err.Error())
	}
	log.Println("Exit the container successfully")
}

func deleteMountPoint(containerDir string) {
	mountPoint := containerDir + MergeLayerName
	if err := syscall.Unmount(mountPoint, syscall.MNT_DETACH); err != nil {
		log.Panic("Fail to unmount the overlay file system: " + err.Error())
	}

	if err := os.RemoveAll(mountPoint); err != nil {
		log.Panic("Fail to delete the mount point: " + err.Error())
	}
	log.Println("Unmount the overlay file system successfully")
}

func deleteWorkSpace(containerDir string) {
	// delete workSpace
	deleteMountPoint(containerDir)
	deleteWriteLayer(containerDir)
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

// GenerateContainerId returns a random string with a fixed length
func GenerateContainerId(n int) string {
	rand.New(rand.NewSource(time.Now().Unix()))
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
