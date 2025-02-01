package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

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

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run 'go mod tidy': %v", err)
	}

	outputByte, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run 'go mod tidy': %v", err)
	}

	fmt.Println(string(outputByte))

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
