package git

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/barancanatbas/v2b/internal/dto"
)

const (
	DefaultGitScheme = "https"
	GitCommand       = "git"
	GitLsRemoteFlag  = "ls-remote"
	RefHeadsPrefix   = "refs/heads/"
)

type GitService struct {
}

func NewGitService() *GitService {
	return &GitService{}
}

func (g *GitService) ensureHTTPS(repoURL string) (string, error) {
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL '%s': %v", repoURL, err)
	}

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = DefaultGitScheme
	}

	return parsedURL.String(), nil
}

func (g *GitService) GetBranch(module *dto.Module) (string, error) {
	if module.CommitHash == "" {
		return "", fmt.Errorf("commit hash is empty")
	}

	// Check if path is a local directory
	if info, err := os.Stat(module.Path); err == nil && info.IsDir() {
		// For local repositories, use git directly
		cmd := exec.Command(GitCommand, "rev-parse", "--abbrev-ref", "HEAD")
		cmd.Dir = module.Path
		out, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("failed to get branch: %w", err)
		}

		return strings.TrimSpace(string(out)), nil
	}

	// For remote repositories, use ls-remote
	repoURL, err := g.ensureHTTPS(module.Path)
	if err != nil {
		return "", fmt.Errorf("failed to ensure HTTPS: %w", err)
	}

	cmd := exec.Command(GitCommand, GitLsRemoteFlag, repoURL)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run 'git ls-remote': %v, stderr: %s", err, errBuf.String())
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, module.CommitHash) {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				return strings.TrimPrefix(parts[1], RefHeadsPrefix), nil
			}
		}
	}

	return "", fmt.Errorf("branch not found for commit hash: %s", module.CommitHash)
}

func (g *GitService) IsSpecialBranch(branch string) bool {
	standardBranches := []string{
		"main",
		"master",
		"develop",
	}

	for _, std := range standardBranches {
		if branch == std {
			return false
		}
	}

	specialPrefixes := []string{
		"feature/",
		"release/",
		"hotfix/",
		"bugfix/",
		"tags/",
		"pull/",
	}

	for _, prefix := range specialPrefixes {
		if strings.HasPrefix(branch, prefix) {
			return true
		}
	}

	return false
}
