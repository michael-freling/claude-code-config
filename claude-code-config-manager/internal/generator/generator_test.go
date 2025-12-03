package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGenerator(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "successfully creates generator",
			config: &Config{
				OutputDir: "/test/output",
				DryRun:    false,
			},
			wantErr: false,
		},
		{
			name: "successfully creates generator with dry run",
			config: &Config{
				OutputDir: "/test/output",
				DryRun:    true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := NewGenerator(tt.config)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, gen)
			assert.NotNil(t, gen.engine)
			assert.NotNil(t, gen.writer)
			assert.Equal(t, tt.config, gen.config)
		})
	}
}

func TestGenerator_Generate_Success(t *testing.T) {
	tests := []struct {
		name         string
		itemType     ItemType
		templateName string
	}{
		{
			name:         "generate skill",
			itemType:     ItemTypeSkill,
			templateName: "coding",
		},
		{
			name:         "generate agent",
			itemType:     ItemTypeAgent,
			templateName: "golang-code-reviewer",
		},
		{
			name:         "generate command",
			itemType:     ItemTypeCommand,
			templateName: "feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			config := &Config{
				OutputDir: tempDir,
				DryRun:    false,
			}

			gen, err := NewGenerator(config)
			require.NoError(t, err)

			err = gen.Generate(tt.itemType, tt.templateName)

			require.NoError(t, err)

			// Verify file was created
			outputPath, err := gen.writer.GetOutputPath(tt.itemType, tt.templateName)
			require.NoError(t, err)

			_, err = os.Stat(outputPath)
			require.NoError(t, err)

			// Verify content is not empty
			content, err := os.ReadFile(outputPath)
			require.NoError(t, err)
			assert.NotEmpty(t, content)
		})
	}
}

func TestGenerator_Generate_Errors(t *testing.T) {
	tests := []struct {
		name         string
		itemType     ItemType
		templateName string
		wantErrMsg   string
	}{
		{
			name:         "non-existent template",
			itemType:     ItemTypeSkill,
			templateName: "nonexistent",
			wantErrMsg:   "failed to generate content",
		},
		{
			name:         "invalid item type",
			itemType:     ItemType("invalid"),
			templateName: "test",
			wantErrMsg:   "failed to generate content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			config := &Config{
				OutputDir: tempDir,
				DryRun:    false,
			}

			gen, err := NewGenerator(config)
			require.NoError(t, err)

			err = gen.Generate(tt.itemType, tt.templateName)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErrMsg)
		})
	}
}

func TestGenerator_GenerateAll_Success(t *testing.T) {
	tests := []struct {
		name     string
		itemType ItemType
		minCount int
	}{
		{
			name:     "generate all skills",
			itemType: ItemTypeSkill,
			minCount: 3,
		},
		{
			name:     "generate all agents",
			itemType: ItemTypeAgent,
			minCount: 3,
		},
		{
			name:     "generate all commands",
			itemType: ItemTypeCommand,
			minCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			config := &Config{
				OutputDir: tempDir,
				DryRun:    false,
			}

			gen, err := NewGenerator(config)
			require.NoError(t, err)

			err = gen.GenerateAll(tt.itemType)

			require.NoError(t, err)

			// Verify files were created
			var subDir string
			switch tt.itemType {
			case ItemTypeSkill:
				subDir = "skills"
			case ItemTypeAgent:
				subDir = "agents"
			case ItemTypeCommand:
				subDir = "commands"
			}

			outputDir := filepath.Join(tempDir, subDir)
			entries, err := os.ReadDir(outputDir)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(entries), tt.minCount)
		})
	}
}

func TestGenerator_GenerateAll_Errors(t *testing.T) {
	tests := []struct {
		name       string
		itemType   ItemType
		outputDir  string
		setupFunc  func(t *testing.T, tempDir string)
		wantErrMsg string
	}{
		{
			name:      "write failure due to read-only directory",
			itemType:  ItemTypeSkill,
			outputDir: "",
			setupFunc: func(t *testing.T, tempDir string) {
				// Create a read-only parent directory
				skillsDir := filepath.Join(tempDir, "skills")
				err := os.MkdirAll(skillsDir, 0555)
				require.NoError(t, err)
			},
			wantErrMsg: "failed to write content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			if tt.setupFunc != nil {
				tt.setupFunc(t, tempDir)
			}

			config := &Config{
				OutputDir: tempDir,
				DryRun:    false,
			}

			gen, err := NewGenerator(config)
			require.NoError(t, err)

			err = gen.GenerateAll(tt.itemType)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErrMsg)
		})
	}
}

func TestGenerator_List(t *testing.T) {
	tests := []struct {
		name         string
		itemType     ItemType
		wantContains []string
	}{
		{
			name:         "list skills",
			itemType:     ItemTypeSkill,
			wantContains: []string{"coding", "docker", "bash"},
		},
		{
			name:         "list agents",
			itemType:     ItemTypeAgent,
			wantContains: []string{"golang-code-reviewer", "golang-engineer"},
		},
		{
			name:         "list commands",
			itemType:     ItemTypeCommand,
			wantContains: []string{"feature", "fix", "refactor"},
		},
		{
			name:         "list invalid type returns empty",
			itemType:     ItemType("invalid"),
			wantContains: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				OutputDir: "/test",
				DryRun:    true,
			}

			gen, err := NewGenerator(config)
			require.NoError(t, err)

			result := gen.List(tt.itemType)

			for _, want := range tt.wantContains {
				assert.Contains(t, result, want)
			}
		})
	}
}

func TestGenerator_Generate_DryRun(t *testing.T) {
	tests := []struct {
		name         string
		itemType     ItemType
		templateName string
	}{
		{
			name:         "dry run does not create file",
			itemType:     ItemTypeSkill,
			templateName: "coding",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			config := &Config{
				OutputDir: tempDir,
				DryRun:    true,
			}

			gen, err := NewGenerator(config)
			require.NoError(t, err)

			err = gen.Generate(tt.itemType, tt.templateName)

			require.NoError(t, err)

			// Verify no file was created
			outputPath, err := gen.writer.GetOutputPath(tt.itemType, tt.templateName)
			require.NoError(t, err)

			_, err = os.Stat(outputPath)
			assert.True(t, os.IsNotExist(err))
		})
	}
}
