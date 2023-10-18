package cmd

import (
	"go_docker_learning/ganker/container"

	"github.com/spf13/cobra"
)

// Define the init command
var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new container",
		Long:  `Initialize a new container`,

		Run: func(cmd *cobra.Command, args []string) {
			err := container.InitRunContainerProcess()
			if err != nil {
				CommandLogger.Infof("init container failed %v", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}
