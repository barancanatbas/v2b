package dto

import "time"

type Module struct {
	ID          int       `json:"id"`
	Path        string    `json:"path"`
	Version     string    `json:"version"`
	CommitHash  string    `json:"commit_hash"`
	Branch      *string   `json:"branch,omitempty"`
	LastUpdated time.Time `json:"last_updated"`
	Size        int64     `json:"size"`
	UsageCount  int       `json:"usage_count"`
	License     string    `json:"license"`
}
