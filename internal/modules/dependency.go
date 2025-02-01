package modules

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/barancanatbas/v2b/internal/dto"
)

func (m *ModuleService) PinModule(module dto.Module) (dto.Module, error) {
	if module.IsPinned {
		return dto.Module{}, fmt.Errorf("module %s is already pinned", module.Path)
	}

	// Get current version
	cmd := exec.Command("go", "list", "-m", "-json", module.Path)
	output, err := cmd.Output()
	if err != nil {
		return dto.Module{}, fmt.Errorf("failed to get module info: %w", err)
	}

	var info struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(output, &info); err != nil {
		return dto.Module{}, fmt.Errorf("failed to parse module info: %w", err)
	}

	// Pin module by adding replace directive
	cmd = exec.Command("go", "mod", "edit", "-replace", fmt.Sprintf("%s=%s@%s", module.Path, module.Path, info.Version))
	if err := cmd.Run(); err != nil {
		return dto.Module{}, fmt.Errorf("failed to pin module: %w", err)
	}

	module.IsPinned = true
	module.Version = info.Version
	return module, nil
}

func (m *ModuleService) UnpinModule(module dto.Module) (dto.Module, error) {
	if !module.IsPinned {
		return dto.Module{}, fmt.Errorf("module %s is not pinned", module.Path)
	}

	// Remove replace directive
	cmd := exec.Command("go", "mod", "edit", "-dropreplace", module.Path)
	if err := cmd.Run(); err != nil {
		return dto.Module{}, fmt.Errorf("failed to unpin module: %w", err)
	}

	// Run go mod tidy to ensure dependencies are correct
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		return dto.Module{}, fmt.Errorf("failed to tidy modules after unpinning: %w", err)
	}

	module.IsPinned = false
	return module, nil
}

func (m *ModuleService) ExplainDependency(module dto.Module) ([]dto.Module, error) {
	// Get module dependencies using go mod why
	cmd := exec.Command("go", "mod", "why", "-m", module.Path)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to explain dependency: %w", err)
	}

	// Parse output to get dependent modules
	var dependencies []dto.Module
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				dependencies = append(dependencies, dto.Module{
					Path: parts[1],
				})
			}
		}
	}

	// Get versions for each dependency
	for i, dep := range dependencies {
		cmd = exec.Command("go", "list", "-m", "-json", dep.Path)
		output, err = cmd.Output()
		if err != nil {
			continue
		}

		var info struct {
			Version string `json:"version"`
		}
		if err := json.Unmarshal(output, &info); err != nil {
			continue
		}

		dependencies[i].Version = info.Version
	}

	return dependencies, nil
}

func (m *ModuleService) GetDependencyGraph(module dto.Module) (map[string][]dto.Module, error) {
	// Get module graph using go mod graph
	cmd := exec.Command("go", "mod", "graph")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get dependency graph: %w", err)
	}

	// Parse output to build dependency graph
	graph := make(map[string][]dto.Module)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		source := parts[0]
		target := parts[1]

		// Extract version from module@version format
		sourceModule := strings.Split(source, "@")
		targetModule := strings.Split(target, "@")
		if len(sourceModule) != 2 || len(targetModule) != 2 {
			continue
		}

		// Add to graph
		if _, exists := graph[sourceModule[0]]; !exists {
			graph[sourceModule[0]] = make([]dto.Module, 0)
		}

		graph[sourceModule[0]] = append(graph[sourceModule[0]], dto.Module{
			Path:    targetModule[0],
			Version: targetModule[1],
		})
	}

	return graph, nil
}
