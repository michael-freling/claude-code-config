package main

import (
	"fmt"
	"os"

	"github.com/michael-freling/claude-code-config/internal/generator"
	"github.com/spf13/cobra"
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

	rootCmd.AddCommand(newAgentsCmd())
	rootCmd.AddCommand(newCommandsCmd())
	rootCmd.AddCommand(newSkillsCmd())

	return rootCmd
}

func newAgentsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "agents",
		Short: "Generate prompts for all agents",
		Long:  `Generate prompts to create all agent definitions.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			gen, err := generator.NewGenerator()
			if err != nil {
				return fmt.Errorf("failed to create generator: %w", err)
			}

			if err := gen.GenerateAll(generator.ItemTypeAgent); err != nil {
				return fmt.Errorf("failed to generate agents: %w", err)
			}

			return nil
		},
	}
}

func newCommandsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "commands",
		Short: "Generate prompts for all commands",
		Long:  `Generate prompts to create all command definitions.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			gen, err := generator.NewGenerator()
			if err != nil {
				return fmt.Errorf("failed to create generator: %w", err)
			}

			if err := gen.GenerateAll(generator.ItemTypeCommand); err != nil {
				return fmt.Errorf("failed to generate commands: %w", err)
			}

			return nil
		},
	}
}

func newSkillsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "skills",
		Short: "Generate prompts for all skills",
		Long:  `Generate prompts to create all skill definitions.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			gen, err := generator.NewGenerator()
			if err != nil {
				return fmt.Errorf("failed to create generator: %w", err)
			}

			if err := gen.GenerateAll(generator.ItemTypeSkill); err != nil {
				return fmt.Errorf("failed to generate skills: %w", err)
			}

			return nil
		},
	}
}
