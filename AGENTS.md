# AGENTS.md

## Repository purpose

This repository is a Go Web learning laboratory.

It contains:

- book-based study projects;
- small Go Web labs;
- framework experiments;
- REST/API practice;
- testing and documentation practice.

The repository is connected to a larger learning and portfolio path for Go backend freelance work.

## Important documents

- Use `PLAN.md` as the source of truth for the current priority and active stage.
- Read the root `README.md` and the nearest project README only when they are relevant to the task.
- Do not load every project or stage document by default.

## Operating contract

The user's latest explicit instruction takes precedence over the defaults in this file.

- Do not rewrite large parts of a project without explicit request.
- Prefer small, reviewable changes.
- Inspect only the files and documents relevant to the requested task.
- Keep study projects close to their source material.
- Do not mix unrelated stages in one task.
- Do not introduce new frameworks unless the task asks for it.
- Keep README files accurate and runnable.
- Prefer standard Go tooling first.
- Do not start a later stage unless the user explicitly requests it or `PLAN.md` identifies it as
  active.
- When the prompt names a working mode, apply the matching contract below.

## Working modes

### Manager Mode

- Inspect the relevant code before editing.
- Make the smallest focused change that satisfies the request.
- Preserve the existing project style and relation to the source material.
- Run the narrowest relevant checks.
- Report the result and remaining risks.

### Documentation Mode

- Change only the requested documentation files.
- Do not change implementation files.
- Preserve technical meaning.
- Keep README commands accurate and runnable.
- Separate current behavior from planned behavior.

### Review Mode

- Inspect the requested change without editing files.
- List findings first and order them by severity.
- Include file and line references when possible.
- Mention test gaps and residual risks.

### Planning Mode

- Inspect only the context needed to prepare the plan.
- Define scope, stages, verification, and the definition of done.
- Do not implement or edit files unless the user explicitly changes modes.

### Book or Source Study Mode

- Summarize concepts in original wording without copying large source fragments.
- Keep study projects close to their source material without copying blindly.
- Separate ideas into applicable now, deferred, or not relevant.
- Do not implement source material until the user explicitly changes modes.

### Tutor Mode

- Treat the request as one interactive learning turn, not a persistent goal.
- Explain the concept and provide one focused exercise.
- Do not write the final implementation first.
- Do not edit files unless the user explicitly changes modes.
- Finish the response after the exercise and wait for the user.
- Waiting for the user is not a blocker.
- Do not poll the repository while waiting.
- Do not create a persistent goal or token budget.

### Pair Programmer Mode

- Treat the request as an interactive session, not a persistent goal.
- Let the user implement the first version.
- Before implementation, provide only scoped guidance and acceptance criteria, then wait.
- Review the diff only after the user says it is ready.
- Suggest minimal targeted fixes instead of rewriting the solution.
- Waiting for the user is not a blocker.
- Do not poll the repository while waiting.
- Do not create a persistent goal or token budget.

## Persistent goals and token budgets

- Create a persistent Codex goal (`/goal`) only when the user explicitly requests one.
- Set a token budget only when the user explicitly specifies it.
- A persistent goal must describe one bounded result that the agent can complete without waiting for
  user work.
- Do not use persistent goals for Tutor Mode or Pair Programmer Mode.
- Do not use pause as a substitute for ending an interactive turn.
- If the task requires user implementation or input, finish the current turn and wait normally.

## Verification contract

- Run checks from the module relevant to the task.
- Run the narrowest relevant check first.
- Run the module's full test suite when the scope or risk justifies it.
- Use a writable Go cache in restricted environments when needed:

```bash
GOCACHE=/tmp/go-web-labs-go-cache go test ./...
```

- Report the exact commands run and their results.
- Separate pre-existing failures from failures caused by the current change.
- If sandbox or environment restrictions prevent a check, explain the limitation and provide the
  exact local verification step.

## Completion contract

After an implementation or documentation change, summarize:

- files changed;
- behavior or documentation changed;
- tests and checks run;
- failures, residual risks, or work intentionally left for later.

Do not commit changes unless the user explicitly requests it.
