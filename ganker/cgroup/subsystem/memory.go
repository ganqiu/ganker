package subsystem

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct{}

const (
	memoryLimit = "memory.limit_in_bytes"
)

// Set the memory limit of the cgroup in the cgrouPath path
func (c *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	if res.Memory == "" {
		return nil
	}
	// Write the memory.limit_in_bytes file
	if err := ioutil.WriteFile(path.Join(subsysCgroupPath, memoryLimit),
		[]byte(res.Memory), 0644); err != nil {
		return fmt.Errorf("" + err.Error())
	}
	return nil
}

// Delete the memory limit of the cgroup in the cgrouPath path
func (c *MemorySubSystem) Delete(cgroupPath string) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	return os.Remove(subsysCgroupPath)
}

// Add the process to the cgroup in the cgrouPath path
func (c *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
	if err = ioutil.WriteFile(path.Join(subsysCgroupPath, processIdPath), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("set cgroup proc fail,error:%v", err.Error())
	}
	return nil
}

// Get the name of the subsystem
func (c *MemorySubSystem) Name() string {
	return "memory"
}
