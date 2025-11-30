package models

import (
	"slices"
	"time"
)

type Tags []string

func (t Tags) Contains(key string) bool {
	return slices.Contains(t, key)
}

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         Tags      `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}
