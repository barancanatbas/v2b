package mocks

import (
	"github.com/barancanatbas/v2b/internal/dto"
)

type MockGitService struct {
	GetBranchFunc       func(mod *dto.Module) (string, error)
	IsSpecialBranchFunc func(branch string) bool
}

func (m *MockGitService) GetBranch(mod *dto.Module) (string, error) {
	if m.GetBranchFunc != nil {
		return m.GetBranchFunc(mod)
	}
	return "", nil
}

func (m *MockGitService) IsSpecialBranch(branch string) bool {
	if m.IsSpecialBranchFunc != nil {
		return m.IsSpecialBranchFunc(branch)
	}
	return false
}
