package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config holds configuration for the Writer
type Config struct {
	OutputDir string // Base output directory (default: ~/.claude)
	DryRun    bool   // If true, write to stdout instead of files
}

// Writer handles writing generated content to files or stdout
type Writer struct {
	config *Config
}

// NewWriter creates a new Writer with the given configuration
func NewWriter(config *Config) *Writer {
	return &Writer{
		config: config,
	}
}

// Write writes the generated content to the appropriate location
func (w *Writer) Write(itemType ItemType, name string, content string) error {
	if w.config.DryRun {
		fmt.Println(content)
		return nil
	}

	outputPath, err := w.GetOutputPath(itemType, name)
	if err != nil {
		return fmt.Errorf("failed to get output path: %w", err)
	}

	// Create parent directories if needed
	parentDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory %s: %w", parentDir, err)
	}

	// Write file with 0644 permissions
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", outputPath, err)
	}

	return nil
}

// GetOutputPath returns the output path for the given item type and name
func (w *Writer) GetOutputPath(itemType ItemType, name string) (string, error) {
	outputDir, err := w.expandHomeDir(w.config.OutputDir)
	if err != nil {
		return "", fmt.Errorf("failed to expand home directory: %w", err)
	}

	var relativePath string
	switch itemType {
	case ItemTypeSkill:
		relativePath = filepath.Join("skills", name, "SKILL.md")
	case ItemTypeAgent:
		relativePath = filepath.Join("agents", fmt.Sprintf("%s.md", name))
	case ItemTypeCommand:
		relativePath = filepath.Join("commands", fmt.Sprintf("%s.md", name))
	default:
		return "", fmt.Errorf("unknown item type: %s", itemType)
	}

	return filepath.Join(outputDir, relativePath), nil
}

// expandHomeDir expands ~ to the user's home directory if present at the start
func (w *Writer) expandHomeDir(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	if path == "~" {
		return homeDir, nil
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:]), nil
	}

	return path, nil
}
