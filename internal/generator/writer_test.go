package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriter_GetOutputPath_Success(t *testing.T) {
	tests := []struct {
		name      string
		outputDir string
		itemType  ItemType
		itemName  string
		wantPath  string
	}{
		{
			name:      "skill output path",
			outputDir: "/test/output",
			itemType:  ItemTypeSkill,
			itemName:  "test-skill",
			wantPath:  "/test/output/skills/test-skill/SKILL.md",
		},
		{
			name:      "agent output path",
			outputDir: "/test/output",
			itemType:  ItemTypeAgent,
			itemName:  "test-agent",
			wantPath:  "/test/output/agents/test-agent.md",
		},
		{
			name:      "command output path",
			outputDir: "/test/output",
			itemType:  ItemTypeCommand,
			itemName:  "test-command",
			wantPath:  "/test/output/commands/test-command.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewWriter(&Config{
				OutputDir: tt.outputDir,
				DryRun:    false,
			})

			got, err := writer.GetOutputPath(tt.itemType, tt.itemName)

			require.NoError(t, err)
			assert.Equal(t, tt.wantPath, got)
		})
	}
}

func TestWriter_GetOutputPath_HomeDirectoryExpansion(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := []struct {
		name      string
		outputDir string
		itemType  ItemType
		itemName  string
		wantPath  string
	}{
		{
			name:      "expand tilde only",
			outputDir: "~",
			itemType:  ItemTypeSkill,
			itemName:  "test",
			wantPath:  filepath.Join(homeDir, "skills/test/SKILL.md"),
		},
		{
			name:      "expand tilde with path",
			outputDir: "~/.claude",
			itemType:  ItemTypeAgent,
			itemName:  "test",
			wantPath:  filepath.Join(homeDir, ".claude/agents/test.md"),
		},
		{
			name:      "no expansion needed",
			outputDir: "/absolute/path",
			itemType:  ItemTypeCommand,
			itemName:  "test",
			wantPath:  "/absolute/path/commands/test.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewWriter(&Config{
				OutputDir: tt.outputDir,
				DryRun:    false,
			})

			got, err := writer.GetOutputPath(tt.itemType, tt.itemName)

			require.NoError(t, err)
			assert.Equal(t, tt.wantPath, got)
		})
	}
}

func TestWriter_GetOutputPath_Errors(t *testing.T) {
	tests := []struct {
		name       string
		itemType   ItemType
		itemName   string
		wantErrMsg string
	}{
		{
			name:       "unknown item type",
			itemType:   ItemType("invalid"),
			itemName:   "test",
			wantErrMsg: "unknown item type: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewWriter(&Config{
				OutputDir: "/test",
				DryRun:    false,
			})

			got, err := writer.GetOutputPath(tt.itemType, tt.itemName)

			require.Error(t, err)
			assert.Empty(t, got)
			assert.Contains(t, err.Error(), tt.wantErrMsg)
		})
	}
}

func TestWriter_Write_DryRun(t *testing.T) {
	tests := []struct {
		name     string
		itemType ItemType
		itemName string
		content  string
	}{
		{
			name:     "dry run prints to stdout without writing file",
			itemType: ItemTypeSkill,
			itemName: "test-skill",
			content:  "test content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewWriter(&Config{
				OutputDir: "/test/output",
				DryRun:    true,
			})

			err := writer.Write(tt.itemType, tt.itemName, tt.content)

			require.NoError(t, err)
		})
	}
}

func TestWriter_Write_ActualFileWriting(t *testing.T) {
	tests := []struct {
		name     string
		itemType ItemType
		itemName string
		content  string
	}{
		{
			name:     "write skill file",
			itemType: ItemTypeSkill,
			itemName: "test-skill",
			content:  "skill content",
		},
		{
			name:     "write agent file",
			itemType: ItemTypeAgent,
			itemName: "test-agent",
			content:  "agent content",
		},
		{
			name:     "write command file",
			itemType: ItemTypeCommand,
			itemName: "test-command",
			content:  "command content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			writer := NewWriter(&Config{
				OutputDir: tempDir,
				DryRun:    false,
			})

			err := writer.Write(tt.itemType, tt.itemName, tt.content)
			require.NoError(t, err)

			// Verify file was written
			outputPath, err := writer.GetOutputPath(tt.itemType, tt.itemName)
			require.NoError(t, err)

			// Check file exists
			_, err = os.Stat(outputPath)
			require.NoError(t, err)

			// Verify content
			gotContent, err := os.ReadFile(outputPath)
			require.NoError(t, err)
			assert.Equal(t, tt.content, string(gotContent))

			// Verify file permissions
			info, err := os.Stat(outputPath)
			require.NoError(t, err)
			assert.Equal(t, os.FileMode(0644), info.Mode().Perm())
		})
	}
}

func TestWriter_Write_CreatesParentDirectories(t *testing.T) {
	tests := []struct {
		name     string
		itemType ItemType
		itemName string
		content  string
	}{
		{
			name:     "creates nested directories for skill",
			itemType: ItemTypeSkill,
			itemName: "deeply/nested/skill",
			content:  "test content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			writer := NewWriter(&Config{
				OutputDir: tempDir,
				DryRun:    false,
			})

			err := writer.Write(tt.itemType, tt.itemName, tt.content)
			require.NoError(t, err)

			// Verify file was written
			outputPath, err := writer.GetOutputPath(tt.itemType, tt.itemName)
			require.NoError(t, err)

			// Check file exists
			_, err = os.Stat(outputPath)
			require.NoError(t, err)

			// Verify parent directory has correct permissions
			parentDir := filepath.Dir(outputPath)
			info, err := os.Stat(parentDir)
			require.NoError(t, err)
			assert.True(t, info.IsDir())
			assert.Equal(t, os.FileMode(0755), info.Mode().Perm())
		})
	}
}

func Test_expandHomeDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := []struct {
		name     string
		path     string
		wantPath string
	}{
		{
			name:     "expand tilde only",
			path:     "~",
			wantPath: homeDir,
		},
		{
			name:     "expand tilde with slash",
			path:     "~/test/path",
			wantPath: filepath.Join(homeDir, "test/path"),
		},
		{
			name:     "no tilde no expansion",
			path:     "/absolute/path",
			wantPath: "/absolute/path",
		},
		{
			name:     "tilde not at start no expansion",
			path:     "/path/~/test",
			wantPath: "/path/~/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewWriter(&Config{})

			got, err := writer.expandHomeDir(tt.path)

			require.NoError(t, err)
			// Normalize paths for comparison on different platforms
			gotNormalized := filepath.Clean(got)
			wantNormalized := filepath.Clean(tt.wantPath)

			// Handle case where paths might have different separators
			gotParts := strings.Split(gotNormalized, string(filepath.Separator))
			wantParts := strings.Split(wantNormalized, string(filepath.Separator))

			assert.Equal(t, wantParts, gotParts)
		})
	}
}
