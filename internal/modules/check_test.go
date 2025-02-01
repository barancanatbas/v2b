package modules

import (
	"testing"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestModuleService_CheckOutdated(t *testing.T) {
	tests := []struct {
		name        string
		modules     []dto.Module
		expectError bool
		validate    func([]dto.Module) bool
	}{
		{
			name: "detect outdated modules",
			modules: []dto.Module{
				{
					Path:    "github.com/test/mod1",
					Version: "v1.0.0",
				},
				{
					Path:    "github.com/test/mod2",
					Version: "v2.0.0",
				},
			},
			validate: func(modules []dto.Module) bool {
				return len(modules) > 0 && modules[0].Version != "latest"
			},
		},
		{
			name:    "all modules up to date",
			modules: []dto.Module{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			outdated, err := service.CheckOutdated()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					assert.True(t, tt.validate(outdated))
				}
			}
		})
	}
}

func TestModuleService_CheckSecurity(t *testing.T) {
	tests := []struct {
		name                    string
		modules                 []dto.Module
		expectError             bool
		expectedVulnerabilities int
	}{
		{
			name: "detect security vulnerabilities",
			modules: []dto.Module{
				{
					Path:    "github.com/test/vulnerable-mod",
					Version: "v1.0.0",
				},
			},
			expectedVulnerabilities: 1,
		},
		{
			name: "no vulnerabilities found",
			modules: []dto.Module{
				{
					Path:    "github.com/test/safe-mod",
					Version: "v1.0.0",
				},
			},
			expectedVulnerabilities: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			vulns, err := service.CheckSecurity()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, vulns, tt.expectedVulnerabilities)
			}
		})
	}
}

func TestModuleService_CheckDeprecated(t *testing.T) {
	tests := []struct {
		name               string
		modules            []dto.Module
		expectError        bool
		expectedDeprecated int
	}{
		{
			name: "detect deprecated modules",
			modules: []dto.Module{
				{
					Path:    "github.com/test/deprecated-mod",
					Version: "v1.0.0",
				},
			},
			expectedDeprecated: 1,
		},
		{
			name: "no deprecated modules",
			modules: []dto.Module{
				{
					Path:    "github.com/test/active-mod",
					Version: "v1.0.0",
				},
			},
			expectedDeprecated: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			deprecated, err := service.CheckDeprecated()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, deprecated, tt.expectedDeprecated)
			}
		})
	}
}

func TestModuleService_CheckCompatibility(t *testing.T) {
	tests := []struct {
		name                 string
		modules              []dto.Module
		expectError          bool
		expectedIncompatible int
	}{
		{
			name: "detect incompatible versions",
			modules: []dto.Module{
				{
					Path:    "github.com/test/mod1",
					Version: "v1.0.0",
				},
				{
					Path:    "github.com/test/mod2",
					Version: "v2.0.0",
				},
			},
			expectedIncompatible: 1,
		},
		{
			name: "all modules compatible",
			modules: []dto.Module{
				{
					Path:    "github.com/test/mod1",
					Version: "v1.0.0",
				},
			},
			expectedIncompatible: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			incompatible, err := service.CheckCompatibility()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, incompatible, tt.expectedIncompatible)
			}
		})
	}
}

func TestModuleService_Cleanup(t *testing.T) {
	tests := []struct {
		name            string
		modules         []dto.Module
		expectError     bool
		expectedRemoved int
	}{
		{
			name: "cleanup unused modules",
			modules: []dto.Module{
				{
					Path:       "github.com/test/unused-mod",
					Version:    "v1.0.0",
					UsageCount: 0,
				},
			},
			expectedRemoved: 1,
		},
		{
			name: "no modules to cleanup",
			modules: []dto.Module{
				{
					Path:       "github.com/test/used-mod",
					Version:    "v1.0.0",
					UsageCount: 5,
				},
			},
			expectedRemoved: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			removed, err := service.Cleanup()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, removed, tt.expectedRemoved)
			}
		})
	}
}
