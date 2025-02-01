package git

import "github.com/barancanatbas/v2b/internal/dto"

type GitServiceInterface interface {
	GetBranch(mod *dto.Module) (string, error)
	IsSpecialBranch(branch string) bool
}
