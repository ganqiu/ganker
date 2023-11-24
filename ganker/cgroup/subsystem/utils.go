package subsystem

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// FindCgroupMountPoint finds the mount point of the cgroup subsystem
func FindCgroupMountPoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	//search for the subsystem, if found, return the mount point
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		if strings.Contains(fields[len(fields)-1], subsystem) {
			if subsystem == "cpu" {
				if strings.Contains(fields[4], "cpu,") {
					return strings.Split(fields[4], ",")[0]
				}
				continue
			}
			return fields[4]
		}
	}

	if err := scanner.Err(); err != nil {
		return ""
	}
	return ""
}

func FindCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountPoint(subsystem)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && errors.Is(err, os.ErrNotExist)) {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(path.Join(cgroupRoot, cgroupPath), 0755); err != nil {
				return "", fmt.Errorf("error create %v cgroup:%v", subsystem, err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("find cgroup in path \"%v\" error:%v", cgroupPath, err)
	}
}
