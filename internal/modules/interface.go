package modules

import "github.com/barancanatbas/v2b/internal/dto"

type ModuleServiceInterface interface {
	GetGoModules(prefix string) ([]dto.Module, error)
	TidyForModule(module dto.Module) error
}
