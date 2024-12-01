package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/barancanatbas/v2b/internal/dto"
	"os/exec"
	"strings"
)

type ModuleService struct {
}

func NewModule() *ModuleService {
	return &ModuleService{}
}

func (m *ModuleService) GetGoModules(prefix string) ([]dto.Module, error) {
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
