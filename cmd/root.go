package cmd

import (
	"github.com/barancanatbas/v2b/internal/checker"
	"github.com/barancanatbas/v2b/internal/git"
	"github.com/barancanatbas/v2b/internal/modules"
	"github.com/barancanatbas/v2b/internal/tidy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logger = logrus.New()

	ignoreErrors bool
	special      bool
	prefix       string
)

var rootCmd = &cobra.Command{
	Use:   "v2b",
	Short: "v2b is a Go modules management tool",
	Long:  "A tool to fetch Go modules and determine their branches based on commit hashes",
	Run: func(cmd *cobra.Command, args []string) {
		_, _, checkerService, _ := initServices()

		err := checkerService.FetchAndDisplayModules()
		if err != nil {
			logger.WithError(err).Fatal("Failed to fetch modules")
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVar(&ignoreErrors, "ignore-errors", false, "Hide error messages from the output")
	rootCmd.Flags().BoolVar(&special, "special", false, "Only display special branches (tags, pull requests, etc.)")
	rootCmd.Flags().StringVar(&prefix, "prefix", "", "Filter modules by prefix")

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

func initServices() (*modules.ModuleService, *git.GitService, *checker.Checker, *tidy.TidyService) {
	moduleService := modules.NewModule()
	gitService := git.NewGitService()
	checkerService := checker.NewChecker(moduleService, gitService, !ignoreErrors, special, prefix)
	tidyService := tidy.NewTidyService(moduleService, gitService, prefix)

	return moduleService, gitService, checkerService, tidyService
}
