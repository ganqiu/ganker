package subsystem

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

type CpuShareSubSystem struct{}

const cpuShareLimit = "cpu.shares"

func (c *CpuShareSubSystem) Name() string {
	return "cpu"
}

func (c *CpuShareSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	// Write the cpu.cfs_quota_us file
	if err := os.WriteFile(path.Join(subsysCgroupPath, cpuShareLimit), []byte(res.CpuShare), 0644); err != nil {
		fmt.Println("error happen", res.CpuShare)
		return fmt.Errorf("set %s cgroup fail %v", cpuShareLimit, err)
	}
	return nil
}

func (c *CpuShareSubSystem) Delete(cgroupPath string) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	if err := os.Remove(subsysCgroupPath); err != nil {
		return fmt.Errorf("remove %s cgroup fail %v", cpuShareLimit, err)
	}
	return nil
}

func (c *CpuShareSubSystem) Apply(cgroupPath string, pid int) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
	if err = os.WriteFile(path.Join(subsysCgroupPath, processIdPath), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("set %s cgroup proc fail,error:%v", cpuShareLimit, err)
	}
	return nil
}
