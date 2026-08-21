-- name: ListRecipes :many
SELECT id, name, created_at, updated_at
FROM recipes
ORDER BY name ASC, id ASC;

-- name: CreateRecipe :one
INSERT INTO recipes (name, created_at, updated_at)
VALUES (?, ?, ?)
RETURNING id, name, created_at, updated_at;

-- name: GetRecipe :one
SELECT id, name, created_at, updated_at
FROM recipes
WHERE id = ?;
