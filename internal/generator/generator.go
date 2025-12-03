package generator

import (
	"fmt"
)

// Generator orchestrates template generation and file writing
type Generator struct {
	config *Config
	engine *Engine
	writer *Writer
}

// NewGenerator creates a new Generator with the given config
func NewGenerator(config *Config) (*Generator, error) {
	engine, err := NewEngine()
	if err != nil {
		return nil, fmt.Errorf("failed to create engine: %w", err)
	}

	writer := NewWriter(config)

	return &Generator{
		config: config,
		engine: engine,
		writer: writer,
	}, nil
}

// Generate generates content for a specific item and writes it to a file
func (g *Generator) Generate(itemType ItemType, name string) error {
	content, err := g.engine.Generate(itemType, name)
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	if err := g.writer.Write(itemType, name, content); err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}

	return nil
}

// GenerateAll generates all templates of the given type
func (g *Generator) GenerateAll(itemType ItemType) error {
	templates := g.engine.List(itemType)

	for _, name := range templates {
		if err := g.Generate(itemType, name); err != nil {
			return err
		}
	}

	return nil
}

// List returns available template names for the given item type
func (g *Generator) List(itemType ItemType) []string {
	return g.engine.List(itemType)
}
