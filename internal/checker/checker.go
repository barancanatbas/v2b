package checker

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/barancanatbas/v2b/internal/git"
	"github.com/barancanatbas/v2b/internal/modules"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

const (
	NotFoundMessage = "Not-Found"
	workerPoolSize  = 10
)

var (
	logger = color.New(color.FgCyan).Add(color.Bold)

	green = color.New(color.FgGreen)
	blue  = color.New(color.FgBlue)
	red   = color.New(color.FgRed)
)

type Checker struct {
	showErrors    bool
	special       bool
	prefix        string
	moduleService modules.ModuleServiceInterface
	gitService    git.GitServiceInterface
}

func NewChecker(moduleService modules.ModuleServiceInterface, gitService git.GitServiceInterface, showErrors, special bool, prefix string) *Checker {
	return &Checker{
		showErrors:    showErrors,
		special:       special,
		prefix:        prefix,
		moduleService: moduleService,
		gitService:    gitService,
	}
}

func (c *Checker) FetchAndDisplayModules() error {
	modules, err := c.moduleService.GetGoModules(c.prefix)
	if err != nil {
		return fmt.Errorf("failed to fetch Go modules: %w", err)
	}

	logger.Println("Fetching special names for each module...")

	wg := sync.WaitGroup{}
	results := sync.Map{}
	versionResults := sync.Map{}
	errResults := sync.Map{}
	moduleChannel := make(chan dto.Module, len(modules))

	for i := 0; i < workerPoolSize; i++ {
		go c.worker(moduleChannel, &wg, &results, &versionResults, &errResults)
	}

	for _, mod := range modules {
		wg.Add(1)
		moduleChannel <- mod
	}

	wg.Wait()
	close(moduleChannel)

	c.generateTable(&results, &versionResults, &errResults)

	return nil
}

func (c *Checker) worker(moduleChannel <-chan dto.Module, wg *sync.WaitGroup, results, versionResults, errResults *sync.Map) {
	for mod := range moduleChannel {
		c.processModule(&mod, results, versionResults, errResults, wg)
	}
}

func (c *Checker) generateTable(specialBranchResults, versionResults, errResults *sync.Map) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Module Path", "Version", "Branch"})
	table.SetBorder(true)
	table.SetCenterSeparator(" ")

	c.appendErrResults(specialBranchResults, errResults)
	c.appendSpecialBranchResults(specialBranchResults, versionResults)
	c.appendVersionResults(specialBranchResults, versionResults)

	specialBranchResults.Range(func(key, value interface{}) bool {
		result := value.([]string)
		table.Append(result)
		return true
	})

	table.Render()
}

func (c *Checker) appendVersionResults(result, versionResults *sync.Map) {
	if c.special {
		return
	}

	versionResults.Range(func(key, value interface{}) bool {
		result.Store(key, value)
		return true
	})
}

func (c *Checker) appendSpecialBranchResults(result, versionResults *sync.Map) {
	result.Range(func(key, value interface{}) bool {
		versionResults.Store(key, value)
		return true
	})
}

func (c *Checker) appendErrResults(result, errResults *sync.Map) {
	if !c.showErrors {
		return
	}

	errResults.Range(func(key, value interface{}) bool {
		result.Store(key, value)
		return true
	})
}

func (c *Checker) processModule(mod *dto.Module, specialBranchResults, versionResults, errResults *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()

	if !strings.Contains(mod.Path, ".") || mod.Version == "" {
		return
	}

	parts := strings.Split(mod.Version, "-")
	mod.CommitHash = parts[len(parts)-1]

	branch, err := c.gitService.GetBranch(mod)
	if err != nil {
		errResults.Store(mod.Path, []string{mod.Path, mod.Version, red.Sprintf(NotFoundMessage)})
		return
	}

	if c.gitService.IsSpecialBranch(branch) {
		specialBranchResults.Store(mod.Path, []string{mod.Path, mod.Version, green.Sprintf(branch)})
	} else {
		versionResults.Store(mod.Path, []string{mod.Path, mod.Version, blue.Sprintf(branch)})
	}
}
