package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/michael-freling/go-prompt-generator/internal/generator"
	"github.com/spf13/cobra"
)

var (
	outputDir string
	dryRun    bool
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "generator",
		Short: "Generate Claude Code prompts for skills, agents, and commands",
		Long:  `A CLI tool to generate Claude Code prompts from templates for skills, agents, and commands.`,
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "~"
	}
	defaultOutputDir := filepath.Join(homeDir, ".claude")

	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", defaultOutputDir, "Output directory")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Print to stdout instead of writing files")

	rootCmd.AddCommand(newGenerateCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newGenerateAllCmd())

	return rootCmd
}

func newGenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a specific item",
		Long:  `Generate a specific skill, agent, or command from a template.`,
	}

	generateCmd.AddCommand(newGenerateItemCmd(generator.ItemTypeSkill))
	generateCmd.AddCommand(newGenerateItemCmd(generator.ItemTypeAgent))
	generateCmd.AddCommand(newGenerateItemCmd(generator.ItemTypeCommand))

	return generateCmd
}

func newGenerateItemCmd(itemType generator.ItemType) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s <name>", itemType),
		Short: fmt.Sprintf("Generate a specific %s", itemType),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			config := &generator.Config{
				OutputDir: outputDir,
				DryRun:    dryRun,
			}

			gen, err := generator.NewGenerator(config)
			if err != nil {
				return fmt.Errorf("failed to create generator: %w", err)
			}

			if err := gen.Generate(itemType, name); err != nil {
				return fmt.Errorf("failed to generate %s %s: %w", itemType, name, err)
			}

			if !dryRun {
				outputPath, err := generator.NewWriter(config).GetOutputPath(itemType, name)
				if err != nil {
					return fmt.Errorf("failed to get output path: %w", err)
				}
				fmt.Printf("Generated %s '%s' at %s\n", itemType, name, outputPath)
			}

			return nil
		},
	}
}

func newListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available items",
		Long:  `List all available skills, agents, or commands.`,
	}

	listCmd.AddCommand(newListItemCmd(generator.ItemTypeSkill))
	listCmd.AddCommand(newListItemCmd(generator.ItemTypeAgent))
	listCmd.AddCommand(newListItemCmd(generator.ItemTypeCommand))

	return listCmd
}

func newListItemCmd(itemType generator.ItemType) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%ss", itemType),
		Short: fmt.Sprintf("List all available %ss", itemType),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := &generator.Config{
				OutputDir: outputDir,
				DryRun:    dryRun,
			}

			gen, err := generator.NewGenerator(config)
			if err != nil {
				return fmt.Errorf("failed to create generator: %w", err)
			}

			items := gen.List(itemType)
			if len(items) == 0 {
				fmt.Printf("No %ss available\n", itemType)
				return nil
			}

			fmt.Printf("Available %ss:\n", itemType)
			for _, item := range items {
				fmt.Printf("  - %s\n", item)
			}

			return nil
		},
	}
}

func newGenerateAllCmd() *cobra.Command {
	generateAllCmd := &cobra.Command{
		Use:   "generate-all [type]",
		Short: "Generate all items",
		Long:  `Generate all items. Optionally specify a type (skills, agents, or commands) to generate only that type.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := &generator.Config{
				OutputDir: outputDir,
				DryRun:    dryRun,
			}

			gen, err := generator.NewGenerator(config)
			if err != nil {
				return fmt.Errorf("failed to create generator: %w", err)
			}

			var itemTypes []generator.ItemType
			if len(args) == 1 {
				switch args[0] {
				case "skills":
					itemTypes = []generator.ItemType{generator.ItemTypeSkill}
				case "agents":
					itemTypes = []generator.ItemType{generator.ItemTypeAgent}
				case "commands":
					itemTypes = []generator.ItemType{generator.ItemTypeCommand}
				default:
					return fmt.Errorf("invalid type: %s (must be skills, agents, or commands)", args[0])
				}
			} else {
				itemTypes = []generator.ItemType{
					generator.ItemTypeSkill,
					generator.ItemTypeAgent,
					generator.ItemTypeCommand,
				}
			}

			for _, itemType := range itemTypes {
				if err := gen.GenerateAll(itemType); err != nil {
					return fmt.Errorf("failed to generate all %ss: %w", itemType, err)
				}

				items := gen.List(itemType)
				if !dryRun && len(items) > 0 {
					fmt.Printf("Generated %d %s(s)\n", len(items), itemType)
				}
			}

			return nil
		},
	}

	return generateAllCmd
}
