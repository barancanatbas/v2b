package modules

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/barancanatbas/v2b/internal/dto"
)

func (m *ModuleService) Search(query, searchType string) ([]dto.Module, error) {
	switch searchType {
	case "content":
		return m.searchContent(query)
	case "version":
		return m.searchVersion(query)
	case "date":
		return m.searchDate(query)
	default:
		return nil, fmt.Errorf("invalid search type: %s", searchType)
	}
}

func (m *ModuleService) searchContent(query string) ([]dto.Module, error) {
	// Get all modules first
	modules, err := m.GetGoModules("")
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}

	var results []dto.Module
	query = strings.ToLower(query)

	for _, mod := range modules {
		// Check in vendor directory
		vendorDir := filepath.Join("vendor", mod.Path)
		if _, err := os.Stat(vendorDir); err == nil {
			// Use grep to search in files
			cmd := exec.Command("grep", "-r", "-i", query, vendorDir)
			if err := cmd.Run(); err == nil {
				results = append(results, mod)
				continue
			}
		}

		// If not found in vendor, try downloading the module
		tmpDir, err := os.MkdirTemp("", "module-search")
		if err != nil {
			continue
		}
		defer os.RemoveAll(tmpDir)

		cmd := exec.Command("go", "get", "-d", mod.Path+"@"+mod.Version)
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			continue
		}

		// Search in downloaded files
		moduleDir := filepath.Join(tmpDir, "pkg", "mod", mod.Path+"@"+mod.Version)
		cmd = exec.Command("grep", "-r", "-i", query, moduleDir)
		if err := cmd.Run(); err == nil {
			results = append(results, mod)
		}
	}

	return results, nil
}

func (m *ModuleService) searchVersion(version string) ([]dto.Module, error) {
	modules, err := m.GetGoModules("")
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}

	var results []dto.Module
	for _, mod := range modules {
		if strings.Contains(mod.Version, version) {
			results = append(results, mod)
		}
	}

	return results, nil
}

func (m *ModuleService) searchDate(dateQuery string) ([]dto.Module, error) {
	modules, err := m.GetGoModules("")
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}

	// Parse date query (e.g., "7d" for 7 days)
	days, err := strconv.Atoi(strings.TrimSuffix(dateQuery, "d"))
	if err != nil {
		return nil, fmt.Errorf("invalid date format, use format like '7d': %w", err)
	}

	cutoffDate := time.Now().AddDate(0, 0, -days)
	var results []dto.Module

	for _, mod := range modules {
		if mod.LastUpdated.After(cutoffDate) {
			results = append(results, mod)
		}
	}

	return results, nil
}
