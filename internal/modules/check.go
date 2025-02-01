package modules

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/barancanatbas/v2b/internal/dto"
)

type VulnerabilityInfo struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	Affected    []string `json:"affected_versions"`
}

func (m *ModuleService) CheckOutdated() ([]dto.Module, error) {
	modules, err := m.GetGoModules("")
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}

	var outdated []dto.Module
	for _, mod := range modules {
		cmd := exec.Command("go", "list", "-m", "-u", "-json", mod.Path)
		output, err := cmd.Output()
		if err != nil {
			continue
		}

		var info struct {
			Update struct {
				Version string `json:"version"`
			} `json:"update"`
		}
		if err := json.Unmarshal(output, &info); err != nil {
			continue
		}

		if info.Update.Version != "" && info.Update.Version != mod.Version {
			outdated = append(outdated, mod)
		}
	}

	return outdated, nil
}

func (m *ModuleService) CheckSecurity() ([]VulnerabilityInfo, error) {
	modules, err := m.GetGoModules("")
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}

	var vulnerabilities []VulnerabilityInfo
	for _, mod := range modules {
		// Check Go vulnerability database
		resp, err := http.Get(fmt.Sprintf("https://vuln.go.dev/ID.json?module=%s", mod.Path))
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var vulns []VulnerabilityInfo
		if err := json.NewDecoder(resp.Body).Decode(&vulns); err != nil {
			continue
		}

		for _, vuln := range vulns {
			for _, affected := range vuln.Affected {
				if strings.Contains(mod.Version, affected) {
					vulnerabilities = append(vulnerabilities, vuln)
					break
				}
			}
		}
	}

	return vulnerabilities, nil
}

func (m *ModuleService) CheckDeprecated() ([]dto.Module, error) {
	modules, err := m.GetGoModules("")
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}

	var deprecated []dto.Module
	for _, mod := range modules {
		cmd := exec.Command("go", "list", "-m", "-json", mod.Path)
		output, err := cmd.Output()
		if err != nil {
			continue
		}

		var info struct {
			Deprecated string `json:"deprecated"`
		}
		if err := json.Unmarshal(output, &info); err != nil {
			continue
		}

		if info.Deprecated != "" {
			deprecated = append(deprecated, mod)
		}
	}

	return deprecated, nil
}

func (m *ModuleService) CheckCompatibility() ([]dto.Module, error) {
	modules, err := m.GetGoModules("")
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}

	var incompatible []dto.Module
	cmd := exec.Command("go", "mod", "verify")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Parse verify output for incompatible modules
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			for _, mod := range modules {
				if strings.Contains(line, mod.Path) {
					incompatible = append(incompatible, mod)
					break
				}
			}
		}
	}

	return incompatible, nil
}

func (m *ModuleService) Cleanup() ([]dto.Module, error) {
	modules, err := m.GetGoModules("")
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}

	var removed []dto.Module
	for _, mod := range modules {
		if mod.UsageCount == 0 {
			cmd := exec.Command("go", "mod", "edit", "-droprequire", mod.Path)
			if err := cmd.Run(); err == nil {
				removed = append(removed, mod)
			}
		}
	}

	// Run tidy after removing modules
	if len(removed) > 0 {
		cmd := exec.Command("go", "mod", "tidy")
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("failed to tidy modules: %w", err)
		}
	}

	return removed, nil
}
