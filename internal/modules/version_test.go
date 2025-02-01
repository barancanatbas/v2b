package modules

import (
	"testing"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestModuleService_UpgradeVersion(t *testing.T) {
	tests := []struct {
		name        string
		module      dto.Module
		newVersion  string
		expectError bool
		validate    func(dto.Module) bool
	}{
		{
			name: "successful version upgrade",
			module: dto.Module{
				Path:    "github.com/test/mod1",
				Version: "v1.0.0",
			},
			newVersion: "v1.1.0",
			validate: func(m dto.Module) bool {
				return m.Version == "v1.1.0"
			},
		},
		{
			name: "upgrade to invalid version",
			module: dto.Module{
				Path:    "github.com/test/mod1",
				Version: "v1.0.0",
			},
			newVersion:  "invalid",
			expectError: true,
		},
		{
			name: "upgrade to same version",
			module: dto.Module{
				Path:    "github.com/test/mod1",
				Version: "v1.0.0",
			},
			newVersion:  "v1.0.0",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			upgraded, err := service.UpgradeVersion(tt.module, tt.newVersion)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					assert.True(t, tt.validate(upgraded))
				}
			}
		})
	}
}

func TestModuleService_RollbackVersion(t *testing.T) {
	tests := []struct {
		name        string
		module      dto.Module
		expectError bool
		validate    func(dto.Module) bool
	}{
		{
			name: "successful version rollback",
			module: dto.Module{
				Path:    "github.com/test/mod1",
				Version: "v1.1.0",
			},
			validate: func(m dto.Module) bool {
				return m.Version == "v1.0.0"
			},
		},
		{
			name: "rollback first version",
			module: dto.Module{
				Path:    "github.com/test/mod1",
				Version: "v1.0.0",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			rolledBack, err := service.RollbackVersion(tt.module)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					assert.True(t, tt.validate(rolledBack))
				}
			}
		})
	}
}

func TestModuleService_ResolveVersionConflict(t *testing.T) {
	tests := []struct {
		name        string
		modules     []dto.Module
		expectError bool
		validate    func([]dto.Module) bool
	}{
		{
			name: "resolve version conflict",
			modules: []dto.Module{
				{
					Path:    "github.com/test/mod1",
					Version: "v1.0.0",
				},
				{
					Path:    "github.com/test/mod2",
					Version: "v2.0.0",
					Dependencies: []dto.Module{
						{
							Path:    "github.com/test/mod1",
							Version: "v1.1.0",
						},
					},
				},
			},
			validate: func(modules []dto.Module) bool {
				return len(modules) > 0 && modules[0].Version == "v1.1.0"
			},
		},
		{
			name: "no conflicts to resolve",
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
				return len(modules) == 2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			resolved, err := service.ResolveVersionConflict(tt.modules)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					assert.True(t, tt.validate(resolved))
				}
			}
		})
	}
}
