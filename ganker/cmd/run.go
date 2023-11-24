package cmd

import (
	"go_docker_learning/ganker/container"
	"strings"

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
	// CgroupName = "GankerCgroup"
	volume string
	image  string
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
			CommandLogger.Infof("Running Container...")
			CommandLogger.Infof("Run command:%v", strings.Join(args, " "))
			container.RunContainer(tty, args, image, volume, ResourceConfig)
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	// add a parameter, used to specify whether to use tty
	runCmd.Flags().BoolVarP(&tty, "it", "i", false, "enable tty")
	runCmd.Flags().StringVar(&ResourceConfig.Memory, "memory-limit", "9223372036854771712", "memory limit")
	runCmd.Flags().StringVar(&ResourceConfig.CpuShare, "cpu-shares", "1024", "cpu-shares limit")
	runCmd.Flags().StringVar(&ResourceConfig.CpuQuota, "cpu-quotas", "-1", "cpuset-cpus limit")
	runCmd.Flags().StringVarP(&volume, "volume", "v", "", "add volume")
	runCmd.Flags().StringVarP(&image, "image", "m", "busybox", "choose image")
}
