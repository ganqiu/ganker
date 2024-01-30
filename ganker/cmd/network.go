package cmd

import (
	"fmt"

	"go_docker_learning/ganker/container"

	"github.com/spf13/cobra"
)

var (
	driver string
	subnet string
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
			// load network driver
			if err := container.InitNet(); err != nil {
				fmt.Printf("Init network failed, err: %v", err)
				return
			}

			// create network
			if err := container.CreateNet(driver, subnet, args[0]); err != nil {
				fmt.Printf("Create network failed, err: %v", err)
				return
			}
		},
	}
)

var (
	networkListCmd = &cobra.Command{
		Use:   "network list [command]",
		Short: "list all networks",
		Long:  `list all networks`,

		Run: func(cmd *cobra.Command, args []string) {
			// load network driver
			if err := container.InitNet(); err != nil {
				fmt.Printf("Init network failed, err: %v", err)
				return
			}

			// list network
			if err := container.ListNet(); err != nil {
				fmt.Printf("List network failed, err: %v", err)
				return
			}
		},
	}
)

var (
	networkDeleteCmd = &cobra.Command{
		Use:   "network delete [command]",
		Short: "delete a network",
		Long:  `delete a network`,

		Run: func(cmd *cobra.Command, args []string) {
			// load network driver
			if err := container.InitNet(); err != nil {
				fmt.Printf("Init network failed, err: %v", err)
				return
			}

			// delete network
			if err := container.DeleteNet(args[0]); err != nil {
				fmt.Printf("Delete network failed, err: %v", err)
				return
			}
		},
	}
)

var (
	networkConnectCmd = &cobra.Command{
		Use:   "network connect [containerId] [networkName]",
		Short: "connect a container to a network",
		Long:  `connect a container to a network`,

		Run: func(cmd *cobra.Command, args []string) {
			// load network driver
			if err := container.InitNet(); err != nil {
				fmt.Printf("Init network failed, err: %v", err)
				return
			}

			// Connect network
			if err := container.Connect(args[0], args[1]); err != nil {
				fmt.Printf("Delete network failed, err: %v", err)
				return
			}
		},
	}
)

var (
	networkDisconnectCmd = &cobra.Command{
		Use:   "network disconnect [command]",
		Short: "disconnect a container from a network",
		Long:  `disconnect a container from a network`,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(cmd.Help())
		},
	}
)

func init() {
	rootCmd.AddCommand(networkCmd)
	networkCmd.AddCommand(networkDeleteCmd, networkCreateCmd, networkDisconnectCmd, networkConnectCmd, networkListCmd)
	networkCreateCmd.Flags().StringVarP(&driver, "driver", "d", "bridge", "network driver")
	networkCreateCmd.Flags().StringVarP(&subnet, "subnet", "s", "", "subnet cidr")
	networkCreateCmd.MarkFlagRequired("subnet")
	networkCreateCmd.MarkFlagRequired("driver")
}
