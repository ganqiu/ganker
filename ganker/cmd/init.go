package cmd

import (
	"go_docker_learning/ganker/container"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Define the init command
var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new container",
		Long:  `Initialize a new container`,

		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Infof("container init")
			command := args[0]
			err := container.InitRunContainerProcess(command, nil)
			return err
		},
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
}
