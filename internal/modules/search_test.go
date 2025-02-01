package modules

import (
	"strings"
	"testing"
	"time"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestModuleService_Search(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		searchType  string
		expectError bool
		expectedLen int
		validate    func([]dto.Module) bool
	}{
		{
			name:        "search by content",
			query:       "logger",
			searchType:  "content",
			expectedLen: 1,
			validate: func(modules []dto.Module) bool {
				return len(modules) > 0
			},
		},
		{
			name:        "search by version",
			query:       "v1.0",
			searchType:  "version",
			expectedLen: 2,
			validate: func(modules []dto.Module) bool {
				for _, mod := range modules {
					if !strings.Contains(mod.Version, "v1.0") {
						return false
					}
				}
				return true
			},
		},
		{
			name:        "search by date range",
			query:       "7d", // last 7 days
			searchType:  "date",
			expectedLen: 1,
			validate: func(modules []dto.Module) bool {
				sevenDaysAgo := time.Now().AddDate(0, 0, -7)
				for _, mod := range modules {
					if mod.LastUpdated.Before(sevenDaysAgo) {
						return false
					}
				}
				return true
			},
		},
		{
			name:        "invalid search type",
			query:       "test",
			searchType:  "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			results, err := service.Search(tt.query, tt.searchType)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, results)
			} else {
				assert.NoError(t, err)
				assert.Len(t, results, tt.expectedLen)
				if tt.validate != nil {
					assert.True(t, tt.validate(results))
				}
			}
		})
	}
}

func TestModuleService_SearchContent(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		expectError bool
		expectedLen int
	}{
		{
			name:        "find existing content",
			query:       "fmt.Println",
			expectedLen: 1,
		},
		{
			name:        "find with case insensitive",
			query:       "LOGGER",
			expectedLen: 1,
		},
		{
			name:        "no results",
			query:       "nonexistentcontent",
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			results, err := service.searchContent(tt.query)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, results)
			} else {
				assert.NoError(t, err)
				assert.Len(t, results, tt.expectedLen)
			}
		})
	}
}

func TestModuleService_SearchVersion(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		expectError bool
		expectedLen int
	}{
		{
			name:        "exact version match",
			version:     "v1.0.0",
			expectedLen: 1,
		},
		{
			name:        "version prefix match",
			version:     "v1",
			expectedLen: 2,
		},
		{
			name:        "no version match",
			version:     "v9.9.9",
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			results, err := service.searchVersion(tt.version)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, results)
			} else {
				assert.NoError(t, err)
				assert.Len(t, results, tt.expectedLen)
			}
		})
	}
}

func TestModuleService_SearchDate(t *testing.T) {
	tests := []struct {
		name        string
		dateQuery   string
		expectError bool
		expectedLen int
	}{
		{
			name:        "last 7 days",
			dateQuery:   "7d",
			expectedLen: 1,
		},
		{
			name:        "last month",
			dateQuery:   "30d",
			expectedLen: 2,
		},
		{
			name:        "invalid date format",
			dateQuery:   "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			results, err := service.searchDate(tt.dateQuery)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, results)
			} else {
				assert.NoError(t, err)
				assert.Len(t, results, tt.expectedLen)
			}
		})
	}
}
