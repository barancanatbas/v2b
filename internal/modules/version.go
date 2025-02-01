package modules

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/barancanatbas/v2b/internal/dto"
)

func (m *ModuleService) UpgradeVersion(module dto.Module, newVersion string) (dto.Module, error) {
	// Validate versions
	current, err := semver.NewVersion(module.Version)
	if err != nil {
		return dto.Module{}, fmt.Errorf("invalid current version: %w", err)
	}

	target, err := semver.NewVersion(newVersion)
	if err != nil {
		return dto.Module{}, fmt.Errorf("invalid target version: %w", err)
	}

	if current.Equal(target) {
		return dto.Module{}, fmt.Errorf("target version is same as current version")
	}

	// Execute go get with specific version
	cmd := exec.Command("go", "get", fmt.Sprintf("%s@%s", module.Path, newVersion))
	if err := cmd.Run(); err != nil {
		return dto.Module{}, fmt.Errorf("failed to upgrade version: %w", err)
	}

	// Run go mod tidy to ensure dependencies are correct
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		return dto.Module{}, fmt.Errorf("failed to tidy modules after upgrade: %w", err)
	}

	// Return updated module info
	module.Version = newVersion
	return module, nil
}

func (m *ModuleService) RollbackVersion(module dto.Module) (dto.Module, error) {
	// Get version history
	cmd := exec.Command("go", "list", "-m", "-versions", module.Path)
	output, err := cmd.Output()
	if err != nil {
		return dto.Module{}, fmt.Errorf("failed to get version history: %w", err)
	}

	versions := strings.Fields(string(output))
	if len(versions) <= 1 {
		return dto.Module{}, fmt.Errorf("no previous version available for rollback")
	}

	// Find current version index
	currentIdx := -1
	for i, v := range versions {
		if v == module.Version {
			currentIdx = i
			break
		}
	}

	if currentIdx <= 0 {
		return dto.Module{}, fmt.Errorf("cannot rollback from first version")
	}

	// Get previous version
	previousVersion := versions[currentIdx-1]

	// Execute rollback
	cmd = exec.Command("go", "get", fmt.Sprintf("%s@%s", module.Path, previousVersion))
	if err := cmd.Run(); err != nil {
		return dto.Module{}, fmt.Errorf("failed to rollback version: %w", err)
	}

	// Run go mod tidy to ensure dependencies are correct
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		return dto.Module{}, fmt.Errorf("failed to tidy modules after rollback: %w", err)
	}

	// Return updated module info
	module.Version = previousVersion
	return module, nil
}

func (m *ModuleService) ResolveVersionConflict(modules []dto.Module) ([]dto.Module, error) {
	// Create version constraint map
	constraints := make(map[string][]string)

	// Collect all version constraints for each module
	for _, mod := range modules {
		if _, exists := constraints[mod.Path]; !exists {
			constraints[mod.Path] = []string{mod.Version}
		} else {
			constraints[mod.Path] = append(constraints[mod.Path], mod.Version)
		}
	}

	// Resolve conflicts for each module
	resolved := make([]dto.Module, 0)
	for path, versions := range constraints {
		if len(versions) <= 1 {
			// No conflict for this module
			continue
		}

		// Parse all versions and find the highest compatible version
		var highest *semver.Version
		for _, v := range versions {
			version, err := semver.NewVersion(v)
			if err != nil {
				continue
			}

			if highest == nil || version.GreaterThan(highest) {
				highest = version
			}
		}

		if highest == nil {
			return nil, fmt.Errorf("failed to resolve version conflict for %s", path)
		}

		// Update to the highest compatible version
		cmd := exec.Command("go", "get", fmt.Sprintf("%s@%s", path, highest.String()))
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("failed to update module %s to version %s: %w", path, highest.String(), err)
		}

		resolved = append(resolved, dto.Module{
			Path:    path,
			Version: highest.String(),
		})
	}

	// Run go mod tidy to ensure all dependencies are correct
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		return nil, fmt.Errorf("failed to tidy modules after resolving conflicts: %w", err)
	}

	return resolved, nil
}
