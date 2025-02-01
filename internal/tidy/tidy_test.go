package tidy

import (
	"testing"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/barancanatbas/v2b/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestTidyService_UpdateModuleBranch(t *testing.T) {
	tests := []struct {
		name        string
		modulePath  string
		branchName  string
		expectError bool
		setupMocks  func(*mocks.MockModuleService, *mocks.MockGitService)
	}{
		{
			name:       "successful module update",
			modulePath: "github.com/test/repo",
			branchName: "feature/test",
			setupMocks: func(mm *mocks.MockModuleService, mg *mocks.MockGitService) {
				mm.GetGoModulesFunc = func(prefix string) ([]dto.Module, error) {
					return []dto.Module{
						{
							Path: "github.com/test/repo",
						},
					}, nil
				}
				mm.TidyForModuleFunc = func(module dto.Module) error {
					assert.Equal(t, "github.com/test/repo", module.Path)
					assert.Equal(t, "feature/test", *module.Branch)
					return nil
				}
			},
		},
		{
			name:        "module not found",
			modulePath:  "github.com/nonexistent/repo",
			branchName:  "main",
			expectError: true,
			setupMocks: func(mm *mocks.MockModuleService, mg *mocks.MockGitService) {
				mm.GetGoModulesFunc = func(prefix string) ([]dto.Module, error) {
					return []dto.Module{}, nil
				}
			},
		},
		{
			name:        "tidy error",
			modulePath:  "github.com/test/repo",
			branchName:  "feature/test",
			expectError: true,
			setupMocks: func(mm *mocks.MockModuleService, mg *mocks.MockGitService) {
				mm.GetGoModulesFunc = func(prefix string) ([]dto.Module, error) {
					return []dto.Module{
						{
							Path: "github.com/test/repo",
						},
					}, nil
				}
				mm.TidyForModuleFunc = func(module dto.Module) error {
					return assert.AnError
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

			service := NewTidyService(mockModule, mockGit, "")
			err := service.UpdateModuleBranch(tt.modulePath, tt.branchName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
