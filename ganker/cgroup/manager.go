package cgroup

import (
	"go_docker_learning/ganker/cgroup/subsystem"

	log "github.com/sirupsen/logrus"
)

type CgroupManager struct {
	// cgroup path in the hierarchy, relative to the each root cgroup
	Path string
	//resource config of cgroup
	Resource *subsystem.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

func (c *CgroupManager) Apply(pid int) error {
	for _, subSysInit := range subsystem.SubsystemsInit {
		err := subSysInit.Apply(c.Path, pid)
		if err != nil {
			log.Errorf("apply %v cgroup fail", subSysInit.Name())
			return err
		}
	}
	return nil
}

func (c *CgroupManager) Set(res *subsystem.ResourceConfig) error {
	for _, subSysInit := range subsystem.SubsystemsInit {
		err := subSysInit.Set(c.Path, res)
		if err != nil {
			log.Errorf("set %v cgroup fail", subSysInit.Name())
			return err
		}
	}
	return nil
}

func (c *CgroupManager) Delete() error {
	for _, subSysInit := range subsystem.SubsystemsInit {
		if err := subSysInit.Delete(c.Path); err != nil {
			return err
		}
	}
	return nil
}
