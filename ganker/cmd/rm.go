package cmd

import (
	"fmt"
	"go_docker_learning/ganker/container"

	"github.com/spf13/cobra"
)

var (
	rmCmd = &cobra.Command{
		Use:   "rm [containerID]",
		Short: "rm container by containerId",
		Long:  `with rm,you can remove a container by containerId,and the container must be stopped`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println(" rm command need a containerId")
				return
			}
			container.DeleteContainer(args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(rmCmd)
}
