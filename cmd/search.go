package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	searchType string
	page       int
	pageSize   int
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for modules",
	Long: `Search for modules using different criteria:
- Content-based search: search in module source code
- Version-based search: search by version numbers
- Date-based search: search by update date

Examples:
  v2b search "logger" --type content
  v2b search "v1.0" --type version
  v2b search "7d" --type date`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleService, _, _, _ := initServices()

		query := args[0]
		results, err := moduleService.Search(query, searchType)
		if err != nil {
			logger.WithError(err).Fatal("Failed to search modules")
		}

		// Apply pagination
		start := (page - 1) * pageSize
		end := start + pageSize
		if start >= len(results) {
			logger.Fatal("Page number exceeds available results")
		}
		if end > len(results) {
			end = len(results)
		}
		results = results[start:end]

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "MODULE PATH", "VERSION", "LAST UPDATED", "LICENSE"})
		table.SetBorder(true)

		for _, mod := range results {
			lastUpdated := mod.LastUpdated.Format("2006-01-02")

			table.Append([]string{
				strconv.Itoa(mod.ID),
				mod.Path,
				mod.Version,
				lastUpdated,
				mod.License,
			})
		}

		fmt.Printf("Showing results %d-%d of %d\n", start+1, end, len(results))
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVar(&searchType, "type", "content", "Search type: content, version, date")
	searchCmd.Flags().IntVar(&page, "page", 1, "Page number for results")
	searchCmd.Flags().IntVar(&pageSize, "page-size", 10, "Number of results per page")
}
