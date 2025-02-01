package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	sortBy string
	filter string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all modules with detailed information",
	Long: `List all Go modules with detailed information including:
- Last update date
- Module size
- Usage count
- License information
You can sort and filter the results using flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		moduleService, _, _, _ := initServices()

		modules, err := moduleService.ListModules(sortBy, filter)
		if err != nil {
			logger.WithError(err).Fatal("Failed to list modules")
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "MODULE PATH", "VERSION", "SIZE", "LAST UPDATED", "LICENSE"})
		table.SetBorder(true)

		for _, mod := range modules {
			size := fmt.Sprintf("%.2f MB", float64(mod.Size)/(1024*1024))
			lastUpdated := mod.LastUpdated.Format("2006-01-02")

			table.Append([]string{
				strconv.Itoa(mod.ID),
				mod.Path,
				mod.Version,
				size,
				lastUpdated,
				mod.License,
			})
		}

		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&sortBy, "sort", "", "Sort by: name, date, size")
	listCmd.Flags().StringVar(&filter, "filter", "", "Filter by: license:MIT, size:1000000, updated:7")
}
