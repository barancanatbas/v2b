package cmd

import (
	"fmt"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/spf13/cobra"
)

var (
	modulePath  string
	newVersion  string
	autoResolve bool
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Manage module versions",
	Long: `Manage Go module versions:
- Upgrade to a specific version
- Rollback to previous version
- Resolve version conflicts

Examples:
  v2b version upgrade --module github.com/example/mod --version v1.2.0
  v2b version rollback --module github.com/example/mod
  v2b version resolve --auto`,
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade module to a specific version",
	Run: func(cmd *cobra.Command, args []string) {
		moduleService, _, _, _ := initServices()

		module := dto.Module{
			Path:    modulePath,
			Version: "", // Current version will be determined by the service
		}

		upgraded, err := moduleService.UpgradeVersion(module, newVersion)
		if err != nil {
			logger.WithError(err).Fatal("Failed to upgrade module version")
		}

		fmt.Printf("Successfully upgraded %s to version %s\n", upgraded.Path, upgraded.Version)
	},
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback module to previous version",
	Run: func(cmd *cobra.Command, args []string) {
		moduleService, _, _, _ := initServices()

		module := dto.Module{
			Path:    modulePath,
			Version: "", // Current version will be determined by the service
		}

		rolledBack, err := moduleService.RollbackVersion(module)
		if err != nil {
			logger.WithError(err).Fatal("Failed to rollback module version")
		}

		fmt.Printf("Successfully rolled back %s to version %s\n", rolledBack.Path, rolledBack.Version)
	},
}

var resolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve version conflicts",
	Run: func(cmd *cobra.Command, args []string) {
		moduleService, _, _, _ := initServices()

		modules, err := moduleService.GetGoModules("")
		if err != nil {
			logger.WithError(err).Fatal("Failed to get modules")
		}

		resolved, err := moduleService.ResolveVersionConflict(modules)
		if err != nil {
			logger.WithError(err).Fatal("Failed to resolve version conflicts")
		}

		if len(resolved) == 0 {
			fmt.Println("No version conflicts found")
			return
		}

		fmt.Println("\nResolved version conflicts:")
		for _, mod := range resolved {
			fmt.Printf("- %s: updated to version %s\n", mod.Path, mod.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.AddCommand(upgradeCmd, rollbackCmd, resolveCmd)

	// Flags for upgrade command
	upgradeCmd.Flags().StringVar(&modulePath, "module", "", "Module path to upgrade")
	upgradeCmd.Flags().StringVar(&newVersion, "version", "", "Target version to upgrade to")
	upgradeCmd.MarkFlagRequired("module")
	upgradeCmd.MarkFlagRequired("version")

	// Flags for rollback command
	rollbackCmd.Flags().StringVar(&modulePath, "module", "", "Module path to rollback")
	rollbackCmd.MarkFlagRequired("module")

	// Flags for resolve command
	resolveCmd.Flags().BoolVar(&autoResolve, "auto", false, "Automatically resolve conflicts")
}
