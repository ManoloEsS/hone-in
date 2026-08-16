# Issues

This document is the local issue backlog for Hone In. Issues are drafted and reviewed here before they are added to GitHub. No issue should be implemented until its scope and acceptance criteria have been approved.

## Issue Statuses

- Draft: Initial issue proposal.
- Discussing: Scope or acceptance criteria are being reviewed.
- Approved: Ready to be added to GitHub and implemented when scheduled.
- In progress: Implementation has started.
- Complete: Implementation and verification are finished.
- Deferred: Identified as a future or nice-to-have feature and not part of the current delivery plan.

## ISSUE-000: Bootstrap The Go HTTP Server

- Type: Foundation
- Related user story: None; prerequisite for `US-001`
- Status: Approved
- GitHub status: Not added

### Goal

Create a reliable Go HTTP server foundation with explicit startup, runtime error handling, signal handling, graceful shutdown, logging, and HTTP server timeouts.

### Scope

- Add the server entry point under `cmd/server`.
- Use a `run() error` pattern so `main()` handles only process-level exit behavior.
- Configure structured logging with Go's `log/slog`.
- Handle `SIGINT` and `SIGTERM`.
- Shut down gracefully with a configurable timeout.
- Return startup and runtime errors instead of silently ignoring them.
- Ignore `http.ErrServerClosed` during normal shutdown.
- Configure HTTP server timeouts:
  - `ReadHeaderTimeout`
  - `ReadTimeout`
  - `WriteTimeout`
  - `IdleTimeout`
  - `MaxHeaderBytes`
- Add a minimal `/healthz` endpoint.
- Add tests for startup, the health endpoint, shutdown, and error handling.

### Acceptance Criteria

- The server starts on a configurable address with a sensible local default.
- A failure to bind the address causes the process to exit with an error.
- `SIGINT` and `SIGTERM` initiate graceful shutdown.
- In-flight requests are allowed to complete within the shutdown timeout.
- Shutdown timeout expiration is logged as an error.
- Normal `http.ErrServerClosed` shutdown is not logged as a server failure.
- Unexpected server errors are returned and logged.
- The process exits nonzero when `run()` returns an error.
- `GET /healthz` returns HTTP `200`.
- The server does not use `log.Fatal` inside application logic.
- Tests do not require sending signals to the actual development process.

### Out Of Scope

- Recipe functionality
- Recipe catalog routes
- SQLite database connection
- Database migrations
- Recipe database tables
- HTML templates beyond what is required for the health endpoint
- HTMX behavior
- Application styling

### Implementation Notes

- Recommended default address: `:8080`.
- Configuration can initially use an `HTTP_ADDR` environment variable.
- Recommended default shutdown timeout: 10 seconds.
- Use `signal.NotifyContext` for signal cancellation.
- Keep the server lifecycle separate from future recipe and database functionality.

### Review Notes

- Approved during planning before being added to GitHub.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.

## ISSUE-001: Display The Recipe Catalog

- Type: Vertical slice
- Related user story: `US-001`
- Prerequisite: `ISSUE-000`
- Status: Approved
- GitHub status: Not added

### Goal

Load recipes from SQLite and display the recipe catalog, including its empty state.

### Scope

- Open the configured SQLite database.
- Enable SQLite foreign-key enforcement.
- Run the initial database migration.
- Create the minimal `recipes` table.
- Add recipe repository logic for listing recipes.
- Add `GET /recipes`.
- Redirect `GET /` to `/recipes`.
- Render the catalog with server-side Go templates.
- Add minimal plain CSS for the catalog and empty state.
- Add tests for an empty catalog and a catalog containing recipes.

### Initial Recipe Schema

The first migration should contain only the fields needed by the catalog:

- `id`
- `name`
- `created_at`
- `updated_at`

Ingredient and preparation-step tables are intentionally deferred to the create-recipe slice.

### Acceptance Criteria

- The application starts with a configured SQLite database path.
- Database migrations run successfully.
- `GET /recipes` returns HTTP `200`.
- `GET /` redirects to `/recipes`.
- Recipes are loaded from SQLite.
- SQLite foreign-key enforcement is enabled.
- Each stored recipe displays its name.
- An empty database displays a clear empty state.
- The page is server-rendered HTML.
- The page works without JavaScript.
- Repository and handler behavior are covered by tests.
- No recipe creation, editing, deletion, search, or scaling is implemented.

### Out Of Scope

- Ingredient persistence
- Preparation-step persistence
- Recipe creation forms
- Recipe detail pages
- Recipe descriptions
- Recipe yields or servings
- Search
- Scaling
- Authentication or ownership

### Review Notes

- Scope agreed during planning before being added to GitHub.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.

## ISSUE-002: Create A Recipe Record

- Type: Vertical slice
- Related user story: `US-003`
- Prerequisite: `ISSUE-001`
- Status: Approved
- GitHub status: Not added

### Goal

Allow a chef to create a recipe record with a name and see it in the catalog.

### Scope

- Add a "New recipe" action to the catalog.
- Add `GET /recipes/new`.
- Add `POST /recipes`.
- Render a server-side recipe form.
- Validate that the recipe name is present.
- Trim surrounding whitespace from the recipe name.
- Persist the recipe name in SQLite.
- Set `created_at` and `updated_at`.
- Redirect to `/recipes` after successful creation.
- Display the newly created recipe in the catalog.
- Return the form with validation errors when submission fails.
- Add handler, service, repository, and persistence tests.

### Acceptance Criteria

- A chef can open the new-recipe form from the catalog.
- The recipe name is required.
- A blank or whitespace-only name is rejected.
- A valid name is persisted.
- A successful submission redirects to `/recipes`.
- The new recipe appears in the catalog.
- Invalid submissions do not create a database record.
- The form works without JavaScript.

### Out Of Scope

- Ingredients
- Preparation steps
- Recipe detail pages
- Editing
- Deleting
- Search
- Scaling

### Review Notes

- The incremental approach was approved during planning.
- Ingredient and preparation-step functionality will be planned as separate vertical slices.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.

## ISSUE-003: View A Recipe Record

- Type: Vertical slice
- Related user story: `US-004`
- Prerequisite: `ISSUE-002`
- Status: Approved
- GitHub status: Not added

### Goal

Allow a chef to open a recipe from the catalog and view its stored information.

### Scope

- Link recipe names in the catalog to `/recipes/:id`.
- Add `GET /recipes/:id`.
- Load a recipe by ID from SQLite.
- Render the recipe name using server-side HTML.
- Add a minimal recipe detail layout and plain CSS.
- Return a proper 404 response for an unknown recipe.
- Add handler and repository tests.

### Acceptance Criteria

- Each recipe in the catalog links to its detail page.
- `GET /recipes/:id` returns HTTP `200` for an existing recipe.
- The recipe name is displayed.
- `GET /recipes/:id` returns HTTP `404` for a missing recipe.
- The page works without JavaScript.
- Existing catalog behavior continues to work.
- No recipe data is changed by viewing the page.

### Out Of Scope

- Ingredients
- Preparation steps
- Recipe descriptions
- Editing
- Deleting
- Search
- Scaling

### Review Notes

- The detail-page-first sequence was approved during planning.
- Ingredient persistence will be planned after this slice.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.

## ISSUE-004: Add Ingredients To A Recipe

- Type: Vertical slice
- Related user story: `US-003`
- Prerequisite: `ISSUE-003`
- Status: Approved
- GitHub status: Not added

### Goal

Allow a chef to add ingredients to an existing recipe and see them displayed on the recipe detail page.

### Scope

- Add the `ingredients` table and migration.
- Add a foreign key from ingredients to recipes.
- Configure the foreign key with `ON DELETE CASCADE`.
- Add `POST /recipes/:id/ingredients`.
- Add an inline ingredient form to the recipe detail page.
- Require a numeric quantity greater than zero.
- Require an ingredient name.
- Store the unit as optional free-form text.
- Append ingredients in entry order.
- Display ingredients on the recipe detail page.
- Update the parent recipe's `updated_at` timestamp.
- Return the updated ingredient list through HTMX.
- Keep normal form submission working without JavaScript.
- Add repository, handler, validation, and persistence tests.

### Initial Ingredient Schema

- `id`
- `recipe_id`
- `position`
- `quantity`
- `unit`
- `name`

### Acceptance Criteria

- A chef can open a recipe detail page and submit an ingredient.
- The ingredient quantity is required.
- The quantity must be numeric and greater than zero.
- The ingredient name is required.
- The unit may be empty.
- A provided unit is stored as user-entered text.
- A valid ingredient is persisted for the selected recipe.
- Ingredients are displayed in the order they were added.
- An HTMX submission updates the ingredient list without a full page reload.
- A normal form submission works without JavaScript.
- Invalid submissions do not create an ingredient record.
- Adding an ingredient to one recipe does not affect another recipe.
- The parent recipe's `updated_at` value changes after a successful addition.

### Out Of Scope

- Editing ingredients
- Removing ingredients
- Reordering ingredients
- Ingredient notes
- Non-numeric quantities
- Quantity ranges
- Unit conversion
- Scaling
- Special ingredient edge cases

### Review Notes

- Units are optional free-form text.
- Quantity validation requires a value greater than zero.
- Edge cases will be handled in later issues as they arise.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.

## ISSUE-005: Add Preparation Steps To A Recipe

- Type: Vertical slice
- Related user story: `US-003`
- Prerequisite: `ISSUE-004`
- Status: Approved
- GitHub status: Not added

### Goal

Allow a chef to add ordered preparation steps to an existing recipe and see them on the recipe detail page.

### Scope

- Add the `recipe_steps` table and migration.
- Add a foreign key from steps to recipes.
- Configure the foreign key with `ON DELETE CASCADE`.
- Add `POST /recipes/:id/steps`.
- Add an inline preparation-step form to the recipe detail page.
- Require non-empty instruction text.
- Trim surrounding whitespace from instructions.
- Append steps in entry order.
- Display steps as a numbered list.
- Update the parent recipe's `updated_at` timestamp.
- Return the updated step list through HTMX.
- Keep normal form submission working without JavaScript.
- Add repository, handler, validation, and persistence tests.

### Initial Step Schema

- `id`
- `recipe_id`
- `position`
- `instruction`

### Acceptance Criteria

- A chef can open a recipe detail page and submit a preparation step.
- The instruction is required.
- A blank or whitespace-only instruction is rejected.
- Surrounding whitespace is removed before persistence.
- A valid step is persisted for the selected recipe.
- Steps are displayed as a numbered list in the order they were added.
- An HTMX submission updates the steps section without a full page reload.
- A normal form submission works without JavaScript.
- Invalid submissions do not create a step record.
- Adding a step to one recipe does not affect another recipe.
- The parent recipe's `updated_at` value changes after a successful addition.

### Out Of Scope

- Editing steps
- Removing steps
- Reordering steps
- Recipe descriptions
- Search
- Scaling
- Special instruction formatting

### Review Notes

- The issue follows the same one-at-a-time submission pattern as ingredient creation.
- Editing, removal, and reordering will be planned as later issues.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.

## ISSUE-006: Edit A Recipe Name

- Type: Vertical slice
- Related user story: `US-005`
- Prerequisite: `ISSUE-005`
- Status: Approved
- GitHub status: Not added

### Goal

Allow a chef to correct or update the name of an existing recipe.

### Scope

- Add an "Edit" action to the recipe detail page.
- Add `GET /recipes/:id/edit`.
- Add `POST /recipes/:id/edit`.
- Prepopulate the form with the current recipe name.
- Require a non-empty name.
- Trim surrounding whitespace from the recipe name.
- Update `updated_at`.
- Redirect to the recipe detail page after a successful update.
- Preserve all ingredients and preparation steps.
- Add handler, validation, repository, and persistence tests.

### Acceptance Criteria

- A chef can open the edit form from a recipe detail page.
- The current recipe name is displayed in the form.
- A blank or whitespace-only name is rejected.
- A valid name is persisted.
- The updated name appears on the detail page and in the catalog.
- Ingredients and preparation steps are unchanged.
- A missing recipe returns HTTP `404`.
- The form works without JavaScript.

### Out Of Scope

- Editing ingredients
- Removing ingredients
- Reordering ingredients
- Editing preparation steps
- Removing preparation steps
- Reordering preparation steps
- Deleting recipes
- Search
- Scaling

### Review Notes

- Recipe editing is intentionally split by recipe component.
- Ingredient and preparation-step editing will be planned as separate issues.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.

## ISSUE-007: Edit An Ingredient

- Type: Vertical slice
- Related user story: `US-005`
- Prerequisite: `ISSUE-006`
- Status: Approved
- GitHub status: Not added

### Goal

Allow a chef to correct an existing ingredient without changing its position in the recipe.

### Scope

- Add an edit action for each displayed ingredient.
- Add `GET /recipes/:id/ingredients/:ingredient_id/edit`.
- Add `POST /recipes/:id/ingredients/:ingredient_id/edit`.
- Prepopulate quantity, unit, and name.
- Require a quantity greater than zero.
- Require an ingredient name.
- Keep the unit as optional free-form text.
- Preserve the ingredient's existing position.
- Update the ingredient through HTMX.
- Support normal form submission without JavaScript.
- Return HTTP `404` when the recipe or ingredient does not exist.
- Prevent an ingredient from being edited through the wrong recipe URL.
- Update the parent recipe's `updated_at` timestamp.
- Add validation, handler, repository, and persistence tests.

### Acceptance Criteria

- A chef can open an edit form for an existing ingredient.
- The form is prepopulated with the ingredient's current quantity, unit, and name.
- A quantity must be greater than zero.
- An ingredient name is required.
- The unit may be empty.
- A valid update is persisted.
- The ingredient remains in its original position.
- The updated ingredient row is returned through HTMX without a full page reload.
- A normal form submission works without JavaScript.
- A missing recipe returns HTTP `404`.
- A missing ingredient returns HTTP `404`.
- An ingredient belonging to another recipe cannot be updated through the selected recipe URL.
- The parent recipe's `updated_at` value changes after a successful update.

### Out Of Scope

- Removing ingredients
- Reordering ingredients
- Adding ingredients
- Editing preparation steps
- Deleting recipes
- Search
- Scaling

### Review Notes

- The issue follows the existing ingredient validation rules.
- Updating an ingredient also updates the parent recipe timestamp.
- Ingredient removal and reordering will be planned as separate issues.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.

## ISSUE-008: Remove An Ingredient

- Type: Vertical slice
- Related user story: `US-005`
- Prerequisite: `ISSUE-007`
- Status: Approved
- GitHub status: Not added

### Goal

Allow a chef to remove an ingredient from a recipe.

### Scope

- Add a remove action for each ingredient.
- Add a server-rendered confirmation step.
- Add `POST /recipes/:id/ingredients/:ingredient_id/delete`.
- Delete the selected ingredient.
- Reorder remaining ingredient positions.
- Update the parent recipe's `updated_at` timestamp.
- Update the ingredient list through HTMX.
- Support normal form submission without JavaScript.
- Prevent deleting an ingredient through the wrong recipe URL.
- Allow the last ingredient to be removed.
- Add repository, handler, transaction, and persistence tests.

### Acceptance Criteria

- A chef can request removal of an ingredient.
- The removal requires server-rendered confirmation.
- The selected ingredient is deleted after confirmation.
- Remaining ingredients retain their relative order.
- Remaining positions are normalized.
- The parent recipe's `updated_at` value changes.
- The updated ingredient list is returned through HTMX.
- A normal form submission works without JavaScript.
- The confirmation flow works without JavaScript.
- A missing recipe returns HTTP `404`.
- A missing ingredient returns HTTP `404`.
- An ingredient belonging to another recipe cannot be deleted.
- Removing the last ingredient leaves the recipe with no ingredients.

### Out Of Scope

- Reordering ingredients manually
- Editing ingredients
- Adding ingredients
- Preparation-step management
- Search
- Scaling

### Review Notes

- The last ingredient may be removed because no minimum ingredient count has been established.
- Remaining ingredient positions are normalized after deletion.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.

## ISSUE-009: Reorder Ingredients

- Type: Nice to have
- Related user story: `US-005`
- Prerequisite: `ISSUE-008`
- Priority: Nice to have
- Status: Deferred
- GitHub status: Not added

### Goal

Allow a chef to manually change the order of ingredients in a recipe.

### Current Decision

Ingredient order does not materially affect the initial cooking workflow. Ingredients will remain in insertion order for the current delivery plan.

### Future Scope

- Add controls to move an ingredient up or down.
- Persist the updated ingredient order.
- Preserve ingredient data while changing positions.
- Support HTMX and normal form submission.
- Add tests for boundary positions and persisted ordering.

### Out Of Scope For Current Delivery

- Drag-and-drop ordering
- Bulk ordering
- Preparation-step ordering

### Review Notes

- This issue is deferred and does not block preparation-step management.
- It should not be added to GitHub until the feature is prioritized.

## ISSUE-010: Edit A Preparation Step

- Type: Vertical slice
- Related user story: `US-005`
- Prerequisite: `ISSUE-005`
- Status: Approved
- GitHub status: Not added

### Goal

Allow a chef to correct an existing preparation step without changing its order.

### Scope

- Add an edit action for each preparation step.
- Add `GET /recipes/:id/steps/:step_id/edit`.
- Add `POST /recipes/:id/steps/:step_id/edit`.
- Prepopulate the instruction field.
- Require non-empty instruction text.
- Trim surrounding whitespace from instructions.
- Preserve the step's existing position.
- Update the step through HTMX.
- Support normal form submission without JavaScript.
- Return HTTP `404` for missing recipes or steps.
- Prevent editing a step through the wrong recipe URL.
- Update the parent recipe's `updated_at` timestamp.
- Add validation, handler, repository, and persistence tests.

### Acceptance Criteria

- A chef can open an edit form for an existing step.
- The current instruction is displayed.
- A blank or whitespace-only instruction is rejected.
- A valid change is persisted.
- The step remains in its original position.
- The updated step row is returned through HTMX without a full page reload.
- A normal form submission works without JavaScript.
- A missing recipe returns HTTP `404`.
- A missing step returns HTTP `404`.
- A step belonging to another recipe cannot be updated through the selected recipe URL.
- The parent recipe's `updated_at` value changes after a successful update.

### Out Of Scope

- Removing steps
- Reordering steps
- Editing ingredients
- Deleting recipes
- Search
- Scaling

### Review Notes

- Step order is preserved during editing.
- Preparation-step removal and reordering will be planned as separate issues.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.

## ISSUE-011: Remove A Preparation Step

- Type: Vertical slice
- Related user story: `US-005`
- Prerequisite: `ISSUE-010`
- Status: Approved
- GitHub status: Not added

### Goal

Allow a chef to remove a preparation step from a recipe while preserving the order of the remaining steps.

### Scope

- Add a remove action for each preparation step.
- Add a server-rendered confirmation step.
- Add `POST /recipes/:id/steps/:step_id/delete`.
- Delete the selected step.
- Normalize the positions of remaining steps.
- Update the parent recipe's `updated_at` timestamp.
- Refresh the steps list through HTMX.
- Support normal form submission without JavaScript.
- Prevent deleting a step through the wrong recipe URL.
- Allow the last preparation step to be removed.
- Add repository, handler, transaction, and persistence tests.

### Acceptance Criteria

- A chef can request removal of a preparation step.
- Removal requires server-rendered confirmation.
- The selected step is deleted after confirmation.
- Remaining steps retain their relative order.
- Remaining positions are normalized.
- The parent recipe's `updated_at` value changes.
- The updated steps list is returned through HTMX.
- A normal form submission works without JavaScript.
- The confirmation flow works without JavaScript.
- A missing recipe returns HTTP `404`.
- A missing step returns HTTP `404`.
- A step belonging to another recipe cannot be deleted.
- Removing the last preparation step leaves the recipe with no steps.

### Out Of Scope

- Reordering preparation steps
- Editing preparation steps
- Ingredient management
- Deleting recipes
- Search
- Scaling

### Review Notes

- The last preparation step may be removed because no minimum step count has been established.
- Remaining step positions are normalized after deletion.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.

## ISSUE-012: Reorder Preparation Steps

- Type: Vertical slice
- Related user story: `US-005`
- Prerequisite: `ISSUE-011`
- Status: Approved
- GitHub status: Not added

### Goal

Allow a chef to change the order of preparation steps so the recipe reflects the correct cooking sequence.

### Scope

- Add "Move up" and "Move down" controls for each step.
- Add an endpoint to move a step one position.
- Swap adjacent step positions transactionally.
- Keep step positions normalized.
- Preserve instruction text.
- Update the parent recipe's `updated_at` timestamp.
- Refresh the steps list through HTMX.
- Support normal form submission without JavaScript.
- Add repository, handler, transaction, and persistence tests.

### Acceptance Criteria

- A chef can move a step up one position.
- A chef can move a step down one position.
- The first step cannot move up.
- The last step cannot move down.
- Moving a step preserves all step data.
- The new order is persisted.
- The updated steps list is returned through HTMX.
- A normal form submission works without JavaScript.
- The parent recipe's `updated_at` value changes when the order changes.
- A step cannot be reordered through the wrong recipe URL.

### Out Of Scope

- Drag-and-drop ordering
- Bulk reordering
- Ingredient reordering
- Recipe deletion
- Search
- Scaling

### Review Notes

- Preparation-step order is part of the core recipe workflow.
- Simple move-up and move-down controls are preferred over drag-and-drop initially.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.

## ISSUE-013: Delete A Recipe

- Type: Vertical slice
- Related user story: `US-006`
- Prerequisite: `ISSUE-012`
- Status: Approved
- GitHub status: Not added

### Goal

Allow a chef to permanently remove a recipe and all of its associated ingredients and preparation steps.

### Scope

- Add a delete action to the recipe detail page.
- Add a server-rendered confirmation step.
- Add `POST /recipes/:id/delete`.
- Delete the selected recipe.
- Cascade-delete its ingredients and preparation steps.
- Redirect to `/recipes` after deletion.
- Return an appropriate HTMX redirect when submitted through HTMX.
- Prevent deleting a recipe through an invalid or mismatched URL.
- Add repository, handler, transaction, and persistence tests.

### Acceptance Criteria

- A chef can request deletion from a recipe detail page.
- Deletion requires server-rendered confirmation.
- The recipe is removed after confirmation.
- Associated ingredients are removed.
- Associated preparation steps are removed.
- The deleted recipe no longer appears in the catalog.
- Viewing the deleted recipe returns HTTP `404`.
- A missing recipe returns HTTP `404`.
- A normal form submission works without JavaScript.
- The confirmation flow works without JavaScript.
- The delete operation does not affect other recipes.

### Out Of Scope

- Search
- Scaling
- Soft deletion
- Undo
- Bulk deletion
- Importing or exporting

### Review Notes

- Permanent deletion is the initial behavior because no recovery or archival requirement exists.
- Associated records must be removed as part of the same delete operation.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.

## ISSUE-014: Search Recipes By Name

- Type: Vertical slice
- Related user story: `US-002`
- Prerequisite: `ISSUE-013`
- Status: Approved
- GitHub status: Not added

### Goal

Allow a chef to search the catalog using independent complete words from a recipe name.

### Scope

- Add a search input to the catalog page.
- Use `GET /recipes?q=...` so search remains part of the catalog route.
- Normalize and split the query into non-empty terms.
- Require every term to match a complete word in the recipe name.
- Make matching case-insensitive.
- Ignore accents during matching.
- Make term order irrelevant.
- Treat punctuation as a word separator.
- Load recipe names through sqlc and filter matches in Go memory.
- Return the recipe list through HTMX.
- Add a short HTMX debounce to avoid a request on every keystroke.
- Preserve normal full-page search behavior without JavaScript.
- Show the full catalog for an empty query.
- Show a clear no-results state.
- Preserve the existing alphabetical result ordering.
- Add search service, repository, handler, and template tests.

### Acceptance Criteria

- Searching for one complete word returns matching recipe names.
- Searching for multiple words requires every word to match.
- Search term order does not affect results.
- Search is case-insensitive.
- Accented and unaccented forms match equivalently.
- A partial word does not match a larger word, such as `cake` matching `cupcake`.
- Punctuation separates words for matching purposes.
- Empty searches return the full catalog.
- No-match searches show a useful empty state.
- HTMX updates only the recipe results.
- The search query remains visible in the input.
- Normal form submission works without JavaScript.
- Recipe descriptions are not searched.
- Existing alphabetical result ordering is preserved.

### Out Of Scope

- Searching descriptions
- Searching ingredients or instructions
- Categories, tags, or labels
- Search ranking
- Pagination
- Advanced search syntax
- Search history

### Review Notes

- Complete-word matching was selected during planning.
- Accent-insensitive matching and punctuation-as-separator behavior were approved during planning.
- Server-side Go filtering was selected instead of FTS5 for the initial implementation.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.

## ISSUE-015: Scale A Recipe With An Integer Multiplier

- Type: Vertical slice
- Related user story: `US-007`
- Prerequisite: `ISSUE-014`
- Status: Approved
- GitHub status: Not added

### Goal

Allow a chef to view scaled ingredient quantities without changing the stored recipe.

### Scope

- Add an integer multiplier control to the recipe detail page.
- Default the multiplier to `x1`.
- Add `GET /recipes/:id/scale?multiplier=N`.
- Require a positive integer multiplier.
- Calculate each scaled quantity server-side:

  ```text
  scaled quantity = original quantity * multiplier
  ```

- Preserve ingredient names and units.
- Return the scaled ingredient list through HTMX.
- Keep normal full-page behavior without JavaScript.
- Do not write scaled values back to SQLite.
- Add scaling service, handler, validation, and calculation tests.

### Acceptance Criteria

- A chef can select a multiplier such as `x1`, `x2`, or `x3`.
- Any positive integer multiplier is accepted.
- Every ingredient quantity is multiplied by the selected integer.
- The original ingredient quantities remain unchanged.
- The scaled ingredient list updates through HTMX.
- Invalid, zero, negative, or non-integer multipliers are rejected.
- A missing recipe returns HTTP `404`.
- A recipe without ingredients displays an empty ingredient section.
- Instructions remain unchanged.
- Scaling works without JavaScript.

### Out Of Scope

- Yield-based scaling
- Fractional multipliers
- Unit conversion
- Quantity formatting
- Special ingredient handling
- Persisting scaled recipes
- Scaling preparation steps

### Review Notes

- Any positive integer is valid initially.
- Scaling is calculated at request time and never persisted.
- The issue should be added to GitHub only after the complete local issue backlog has been reviewed.
