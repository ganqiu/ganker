package subsystem

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct{}

const (
	memoryLimit = "memory.limit_in_bytes"
)

// Get the name of the subsystem
func (c *MemorySubSystem) Name() string {
	return "memory"
}

// Set the memory limit of the cgroup in the cgrouPath path
func (c *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}

	// Write the memory.limit_in_bytes file
	if err := os.WriteFile(path.Join(subsysCgroupPath, memoryLimit), []byte(res.Memory), 0644); err != nil {
		return fmt.Errorf("set %s cgroup fail %v", memoryLimit, err)
	}

	return nil
}

// Delete the memory limit of the cgroup in the cgrouPath path
func (c *MemorySubSystem) Delete(cgroupPath string) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	if err := os.Remove(subsysCgroupPath); err != nil {
		return fmt.Errorf("remove %s cgroup fail %v", memoryLimit, err)
	}
	return nil
}

// Add the process to the cgroup in the cgrouPath path
func (c *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
	if err = os.WriteFile(path.Join(subsysCgroupPath, processIdPath), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("set %s cgroup proc fail,error:%v", memoryLimit, err)
	}
	return nil
}
