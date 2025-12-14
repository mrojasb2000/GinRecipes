package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrojasb2000/GinRecipes/models"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/recipes", NewRecipeHandler)
	router.GET("/api/v1/recipes", ListRecipesHandler)
	router.PUT("/api/v1/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/api/v1/recipes/:id", DeleteRecipeHandler)
	router.GET("/api/v1/recipes/search", SearchRecipesHandler)
	return router
}

func setupTestData() {
	recipes = []models.Recipe{
		{
			ID:           "test1",
			Name:         "Test Pizza",
			Tags:         models.Tags{"italian", "pizza"},
			Ingredients:  []string{"dough", "tomato", "cheese"},
			Instructions: []string{"prepare dough", "add toppings", "bake"},
			PublishedAt:  time.Now(),
		},
		{
			ID:           "test2",
			Name:         "Test Pasta",
			Tags:         models.Tags{"italian", "pasta"},
			Ingredients:  []string{"pasta", "sauce"},
			Instructions: []string{"boil pasta", "add sauce"},
			PublishedAt:  time.Now(),
		},
	}
}

func TestListRecipesHandler(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/recipes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Recipe
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response), 2)
}

func TestNewRecipeHandler(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	newRecipe := models.Recipe{
		Name:         "New Test Recipe",
		Tags:         models.Tags{"test", "new"},
		Ingredients:  []string{"ingredient1", "ingredient2"},
		Instructions: []string{"step1", "step2"},
	}

	jsonData, _ := json.Marshal(newRecipe)
	req, _ := http.NewRequest("POST", "/api/v1/recipes", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Recipe
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, newRecipe.Name, response.Name)
	assert.Equal(t, newRecipe.Tags, response.Tags)
	assert.NotZero(t, response.PublishedAt)
}

func TestNewRecipeHandler_InvalidJSON(t *testing.T) {
	router := setupTestRouter()

	invalidJSON := []byte(`{"name": "Test", "tags": [invalid]}`)
	req, _ := http.NewRequest("POST", "/api/v1/recipes", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
}

func TestUpdateRecipeHandler(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	updatedRecipe := models.Recipe{
		Name:         "Updated Pizza",
		Tags:         models.Tags{"italian", "updated"},
		Ingredients:  []string{"new dough", "new cheese"},
		Instructions: []string{"new instructions"},
	}

	jsonData, _ := json.Marshal(updatedRecipe)
	req, _ := http.NewRequest("PUT", "/api/v1/recipes/test1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Recipe
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test1", response.ID)
	assert.Equal(t, updatedRecipe.Name, response.Name)
	assert.Equal(t, updatedRecipe.Tags, response.Tags)
}

func TestUpdateRecipeHandler_NotFound(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	updatedRecipe := models.Recipe{
		Name: "Non-existent Recipe",
		Tags: models.Tags{"test"},
	}

	jsonData, _ := json.Marshal(updatedRecipe)
	req, _ := http.NewRequest("PUT", "/api/v1/recipes/nonexistent", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Recipe not found", response["error"])
}

func TestUpdateRecipeHandler_InvalidJSON(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	invalidJSON := []byte(`{"name": invalid}`)
	req, _ := http.NewRequest("PUT", "/api/v1/recipes/test1", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteRecipeHandler(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	initialCount := len(recipes)

	req, _ := http.NewRequest("DELETE", "/api/v1/recipes/test1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Recipe deleted", response["message"])
	assert.Equal(t, initialCount-1, len(recipes))
}

func TestDeleteRecipeHandler_NotFound(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	req, _ := http.NewRequest("DELETE", "/api/v1/recipes/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Recipe not found", response["error"])
}

func TestSearchRecipesHandler(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/recipes/search?tag=pizza", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Recipe
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Greater(t, len(response), 0)

	// Verify all returned recipes have the pizza tag
	for _, recipe := range response {
		assert.True(t, recipe.Tags.Contains("pizza"))
	}
}

func TestSearchRecipesHandler_MultipleResults(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/recipes/search?tag=italian", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Recipe
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response))
}

func TestSearchRecipesHandler_NotFound(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/recipes/search?tag=nonexistenttag", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Recipe not found", response["error"])
}

func TestSearchRecipesHandler_EmptyTag(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/recipes/search?tag=", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
