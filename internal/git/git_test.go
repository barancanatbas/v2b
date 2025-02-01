package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	// Create a temporary directory for test repository
	tmpDir, err := os.MkdirTemp("", "git-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize a test git repository
	if err := exec.Command("git", "init", tmpDir).Run(); err != nil {
		t.Fatal(err)
	}

	// Set git config for test repository
	cmd := exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	// Create and commit a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command("git", "commit", "-m", "test commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	// Get the commit hash
	cmd = exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = tmpDir
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	commitHash := strings.TrimSpace(string(out))

	tests := []struct {
		name        string
		module      *dto.Module
		expectError bool
	}{
		{
			name: "valid module with commit hash",
			module: &dto.Module{
				Path:       tmpDir,
				Version:    "v1.0.0-" + commitHash[:7],
				CommitHash: commitHash[:7],
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
				Path:       tmpDir,
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
				assert.NoError(t, err)
				assert.NotEmpty(t, branch)
			}
		})
	}
}
