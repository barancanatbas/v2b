package mocks

import (
	"github.com/barancanatbas/v2b/internal/dto"
)

type MockModuleService struct {
	GetGoModulesFunc  func(prefix string) ([]dto.Module, error)
	TidyForModuleFunc func(module dto.Module) error
}

func (m *MockModuleService) GetGoModules(prefix string) ([]dto.Module, error) {
	if m.GetGoModulesFunc != nil {
		return m.GetGoModulesFunc(prefix)
	}
	return nil, nil
}

func (m *MockModuleService) TidyForModule(module dto.Module) error {
	if m.TidyForModuleFunc != nil {
		return m.TidyForModuleFunc(module)
	}
	return nil
}
