package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [module-id]",
	Short: "Update a specific module by its ID",
	Long:  "Update a specific Go module by its ID using 'go get'",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, _, _, tidyService := initServices()

		moduleID, err := strconv.Atoi(args[0])
		if err != nil {
			logger.WithError(err).Fatal("Invalid module ID")
		}

		err = tidyService.UpdateModuleByID(moduleID)
		if err != nil {
			logger.WithError(err).Fatalf("Failed to update module with ID %d", moduleID)
		}

		logger.Printf("Successfully updated module with ID %d", moduleID)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
