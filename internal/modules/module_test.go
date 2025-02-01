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

	tests := []struct {
		name          string
		goModContent  string
		prefix        string
		expectError   bool
		expectedCount int
	}{
		{
			name: "valid go.mod with multiple modules",
			goModContent: `module test

go 1.23.1

require (
	github.com/test/mod1 v1.0.0-abc123
	github.com/test/mod2 v2.0.0-def456
)
`,
			prefix:        "",
			expectError:   false,
			expectedCount: 2,
		},
		{
			name: "valid go.mod with prefix filter",
			goModContent: `module test

go 1.23.1

require (
	github.com/test/mod1 v1.0.0-abc123
	github.com/other/mod2 v2.0.0-def456
)
`,
			prefix:        "github.com/test",
			expectError:   false,
			expectedCount: 1,
		},
		{
			name: "invalid go.mod content",
			goModContent: `invalid content
not a valid go.mod file
`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary go.mod file
			goModPath := filepath.Join(tmpDir, "go.mod")
			err := os.WriteFile(goModPath, []byte(tt.goModContent), 0644)
			if err != nil {
				t.Fatal(err)
			}

			// Change working directory to temp directory
			originalWd, _ := os.Getwd()
			err = os.Chdir(tmpDir)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(originalWd)

			moduleService := NewModule()
			modules, err := moduleService.GetGoModules(tt.prefix)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, modules, tt.expectedCount)

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
