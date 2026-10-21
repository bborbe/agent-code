---
status: completed
spec: [001-package-layout-docs]
summary: Added a Package layout section to README describing pkg/factory and pkg/steps, plus an Unreleased changelog entry, with all verification checks passing.
execution_id: be9d3b2c-65e2-5dbc-a8d8-b94acb40c014-exec-001-spec-001-package-layout-docs
dark-factory-version: dev
created: "2026-10-21T14:50:08Z"
queued: "2026-10-21T14:52:50Z"
started: "2026-10-21T14:52:51Z"
completed: "2026-10-21T15:12:20Z"
branch: dark-factory/package-layout-docs
---

<summary>
- The project README gains a short "Package layout" section that says, in one line each, what the two library packages do.
- A reader can name the responsibility of each library package without opening a Go file.
- The section sits immediately before the existing "Related" pointer list.
- No existing README section is reordered, reworded, or removed.
- No Go source file, doc comment, or test is changed.
- The GoDoc coverage that exists today is proven unchanged by a verification command that fails if any exported declaration loses its comment or a package is emptied.
- The changelog records the new README section under an Unreleased heading.
- The repository's standard pre-commit gate still passes.
</summary>

<objective>
Document the two library packages in the README so a new agent author can understand the repository layout without reading Go source, and prove at verification time that the existing GoDoc coverage on every exported declaration is intact. This is a documentation-only change: one inserted README section, one changelog entry, and no code changes.
</objective>

<context>
Read these before making changes:

- `specs/in-progress/001-package-layout-docs.md` — the approved spec, including its Acceptance Criteria, Constraints, and Failure Modes. This prompt implements it.
- `README.md` — note the current top-level section order: `## How It Works`, `## Env Vars`, `## Creating a New Agent`, `## Local Quick Test`, `## Links`, `## Related`, `## License`. The new section goes between `## Links` and `## Related`.
- `CHANGELOG.md` — note the preamble at the top and that the newest released section is `## v0.2.9`; `## Unreleased` goes after the preamble and before `## v0.2.9`.
- `pkg/factory/factory.go` — confirms what the factory package actually wires: `CreateKafkaResultDeliverer`, `CreateFileResultDeliverer`, `CreateAgent`, `CreateSyncProducerProvider`, `CreateAgentProvider`, and the `SyncProducerProvider` interface.
- `pkg/steps/steps.go` — confirms the three steps: `PlanStep` (`Name()` returns `"plan"`), `ExecuteStep` (`"execute"`), `VerifyStep` (`"verify"`), all pure Go with no AI call.

Coding guides (in-container paths):
- `/home/node/.claude/plugins/marketplaces/coding/docs/readme-guide.md` — README audience and structure.
- `/home/node/.claude/plugins/marketplaces/coding/docs/changelog-guide.md` — `## Unreleased` placement and conventional entry prefixes.
- `/home/node/.claude/plugins/marketplaces/coding/docs/go-doc-best-practices.md` — why the existing doc comments must not be altered.
</context>

<requirements>
1. Insert exactly one `## Package layout` section into `README.md`, positioned immediately before the existing `## Related` heading (i.e. after the `## Links` section and its admin-endpoint bullet list). Do not insert it anywhere else and do not add a second copy of the heading.

2. The section body contains one prose line for `pkg/factory` and one prose line for `pkg/steps`. A two-item bullet list is the expected form. Do not use nested bullets and do not use a fenced code block. Suggested content (word the lines so the AC1 assertions in `<verification>` still match):
   - `` `pkg/factory/` — wires the binary's concrete dependencies: the Kafka and file result deliverers, the Kafka sync-producer provider, and the agent provider. ``
   - `` `pkg/steps/` — the three pure-Go phase steps of the reference math agent (plan, execute, verify); none of them makes an AI call. ``

3. The `pkg/factory` line must name at least two of: result deliverer(s), sync-producer provider, agent provider. The `pkg/steps` line must name at least two of: plan, execute, verify. The AC1 assertion regex in `<verification>` is case-sensitive on the token `pkg/factory` / `pkg/steps` and accepts `Deliverer`/`deliverer`, `Sync`/`sync`, `Agent`/`agent`, `Plan`/`plan`, `Execute`/`execute`, `Verify`/`verify`.

4. Do NOT change any existing README content: do not reorder sections, do not reword `## How It Works` or `## Creating a New Agent`, and do not fix the stale paths in `## Related`. The only README change is the inserted section.

5. Do NOT modify any Go file, doc comment, test, mock, `go.mod`, or `go.sum`. Every exported top-level declaration and every exported method under `pkg/`, `cmd/`, and the root `main.go` must keep the `//` comment line it carries today, attached with no blank line inserted between comment and declaration. Do not add doc comments to unexported declarations.

6. Add an `## Unreleased` section to `CHANGELOG.md` immediately after the preamble block (`All notable changes ...` plus any blank lines) and before `## v0.2.9`. Under it add exactly one entry naming the new README section, in conventional-prefix form, for example:
   ```markdown
   ## Unreleased

   - docs: add "Package layout" section to README describing pkg/factory and pkg/steps
   ```
   The text `Package layout` must appear within the five lines following the `## Unreleased` heading. Do not modify the preamble or any released section.

7. Failure handling — apply these if the corresponding situation occurs:
   - If `make precommit`'s `format` target rewrites a Go file under `pkg/`, `cmd/`, or the root `main.go` (import order, line wrapping, or a comment detached by a blank line), discard that formatting-only rewrite by restoring the file to its committed content (`git checkout -- <path>`, or re-check out the file from the merge-base if `.git` is unavailable) and re-run `make precommit`.
   - If `make precommit`'s `generate` target rewrites `mocks/` or `go mod tidy` touches `go.mod`/`go.sum`, restore those paths (`git checkout -- mocks go.mod go.sum`) and re-run `make precommit`.
   - If a `make precommit` target fails for a network-dependent reason (tool download, vulnerability database), re-run `make precommit` once. If it fails again, report the failing target and its output; do not edit Go source to work around it.
   - If you are tempted to "improve" the stale `## Related` paths or the "How It Works" narrative while editing nearby, do not: revert any such line and keep only the inserted section.
   - If the AC1 assertions report `FAIL: layout heading missing, duplicated, or after Related` or `FAIL: layout section is not immediately before Related`, move the section to immediately before `## Related` and re-run until the block prints `OK` with exit 0.
   - If the AC1 assertions report `FAIL: pkg/factory line missing wiring description` or its `pkg/steps` variant, expand that package's line to name the required responsibilities (two tokens minimum) and re-run until `OK`.
</requirements>

<constraints>
- Documentation-only change: no exported API change, no behavior change, no test change.
- Do NOT add, reword, or move GoDoc comments on identifiers that already carry one — coverage is complete; this work only locks it.
- Do NOT add doc comments to unexported declarations (the `application` structs in the two entry points, methods on unexported receiver types).
- Do NOT fix the stale paths in the README's "Related" section, and do NOT rewrite the stale "How It Works" / "Creating a New Agent" narrative.
- Do NOT document the command entry points or any other directory in the new section — it covers the two library packages only.
- Do NOT add a lint rule, CI job, Makefile target, or committed script for the doc-comment check. The check is a verification-time command only.
- Do NOT restructure the README — no reordering or rewriting of existing sections; the change is one inserted section.
- The new heading uses `##`, the README's existing top-level section level.
- `README.md` stays in English.
- The change touches no file other than `README.md`, `CHANGELOG.md`, and the spec file itself.
- Do NOT commit — dark-factory handles git.
</constraints>

<verification>
Run all of the following from the repository root and record the exit code of each.

1. Baseline gate (AC3) — formatting, module hygiene, tests, lint, vet, vulnerability and license checks:

```bash
make precommit
```

2. GoDoc regression lock (AC2). It inspects every non-test Go file under `pkg/` and `cmd/` plus root `*.go`, requires each exported declaration to sit directly under a `//` comment line, counts what it inspected, and fails if the count drops below the 24 present at v0.2.9:

```bash
cat > /tmp/godoc-lock.awk <<'AWK'
function recv(line,   r,t,n,a) {
  match(line, /^func[ \t]+\([^)]*\)/)
  r = substr(line, RSTART, RLENGTH)
  sub(/^func[ \t]+\(/, "", r); sub(/\)$/, "", r)
  gsub(/^[ \t]+|[ \t]+$/, "", r)
  n = split(r, a, /[ \t]+/)
  t = a[n]; sub(/^\*/, "", t)
  return t
}
FNR==1 { prev="" }
/^func[ \t]+[A-Z]/ { n++; if (prev !~ /^[[:space:]]*\/\//) { print FILENAME ":" FNR ": " $0; bad=1 } }
/^func[ \t]+\(/ { if (recv($0) ~ /^[A-Z]/) { n++; if (prev !~ /^[[:space:]]*\/\//) { print FILENAME ":" FNR ": " $0; bad=1 } } }
/^(type|var|const)[ \t]+[A-Z]/ { n++; if (prev !~ /^[[:space:]]*\/\//) { print FILENAME ":" FNR ": " $0; bad=1 } }
{ prev=$0 }
END { printf "documented_exported=%d\n", n; exit (bad || n < 24) }
AWK
awk -f /tmp/godoc-lock.awk $(find pkg cmd -name '*.go' -not -name '*_test.go' | sort) $(find . -maxdepth 1 -name '*.go' -not -name '*_test.go')
```

Expected on success: the single line `documented_exported=24`, exit 0. A removed or blank-line-detached doc comment prints the offending `file:line` and exits 1; fewer than 24 exported declarations prints `documented_exported=<n>` and exits 1.

3. README section assertions (AC1):

```bash
P=$(grep -n '^## Package layout' README.md | cut -d: -f1)
R=$(grep -n '^## Related' README.md | cut -d: -f1)
test "$(grep -c '^## Package layout' README.md)" -eq 1 && test -n "$P" && test -n "$R" && test "$P" -lt "$R" || { echo "FAIL: layout heading missing, duplicated, or after Related"; exit 1; }
test "$(awk -v p="$P" 'NR>p && /^## /{print NR; exit}' README.md)" -eq "$R" || { echo "FAIL: layout section is not immediately before Related"; exit 1; }
sed -n "$((P+1)),$((R-1))p" README.md | grep -qE 'pkg/factory.*([Dd]eliverer).*([Ss]ync|[Aa]gent)|pkg/factory.*([Ss]ync|[Aa]gent).*([Dd]eliverer)' || { echo "FAIL: pkg/factory line missing wiring description"; exit 1; }
sed -n "$((P+1)),$((R-1))p" README.md | grep -qE 'pkg/steps.*([Pp]lan).*([Ee]xecute|[Vv]erify)|pkg/steps.*([Ee]xecute).*([Vv]erify)' || { echo "FAIL: pkg/steps line missing phase description"; exit 1; }
echo "OK"
```

Expected on success: `OK`, exit 0.

4. CHANGELOG assertion (AC5):

```bash
grep -q '^## Unreleased' CHANGELOG.md && grep -A5 '^## Unreleased' CHANGELOG.md | grep -q 'Package layout' || { echo "FAIL: no Unreleased entry naming the package-layout section"; exit 1; }
```

Expected on success: exit 0.

5. File-scope check (AC4). This repo's `.dark-factory.yaml` sets `workflow: direct` and does not set `hideGit: true`, so `.git` is available in the container and `git status` is a valid check:

```bash
git status --porcelain
```

Expected: only `README.md`, `CHANGELOG.md`, and the spec file appear as modified; no Go, mock, `go.mod`, or `go.sum` changes. Note that untracked daemon artifacts (`.dark-factory.lock`, `.dark-factory.log`, `specs/log/`) are pre-existing and out of scope.
</verification>
