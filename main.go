package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	docs "github.com/mrojasb2000/GinRecipes/docs"
	"github.com/mrojasb2000/GinRecipes/models"
	"github.com/rs/xid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var recipes []models.Recipe
var ctx context.Context
var err error
var client *mongo.Client

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("Connected to MongoDB!")
}

// Add new Recipe
//
//	@Summary		Operation POST /recipes recipes.
//	@Description	Add a new recipe.
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			models.Recipe	body		models.Recipe	true	"Add recipe"
//	@Success		200	{object}	models.Recipe
//	@Failure		400	{object}	httputil.HTTPError
//	@Failure		404	{object}	httputil.HTTPError
//	@Failure		500	{object}	httputil.HTTPError
//	@Router			/recipes [post]
//
// NewRecipeHandler handles POST requests to create a new recipe.
// It binds the JSON request body to a Recipe model, validates the input,
// generates a unique ID using xid, sets the published timestamp to the current time,
// appends the recipe to the recipes slice, and returns the created recipe with HTTP 201 status.
// If the JSON binding fails, it returns an HTTP 400 error with the validation error message.
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

// Recipes list
// @Summary      Operation GET /recipes returns a list of recipes.
// @Description  Return a recipes list
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Recipe
// @Failure		400	{object}	httputil.HTTPError
// @Failure		404	{object}	httputil.HTTPError
// @Failure		500	{object}	httputil.HTTPError
// @Router       /recipes [get]
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// Update Recipe
//
//	@Summary		Operation PUT /recipes/{id} recipes.
//	@Description	Update an existing recipe.
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Recipe ID"
//	@Param			models.Recipe	body		models.Recipe	true	"Update recipe"
//	@Success		200	{object}	models.Recipe
//	@Failure		400	{object}	httputil.HTTPError
//	@Failure		404	{object}	httputil.HTTPError
//	@Failure		500	{object}	httputil.HTTPError
//	@Router			/recipes/{id} [put]
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

// Delete Recipe
//
//	@Summary		Operation DELETE /recipes/{id} recipes.
//	@Description	Delete an existing recipe.
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Recipe ID"
//	@Success		200	{object}	models.Recipe
//	@Failure		400	{object}	httputil.HTTPError
//	@Failure		404	{object}	httputil.HTTPError
//	@Failure		500	{object}	httputil.HTTPError
//	@Router			/recipes/{id} [delete]
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

// Search Recipes
//
//	@Summary		Operation Search Recipe GET /recipes/search?={tag} recipes.
//	@Description	Search an existing recipe.
//	@Tags			recipes
//	@Accept			json
//	@Produce		json
//	@Param			tag	query		string	true	"Tag Recipe"
//	@Success		200	{object}	models.Recipe
//	@Failure		400	{object}	httputil.HTTPError
//	@Failure		404	{object}	httputil.HTTPError
//	@Failure		500	{object}	httputil.HTTPError
//	@Router			/recipes/search [get]
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

// @title           Recipes Example API.
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
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.POST("api/v1/recipes", NewRecipeHandler)
	router.GET("api/v1/recipes", ListRecipesHandler)
	router.PUT("api/v1/recipes/:id", UpdateRecipeHandler)
	router.DELETE("api/v1/recipes/:id", DeleteRecipeHandler)
	router.GET("api/v1/recipes/search", SearchRecipesHandler)

	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run()
}
