package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/barancanatbas/v2b/internal/dto"
)

type ModuleService struct {
}

func NewModule() *ModuleService {
	return &ModuleService{}
}

func (m *ModuleService) GetGoModules(prefix string) ([]dto.Module, error) {
	// Try to read go.mod file directly first
	content, err := os.ReadFile("go.mod")
	if err == nil {
		modules, err := m.parseGoMod(content)
		if err == nil {
			// Assign IDs to modules
			for i := range modules {
				modules[i].ID = i + 1
			}

			// Filter by prefix if specified
			if prefix != "" {
				var filtered []dto.Module
				for _, mod := range modules {
					if strings.HasPrefix(mod.Path, prefix) {
						filtered = append(filtered, mod)
					}
				}
				return filtered, nil
			}
			return modules, nil
		}
	}

	// Fallback to go list command
	cmd := exec.Command("go", "list", "-m", "-json", "all")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run 'go list': %v", err)
	}

	var modules []dto.Module
	decoder := json.NewDecoder(&out)
	id := 1
	for decoder.More() {
		var mod dto.Module
		if err := decoder.Decode(&mod); err != nil {
			return nil, fmt.Errorf("failed to decode module data: %v", err)
		}

		if prefix != "" {
			if !strings.HasPrefix(mod.Path, prefix) {
				continue
			}
		}

		mod.ID = id
		id++
		modules = append(modules, mod)
	}

	return modules, nil
}

func (m *ModuleService) TidyForModule(module dto.Module) error {
	if module.Branch == nil {
		return nil
	}

	cmd := exec.Command("go", "get", module.Path+"@"+*module.Branch)
	cmd.Dir = module.Path

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run 'go get': %v, output: %s", err, string(output))
	}

	fmt.Println(string(output))
	return nil
}

func (m *ModuleService) parseGoMod(content []byte) ([]dto.Module, error) {
	var modules []dto.Module
	lines := strings.Split(string(content), "\n")

	inRequireBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if line == "require (" {
			inRequireBlock = true
			continue
		}

		if line == ")" {
			inRequireBlock = false
			continue
		}

		if inRequireBlock || strings.HasPrefix(line, "require ") {
			// Remove "require " prefix if it exists
			if strings.HasPrefix(line, "require ") {
				line = strings.TrimPrefix(line, "require ")
			}

			// Split the line into parts
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				// Remove indirect comment if exists
				version := parts[1]
				if strings.Contains(version, "//") {
					version = strings.Fields(version)[0]
				}

				modules = append(modules, dto.Module{
					Path:    parts[0],
					Version: version,
				})
			}
		}
	}

	if len(modules) == 0 {
		return nil, fmt.Errorf("no valid modules found in go.mod")
	}

	return modules, nil
}

func (m *ModuleService) GetModuleDetails(modulePath string) (*dto.Module, error) {
	// Get basic module info
	modules, err := m.GetGoModules("")
	if err != nil {
		return nil, fmt.Errorf("failed to get module info: %w", err)
	}

	// Find the specific module
	var targetModule *dto.Module
	for _, mod := range modules {
		if mod.Path == modulePath {
			targetModule = &mod
			break
		}
	}

	if targetModule == nil {
		return nil, fmt.Errorf("module %s not found", modulePath)
	}

	// Get module directory
	moduleDir := filepath.Join("vendor", targetModule.Path)
	if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
		// Try go list to get module info
		cmd := exec.Command("go", "list", "-m", "-json", modulePath)
		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get module details: %w", err)
		}

		if err := json.Unmarshal(output, targetModule); err != nil {
			return nil, fmt.Errorf("failed to parse module details: %w", err)
		}
	} else {
		// Get size
		var size int64
		err := filepath.Walk(moduleDir, func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				size += info.Size()
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to calculate module size: %w", err)
		}
		targetModule.Size = size

		// Get last updated time
		info, err := os.Stat(moduleDir)
		if err != nil {
			return nil, fmt.Errorf("failed to get module info: %w", err)
		}
		targetModule.LastUpdated = info.ModTime()

		// Try to find license
		licensePaths := []string{
			filepath.Join(moduleDir, "LICENSE"),
			filepath.Join(moduleDir, "LICENSE.md"),
			filepath.Join(moduleDir, "LICENSE.txt"),
		}
		for _, licensePath := range licensePaths {
			if _, err := os.Stat(licensePath); err == nil {
				targetModule.License = "Found"
				break
			}
		}
	}

	return targetModule, nil
}

func (m *ModuleService) ListModules(sortBy, filter string) ([]dto.Module, error) {
	// Get all modules
	modules, err := m.GetGoModules("")
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}

	// Apply filter
	if filter != "" {
		filtered := make([]dto.Module, 0)
		for _, mod := range modules {
			if m.matchesFilter(mod, filter) {
				filtered = append(filtered, mod)
			}
		}
		modules = filtered
	}

	// Sort modules
	switch sortBy {
	case "name":
		sort.Slice(modules, func(i, j int) bool {
			return modules[i].Path < modules[j].Path
		})
	case "date":
		sort.Slice(modules, func(i, j int) bool {
			return modules[i].LastUpdated.After(modules[j].LastUpdated)
		})
	case "size":
		sort.Slice(modules, func(i, j int) bool {
			return modules[i].Size > modules[j].Size
		})
	case "":
		// No sorting needed
	default:
		return nil, fmt.Errorf("invalid sort field: %s", sortBy)
	}

	return modules, nil
}

func (m *ModuleService) matchesFilter(mod dto.Module, filter string) bool {
	parts := strings.SplitN(filter, ":", 2)
	if len(parts) != 2 {
		return false
	}

	key := parts[0]
	value := parts[1]

	switch key {
	case "license":
		return mod.License == value
	case "size":
		size, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return false
		}
		return mod.Size >= size
	case "updated":
		days, err := strconv.Atoi(value)
		if err != nil {
			return false
		}
		return time.Since(mod.LastUpdated) <= time.Duration(days)*24*time.Hour
	default:
		return false
	}
}
