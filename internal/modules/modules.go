package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/barancanatbas/v2b/internal/git"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"os"
	"os/exec"
	"strings"
	"sync"
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

type ModuleService struct {
	printErrors bool
	special     bool
	prefix      string
}

func NewModule(printErrors bool, special bool, prefix string) *ModuleService {
	return &ModuleService{
		printErrors: printErrors,
		special:     special,
		prefix:      prefix,
	}
}

func (m *ModuleService) FetchAndDisplayModules() error {
	modules, err := m.getGoModules()
	if err != nil {
		return fmt.Errorf("failed to fetch Go modules: %w", err)
	}

	logger.Println("Fetching branch names for each module...")

	wg := sync.WaitGroup{}
	results := sync.Map{}
	versionResults := sync.Map{}
	errResults := sync.Map{}
	moduleChannel := make(chan Module, len(modules))

	for i := 0; i < workerPoolSize; i++ {
		go m.worker(moduleChannel, &wg, &results, &versionResults, &errResults)
	}

	for _, mod := range modules {
		wg.Add(1)
		moduleChannel <- mod
	}

	wg.Wait()
	close(moduleChannel)

	m.generateTable(&results, &versionResults, &errResults)

	return nil
}

func (m *ModuleService) worker(moduleChannel <-chan Module, wg *sync.WaitGroup, results, versionResults, errResults *sync.Map) {
	for mod := range moduleChannel {
		m.processModule(mod, results, versionResults, errResults, wg)
	}
}

func (m *ModuleService) generateTable(specialBranchResults, versionResults, errResults *sync.Map) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Module Path", "Version", "Branch"})
	table.SetBorder(true)
	table.SetCenterSeparator(" ")

	m.appendErrResults(specialBranchResults, errResults)
	m.appendSpecialBranchResults(specialBranchResults, versionResults)
	m.appendVersionResults(specialBranchResults, versionResults)

	specialBranchResults.Range(func(key, value interface{}) bool {
		result := value.([]string)
		table.Append(result)
		return true
	})

	table.Render()
}

func (m *ModuleService) appendVersionResults(result, versionResults *sync.Map) {
	if m.special {
		return
	}

	versionResults.Range(func(key, value interface{}) bool {
		result.Store(key, value)
		return true
	})
}

func (m *ModuleService) appendSpecialBranchResults(result, versionResults *sync.Map) {
	result.Range(func(key, value interface{}) bool {
		versionResults.Store(key, value)
		return true
	})
}

func (m *ModuleService) appendErrResults(result, errResults *sync.Map) {
	if !m.printErrors {
		return
	}

	errResults.Range(func(key, value interface{}) bool {
		result.Store(key, value)
		return true
	})
}

func (m *ModuleService) getGoModules() ([]Module, error) {
	cmd := exec.Command("go", "list", "-m", "-json", "all")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run 'go list': %v", err)
	}

	var modules []Module
	decoder := json.NewDecoder(&out)
	for decoder.More() {
		var mod Module
		if err := decoder.Decode(&mod); err != nil {
			return nil, fmt.Errorf("failed to decode module data: %v", err)
		}

		if m.prefix != "" {
			if !strings.HasPrefix(mod.Path, m.prefix) {
				continue
			}
		}

		modules = append(modules, mod)
	}

	return modules, nil
}

func (m *ModuleService) processModule(mod Module, specialBranchResults, versionResults, errResults *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()

	if !strings.Contains(mod.Path, ".") || mod.Version == "" {
		return
	}

	parts := strings.Split(mod.Version, "-")
	commitHash := parts[len(parts)-1]

	branch, err := git.GetBranchFromCommit(mod.Path, commitHash)
	if err != nil {
		errResults.Store(mod.Path, []string{mod.Path, mod.Version, red.Sprintf(NotFoundMessage)})
		return
	}

	if m.isSpecialBranch(branch) {
		specialBranchResults.Store(mod.Path, []string{mod.Path, mod.Version, green.Sprintf(branch)})
	} else {
		versionResults.Store(mod.Path, []string{mod.Path, mod.Version, blue.Sprintf(branch)})
	}
}

func (m *ModuleService) isSpecialBranch(branch string) bool {
	nonSpecialPrefixes := []string{
		"refs/tags/",
		"refs/pull/",
		"refs/heads/",
		"refs/remotes/",
		"refs/merge-requests/",
		"refs/stash",
		"HEAD",
	}

	for _, prefix := range nonSpecialPrefixes {
		if strings.HasPrefix(branch, prefix) {
			return false
		}
	}

	return true
}
