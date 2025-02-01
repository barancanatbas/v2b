package tidy

import (
	"testing"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/barancanatbas/v2b/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestTidyService_UpdateModuleByID(t *testing.T) {
	tests := []struct {
		name        string
		moduleID    int
		expectError bool
		setupMocks  func(*mocks.MockModuleService, *mocks.MockGitService)
	}{
		{
			name:     "successful module update by ID",
			moduleID: 1,
			setupMocks: func(mm *mocks.MockModuleService, mg *mocks.MockGitService) {
				mm.GetGoModulesFunc = func(prefix string) ([]dto.Module, error) {
					return []dto.Module{
						{
							ID:   1,
							Path: "github.com/test/repo",
						},
					}, nil
				}
				mm.TidyForModuleFunc = func(module dto.Module) error {
					assert.Equal(t, 1, module.ID)
					assert.Equal(t, "github.com/test/repo", module.Path)
					return nil
				}
			},
		},
		{
			name:        "module ID not found",
			moduleID:    999,
			expectError: true,
			setupMocks: func(mm *mocks.MockModuleService, mg *mocks.MockGitService) {
				mm.GetGoModulesFunc = func(prefix string) ([]dto.Module, error) {
					return []dto.Module{
						{
							ID:   1,
							Path: "github.com/test/repo",
						},
					}, nil
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
			err := service.UpdateModuleByID(tt.moduleID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
