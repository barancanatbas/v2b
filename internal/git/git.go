package git

import (
	"bytes"
	"fmt"
	"github.com/barancanatbas/v2b/internal/dto"
	"net/url"
	"os/exec"
	"strings"
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
