package cmd

import (
	"fmt"
	"go_docker_learning/ganker/container"

	"github.com/spf13/cobra"
)

var tty bool

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Create a container ",
		Long:  `Create a container by ganker run  [arg]`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("Missing container command")
			}
			command := args[0]
			container.RunContainer(tty, command)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&tty, "it", "i", false, "enable tty")
}
