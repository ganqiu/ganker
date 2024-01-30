package cmd

import (
	"go_docker_learning/ganker/container"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logCmd = &cobra.Command{
		Use:   "log [containerId]",
		Short: "print logs of a container",
		Long:  `print logs of a container`,

		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Errorf("Missing container id")
				return
			}
			container.ShowContainerLog(args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(logCmd)
}
