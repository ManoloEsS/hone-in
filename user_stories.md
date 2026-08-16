# User Stories

These stories describe user-visible behavior for Hone In. They are not implementation tasks. Before work begins on a story, we will discuss its acceptance criteria, split it into one or more vertical slices when necessary, and draft the related GitHub issues.

## Story Statuses

- Draft: Identified but not yet discussed in detail.
- Ready: Acceptance criteria and scope are agreed upon.
- In progress: An approved GitHub issue is being implemented.
- Complete: The behavior is implemented and verified.

## MVP Story Map

### Catalog

#### US-001: View The Recipe Catalog

As a chef, I want to view my recipe catalog so that I can choose a recipe to work with.

Acceptance criteria:

- The application has a recipe catalog page.
- Stored recipes are displayed in the catalog.
- Each recipe links to its detail page.
- An empty catalog has a useful empty state.

Issues: `ISSUE-001`, `ISSUE-003`

Status: Ready

#### US-002: Search Recipes By Name

As a chef, I want to search recipes by independent words in their names so that I can quickly find a recipe.

Acceptance criteria:

- The catalog provides a search input.
- Search is case-insensitive.
- Each complete word in a multi-word query must match the recipe name.
- Search term order does not affect results.
- Accented and unaccented forms match equivalently.
- Punctuation separates words.
- Search results update without a full page reload using HTMX.
- An empty query returns the full catalog.
- A query with no matches displays a useful empty state.

Search does not include recipe descriptions or other fields in the initial version.

Issues: `ISSUE-014`

Status: Ready

### Recipe Management

#### US-003: Create A Recipe

As a chef, I want to create a recipe so that I can store it in my catalog.

Acceptance criteria:

- A chef can open a create-recipe form from the catalog.
- A recipe has a required name.
- A recipe may be created before ingredients or preparation steps are added.
- Ingredients have required numeric quantities greater than zero.
- Ingredient units are optional free-form text.
- A recipe can contain ordered preparation steps.
- Invalid input is reported clearly.
- A valid recipe is persisted and can be viewed afterward.

This story is delivered incrementally through recipe creation, ingredient creation, and preparation-step creation issues.

Issues: `ISSUE-002`, `ISSUE-004`, `ISSUE-005`

Status: Ready

#### US-004: View A Recipe

As a chef, I want to view a complete recipe so that I can use it while cooking.

Acceptance criteria:

- The recipe page displays the recipe name.
- Ingredients are displayed in their stored order.
- Ingredient quantities and optional units are displayed with their names.
- Preparation steps are displayed in their stored order.
- A missing recipe returns a not-found response.

This story is delivered incrementally as recipe details, ingredients, and preparation steps become available.

Issues: `ISSUE-003`, `ISSUE-004`, `ISSUE-005`

Status: Ready

#### US-005: Edit A Recipe

As a chef, I want to edit a recipe so that I can correct or improve it.

Acceptance criteria:

- Existing recipe data is loaded into an edit form.
- The recipe name can be changed.
- Existing ingredients can be changed or removed.
- Preparation steps can be changed, removed, and reordered.
- Ingredient quantities remain numeric, required, and greater than zero.
- A valid update is persisted and visible on the recipe page.

Ingredient reordering is deferred as `ISSUE-009` because ingredient order does not materially affect the initial cooking workflow.

Issues: `ISSUE-006`, `ISSUE-007`, `ISSUE-008`, `ISSUE-010`, `ISSUE-011`, `ISSUE-012`

Status: Ready

#### US-006: Delete A Recipe

As a chef, I want to delete a recipe so that I can remove recipes I no longer need.

Acceptance criteria:

- The delete action requires server-rendered confirmation.
- Deleting a recipe removes it from the catalog.
- Deleting a recipe also removes its ingredients and preparation steps.
- A deleted recipe can no longer be viewed.

Issues: `ISSUE-013`

Status: Ready

### Recipe Scaling

#### US-007: Scale A Recipe With An Integer Multiplier

As a chef, I want to multiply a recipe by an integer so that I can prepare more or less of it without changing the stored recipe.

Acceptance criteria:

- The recipe page provides an integer multiplier such as x1, x2, or x3.
- Any positive integer multiplier is valid.
- The multiplier applies to every numeric ingredient quantity.
- Scaling is linear:

  ```text
  scaled quantity = original quantity * integer multiplier
  ```

- The scaled ingredient list updates without a full page reload using HTMX.
- The original ingredient quantities remain unchanged in storage.
- Yield-based scaling is not part of this story.

Issues: `ISSUE-015`

Status: Ready

## Approved Delivery Order

The approved issue delivery order is:

1. `ISSUE-000`: Bootstrap the Go HTTP server
2. `ISSUE-001`: Display the recipe catalog
3. `ISSUE-002`: Create a recipe record
4. `ISSUE-003`: View a recipe record
5. `ISSUE-004`: Add ingredients to a recipe
6. `ISSUE-005`: Add preparation steps to a recipe
7. `ISSUE-006`: Edit a recipe name
8. `ISSUE-007`: Edit an ingredient
9. `ISSUE-008`: Remove an ingredient
10. `ISSUE-010`: Edit a preparation step
11. `ISSUE-011`: Remove a preparation step
12. `ISSUE-012`: Reorder preparation steps
13. `ISSUE-013`: Delete a recipe
14. `ISSUE-014`: Search recipes by name
15. `ISSUE-015`: Scale a recipe with an integer multiplier

`ISSUE-009`, ingredient reordering, is deferred and does not block the MVP delivery order.

## Vertical Slice Workflow

For each story, we will:

1. Discuss the user outcome and edge cases.
2. Agree on acceptance criteria.
3. Decide whether the story is one slice or needs multiple user-visible slices.
4. Define the necessary database, domain, HTTP, template, and test work.
5. Draft the GitHub issue or issues.
6. Review and approve the issues before adding them to GitHub.
7. Implement one approved issue at a time.
8. Verify the slice and update its status here.
