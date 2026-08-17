---
name: hone-in-development
description: Use for Hone In implementation work when changes must follow the approved issue backlog, Jujutsu issue roots, complete runnable child edits, continuous within-issue review, and user-controlled squashing.
---

# Hone In Development Workflow

This skill defines the implementation workflow for the Hone In repository. Follow it for application code, tests, templates, static assets, migrations, generated code, and implementation-related documentation.

## Non-Negotiable Rules

- Work on one approved issue at a time.
- Do not implement an issue whose status is not `Approved` or `In progress` in `issues.md`.
- Do not move to a later issue until the current issue's child edits are complete, the user has reviewed the issue's work, and the user has squashed them into the current issue root.
- Require explicit user confirmation before creating a new issue root. No confirmation is required between child edits within the current issue.
- The user controls Jujutsu squashing. Never run `jj squash`, `jj abandon`, rebase operations, or equivalent history-rewriting commands unless the user explicitly asks for that operation.
- Never edit the issue root after it is created unless the user explicitly permits it.
- Never edit a parent or ancestor of the issue root unless the user explicitly permits it.
- Keep every implementation edit limited to one file by default. Allow multiple files in one child when they are required together for one coherent, complete, runnable behavior, such as implementation and its tests.
- Follow test-driven development for every vertical slice: define the tests first, then complete the smallest child that makes the relevant tests pass.
- Never complete or report a code child that does not compile, has failing relevant tests, or leaves unfinished or unrunnable functionality.
- Do not make unrelated cleanup, refactoring, formatting, dependency, or documentation changes in an issue child edit.

## Jujutsu Model

An issue root is the Jujutsu change named for the issue being implemented. It is the target into which the user will eventually squash the reviewed child edits. The root is an anchor, not a working change.

Create the root only after the user identifies or approves the issue:

```text
jj new -m "ISSUE-000: Bootstrap The Go HTTP Server"
```

Use the exact issue identifier and the issue title from `issues.md`. Do not put implementation files directly in this root.

Each functional change is a new child edit on top of the current edit:

```text
jj new -m "ISSUE-000: Add the server entry point"
```

The child edit should change one file when possible. It may change multiple files when those files are necessary to deliver one complete, runnable behavior. Its description must explain the change in enough detail for review:

- `Files`: every file changed by the child and why each belongs in the same scope.
- `Behavior`: what the change does for the application or test suite.
- `Issue`: the issue acceptance criterion or scope item it implements.
- `Decisions`: the relevant decision IDs from `decisions.md`, or `None`.
- `Validation`: the targeted command or review used to verify it.

The description can be a multi-line Jujutsu description when useful. The change identifier should remain concise, while the body should make the rationale and boundaries clear.

## Test-Driven Slice Workflow

Every vertical slice follows a test-first, green-child sequence:

1. Translate the approved issue acceptance criteria into a small end-to-end test suite.
2. Define focused unit or integration tests for logic and boundaries when they provide useful diagnosis.
3. Decide whether the tests can pass against existing behavior. If not, include the tests and the smallest required implementation in the same coherent child edit.
4. Construct the child test-first and use a temporary red state if useful, but do not leave the child as a failing, uncompilable, or unfinished result.
5. Run the relevant focused tests, formatting checks, and compilation checks before completing the child.
6. Run the complete repository check suite before moving to the next child when the repository supports it.
7. Refactor only while the tests remain green and the change stays within the approved issue.

End-to-end tests are the primary acceptance signal for each vertical slice. They should exercise the user-visible flow through the application's real HTTP boundary and an isolated test data store or equivalent test environment. Use the simplest reliable test harness; do not build a speculative testing framework.

Unit and integration tests complement, rather than replace, end-to-end coverage. Do not duplicate every assertion at every level. Put user outcomes in end-to-end tests, detailed edge cases in unit tests, and wiring or persistence behavior in integration tests.

Tests may be separate child edits when they pass against existing behavior. When new tests require new implementation, include the test and implementation files in the same child edit. Never include unrelated implementation in a test child merely to avoid a separate scope decision.

If a required test cannot be written or fails for a reason unrelated to the approved behavior, stop and ask the user rather than weakening the test or proceeding without coverage.

## Issue Lifecycle

### 1. Establish Context

Before creating an issue root:

1. Read the issue in `issues.md`.
2. Read the related user story in `user_stories.md`.
3. Read all relevant decisions in `decisions.md`.
4. Inspect the existing implementation and tests.
5. Run `jj status` and `jj log` to confirm the current history.
6. Check issue dependencies and confirm that the user wants to start this issue.
7. Identify the first end-to-end behavior and the smallest test suite that will express it.

If there are unrelated working-copy changes, an existing issue root for a different issue, unresolved child-edit changes, or unclear requirements, stop and ask the user before creating or modifying anything.

### 2. Create the Issue Root

Create one empty issue root with `jj new -m`. Confirm its description and parent with read-only Jujutsu commands. Do not modify files in the root.

The root remains the issue's stable anchor throughout implementation. All code, tests, migrations, templates, assets, generated output, and issue-status documentation changes belong in child edits.

### 3. Propose One Child Edit

Before each child edit, describe:

- Every file to change and why the files belong together.
- The smallest coherent behavior being added or changed.
- What the user will observe or what the test will prove.
- The exact issue acceptance criteria covered.
- The relevant decision IDs.
- The validation to run.

Bundle source code and tests into one child edit when that is necessary to leave the child compiling and passing. Keep the child limited to the smallest complete behavior.

### 4. Make and Verify the Child Edit

Create a new child edit on top of the current edit, then make the file change. Use `apply_patch` for manual edits. After editing:

1. Check `jj status` and confirm only the described files changed in the child.
2. Inspect `jj diff` for scope, correctness, and accidental edits.
3. Run formatting checks for changed source files.
4. Run the narrowest relevant tests and compilation checks.
5. Run broader verification when the issue requires it, including `make check` when available.
6. Confirm that no unfinished or unrunnable functionality remains in the child.
7. Report every changed file, detailed change, issue linkage, decisions, validation, and any residual risk.

If a command changes more files than intended, stop and determine whether those files are required for the same complete behavior. Include them only when they are coherent and described; otherwise separate the work into child edits.

Do not create the next child until the current code child compiles and its relevant tests pass. Documentation-only children must still be complete and internally consistent.

### 5. Continue Within The Issue

After every child edit:

- Report the file, detailed change, issue linkage, decisions, validation, and any residual risk.
- Leave the child edit available for the user to inspect and discuss.
- Continue to the next planned child edit within the same issue without waiting for explicit approval.
- Pause if the user asks to stop, requests a change, identifies a conflict, or the implementation exposes an unresolved requirement.
- A review report is not a green-status exception: resolve compilation or test failures before continuing.

The user may review and discuss each child edit independently. Discussion does not authorize changing the issue root or moving to another issue.

### 6. Complete The Current Issue

After all child edits for the issue are complete and the issue has been verified against its acceptance criteria, tell the user that the issue is ready to review and squash. Do not squash the edits yourself. Do not edit the root to imitate a squash. Wait for the user to confirm that the child edits have been squashed into the issue root.

Only after that confirmation may the issue be considered complete for workflow purposes or may work begin on a new issue. Any status update in `issues.md`, `user_stories.md`, or another documentation file must itself be a separately described, single-file child edit unless the user explicitly directs otherwise.

### 7. Move to the Next Issue

Before starting another issue:

- Confirm the prior issue's child edits were reviewed and all requested changes are resolved.
- Confirm the user squashed them into the issue root.
- Confirm the issue root remains the final issue change.
- Confirm the issue is verified against its acceptance criteria.
- Confirm the next issue is approved and its dependencies are satisfied.
- Obtain explicit user confirmation before creating the next issue root.

If any condition is missing, remain on the current issue and ask the user what to do.

## Scope And Design Rules

- Implement the smallest change that satisfies the approved issue.
- Preserve the decisions in `decisions.md`; if implementation exposes a conflict, stop and ask rather than silently changing the decision.
- Prefer the existing project structure and patterns.
- Keep server-rendered HTML and HTMX behavior aligned with the approved stories.
- Start each slice with end-to-end tests and add unit or integration tests where they provide focused coverage.
- Keep tests in separate children when they pass independently; otherwise include them with the required implementation in one coherent child.
- Never leave a child with unfinished, uncompilable, failing, or unrunnable code.
- Treat generated files as ordinary files for scope and review; generated changes must still be isolated to one child edit.
- Do not add compatibility layers, abstractions, or future features without an approved issue or explicit user approval.

## Safe Read-Only Commands

These commands are useful for review and do not alter Jujutsu history:

```text
jj status
jj log
jj diff
jj show
```

Do not use a read-only command as a reason to edit an ancestor or the issue root. History changes and squashing remain user-controlled.
