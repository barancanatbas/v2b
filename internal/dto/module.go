// Package dto contains data transfer objects for the application
package dto

import "time"

// Module represents a Go module with its metadata
type Module struct {
	ID           int       `json:"id"`
	Path         string    `json:"path"`
	Version      string    `json:"version"`
	CommitHash   string    `json:"commit_hash"`
	Branch       *string   `json:"branch,omitempty"`
	LastUpdated  time.Time `json:"last_updated"`
	Size         int64     `json:"size"`
	UsageCount   int       `json:"usage_count"`
	License      string    `json:"license"`
	IsPinned     bool      `json:"is_pinned"`
	Dependencies []Module  `json:"dependencies,omitempty"`
}
