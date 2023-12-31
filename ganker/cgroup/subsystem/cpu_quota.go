package subsystem

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

type CpuQuotaSubSystem struct{}

const (
	cpuQuotaLimit = "cpu.cfs_quota_us"
)

func (c *CpuQuotaSubSystem) Name() string {
	return "cpu"
}

func (c *CpuQuotaSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	// Write the cpu.cfs_quota_us file
	if err := os.WriteFile(path.Join(subsysCgroupPath, cpuQuotaLimit), []byte(res.CpuQuota), 0644); err != nil {
		return fmt.Errorf("set %s cgroup fail %v", cpuQuotaLimit, err)
	}
	return nil
}

func (c *CpuQuotaSubSystem) Delete(cgroupPath string) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	if err := os.Remove(subsysCgroupPath); err != nil {
		return fmt.Errorf("remove %s cgroup fail %v", cpuQuotaLimit, err)
	}
	return nil
}

func (c *CpuQuotaSubSystem) Apply(cgroupPath string, pid int) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
	if err = os.WriteFile(path.Join(subsysCgroupPath, processIdPath), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("set %s cgroup proc fail,error:%v", cpuQuotaLimit, err)
	}
	return nil
}
