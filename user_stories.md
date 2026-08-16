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

Initial acceptance criteria to discuss:

- The application has a recipe catalog page.
- Stored recipes are displayed in the catalog.
- Each recipe links to its detail page.
- An empty catalog has a useful empty state.

Status: Draft

#### US-002: Search Recipes By Name

As a chef, I want to search recipes by independent words in their names so that I can quickly find a recipe.

Initial acceptance criteria to discuss:

- The catalog provides a search input.
- Search is case-insensitive.
- Each meaningful word in a multi-word query must match the recipe name.
- Search results update without a full page reload using HTMX.
- An empty query returns the full catalog.
- A query with no matches displays a useful empty state.

Search does not include recipe descriptions or other fields in the initial version.

Status: Draft

### Recipe Management

#### US-003: Create A Recipe

As a chef, I want to create a recipe so that I can store it in my catalog.

Initial acceptance criteria to discuss:

- A chef can open a create-recipe form from the catalog.
- A recipe has a required name.
- A recipe can contain ingredients with required numeric quantities.
- A recipe can contain ordered preparation steps.
- Invalid input is reported clearly.
- A valid recipe is persisted and can be viewed afterward.

Status: Draft

#### US-004: View A Recipe

As a chef, I want to view a complete recipe so that I can use it while cooking.

Initial acceptance criteria to discuss:

- The recipe page displays the recipe name.
- Ingredients are displayed in their stored order.
- Ingredient quantities and units are displayed with their names.
- Preparation steps are displayed in their stored order.
- A missing recipe returns a not-found response.

Status: Draft

#### US-005: Edit A Recipe

As a chef, I want to edit a recipe so that I can correct or improve it.

Initial acceptance criteria to discuss:

- Existing recipe data is loaded into an edit form.
- The recipe name can be changed.
- Ingredients can be added, changed, removed, and reordered.
- Preparation steps can be added, changed, removed, and reordered.
- Ingredient quantities remain numeric and required.
- A valid update is persisted and visible on the recipe page.

Status: Draft

#### US-006: Delete A Recipe

As a chef, I want to delete a recipe so that I can remove recipes I no longer need.

Initial acceptance criteria to discuss:

- The delete action requires confirmation.
- Deleting a recipe removes it from the catalog.
- Deleting a recipe also removes its ingredients and preparation steps.
- A deleted recipe can no longer be viewed.

Status: Draft

### Recipe Scaling

#### US-007: Scale A Recipe With An Integer Multiplier

As a chef, I want to multiply a recipe by an integer so that I can prepare more or less of it without changing the stored recipe.

Initial acceptance criteria to discuss:

- The recipe page provides an integer multiplier such as x1, x2, or x3.
- The multiplier applies to every numeric ingredient quantity.
- Scaling is linear:

  ```text
  scaled quantity = original quantity * integer multiplier
  ```

- The scaled ingredient list updates without a full page reload using HTMX.
- The original ingredient quantities remain unchanged in storage.
- Yield-based scaling is not part of this story.

Status: Draft

## Proposed Delivery Order

This is a starting point, not an approved implementation plan:

1. US-001: View the recipe catalog
2. US-003: Create a recipe
3. US-004: View a recipe
4. US-005: Edit a recipe
5. US-006: Delete a recipe
6. US-002: Search recipes by name
7. US-007: Scale a recipe with an integer multiplier

The order may change after discussing the stories and identifying the smallest useful vertical slices.

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
