package subsystem

import (
	"fmt"
	"io/ioutil"
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
	if res.CpuQuota == "" {
		return nil
	}
	// Write the cpu.cfs_quota_us file
	if err := ioutil.WriteFile(path.Join(subsysCgroupPath, cpuQuotaLimit), []byte(res.CpuQuota), 0644); err != nil {
		return fmt.Errorf("set cgroup cpu quota fail %v", err)
	}
	return nil
}

func (c *CpuQuotaSubSystem) Delete(cgroupPath string) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	return os.Remove(subsysCgroupPath)
}

func (c *CpuQuotaSubSystem) Apply(cgroupPath string, pid int) error {
	subsysCgroupPath, err := FindCgroupPath(c.Name(), cgroupPath, true)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
	if err = ioutil.WriteFile(path.Join(subsysCgroupPath, processIdPath), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("set cgroup proc fail,error:%v", err.Error())
	}
	return nil
}
