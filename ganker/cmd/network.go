package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	networkCmd = &cobra.Command{
		Use:   "network [command]",
		Short: "a series of network commands for container",
		Long:  `it is a series of network commands for container, such as create, delete, connect, disconnect`,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Help())
		},
	}
)

var (
	networkCreateCmd = &cobra.Command{
		Use:   "network create [command]",
		Short: "create a network for container",
		Long:  `create a network for container`,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Help())
		},
	}
)

func init() {
	rootCmd.AddCommand(networkCmd)
	networkCmd.AddCommand()
}
