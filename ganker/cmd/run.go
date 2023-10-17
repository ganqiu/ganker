package cmd

import (
	"go_docker_learning/ganker/container"

	"go_docker_learning/ganker/cgroup/subsystem"

	"github.com/spf13/cobra"
)

var (
	tty            bool
	ResourceConfig = &subsystem.ResourceConfig{
		Memory:   "9223372036854771712",
		CpuShare: "1024",
		CpuQuota: "-1",
	}
	CgroupName = "GankeCgroup"
)

// Define the run command
var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Create a container ",
		Long:  `Create a container by ganker run  [arg]`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				CommandLogger.Info("missing container command")
			}
			command := args[0]
			CommandLogger.Infof("Run command %s", command)
			container.RunContainer(tty, command, ResourceConfig)
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	// add a parameter, used to specify whether to use tty
	runCmd.Flags().BoolVarP(&tty, "it", "i", false, "enable tty")
	runCmd.Flags().StringVarP(&ResourceConfig.Memory, "memory-limit", "m", "9223372036854771712", "memory limit")
	runCmd.Flags().StringVarP(&ResourceConfig.CpuShare, "cpu-shares", "cs", "1024", "cpu-shares limit")
	runCmd.Flags().StringVarP(&ResourceConfig.CpuQuota, "cpu-quotas", "cq", "-1", "cpuset-cpus limit")
}
