package checker

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/barancanatbas/v2b/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestChecker_FetchAndDisplayModules(t *testing.T) {
	tests := []struct {
		name          string
		modules       []dto.Module
		moduleErr     error
		getBranchFunc func(mod *dto.Module) (string, error)
		isSpecialFunc func(branch string) bool
		showErrors    bool
		special       bool
		prefix        string
		expectedError bool
	}{
		{
			name: "successful fetch and display",
			modules: []dto.Module{
				{Path: "github.com/test/mod1", Version: "v1.0.0-abc123"},
				{Path: "github.com/test/mod2", Version: "v2.0.0-def456"},
			},
			getBranchFunc: func(mod *dto.Module) (string, error) {
				return "main", nil
			},
			isSpecialFunc: func(branch string) bool {
				return false
			},
			showErrors: true,
			special:    false,
			prefix:     "",
		},
		{
			name:          "module service error",
			moduleErr:     errors.New("failed to fetch modules"),
			expectedError: true,
		},
		{
			name: "git service error with show errors",
			modules: []dto.Module{
				{Path: "github.com/test/mod1", Version: "v1.0.0-abc123"},
			},
			getBranchFunc: func(mod *dto.Module) (string, error) {
				return "", errors.New("git error")
			},
			showErrors: true,
		},
		{
			name: "special branches only",
			modules: []dto.Module{
				{Path: "github.com/test/mod1", Version: "v1.0.0-abc123"},
				{Path: "github.com/test/mod2", Version: "v2.0.0-def456"},
			},
			getBranchFunc: func(mod *dto.Module) (string, error) {
				return "feature/test", nil
			},
			isSpecialFunc: func(branch string) bool {
				return strings.HasPrefix(branch, "feature/")
			},
			special: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModule := &mocks.MockModuleService{
				GetGoModulesFunc: func(prefix string) ([]dto.Module, error) {
					return tt.modules, tt.moduleErr
				},
			}

			mockGit := &mocks.MockGitService{
				GetBranchFunc:       tt.getBranchFunc,
				IsSpecialBranchFunc: tt.isSpecialFunc,
			}

			checker := NewChecker(mockModule, mockGit, tt.showErrors, tt.special, tt.prefix)
			err := checker.FetchAndDisplayModules()

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestChecker_ProcessModule(t *testing.T) {
	tests := []struct {
		name          string
		module        dto.Module
		getBranchFunc func(mod *dto.Module) (string, error)
		isSpecialFunc func(branch string) bool
		expectError   bool
		expectSpecial bool
	}{
		{
			name: "process regular branch",
			module: dto.Module{
				Path:    "github.com/test/mod",
				Version: "v1.0.0-abc123",
			},
			getBranchFunc: func(mod *dto.Module) (string, error) {
				return "main", nil
			},
			isSpecialFunc: func(branch string) bool {
				return false
			},
			expectError:   false,
			expectSpecial: false,
		},
		{
			name: "process special branch",
			module: dto.Module{
				Path:    "github.com/test/mod",
				Version: "v1.0.0-abc123",
			},
			getBranchFunc: func(mod *dto.Module) (string, error) {
				return "feature/test", nil
			},
			isSpecialFunc: func(branch string) bool {
				return true
			},
			expectError:   false,
			expectSpecial: true,
		},
		{
			name: "process error case",
			module: dto.Module{
				Path:    "github.com/test/mod",
				Version: "v1.0.0-abc123",
			},
			getBranchFunc: func(mod *dto.Module) (string, error) {
				return "", errors.New("git error")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGit := &mocks.MockGitService{
				GetBranchFunc:       tt.getBranchFunc,
				IsSpecialBranchFunc: tt.isSpecialFunc,
			}

			checker := NewChecker(nil, mockGit, true, false, "")

			var wg sync.WaitGroup
			specialBranchResults := &sync.Map{}
			versionResults := &sync.Map{}
			errResults := &sync.Map{}

			wg.Add(1)
			checker.processModule(&tt.module, specialBranchResults, versionResults, errResults, &wg)

			if tt.expectError {
				_, exists := errResults.Load(tt.module.Path)
				assert.True(t, exists)
			} else if tt.expectSpecial {
				_, exists := specialBranchResults.Load(tt.module.Path)
				assert.True(t, exists)
			} else {
				_, exists := versionResults.Load(tt.module.Path)
				assert.True(t, exists)
			}
		})
	}
}
