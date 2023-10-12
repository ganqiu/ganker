package cmd

import (
	"fmt"
	"go_docker_learning/ganker/container"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tty bool

// Define the run command
var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Create a container ",
		Long:  `Create a container by ganker run  [arg]`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("missing container command")
			}
			command := args[0]
			log.Infof("Run command %s", command)
			container.RunContainer(tty, command)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	// add a parameter, used to specify whether to use tty
	runCmd.Flags().BoolVarP(&tty, "it", "i", false, "enable tty")
}
