package cmd

import (
	"fmt"
	"go_docker_learning/ganker/container"

	"github.com/spf13/cobra"
)

var (
	stopCmd = &cobra.Command{ /**/
		Use:   "stop [containerID]",
		Short: "stop a container",
		Long:  `it is used to stop a container that is running`,

		Run: func(cmd *cobra.Command, args []string) {

			if len(args) < 1 {
				fmt.Println("Please input containerID")
				return
			}
			containerID := args[0]
			container.StopContainer(containerID)
		},
	}
)

func init() {
	rootCmd.AddCommand(stopCmd)
}
