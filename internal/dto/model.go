package dto

type Module struct {
	ID         int     `json:"id"`
	Path       string  `json:"path"`
	Version    string  `json:"version"`
	CommitHash string  `json:"commit_hash"`
	Branch     *string `json:"branch,omitempty"`
}
