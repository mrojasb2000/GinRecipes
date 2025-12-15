package models

import (
	"slices"
	"time"
)

type (
	Tags         []string
	Ingredients  []string
	Instructions []string
)

func (t Tags) Contains(key string) bool {
	return slices.Contains(t, key)
}

type Recipe struct {
	ID           string       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
	Name         string       `json:"name" example:"Recipe name"`
	Tags         Tags         `json:"tags" example:"value1,value2"`
	Ingredients  Ingredients  `json:"ingredients" example:"ingredient1,ingredient2"`
	Instructions Instructions `json:"instructions" example:"instruction1,instruction2"`
	PublishedAt  time.Time    `json:"publishedAt" example:"2024-01-01T00:00:00Z"`
}
