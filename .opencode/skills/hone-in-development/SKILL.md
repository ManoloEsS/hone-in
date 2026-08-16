---
name: hone-in-development
description: Use for Hone In implementation work when changes must follow the approved issue backlog, Jujutsu issue roots, single-file child edits, explicit review, and user-controlled squashing.
---

# Hone In Development Workflow

This skill defines the implementation workflow for the Hone In repository. Follow it for application code, tests, templates, static assets, migrations, generated code, and implementation-related documentation.

## Non-Negotiable Rules

- Work on one approved issue at a time.
- Do not implement an issue whose status is not `Approved` or `In progress` in `issues.md`.
- Do not move to a later issue until the current issue's child edits have all been explicitly approved and the user has squashed them into the current issue root.
- The user controls Jujutsu squashing. Never run `jj squash`, `jj abandon`, rebase operations, or equivalent history-rewriting commands unless the user explicitly asks for that operation.
- Never edit the issue root after it is created unless the user explicitly permits it.
- Never edit a parent or ancestor of the issue root unless the user explicitly permits it.
- Keep every implementation edit limited to exactly one file. If a coherent change requires multiple files, split it into multiple child edits.
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

The child edit may change one file only. Its description must explain the change in enough detail for review:

- `Change`: the exact file-level change being made.
- `Behavior`: what the change does for the application or test suite.
- `Issue`: the issue acceptance criterion or scope item it implements.
- `Decisions`: the relevant decision IDs from `decisions.md`, or `None`.
- `Validation`: the targeted command or review used to verify it.

The description can be a multi-line Jujutsu description when useful. The change identifier should remain concise, while the body should make the rationale and boundaries clear.

## Issue Lifecycle

### 1. Establish Context

Before creating an issue root:

1. Read the issue in `issues.md`.
2. Read the related user story in `user_stories.md`.
3. Read all relevant decisions in `decisions.md`.
4. Inspect the existing implementation and tests.
5. Run `jj status` and `jj log` to confirm the current history.
6. Check issue dependencies and confirm that the user wants to start this issue.

If there are unrelated working-copy changes, an existing issue root, pending unapproved child edits, or unclear requirements, stop and ask the user before creating or modifying anything.

### 2. Create the Issue Root

Create one empty issue root with `jj new -m`. Confirm its description and parent with read-only Jujutsu commands. Do not modify files in the root.

The root remains the issue's stable anchor throughout implementation. All code, tests, migrations, templates, assets, generated output, and issue-status documentation changes belong in child edits.

### 3. Propose One Child Edit

Before each child edit, describe:

- The single file to change.
- The smallest coherent behavior being added or changed.
- What the user will observe or what the test will prove.
- The exact issue acceptance criteria covered.
- The relevant decision IDs.
- The validation to run.

Do not bundle source code and tests into one child edit. Create separate child edits when both files need changes.

### 4. Make and Verify the Child Edit

Create a new child edit on top of the current edit, then make the file change. Use `apply_patch` for manual edits. After editing:

1. Check `jj status` and confirm exactly one file changed in the child.
2. Inspect `jj diff` for scope, correctness, and accidental edits.
3. Run the narrowest relevant test, build, formatter, or static check.
4. Run broader verification when the issue requires it.
5. Report the file, detailed change, issue linkage, decisions, validation, and any residual risk.

If a command changes more than the intended file, stop and cleanly separate the work into child edits rather than accepting a multi-file child.

### 5. Wait for Review

After every child edit, stop and wait for the user's explicit review. Do not create the next child edit after the user rejects or requests changes until the current child has been addressed.

All child edits for the issue must be reviewed and explicitly approved. Approval of one child does not approve the issue root or authorize moving to another issue.

### 6. User Squashes the Issue

After all child edits for the issue are approved, tell the user that the issue is ready to squash. Do not squash the edits yourself. Do not edit the root to imitate a squash. Wait for the user to confirm that the approved child edits have been squashed into the issue root.

Only after that confirmation may the issue be considered complete for workflow purposes. Any status update in `issues.md`, `user_stories.md`, or another documentation file must itself be a separately described, single-file child edit unless the user explicitly directs otherwise.

### 7. Move to the Next Issue

Before starting another issue:

- Confirm the prior issue's child edits were approved.
- Confirm the user squashed them into the issue root.
- Confirm the issue root remains the final issue change.
- Confirm the issue is verified against its acceptance criteria.
- Confirm the next issue is approved and its dependencies are satisfied.

If any condition is missing, remain on the current issue and ask the user what to do.

## Scope And Design Rules

- Implement the smallest change that satisfies the approved issue.
- Preserve the decisions in `decisions.md`; if implementation exposes a conflict, stop and ask rather than silently changing the decision.
- Prefer the existing project structure and patterns.
- Keep server-rendered HTML and HTMX behavior aligned with the approved stories.
- Add tests as separate one-file child edits when behavior needs coverage.
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
