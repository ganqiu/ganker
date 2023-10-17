package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var CommandLogger = initCommandLogger()

var (
	rootCmd = &cobra.Command{
		Use: "ganker",

		Short: "ganker is a docker-like tool for me to zaolunzi",

		Long: `ganker is a docker-like tool , The purpose of this tool is to 
			learn how docker works and how to implement it`,

		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
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

func initCommandLogger() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		PadLevelText:    true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.SetLevel(log.InfoLevel)
	logger.SetOutput(os.Stdout)
	logger.WithField("part", "command")
	return logger
}
