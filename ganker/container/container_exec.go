package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

const EnvExecPid = "container_pid"
const EnvExecCmd = "container_cmd"

// ExecContainer exec container by containerId
func ExecContainer(containerId string, cmdArray []string) {

	containerInfo, err := getContainerInfo(containerId)
	if err != nil {
		log.Errorf("Exec container getContainerInfo %s error %v", containerId, err)
		return
	}

	if containerInfo.Status != RUNNING {
		log.Errorf("Exec container %s not running", containerId)
		return
	}

	cmdStr := strings.Join(cmdArray, " ")
	log.Infof("container %v exec cmd %s", containerInfo.ContainerId, cmdStr)

	// logOut := recordContainerLog(containerId)
	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdin
	cmd.Stderr = os.Stdout

	if err := os.Setenv(EnvExecPid, containerInfo.Pid); err != nil {
		log.Errorf("Exec container setenv %s error %v", containerId, err)
		return
	}

	if err := os.Setenv(EnvExecCmd, cmdStr); err != nil {
		log.Errorf("Exec container setenv %s error %v", containerId, err)
		return
	}

	// set env that container process can inherit
	cmd.Env = append(os.Environ(), getEnvByPid(containerInfo.Pid)...)

	if err := cmd.Run(); err != nil {
		log.Errorf("Exec container %s error %v", containerId, err)
		return
	}
}

// getEnvPid get env from process pid
func getEnvByPid(pid string) []string {

	// process storage env in /proc/pid/environ
	processPath := fmt.Sprintf("/proc/%s/environ", pid)

	content, err := os.ReadFile(processPath)
	if err != nil {
		log.Errorf("getEnvPid read file %s error %v", processPath, err)
		return nil
	}

	envSlice := strings.Split(string(content), "\u0000")
	return envSlice
}
