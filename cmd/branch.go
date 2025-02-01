package cmd

import (
	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use:   "branch [module-path] [branch-name]",
	Short: "Update a module to use a specific branch",
	Long:  "Update a Go module to use a specific branch version using 'go get module@branch'",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		_, _, _, tidyService := initServices()

		modulePath := args[0]
		branchName := args[1]

		err := tidyService.UpdateModuleBranch(modulePath, branchName)
		if err != nil {
			logger.WithError(err).Fatalf("Failed to update module %s to branch %s", modulePath, branchName)
		}

		logger.Printf("Successfully updated module %s to branch %s", modulePath, branchName)
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
}
