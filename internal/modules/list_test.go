package modules

import (
	"testing"
	"time"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestModuleService_GetModuleDetails(t *testing.T) {
	tests := []struct {
		name        string
		modulePath  string
		expectError bool
		expected    *dto.Module
	}{
		{
			name:       "valid module",
			modulePath: "github.com/test/mod1",
			expected: &dto.Module{
				Path:        "github.com/test/mod1",
				Version:     "v1.0.0",
				LastUpdated: time.Now().Truncate(24 * time.Hour),
				Size:        1024,
				UsageCount:  5,
				License:     "MIT",
			},
		},
		{
			name:        "invalid module path",
			modulePath:  "invalid/path",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			result, err := service.GetModuleDetails(tt.modulePath)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.Path, result.Path)
				assert.Equal(t, tt.expected.Version, result.Version)
				assert.Equal(t, tt.expected.LastUpdated.Day(), result.LastUpdated.Day())
				assert.Equal(t, tt.expected.Size, result.Size)
				assert.Equal(t, tt.expected.UsageCount, result.UsageCount)
				assert.Equal(t, tt.expected.License, result.License)
			}
		})
	}
}

func TestModuleService_ListModules(t *testing.T) {
	tests := []struct {
		name        string
		sortBy      string
		filter      string
		expectError bool
		expectedLen int
	}{
		{
			name:        "list all modules",
			sortBy:      "name",
			expectedLen: 2,
		},
		{
			name:        "sort by date",
			sortBy:      "date",
			expectedLen: 2,
		},
		{
			name:        "filter by license",
			filter:      "license:MIT",
			expectedLen: 1,
		},
		{
			name:        "invalid sort field",
			sortBy:      "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			modules, err := service.ListModules(tt.sortBy, tt.filter)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, modules)
			} else {
				assert.NoError(t, err)
				assert.Len(t, modules, tt.expectedLen)

				if tt.sortBy == "date" {
					// Check if sorted by date
					for i := 1; i < len(modules); i++ {
						assert.True(t, modules[i].LastUpdated.After(modules[i-1].LastUpdated) ||
							modules[i].LastUpdated.Equal(modules[i-1].LastUpdated))
					}
				}

				if tt.filter != "" {
					// Check if filter is applied
					for _, mod := range modules {
						if tt.filter == "license:MIT" {
							assert.Equal(t, "MIT", mod.License)
						}
					}
				}
			}
		})
	}
}
