package cmd

import (
	"go_docker_learning/ganker/container"

	"github.com/spf13/cobra"
)

// Define the commit command
var (
	commitCmd = &cobra.Command{
		Use:   "commit",
		Short: "package container into image",
		Long:  `package container into image`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				CommandLogger.Info("missing container command")
			}
			container.CommitContainer(args[0], args[1])
		},
	}
)

func init() {
	rootCmd.AddCommand(commitCmd)
}
