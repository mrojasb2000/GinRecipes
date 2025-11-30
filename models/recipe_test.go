package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTags_Contains(t *testing.T) {
	tests := []struct {
		name     string
		tags     Tags
		key      string
		expected bool
	}{
		{
			name:     "Tag exists in list",
			tags:     Tags{"italian", "pizza", "dinner"},
			key:      "pizza",
			expected: true,
		},
		{
			name:     "Tag does not exist in list",
			tags:     Tags{"italian", "pizza", "dinner"},
			key:      "dessert",
			expected: false,
		},
		{
			name:     "Empty tags list",
			tags:     Tags{},
			key:      "pizza",
			expected: false,
		},
		{
			name:     "Tag exists at beginning",
			tags:     Tags{"italian", "pizza", "dinner"},
			key:      "italian",
			expected: true,
		},
		{
			name:     "Tag exists at end",
			tags:     Tags{"italian", "pizza", "dinner"},
			key:      "dinner",
			expected: true,
		},
		{
			name:     "Case sensitive - different case",
			tags:     Tags{"Italian", "Pizza"},
			key:      "italian",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tags.Contains(tt.key)
			if result != tt.expected {
				t.Errorf("Tags.Contains() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestRecipe_JSONMarshaling(t *testing.T) {
	publishedAt := time.Date(2021, 1, 17, 19, 28, 52, 0, time.UTC)

	recipe := Recipe{
		ID:           "test123",
		Name:         "Test Recipe",
		Tags:         Tags{"italian", "dinner"},
		Ingredients:  []string{"ingredient1", "ingredient2"},
		Instructions: []string{"step1", "step2"},
		PublishedAt:  publishedAt,
	}

	// Test marshaling
	jsonData, err := json.Marshal(recipe)
	if err != nil {
		t.Fatalf("Failed to marshal recipe: %v", err)
	}

	// Test unmarshaling
	var unmarshaledRecipe Recipe
	err = json.Unmarshal(jsonData, &unmarshaledRecipe)
	if err != nil {
		t.Fatalf("Failed to unmarshal recipe: %v", err)
	}

	// Verify fields
	if unmarshaledRecipe.ID != recipe.ID {
		t.Errorf("ID mismatch: got %v, want %v", unmarshaledRecipe.ID, recipe.ID)
	}
	if unmarshaledRecipe.Name != recipe.Name {
		t.Errorf("Name mismatch: got %v, want %v", unmarshaledRecipe.Name, recipe.Name)
	}
	if len(unmarshaledRecipe.Tags) != len(recipe.Tags) {
		t.Errorf("Tags length mismatch: got %v, want %v", len(unmarshaledRecipe.Tags), len(recipe.Tags))
	}
	if len(unmarshaledRecipe.Ingredients) != len(recipe.Ingredients) {
		t.Errorf("Ingredients length mismatch: got %v, want %v", len(unmarshaledRecipe.Ingredients), len(recipe.Ingredients))
	}
	if len(unmarshaledRecipe.Instructions) != len(recipe.Instructions) {
		t.Errorf("Instructions length mismatch: got %v, want %v", len(unmarshaledRecipe.Instructions), len(recipe.Instructions))
	}
}

func TestRecipe_EmptyFields(t *testing.T) {
	recipe := Recipe{}

	if recipe.ID != "" {
		t.Errorf("Expected empty ID, got %v", recipe.ID)
	}
	if recipe.Name != "" {
		t.Errorf("Expected empty Name, got %v", recipe.Name)
	}
	if recipe.Tags != nil {
		t.Error("Expected Tags to be nil when not initialized")
	}
	if recipe.Ingredients != nil {
		t.Error("Expected Ingredients to be nil when not initialized")
	}
	if recipe.Instructions != nil {
		t.Error("Expected Instructions to be nil when not initialized")
	}
	if !recipe.PublishedAt.IsZero() {
		t.Error("Expected PublishedAt to be zero time")
	}
}
