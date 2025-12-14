package main

import (
	"encoding/json"
	"net/http"
	"time"

	"os"

	"example.org/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]models.Recipe, 0)
	for _, recipe := range recipes {
		if recipe.Tags.Contains(tag) {
			listOfRecipes = append(listOfRecipes, recipe)
		}
	}
	if len(listOfRecipes) > 0 {
		c.JSON(http.StatusOK, listOfRecipes)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
}

// @title           Recipes Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  https://bramworks.com/terms/

// @contact.name   API Support
// @contact.url    https://bramworks.com/support
// @contact.email  support@bramworks.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://bramworks.com/resources/open-api/
func main() {
	router := gin.Default()
	router.POST("api/v1/recipes", NewRecipeHandler)
	router.GET("api/v1/recipes", ListRecipesHandler)
	router.PUT("api/v1/recipes/:id", UpdateRecipeHandler)
	router.DELETE("api/v1/recipes/:id", DeleteRecipeHandler)
	router.GET("api/v1/recipes/search", SearchRecipesHandler)

	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run()
}
