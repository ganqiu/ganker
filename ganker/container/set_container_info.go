package container

import (
	"encoding/json"
	"fmt"
	"go_docker_learning/ganker/cgroup"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	RUNNING      = "Up"
	EXIT         = "Exited"
	InfoFileName = "containerInfo.json"
)

type Info struct {
	Pid         string `json:"pid"`          // 容器的init进程在宿主机上的 PID
	Image       string `json:"image"`        // 容器所用镜像的名称
	ContainerId string `json:"container id"` // 容器Id
	Name        string `json:"name"`         // 容器名
	Command     string `json:"command"`      // 容器内init运行命令
	Created     string `json:"created"`      // 创建时间
	Status      string `json:"status"`       // 容器的状态
	Volume      string `json:"volume"`       // 容器的数据卷
}

func recordContainerInfo(cPid, containerId, containerName, imageName, volume string, commandArray []string) error {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, "")
	containerInfo := &Info{
		Pid:         cPid,
		Image:       imageName,
		ContainerId: containerId,
		Name:        containerName,
		Command:     command,
		Created:     createTime,
		Status:      RUNNING,
		Volume:      volume,
	}

	jsonBody, err := json.Marshal(containerInfo)
	if err != nil {
		return err
	}

	containerPath := ContainerRootPath + containerId

	if err := os.MkdirAll(containerPath, 0622); err != nil {
		return err
	}

	jsonFile := filepath.Join(containerPath, InfoFileName)
	file, err1 := os.OpenFile(jsonFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err1 != nil {
		return err
	}

	defer file.Close()

	if _, err := file.Write(jsonBody); err != nil {
		return err
	}

	return nil
}

func quitContainer(containerId string) *Info {
	containerInfoPath := filepath.Join(ContainerRootPath+containerId, InfoFileName)

	file, err := os.OpenFile(containerInfoPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Errorf("Open file %s error %v", containerInfoPath, err)
		return nil
	}

	defer file.Close()

	// 读取文件内容
	var containerInfo Info

	if err := json.NewDecoder(file).Decode(&containerInfo); err != nil {
		log.Errorf("Decode file %s error: %v\n", containerInfoPath, err)
		return nil
	}
	containerInfo.Status = EXIT
	containerInfo.Pid = " "

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Json marshal error %v", err)
		return nil
	}

	if _, err := file.WriteAt(jsonBytes, 0); err != nil {
		log.Errorf("Write string to file error %v", err)
		return nil
	}
	return &containerInfo
}

func deleteContainerInfo(containerID string) {
	dirUrl := ContainerRootPath + containerID
	if err := os.RemoveAll(dirUrl); err != nil {
		log.Errorf("Remove dir %s error %v", dirUrl, err)
	}
}

func ShowContainersInfo(all bool) {
	infoPath := ContainerRootPath
	files, err := os.ReadDir(infoPath)
	if err != nil {
		log.Errorf("Read dir %s error %v", infoPath, err)
		return
	}
	var containerInfos []*Info
	for _, file := range files {
		containerInfo, err := getContainerInfo(file.Name())
		if err != nil {
			log.Errorf("Get container info error %v", err)
			continue
		}
		if all {
			containerInfos = append(containerInfos, containerInfo)
		} else {
			if containerInfo.Status == RUNNING {
				containerInfos = append(containerInfos, containerInfo)
			}
		}
	}
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprintf(w, "ID\tPID\tSTATUS\tIMAGE\tCREATED\tNAME\tCOMMAND\tVOLUME\n")
	for _, item := range containerInfos {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			item.ContainerId,
			item.Pid,
			item.Status,
			item.Image,
			item.Created,
			item.Name,
			item.Command,
			item.Volume)
	}
	if err := w.Flush(); err != nil {
		log.Errorf("Flush error %v", err)
		return
	}
}

// getContainerInfo check if container is running, if not, quit container and return container info
func getContainerInfo(containerId string) (*Info, error) {
	containerInfo, err := getContainerFileInfo(containerId)
	if err != nil {
		log.Errorf("Get container info error %v", err)
		return nil, err
	}

	if containerInfo.Pid == " " {
		return containerInfo, nil
	}

	pid, err1 := strconv.Atoi(containerInfo.Pid)
	if err1 != nil {
		log.Errorf("Parse Pid to int failed :  %v", err)
		return nil, err
	}

	if containerInfo.Status == RUNNING && !checkProcessIsAlive(pid) {

		log.Infof("container %s process is not running now", containerId)

		cGroupManager := cgroup.NewCgroupManager("GankerCgroup" + "/" + containerId)
		if err := cGroupManager.Delete(); err != nil {
			log.Errorf("Destroy cgroup fail %v", err)
			return nil, err
		}

		containerDir := StorageRootPath + containerId + "/"
		deleteWorkSpace(containerDir, containerInfo.Volume)

		containerInfo = quitContainer(containerId)
		if containerInfo == nil {
			return nil, fmt.Errorf("quit container error")
		}
		return containerInfo, nil
	}
	return containerInfo, nil
}

// getContainerFileInfo get container Info struct from file
func getContainerFileInfo(containerId string) (*Info, error) {

	if containerId == "" {
		return nil, fmt.Errorf("container name is empty")
	}

	containerInfoPath := filepath.Join(ContainerRootPath+containerId, InfoFileName)
	file, err := os.OpenFile(containerInfoPath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var containerInfo Info
	if err := json.NewDecoder(file).Decode(&containerInfo); err != nil {
		return nil, err
	}
	return &containerInfo, nil
}
