-- name: ListRecipes :many
SELECT id, name, created_at, updated_at
FROM recipes
ORDER BY name ASC, id ASC;
