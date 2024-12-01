package cmd

import (
	"github.com/spf13/cobra"
)

var tidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "Fetch Go modules and determine their branches",
	Long:  "This command fetches the branches of Go module dependencies from go.mod and updates them using 'go get' based on those branches.",
	Run: func(cmd *cobra.Command, args []string) {
		_, _, _, tidyService := initServices()

		tidyService.ModTidy()
	},
}

func init() {
	tidyCmd.Flags().StringVar(&prefix, "prefix", "", "Filter modules by prefix")

	rootCmd.AddCommand(tidyCmd)
}
