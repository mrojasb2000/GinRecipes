# Tests Unitarios - API de Recetas

Este documento describe los tests unitarios implementados para la API de recetas en Go.

## Estructura de Tests

El proyecto contiene dos archivos de tests:

- `models/recipe_test.go` - Tests para el modelo Recipe y el tipo Tags
- `main_test.go` - Tests para los handlers HTTP de la API

## Cobertura de Tests

- **Cobertura total**: 85.7%
- **Modelo (models)**: 100%
- **Handlers (main)**: 100% (excepto la función `main()`)

## Ejecutar los Tests

### Ejecutar todos los tests
```bash
go test ./... -v
```

### Ejecutar tests con cobertura
```bash
go test ./... -coverprofile=coverage.out
```

### Ver reporte de cobertura
```bash
go tool cover -func=coverage.out
```

### Ver reporte de cobertura en HTML
```bash
go tool cover -html=coverage.out
```

### Ejecutar tests de un paquete específico
```bash
# Solo tests del modelo
go test ./models -v

# Solo tests de handlers
go test . -v
```

## Tests Implementados

### Tests del Modelo (`models/recipe_test.go`)

#### `TestTags_Contains`
Prueba el método `Contains` del tipo `Tags` con diferentes escenarios:
- Tag existe en la lista
- Tag no existe en la lista
- Lista de tags vacía
- Tag al inicio de la lista
- Tag al final de la lista
- Sensibilidad a mayúsculas/minúsculas

#### `TestRecipe_JSONMarshaling`
Verifica que el struct `Recipe` se pueda serializar y deserializar correctamente a/desde JSON.

#### `TestRecipe_EmptyFields`
Prueba el comportamiento de un `Recipe` con valores por defecto.

### Tests de Handlers (`main_test.go`)

#### `TestListRecipesHandler`
Prueba el endpoint `GET /api/v1/recipes` que lista todas las recetas.
- **Código esperado**: 200 OK
- **Validación**: Respuesta contiene al menos 2 recetas

#### `TestNewRecipeHandler`
Prueba el endpoint `POST /api/v1/recipes` para crear una nueva receta.
- **Código esperado**: 201 Created
- **Validación**: 
  - ID generado automáticamente
  - PublishedAt se establece
  - Datos de la receta se guardan correctamente

#### `TestNewRecipeHandler_InvalidJSON`
Prueba el manejo de errores cuando se envía JSON inválido.
- **Código esperado**: 400 Bad Request
- **Validación**: Respuesta contiene mensaje de error

#### `TestUpdateRecipeHandler`
Prueba el endpoint `PUT /api/v1/recipes/:id` para actualizar una receta existente.
- **Código esperado**: 200 OK
- **Validación**: 
  - ID se mantiene
  - Datos se actualizan correctamente
  - PublishedAt se actualiza

#### `TestUpdateRecipeHandler_NotFound`
Prueba la actualización de una receta que no existe.
- **Código esperado**: 404 Not Found
- **Validación**: Mensaje de error apropiado

#### `TestUpdateRecipeHandler_InvalidJSON`
Prueba el manejo de errores en actualización con JSON inválido.
- **Código esperado**: 400 Bad Request

#### `TestDeleteRecipeHandler`
Prueba el endpoint `DELETE /api/v1/recipes/:id` para eliminar una receta.
- **Código esperado**: 200 OK
- **Validación**: 
  - Mensaje de confirmación
  - Receta se elimina del slice

#### `TestDeleteRecipeHandler_NotFound`
Prueba la eliminación de una receta que no existe.
- **Código esperado**: 404 Not Found
- **Validación**: Mensaje de error apropiado

#### `TestSearchRecipesHandler`
Prueba el endpoint `GET /api/v1/recipes/search?tag=<tag>` para buscar recetas por tag.
- **Código esperado**: 200 OK
- **Validación**: Todas las recetas devueltas contienen el tag buscado

#### `TestSearchRecipesHandler_MultipleResults`
Prueba la búsqueda que devuelve múltiples resultados.
- **Código esperado**: 200 OK
- **Validación**: Se devuelven todas las recetas que coinciden

#### `TestSearchRecipesHandler_NotFound`
Prueba la búsqueda de un tag que no existe.
- **Código esperado**: 404 Not Found
- **Validación**: Mensaje de error apropiado

#### `TestSearchRecipesHandler_EmptyTag`
Prueba la búsqueda con un tag vacío.
- **Código esperado**: 404 Not Found

## Funciones Helper

### `setupTestRouter()`
Configura un router de Gin en modo test con todas las rutas de la API.

### `setupTestData()`
Inicializa el slice `recipes` con datos de prueba:
- Test Pizza (tags: italian, pizza)
- Test Pasta (tags: italian, pasta)

## Dependencias de Testing

- `github.com/stretchr/testify/assert` - Para aserciones más expresivas
- `github.com/gin-gonic/gin` - Framework web (modo test)
- `net/http/httptest` - Para crear servidores HTTP de prueba

## Mejores Prácticas Implementadas

1. **Aislamiento**: Cada test configura sus propios datos de prueba
2. **Nomenclatura clara**: Los nombres de los tests describen exactamente qué prueban
3. **Table-driven tests**: Usado en `TestTags_Contains` para múltiples casos
4. **Test de casos límite**: Se prueban casos de error, datos vacíos, y casos normales
5. **Assertions claras**: Uso de testify para mensajes de error descriptivos
6. **Modo test de Gin**: Desactiva el logging innecesario durante los tests

## Ejecutar Tests Continuamente

Para desarrollo, puedes usar herramientas como `gotestsum` o `entr`:

```bash
# Instalar gotestsum
go install gotest.tools/gotestsum@latest

# Ejecutar con mejor output
gotestsum --watch
```

## Notas

- La función `main()` no está cubierta por tests ya que inicia el servidor
- Los tests usan datos en memoria, no acceden al archivo `recipes.json`
- Cada test es independiente y puede ejecutarse en cualquier orden