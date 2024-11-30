package cmd

import (
	"github.com/barancanatbas/v2b/internal/modules"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logger = logrus.New()

	errFlag     bool
	specialFlag bool
	prefixFlag  string
)

var rootCmd = &cobra.Command{
	Use:   "v2b",
	Short: "v2b is a Go modules management tool",
	Long:  "A tool to fetch Go modules and determine their branches based on commit hashes",
	Run: func(cmd *cobra.Command, args []string) {
		moduleService := modules.NewModule(errFlag, specialFlag, prefixFlag)
		err := moduleService.FetchAndDisplayModules()
		if err != nil {
			logger.WithError(err).Fatal("Failed to fetch modules")
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVar(&errFlag, "err", false, "Print error messages along with successful results")
	rootCmd.Flags().BoolVar(&specialFlag, "special", false, "Only display special branches (tags, pull requests, etc.)")
	rootCmd.Flags().StringVar(&prefixFlag, "prefix", "", "Filter modules by prefix")

	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "Set the log level (debug, info, warn, error, fatal, panic)")
	cobra.OnInitialize(initLogger)
}

func initLogger() {
	logLevel, _ := rootCmd.Flags().GetString("log-level")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.WithError(err).Fatal("Invalid log level specified")
	}

	logger.SetLevel(level)
}
