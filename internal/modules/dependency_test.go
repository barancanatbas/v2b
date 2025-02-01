package modules

import (
	"testing"

	"github.com/barancanatbas/v2b/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestModuleService_PinModule(t *testing.T) {
	tests := []struct {
		name        string
		module      dto.Module
		expectError bool
		validate    func(dto.Module) bool
	}{
		{
			name: "successfully pin module",
			module: dto.Module{
				Path:    "github.com/test/mod1",
				Version: "v1.0.0",
			},
			validate: func(m dto.Module) bool {
				return m.IsPinned
			},
		},
		{
			name: "pin already pinned module",
			module: dto.Module{
				Path:     "github.com/test/mod1",
				Version:  "v1.0.0",
				IsPinned: true,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			pinned, err := service.PinModule(tt.module)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					assert.True(t, tt.validate(pinned))
				}
			}
		})
	}
}

func TestModuleService_UnpinModule(t *testing.T) {
	tests := []struct {
		name        string
		module      dto.Module
		expectError bool
		validate    func(dto.Module) bool
	}{
		{
			name: "successfully unpin module",
			module: dto.Module{
				Path:     "github.com/test/mod1",
				Version:  "v1.0.0",
				IsPinned: true,
			},
			validate: func(m dto.Module) bool {
				return !m.IsPinned
			},
		},
		{
			name: "unpin not pinned module",
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
			unpinned, err := service.UnpinModule(tt.module)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					assert.True(t, tt.validate(unpinned))
				}
			}
		})
	}
}

func TestModuleService_ExplainDependency(t *testing.T) {
	tests := []struct {
		name        string
		module      dto.Module
		expectError bool
		validate    func([]dto.Module) bool
	}{
		{
			name: "explain direct dependencies",
			module: dto.Module{
				Path:    "github.com/test/mod1",
				Version: "v1.0.0",
			},
			validate: func(deps []dto.Module) bool {
				return len(deps) > 0
			},
		},
		{
			name: "explain module with no dependencies",
			module: dto.Module{
				Path:    "github.com/test/mod2",
				Version: "v1.0.0",
			},
			validate: func(deps []dto.Module) bool {
				return len(deps) == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			deps, err := service.ExplainDependency(tt.module)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					assert.True(t, tt.validate(deps))
				}
			}
		})
	}
}

func TestModuleService_GetDependencyGraph(t *testing.T) {
	tests := []struct {
		name        string
		module      dto.Module
		expectError bool
		validate    func(map[string][]dto.Module) bool
	}{
		{
			name: "get dependency graph",
			module: dto.Module{
				Path:    "github.com/test/mod1",
				Version: "v1.0.0",
			},
			validate: func(graph map[string][]dto.Module) bool {
				return len(graph) > 0
			},
		},
		{
			name: "get graph for module with no dependencies",
			module: dto.Module{
				Path:    "github.com/test/mod2",
				Version: "v1.0.0",
			},
			validate: func(graph map[string][]dto.Module) bool {
				deps, exists := graph["github.com/test/mod2"]
				return exists && len(deps) == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewModule()
			graph, err := service.GetDependencyGraph(tt.module)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					assert.True(t, tt.validate(graph))
				}
			}
		})
	}
}
