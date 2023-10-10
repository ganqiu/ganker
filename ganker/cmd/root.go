package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "ganker",

		Short: "ganker is a docker-like tool for me to zaolunzi",

		Long: `ganker is a docker-like tool , The purpose of this tool is to 
			learn how docker works and how to implement it`,

		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},

		PreRunE: func(cmd *cobra.Command, args []string) error {
			log.SetFormatter(&log.JSONFormatter{})
			log.SetOutput(os.Stdout)
			return nil
		},
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}
