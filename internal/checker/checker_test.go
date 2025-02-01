package checker

import (
	"errors"
	"sync"
	"testing"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/barancanatbas/v2b/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestChecker_FetchAndDisplayModules(t *testing.T) {
	tests := []struct {
		name        string
		showErrors  bool
		special     bool
		prefix      string
		setupMocks  func(*mocks.MockModuleService, *mocks.MockGitService)
		expectError bool
	}{
		{
			name: "successful_fetch_and_display",
			setupMocks: func(mm *mocks.MockModuleService, mg *mocks.MockGitService) {
				mm.GetGoModulesFunc = func(prefix string) ([]dto.Module, error) {
					return []dto.Module{
						{
							ID:         1,
							Path:       "github.com/test/mod1",
							Version:    "v1.0.0-abc123",
							CommitHash: "abc123",
						},
						{
							ID:         2,
							Path:       "github.com/test/mod2",
							Version:    "v2.0.0-def456",
							CommitHash: "def456",
						},
					}, nil
				}
				mg.GetBranchFunc = func(module *dto.Module) (string, error) {
					return "main", nil
				}
			},
		},
		{
			name: "module_service_error",
			setupMocks: func(mm *mocks.MockModuleService, mg *mocks.MockGitService) {
				mm.GetGoModulesFunc = func(prefix string) ([]dto.Module, error) {
					return nil, assert.AnError
				}
			},
			expectError: true,
		},
		{
			name:       "git_service_error_with_show_errors",
			showErrors: true,
			setupMocks: func(mm *mocks.MockModuleService, mg *mocks.MockGitService) {
				mm.GetGoModulesFunc = func(prefix string) ([]dto.Module, error) {
					return []dto.Module{
						{
							ID:         1,
							Path:       "github.com/test/mod1",
							Version:    "v1.0.0-abc123",
							CommitHash: "abc123",
						},
					}, nil
				}
				mg.GetBranchFunc = func(module *dto.Module) (string, error) {
					return "", assert.AnError
				}
			},
		},
		{
			name:    "special_branches_only",
			special: true,
			setupMocks: func(mm *mocks.MockModuleService, mg *mocks.MockGitService) {
				mm.GetGoModulesFunc = func(prefix string) ([]dto.Module, error) {
					return []dto.Module{
						{
							ID:         1,
							Path:       "github.com/test/mod1",
							Version:    "v1.0.0-abc123",
							CommitHash: "abc123",
						},
						{
							ID:         2,
							Path:       "github.com/test/mod2",
							Version:    "v2.0.0-def456",
							CommitHash: "def456",
						},
					}, nil
				}
				mg.GetBranchFunc = func(module *dto.Module) (string, error) {
					return "feature/test", nil
				}
				mg.IsSpecialBranchFunc = func(branch string) bool {
					return true
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModule := &mocks.MockModuleService{}
			mockGit := &mocks.MockGitService{}

			if tt.setupMocks != nil {
				tt.setupMocks(mockModule, mockGit)
			}

			checker := NewChecker(mockModule, mockGit, tt.showErrors, tt.special, tt.prefix)
			err := checker.FetchAndDisplayModules()

			if tt.expectError {
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
