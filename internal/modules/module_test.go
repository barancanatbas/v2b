package modules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModuleService_GetGoModules(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "module-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a mock go.mod file
	goModContent := `module test

go 1.23.1

require (
	github.com/test/mod1 v1.0.0
	github.com/test/mod2 v2.0.0
)
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		prefix        string
		expectedCount int
		expectError   bool
	}{
		{
			name:          "list all modules",
			prefix:        "",
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:          "filter by prefix",
			prefix:        "github.com/test/mod1",
			expectedCount: 1,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Change to test directory
			originalWd, _ := os.Getwd()
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(originalWd)

			moduleService := NewModule()
			modules, err := moduleService.GetGoModules(tt.prefix)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if len(modules) != tt.expectedCount {
					t.Errorf("expected %d modules, got %d", tt.expectedCount, len(modules))
					for _, mod := range modules {
						t.Logf("module: %+v", mod)
					}
				}

				if tt.prefix != "" {
					for _, mod := range modules {
						assert.True(t, strings.HasPrefix(mod.Path, tt.prefix))
					}
				}
			}
		})
	}
}

func TestModuleService_parseGoMod(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
		expectLen   int
	}{
		{
			name: "valid module with direct dependencies",
			content: `module test
require (
	github.com/test/mod1 v1.0.0
	github.com/test/mod2 v2.0.0
)`,
			expectError: false,
			expectLen:   2,
		},
		{
			name: "valid module with indirect dependencies",
			content: `module test
require (
	github.com/test/mod1 v1.0.0 // indirect
)`,
			expectError: false,
			expectLen:   1,
		},
		{
			name:        "invalid module file",
			content:     "invalid content",
			expectError: true,
		},
	}

	moduleService := NewModule()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modules, err := moduleService.parseGoMod([]byte(tt.content))

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, modules, tt.expectLen)
			}
		})
	}
}
