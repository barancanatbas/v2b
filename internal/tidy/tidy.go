package tidy

import (
	"fmt"
	"sync"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/barancanatbas/v2b/internal/git"
	"github.com/barancanatbas/v2b/internal/modules"
)

const (
	WorkerPoolSize = 10
)

type TidyService struct {
	moduleService modules.ModuleServiceInterface
	gitService    git.GitServiceInterface
	prefix        string
}

func NewTidyService(moduleService modules.ModuleServiceInterface, gitService git.GitServiceInterface, prefix string) *TidyService {
	return &TidyService{
		moduleService: moduleService,
		gitService:    gitService,
		prefix:        prefix,
	}
}

func (t *TidyService) ModTidy() {
	modules, err := t.moduleService.GetGoModules(t.prefix)
	if err != nil {
		fmt.Println("Failed to fetch Go modules:", err)
		return
	}

	var wg sync.WaitGroup
	modChannel := make(chan dto.Module, len(modules))

	for i := 0; i < WorkerPoolSize; i++ {
		go t.worker(modChannel, &wg)
	}

	for _, mod := range modules {
		wg.Add(1)
		modChannel <- mod
	}

	wg.Wait()
	close(modChannel)
}

func (t *TidyService) worker(modChannel <-chan dto.Module, wg *sync.WaitGroup) {
	for mod := range modChannel {
		t.processModule(&mod, wg)
	}
}

func (t *TidyService) processModule(mod *dto.Module, wg *sync.WaitGroup) {
	defer wg.Done()

	branch, err := t.gitService.GetBranch(mod)
	if err != nil {
		fmt.Println("Failed to get branch:", err)
		return
	}

	if !t.gitService.IsSpecialBranch(branch) {
		fmt.Println("Skipping special branch:", branch)
		return
	}

	mod.Branch = &branch

	_ = t.moduleService.TidyForModule(*mod)
}

func (t *TidyService) UpdateModuleBranch(modulePath, branchName string) error {
	// Get all modules
	modules, err := t.moduleService.GetGoModules("")
	if err != nil {
		return fmt.Errorf("failed to get modules: %w", err)
	}

	// Find the target module
	var targetModule *dto.Module
	for _, mod := range modules {
		if mod.Path == modulePath {
			targetModule = &mod
			break
		}
	}

	if targetModule == nil {
		return fmt.Errorf("module %s not found", modulePath)
	}

	// Update the module's branch
	targetModule.Branch = &branchName

	// Run go get with the specified branch
	if err := t.moduleService.TidyForModule(*targetModule); err != nil {
		return fmt.Errorf("failed to update module %s to branch %s: %w", modulePath, branchName, err)
	}

	return nil
}

func (t *TidyService) UpdateModuleByID(moduleID int) error {
	// Get all modules
	modules, err := t.moduleService.GetGoModules("")
	if err != nil {
		return fmt.Errorf("failed to get modules: %w", err)
	}

	// Find the target module
	var targetModule *dto.Module
	for _, mod := range modules {
		if mod.ID == moduleID {
			targetModule = &mod
			break
		}
	}

	if targetModule == nil {
		return fmt.Errorf("module with ID %d not found", moduleID)
	}

	// Run go get for the module
	if err := t.moduleService.TidyForModule(*targetModule); err != nil {
		return fmt.Errorf("failed to update module %s (ID: %d): %w", targetModule.Path, moduleID, err)
	}

	return nil
}
