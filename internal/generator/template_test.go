package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEngine(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successfully creates engine and loads templates",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine()

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, engine)
			assert.NotNil(t, engine.templates)
			assert.NotNil(t, engine.templateNames)

			// Verify all three item types are loaded
			assert.Contains(t, engine.templates, ItemTypeSkill)
			assert.Contains(t, engine.templates, ItemTypeAgent)
			assert.Contains(t, engine.templates, ItemTypeCommand)
		})
	}
}

func TestEngine_Generate_Success(t *testing.T) {
	tests := []struct {
		name         string
		itemType     ItemType
		templateName string
	}{
		{
			name:         "generate skill template",
			itemType:     ItemTypeSkill,
			templateName: "coding",
		},
		{
			name:         "generate agent template",
			itemType:     ItemTypeAgent,
			templateName: "golang-code-reviewer",
		},
		{
			name:         "generate command template",
			itemType:     ItemTypeCommand,
			templateName: "feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine()
			require.NoError(t, err)

			result, err := engine.Generate(tt.itemType, tt.templateName)

			require.NoError(t, err)
			assert.NotEmpty(t, result)
			assert.Contains(t, result, tt.templateName)
		})
	}
}

func TestEngine_Generate_Errors(t *testing.T) {
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
			wantErrMsg:   "template nonexistent not found for type skill",
		},
		{
			name:         "invalid item type",
			itemType:     ItemType("invalid"),
			templateName: "test",
			wantErrMsg:   "no templates found for type: invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine()
			require.NoError(t, err)

			result, err := engine.Generate(tt.itemType, tt.templateName)

			require.Error(t, err)
			assert.Empty(t, result)
			assert.Contains(t, err.Error(), tt.wantErrMsg)
		})
	}
}

func TestEngine_List(t *testing.T) {
	tests := []struct {
		name         string
		itemType     ItemType
		wantContains []string
		wantNotEmpty bool
	}{
		{
			name:         "list skills",
			itemType:     ItemTypeSkill,
			wantContains: []string{"coding", "docker", "bash"},
			wantNotEmpty: true,
		},
		{
			name:         "list agents",
			itemType:     ItemTypeAgent,
			wantContains: []string{"golang-code-reviewer", "golang-engineer"},
			wantNotEmpty: true,
		},
		{
			name:         "list commands",
			itemType:     ItemTypeCommand,
			wantContains: []string{"feature", "fix", "refactor"},
			wantNotEmpty: true,
		},
		{
			name:         "list invalid type returns empty",
			itemType:     ItemType("invalid"),
			wantContains: []string{},
			wantNotEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine()
			require.NoError(t, err)

			result := engine.List(tt.itemType)

			if tt.wantNotEmpty {
				assert.NotEmpty(t, result)
			} else {
				assert.Empty(t, result)
			}

			for _, want := range tt.wantContains {
				assert.Contains(t, result, want)
			}
		})
	}
}
