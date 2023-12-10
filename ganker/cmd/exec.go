package cmd

import (
	"fmt"
	"go_docker_learning/ganker/container"
	"go_docker_learning/ganker/nsenter"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const EnvExecPid = "container_pid"

var (
	execCmd = &cobra.Command{ /**/
		Use:   "exec [containerId] [command]",
		Short: "exec container and execute command",
		Long:  `exec container and execute command`,

		Run: func(cmd *cobra.Command, args []string) {

			if os.Getenv(EnvExecPid) != "" {
				log.Infof("pid:%s", os.Getenv(EnvExecPid))
				nsenter.EnterNamespace()
				return
			}
			fmt.Println("exec container and execute command")
			if len(args) < 2 {
				log.Errorf("Missing container id or command")
				return
			}

			containerId, cmdArray := args[0], args[1:]
			container.ExecContainer(containerId, cmdArray)
		},
	}
)

func init() {
	rootCmd.AddCommand(execCmd)
}
