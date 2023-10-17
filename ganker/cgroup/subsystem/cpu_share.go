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
	if res.CpuQuota == "" {
		return nil
	}
	// Write the cpu.cfs_quota_us file
	if err := os.WriteFile(path.Join(subsysCgroupPath, cpuShareLimit), []byte(res.CpuQuota), 0644); err != nil {
		return fmt.Errorf("set cgroup cpu quota fail %v", err)
	}
	return nil
}

func (c *CpuShareSubSystem) Delete(cgroupPath string) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	return os.Remove(subsysCgroupPath)
}

func (c *CpuShareSubSystem) Apply(cgroupPath string, pid int) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
	if err = os.WriteFile(path.Join(subsysCgroupPath, processIdPath), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("set cgroup proc fail,error:%v", err.Error())
	}
	return nil
}
