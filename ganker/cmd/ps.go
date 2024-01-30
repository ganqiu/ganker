package cmd

import (
	"go_docker_learning/ganker/container"

	"github.com/spf13/cobra"
)

var all bool
var (
	psCmd = &cobra.Command{
		Use:   "ps",
		Short: "show container list",
		Long:  `Show container list and brief info`,

		Run: func(cmd *cobra.Command, args []string) {
			container.ShowContainersInfo(all)

		},
	}
)

func init() {
	rootCmd.AddCommand(psCmd)
	psCmd.Flags().BoolVarP(&all, "all", "a", false, "show all containers status")
}
