# Decisions

This document records product and architecture decisions for Hone In. Decisions are intentionally small and can be revisited when new requirements or implementation experience justify a change.

## D-001: Use Go For The Application

- Status: Accepted
- Date: 2026-08-15
- Decision: Use Go for the backend and application server.
- Rationale: Go provides a small, deployable server with a strong standard library and good support for HTML rendering and SQLite.
- Consequences: The application will initially be a single Go application rather than separate frontend and backend services.

## D-002: Use Server-Rendered HTML With HTMX

- Status: Accepted
- Date: 2026-08-15
- Decision: Render HTML on the server and use HTMX for interactive updates.
- Rationale: The application does not need a separate single-page frontend. Server-rendered HTML keeps the frontend and backend simple while HTMX provides partial-page interactions.
- Consequences: Go handlers will return complete pages for normal requests and HTML fragments for HTMX requests. There will be no JSON API initially.

## D-003: Use Plain CSS

- Status: Accepted
- Date: 2026-08-15
- Decision: Use plain CSS without a frontend framework or frontend build system.
- Rationale: The initial interface is small and does not need the complexity of a CSS framework or JavaScript bundler.
- Consequences: Styles will be maintained as static application assets.

## D-004: Use An Embedded SQLite Database

- Status: Accepted
- Date: 2026-08-15
- Decision: Use SQLite as the application's embedded database.
- Rationale: The application is initially single-user and local, so a separate database server is unnecessary.
- Consequences: The application will use a local SQLite database file and database migrations. The SQLite driver choice remains an implementation detail to confirm before coding.

## D-005: Start As A Single-User Application

- Status: Accepted
- Date: 2026-08-15
- Decision: Do not implement authentication, users, ownership, or sharing initially.
- Rationale: The first version is intended for one chef and should focus on recipe management.
- Consequences: The initial schema will not need a users table or recipe ownership field. Authentication and sharing can be added later if needed.

## D-006: Keep The Repository As One Application

- Status: Accepted
- Date: 2026-08-15
- Decision: Keep the Go code, templates, static assets, and database migrations in one repository and one deployable application.
- Rationale: Separate frontend and backend projects would add complexity without providing value for a server-rendered HTMX application.
- Consequences: Concerns will be separated into packages and directories, but they will not be separate services.

## D-007: Keep The Initial Recipe Model Small

- Status: Accepted
- Date: 2026-08-15
- Decision: The initial recipe model will focus on a recipe name, ingredients, and ordered preparation steps.
- Rationale: These are the core capabilities needed to store and use a recipe.
- Consequences: Recipe descriptions, categories, tags, labels, importing, printing, and exporting are deferred. Recipe yield and serving information are also deferred.

## D-008: Require Numeric Ingredient Quantities

- Status: Accepted
- Date: 2026-08-15
- Decision: Ingredient quantities are required and numeric in the initial version.
- Rationale: Numeric values allow the first scaling implementation to remain predictable and simple.
- Consequences: Edge cases such as quantities described as "to taste," ranges, or non-numeric values are deferred. The behavior of ingredient units remains to be defined separately.

## D-009: Scale With Integer Multipliers

- Status: Accepted
- Date: 2026-08-15
- Decision: Recipes initially scale linearly using positive integer multipliers such as x1, x2, and x3.
- Formula:

  ```text
  scaled quantity = original quantity * integer multiplier
  ```

- Rationale: Integer multiplication provides useful scaling without requiring yield or serving modeling.
- Consequences: Yield-based scaling, fractional multipliers, unit conversion, and quantity formatting rules are deferred.

## D-010: Search Recipe Names By Independent Words

- Status: Accepted
- Date: 2026-08-15
- Decision: Search terms match independently against words in the recipe name. A multi-word query should require each word to match.
- Rationale: Chefs should be able to find a recipe without entering the exact phrase or word order.
- Consequences: Search initially covers recipe names only. Recipe descriptions and other fields may be added later. Exact tokenization, punctuation handling, ranking, and partial-word behavior will be defined when the search story is discussed.

## D-011: Review Work Before Implementation

- Status: Accepted
- Date: 2026-08-15
- Decision: User stories, acceptance criteria, and implementation issues will be discussed and approved before code is written.
- Rationale: The application is being designed incrementally through vertical slices.
- Consequences: No implementation should be started solely from an unreviewed plan. GitHub issues should be drafted and reviewed before work begins on them.

## Open Questions

- Should an ingredient unit be optional or required?
- What is the minimum content required to save a recipe?
- Should a recipe be allowed to have no ingredients or no preparation steps?
- What range of integer multipliers should the initial interface offer?
- How should search handle punctuation, accents, and partial words?
