package main

import (
	"encoding/json"
	"net/http"
	"time"

	"os"

	"example.org/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var recipes []models.Recipe

func init() {
	recipes = make([]models.Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal(file, &recipes)
}

func NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusCreated, recipe)
}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i, r := range recipes {
		if r.ID == id {
			recipe.ID = r.ID
			recipe.PublishedAt = time.Now()
			recipes[i] = recipe
			c.JSON(http.StatusOK, recipe)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	for i, r := range recipes {
		if r.ID == id {
			recipes = append(recipes[:i], recipes[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
}

func main() {
	router := gin.Default()
	router.POST("api/v1/recipes", NewRecipeHandler)
	router.GET("api/v1/recipes", ListRecipesHandler)
	router.PUT("api/v1/recipes/:id", UpdateRecipeHandler)
	router.DELETE("api/v1/recipes/:id", DeleteRecipeHandler)
	router.Run()
}
