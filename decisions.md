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
- Consequences: The application will use a local SQLite database file and database migrations. The SQLite driver is selected in D-014.

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
- Consequences: Edge cases such as quantities described as "to taste," ranges, or non-numeric values are deferred. Ingredient units are optional free-form text.

## D-009: Scale With Integer Multipliers

- Status: Accepted
- Date: 2026-08-15
- Decision: Recipes initially scale linearly using positive integer multipliers such as x1, x2, and x3.
- Formula:

  ```text
  scaled quantity = original quantity * integer multiplier
  ```

- Rationale: Integer multiplication provides useful scaling without requiring yield or serving modeling.
- Consequences: Any positive integer is valid initially. Yield-based scaling, fractional multipliers, unit conversion, and quantity formatting rules are deferred.

## D-010: Search Recipe Names By Independent Words

- Status: Accepted
- Date: 2026-08-15
- Decision: Search terms match independently against complete words in the recipe name. A multi-word query should require each word to match, regardless of term order.
- Rationale: Chefs should be able to find a recipe without entering the exact phrase or word order.
- Consequences: Search initially covers recipe names only. A partial term must not match a larger word, such as `cake` matching `cupcake`. Punctuation separates words, and accents are ignored during matching. Recipe descriptions and other fields may be added later.

## D-011: Review Work Before Implementation

- Status: Accepted
- Date: 2026-08-15
- Decision: User stories, acceptance criteria, and implementation issues will be discussed and approved before code is written.
- Rationale: The application is being designed incrementally through vertical slices.
- Consequences: No implementation should be started solely from an unreviewed plan. GitHub issues should be drafted and reviewed before work begins on them.

## D-012: Use Goose For Database Migrations

- Status: Accepted
- Date: 2026-08-16
- Decision: Use Goose to manage versioned SQLite database migrations.
- Rationale: Goose provides a clear migration history and a repeatable way to initialize and evolve the embedded database.
- Consequences: Migration files will be stored in the repository and applied by the application during startup or through an explicit development command. Migration ownership will remain separate from generated query code.

## D-013: Use Sqlc For Type-Safe Database Access

- Status: Accepted
- Date: 2026-08-16
- Decision: Use sqlc to generate type-safe Go database-access code from SQL queries.
- Rationale: SQL remains explicit while generated types reduce manual row mapping and query boilerplate.
- Consequences: SQL query files and sqlc configuration become part of the repository. Generated code will be treated as a build artifact and regenerated when queries or schema change.

## D-014: Use Modernc Sqlite

- Status: Accepted
- Date: 2026-08-16
- Decision: Use `modernc.org/sqlite` as the SQLite driver.
- Rationale: It is a pure-Go driver that avoids a CGO dependency and keeps local and production builds simple.
- Consequences: SQLite access will use the `database/sql` interfaces supported by the driver.

## D-015: Filter Recipe Search In Go Initially

- Status: Accepted
- Date: 2026-08-16
- Decision: Initially load recipe names through sqlc and filter search results in Go memory instead of adding SQLite FTS5.
- Rationale: The application is single-user and the initial catalog is expected to be small. This avoids specialized SQLite search schema and keeps the first implementation simple.
- Consequences: Search will not require an FTS5 table or search index migrations. A database-backed search strategy can replace the in-memory filter later without changing the search behavior or UI contract.

## D-016: Allow Incomplete Recipe Content Initially

- Status: Accepted
- Date: 2026-08-16
- Decision: A recipe requires a name, but ingredients and preparation steps are optional.
- Rationale: The application is being built incrementally, and a recipe record can be created before its content is added.
- Consequences: A recipe may temporarily have no ingredients or no preparation steps. No minimum content rule will be enforced initially.

## D-017: Update Recipe Timestamps For Content Changes

- Status: Accepted
- Date: 2026-08-16
- Decision: Update a recipe's `updated_at` timestamp whenever its name, ingredients, or preparation steps change.
- Rationale: The timestamp should represent the latest change to the complete recipe, not only changes to the recipe row itself.
- Consequences: Ingredient and preparation-step create, edit, delete, and reorder operations must update the parent recipe timestamp in the same transaction.

## D-018: Use Server-Rendered Destructive Confirmations

- Status: Accepted
- Date: 2026-08-16
- Decision: Destructive operations must use a server-rendered confirmation step rather than relying only on browser JavaScript confirmation dialogs.
- Rationale: Forms must remain safe and usable when JavaScript is unavailable.
- Consequences: Ingredient removal, preparation-step removal, and recipe deletion require a confirmation page or equivalent server-rendered confirmation flow before the final POST.

## D-019: Enforce SQLite Cascading Deletes

- Status: Accepted
- Date: 2026-08-16
- Decision: Enable SQLite foreign-key enforcement and define ingredient and preparation-step foreign keys with `ON DELETE CASCADE`.
- Rationale: Deleting a recipe should reliably remove its dependent content without orphaned records.
- Consequences: The SQLite connection must enable foreign keys, and child-table migrations must define the cascade behavior.

## D-020: Use TDD For Vertical Slices

- Status: Accepted
- Date: 2026-08-16
- Decision: Develop each vertical slice test-first. Define a small end-to-end test suite that expresses the user-visible behavior before implementation, then implement the slice against those tests. Add unit tests for isolated domain or transformation logic and integration tests for database, repository, HTTP, or other boundary behavior when they provide focused coverage. Tests may be delivered in the same child edit as their implementation when a separate test edit would not remain runnable.
- Rationale: End-to-end tests keep each slice aligned with its user outcome, while unit and integration tests provide fast, precise feedback for logic and system boundaries.
- Consequences: Tests are designed before implementation. A test-only child edit is allowed when it passes against existing behavior; otherwise tests and implementation belong in the same coherent child edit. A slice is not complete until its end-to-end behavior and relevant lower-level tests pass. Tests remain isolated, deterministic, and limited to the behavior in the approved issue rather than becoming a speculative test framework.

## D-021: Expand CI With Project Capabilities

- Status: Accepted
- Date: 2026-08-16
- Decision: Keep the initial pull-request CI pipeline intentionally small and expand it as the repository gains application code, features, test levels, dependencies, runtime services, and quality requirements.
- Rationale: A minimal pipeline is appropriate for the empty initial repository, while evolving CI with the application keeps checks relevant without adding speculative infrastructure.
- Consequences: New capabilities should add or update the corresponding Makefile targets and CI steps in the same local/CI contract. Future database, integration, end-to-end, race, coverage, linting, or build checks should be introduced when their related implementation exists and should remain runnable locally where practical.

## D-022: Require Runnable Child Edits

- Status: Accepted
- Date: 2026-08-16
- Decision: Every completed child edit must be a complete, compilable, and runnable increment for its scope. A child defaults to one file, but may include multiple files when they are required together for one coherent behavior, such as implementation and its tests.
- Rationale: Child edits are the units reviewed by the user and should never leave unfinished functionality, missing symbols, failing tests, or an unrunnable project state for the next edit.
- Consequences: Child descriptions must list every changed file and explain why they belong together. Relevant tests must pass before the next child edit is created. A test-only child is valid only when it passes against existing behavior; otherwise the test and implementation must be completed in the same child. Temporary red states may be used while constructing a child but must not be left as its completed result.
