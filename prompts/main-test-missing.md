---
status: draft
created: "2026-08-09T16:35:00Z"
---

<summary>
- The root command package has no test file at all, unlike the sibling command package
- A test suite bootstrap is added so the package is covered by the standard test runner
- The sibling package already solved the same problem and its approach is followed
- A known CI crash makes one common technique unsafe here, and that is respected
- No production code changes
</summary>

<objective>
The root command package has a test suite that runs under the project's standard test tooling, matching the pattern already used by the other command package.
</objective>

<context>
Read `CLAUDE.md` for project conventions if present (this repo has none; conventions come from sibling bborbe repos and the existing test files).

Files to read before making changes (read ALL first):
- `cmd/run-task/main_test.go` — the sibling command package's test file. Its header comment documents why the compile-check spec was deliberately removed: it segfaults under `-race` on GitHub Actions, and because the test binary is itself in the command package, `go test` already fails when the package does not compile.
- `pkg/steps/steps_suite_test.go` — the canonical suite bootstrap shape used in this repo (fail handler, time and formatting setup, spec runner)
- `main.go` — the package under test

This repo uses Ginkgo and Gomega. Match the existing suite files rather than inventing a structure.
</context>

<requirements>
1. Create `main_test.go` in the repository root, in the external test package (`package main_test`).
2. Give it a suite bootstrap following the same shape as `pkg/steps/steps_suite_test.go` — register the fail handler and run the specs.
3. Do NOT add a `gexec.Build` compile-check spec. Reproduce the reasoning from `cmd/run-task/main_test.go` as a comment: it segfaults under `-race` on CI, and the compile check is redundant because the test binary lives in the same package.
4. Keep the suite minimal — the bootstrap plus any trivially true spec needed for the suite to be non-empty. Do not write speculative tests against `application` fields or the `Run` method; broadening coverage is a separate concern.
5. Add a bullet under `## Unreleased` in `CHANGELOG.md` using a conventional prefix.
</requirements>

<constraints>
- Only add `main_test.go` and edit `CHANGELOG.md`
- Do NOT modify `main.go`
- Do NOT commit — dark-factory handles git
- Existing tests must still pass
- Do not add a `gexec.Build` spec, for the documented CI reason
</constraints>

<verification>
make precommit
</verification>
