package models

import (
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Tags         []string
	Ingredients  []string
	Instructions []string
)

func (t Tags) Contains(key string) bool {
	return slices.Contains(t, key)
}

// swagger:parameters recipes newRecipe
type Recipe struct {
	//swagger:ignore
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name" example:"Chocolate Cake"`
	Tags         Tags               `json:"tags" bson:"tags" example:"dessert,sweet"`
	Ingredients  Ingredients        `json:"ingredients" bson:"ingredients" example:"ingredient1,ingredient2"`
	Instructions Instructions       `json:"instructions" bson:"instructions" example:"instruction1,instruction2"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt" example:"2024-01-01T00:00:00Z"`
}
