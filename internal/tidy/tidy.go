package tidy

import (
	"fmt"
	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/barancanatbas/v2b/internal/git"
	"github.com/barancanatbas/v2b/internal/modules"
	"sync"
)

const (
	WorkerPoolSize = 10
)

type TidyService struct {
	moduleService *modules.ModuleService
	gitService    *git.GitService
	prefix        string
}

func NewTidyService(moduleService *modules.ModuleService, gitService *git.GitService, prefix string) *TidyService {
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
