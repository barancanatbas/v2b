package git

import (
	"testing"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestGitService_IsSpecialBranch(t *testing.T) {
	tests := []struct {
		name     string
		branch   string
		expected bool
	}{
		{
			name:     "feature branch",
			branch:   "feature/test",
			expected: true,
		},
		{
			name:     "release branch",
			branch:   "release/1.0.0",
			expected: true,
		},
		{
			name:     "tag reference",
			branch:   "tags/v1.0.0",
			expected: true,
		},
		{
			name:     "pull request",
			branch:   "pull/123",
			expected: true,
		},
		{
			name:     "main branch",
			branch:   "main",
			expected: false,
		},
		{
			name:     "master branch",
			branch:   "master",
			expected: false,
		},
		{
			name:     "develop branch",
			branch:   "develop",
			expected: false,
		},
	}

	gitService := NewGitService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gitService.IsSpecialBranch(tt.branch)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGitService_GetBranch(t *testing.T) {
	tests := []struct {
		name        string
		module      *dto.Module
		expectError bool
	}{
		{
			name: "valid module with commit hash",
			module: &dto.Module{
				Path:       "github.com/test/repo",
				Version:    "v1.0.0-abc123",
				CommitHash: "abc123",
			},
			expectError: false,
		},
		{
			name: "invalid module path",
			module: &dto.Module{
				Path:       "invalid-path",
				Version:    "v1.0.0-abc123",
				CommitHash: "abc123",
			},
			expectError: true,
		},
		{
			name: "empty commit hash",
			module: &dto.Module{
				Path:       "github.com/test/repo",
				Version:    "v1.0.0",
				CommitHash: "",
			},
			expectError: true,
		},
	}

	gitService := NewGitService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			branch, err := gitService.GetBranch(tt.module)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, branch)
			} else {
				// Note: Since actual git operations depend on external state,
				// we can only verify that the function returns without error
				assert.NoError(t, err)
			}
		})
	}
}
