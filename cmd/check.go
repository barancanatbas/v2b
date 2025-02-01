package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/barancanatbas/v2b/internal/modules"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	checkType string
	autoFix   bool
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check modules for various issues",
	Long: `Check Go modules for various issues:
- Outdated versions
- Security vulnerabilities
- Deprecated modules
- Version compatibility
- Unused modules

Examples:
  v2b check --type outdated
  v2b check --type security
  v2b check --type deprecated
  v2b check --type compatibility
  v2b check --type cleanup --auto-fix`,
	Run: func(cmd *cobra.Command, args []string) {
		moduleService, _, _, _ := initServices()

		switch checkType {
		case "outdated":
			outdated, err := moduleService.CheckOutdated()
			if err != nil {
				logger.WithError(err).Fatal("Failed to check outdated modules")
			}
			displayModules("Outdated Modules", outdated)

		case "security":
			vulns, err := moduleService.CheckSecurity()
			if err != nil {
				logger.WithError(err).Fatal("Failed to check security vulnerabilities")
			}
			displayVulnerabilities(vulns)

		case "deprecated":
			deprecated, err := moduleService.CheckDeprecated()
			if err != nil {
				logger.WithError(err).Fatal("Failed to check deprecated modules")
			}
			displayModules("Deprecated Modules", deprecated)

		case "compatibility":
			incompatible, err := moduleService.CheckCompatibility()
			if err != nil {
				logger.WithError(err).Fatal("Failed to check module compatibility")
			}
			displayModules("Incompatible Modules", incompatible)

		case "cleanup":
			removed, err := moduleService.Cleanup()
			if err != nil {
				logger.WithError(err).Fatal("Failed to cleanup modules")
			}
			if len(removed) > 0 {
				displayModules("Removed Modules", removed)
			} else {
				fmt.Println("No unused modules found")
			}

		default:
			logger.Fatalf("Invalid check type: %s", checkType)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringVar(&checkType, "type", "", "Check type: outdated, security, deprecated, compatibility, cleanup")
	checkCmd.Flags().BoolVar(&autoFix, "auto-fix", false, "Automatically fix issues (only applicable for cleanup)")
	checkCmd.MarkFlagRequired("type")
}

func displayModules(title string, modules []dto.Module) {
	if len(modules) == 0 {
		fmt.Printf("No %s found\n", title)
		return
	}

	fmt.Printf("\n%s:\n", title)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "MODULE PATH", "VERSION", "LICENSE"})
	table.SetBorder(true)

	for _, mod := range modules {
		table.Append([]string{
			strconv.Itoa(mod.ID),
			mod.Path,
			mod.Version,
			mod.License,
		})
	}

	table.Render()
}

func displayVulnerabilities(vulns []modules.VulnerabilityInfo) {
	if len(vulns) == 0 {
		fmt.Println("No security vulnerabilities found")
		return
	}

	fmt.Println("\nSecurity Vulnerabilities:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "SEVERITY", "DESCRIPTION", "AFFECTED VERSIONS"})
	table.SetBorder(true)

	for _, vuln := range vulns {
		table.Append([]string{
			vuln.ID,
			vuln.Severity,
			vuln.Description,
			strings.Join(vuln.Affected, ", "),
		})
	}

	table.Render()
}
