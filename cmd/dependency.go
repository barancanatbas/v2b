package cmd

import (
	"fmt"
	"os"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var dependencyCmd = &cobra.Command{
	Use:   "dependency",
	Short: "Manage module dependencies",
	Long: `Manage Go module dependencies:
- Pin modules to specific versions
- Unpin modules
- Explain why a module is needed
- Show dependency graph

Examples:
  v2b dependency pin --module github.com/example/mod
  v2b dependency unpin --module github.com/example/mod
  v2b dependency why --module github.com/example/mod
  v2b dependency graph --module github.com/example/mod`,
}

var pinCmd = &cobra.Command{
	Use:   "pin",
	Short: "Pin a module to its current version",
	Run: func(cmd *cobra.Command, args []string) {
		moduleService, _, _, _ := initServices()

		module := dto.Module{
			Path: modulePath,
		}

		pinned, err := moduleService.PinModule(module)
		if err != nil {
			logger.WithError(err).Fatal("Failed to pin module")
		}

		fmt.Printf("Successfully pinned %s to version %s\n", pinned.Path, pinned.Version)
	},
}

var unpinCmd = &cobra.Command{
	Use:   "unpin",
	Short: "Unpin a module",
	Run: func(cmd *cobra.Command, args []string) {
		moduleService, _, _, _ := initServices()

		module := dto.Module{
			Path: modulePath,
		}

		unpinned, err := moduleService.UnpinModule(module)
		if err != nil {
			logger.WithError(err).Fatal("Failed to unpin module")
		}

		fmt.Printf("Successfully unpinned %s\n", unpinned.Path)
	},
}

var whyCmd = &cobra.Command{
	Use:   "why",
	Short: "Explain why a module is needed",
	Run: func(cmd *cobra.Command, args []string) {
		moduleService, _, _, _ := initServices()

		module := dto.Module{
			Path: modulePath,
		}

		deps, err := moduleService.ExplainDependency(module)
		if err != nil {
			logger.WithError(err).Fatal("Failed to explain dependency")
		}

		if len(deps) == 0 {
			fmt.Printf("No dependencies found for %s\n", module.Path)
			return
		}

		fmt.Printf("\nDependencies for %s:\n", module.Path)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"MODULE", "VERSION", "REASON"})
		table.SetBorder(true)

		for _, dep := range deps {
			table.Append([]string{
				dep.Path,
				dep.Version,
				"Required by project",
			})
		}

		table.Render()
	},
}

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Show dependency graph",
	Run: func(cmd *cobra.Command, args []string) {
		moduleService, _, _, _ := initServices()

		module := dto.Module{
			Path: modulePath,
		}

		graph, err := moduleService.GetDependencyGraph(module)
		if err != nil {
			logger.WithError(err).Fatal("Failed to get dependency graph")
		}

		if len(graph) == 0 {
			fmt.Printf("No dependencies found for %s\n", module.Path)
			return
		}

		fmt.Printf("\nDependency graph for %s:\n", module.Path)
		for source, targets := range graph {
			fmt.Printf("\n%s\n", source)
			for _, target := range targets {
				fmt.Printf("  └── %s@%s\n", target.Path, target.Version)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(dependencyCmd)
	dependencyCmd.AddCommand(pinCmd, unpinCmd, whyCmd, graphCmd)

	// Flags for all commands
	for _, cmd := range []*cobra.Command{pinCmd, unpinCmd, whyCmd, graphCmd} {
		cmd.Flags().StringVar(&modulePath, "module", "", "Module path to operate on")
		cmd.MarkFlagRequired("module")
	}
}
