---
status: completed
approved: "2026-10-20T19:18:11Z"
generating: "2026-10-21T14:49:33Z"
prompted: "2026-10-21T14:51:20Z"
verifying: "2026-10-21T15:12:20Z"
completed: "2026-10-21T15:12:26Z"
branch: dark-factory/package-layout-docs
---

## Summary

- The project README gains a short "Package layout" section that says, in one line each, what the two library packages do.
- The section sits immediately before the existing "Related" pointer list, so a new agent author meets it before a list of pointers that mostly leave this repository.
- The GoDoc half of the original request is already satisfied at the current release: every exported top-level declaration and every method on an exported receiver under the two library packages and the command entry point carries a doc comment. Nothing is added there.
- Instead, this work records that coverage as a verification-time check that fails if any such declaration loses its doc comment or a package is emptied. The check is not wired into the build, so it fires only when a verifier runs it.
- No runtime behavior, exported API, or test changes.

## Problem

Someone copying this reference agent as a template has to open the Go source to learn that one package wires the concrete dependencies and the other holds the phase steps. The README's only package pointers sit in a "Related" list whose entries (a shared prompt library, a shared delivery library, a controller, an executor) do not exist in this repository, so the two packages that do exist are invisible from the README. The original request also asked for GoDoc comments on every exported identifier; an audit at the current release shows that work is already complete — all twenty-four exported top-level declarations and exported methods under `pkg/` carry a comment beginning with the identifier name, and `cmd/` contributes no exported declaration at all; the only undocumented declarations are unexported application structs and methods on unexported receiver types, which is correct. Re-adding those comments would be busywork; the live risk is that coverage erodes silently. The request therefore becomes: describe the layout, and prove the coverage stays intact.

## Goal

A reader of the README can name the responsibility of each of the two library packages without opening a Go file, and the GoDoc coverage that exists today is provably unchanged — any edit that removes a doc comment from an exported declaration or exported method, or empties a package, makes the documented check exit non-zero.

## Non-goals

- Do NOT add, reword, or move GoDoc comments on identifiers that already carry one — coverage is complete; this work only locks it.
- Do NOT add doc comments to unexported declarations (the `application` structs in the two entry points, methods on unexported receiver types). They are correctly undocumented and stay that way.
- Do NOT fix the stale paths in the README's "Related" section, and do NOT rewrite the stale "How It Works" / "Creating a New Agent" narrative that still describes a Claude runner while this repository ships a pure-code agent. Both are separate concerns; this spec only inserts the new section.
- Do NOT document the command entry points or any other directory in the new section — it covers the two library packages only.
- Do NOT add a lint rule, CI job, Makefile target, or committed script for the doc-comment check. Enabling the linter's exported-comment rule is a separate decision; the check here is a verification-time command.
- Do NOT restructure the README — no reordering or rewriting of existing sections; the change is one inserted section.

## Acceptance Criteria

- [ ] AC1 — README contains exactly one `## Package layout` heading whose next `##` heading is `## Related`, and the section body carries a `pkg/factory` line naming at least two of the deliverer/sync-producer/agent-provider responsibilities and a `pkg/steps` line naming at least two of plan/execute/verify. Evidence: file-content grep — the AC1 assertion block in Verification prints `OK` and exits 0; a bare `touch`, a bare path list, a one-token stub line, a section placed after `## Related`, or a section inserted anywhere earlier than immediately before `## Related` prints `FAIL` and exits 1.
- [ ] AC2 — the GoDoc coverage check passes: every exported top-level declaration and every method on an exported receiver type under `pkg/`, `cmd/`, and the root `main.go` (test files excluded) is immediately preceded by a `//` comment line. Evidence: exit code 0 plus the stdout line `documented_exported=24` (or a larger count); a removed or blank-line-detached doc comment exits 1 and prints the offending `file:line`, while a scan finding fewer than 24 exported declarations exits 1 with only the `documented_exported=<n>` line.
- [ ] AC3 — `make precommit` exits 0. Evidence: exit code.
- [ ] AC4 — the change touches no file other than `README.md`, `CHANGELOG.md`, and this spec. Evidence: negative file diff — `git status --porcelain` lists only those paths.
- [ ] AC5 — `CHANGELOG.md` carries an `## Unreleased` section with one entry naming the README package-layout section. Evidence: `grep -n '^## Unreleased' CHANGELOG.md` returns ≥1 match and the entry below it names `Package layout`.

## Verification

### Container-executable (runs inside the YOLO container at prompt time)

Baseline gate — formatting, module hygiene, tests, lint, vet, vulnerability and license checks:

```bash
make precommit
```

GoDoc regression lock (AC2). Run from the repository root. It inspects every non-test Go file under `pkg/` and `cmd/`, requires each exported declaration to sit directly under a `//` comment line, counts what it inspected, and fails if the count drops below the 24 present at v0.2.9 — so deleting or emptying a package cannot pass it:

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

Scope note: the awk patterns see single-line top-level declarations only. An exported name declared inside a grouped `type (` / `const (` / `var (` block, and exported struct fields, are outside this check — no such declaration exists at v0.2.9, and adding one is out of scope for this change.

Expected on success: the single line `documented_exported=24`, exit 0.

README section assertions (AC1):

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

CHANGELOG assertion (AC5):

```bash
grep -q '^## Unreleased' CHANGELOG.md && grep -A5 '^## Unreleased' CHANGELOG.md | grep -q 'Package layout' || { echo "FAIL: no Unreleased entry naming the package-layout section"; exit 1; }
```

Expected on success: exit 0.

No operator-executable rung: this change ships a merged README diff and a merged spec record, and nothing runs in a cluster.

## Desired Behavior

1. `README.md` gains exactly one `## Package layout` heading, placed immediately before `## Related` — that is, after the existing admin-endpoints `## Links` section.
2. Under that heading, one line for `pkg/factory/` describing it as where the binary's concrete dependencies are wired: the Kafka and file result deliverers, the Kafka sync-producer provider, and the agent provider.
3. Under that heading, one line for `pkg/steps/` describing it as the three pure-Go phase steps of the reference math agent — plan, execute, verify — with no AI call in any of them.
4. Every exported top-level declaration and every exported method under `pkg/`, `cmd/`, and the root `main.go` keeps the doc comment it carries at v0.2.9: no comment is removed, no comment is detached from its declaration by an inserted blank line, and no new exported declaration appears without one.
5. `CHANGELOG.md` gains an `## Unreleased` section whose single entry names the new README package-layout section.

## Constraints

- Doc-comment coverage at v0.2.9 is complete and is the baseline the check enforces: 24 exported declarations — 18 in the step package (three result types, three step types, three constructors, and nine interface methods `Name`/`ShouldRun`/`Run`) and 6 in the factory package.
- Coverage means exported top-level declarations and methods on exported receiver types. Exported struct fields (`pkg/steps/steps.go` fields on `Plan`, `Result`, `Review`) carry no doc comment and are deliberately outside the check.
- The two `application` structs and the sync-producer provider's `Get` method remain undocumented: unexported receivers are outside GoDoc's contract.
- No exported API change, no behavior change, no test change.
- The new heading uses `##`, the README's existing top-level section level. Each package gets one line of prose — no nested bullets, no fenced block.
- Whether the section is a bullet list or a table row per package is the implementer's choice, constrained by the AC1 assertion block.
- The doc-comment check is a verification-time command. It is not committed as a script and not wired into the build.
- `README.md` stays in English and keeps its existing sections (`How It Works`, `Env Vars`, `Creating a New Agent`, `Local Quick Test`, `Links`, `Related`, `License`) in their current order and content.

## Failure Modes

| Trigger | Expected behavior | Recovery |
|---------|-------------------|----------|
| The section is appended after `## Related`, placed at the end of the README, or inserted anywhere earlier than immediately before `## Related` | AC1 exits 1 with `FAIL: layout heading missing, duplicated, or after Related` or `FAIL: layout section is not immediately before Related` | Move the section to immediately before `## Related`; re-run the AC1 block and confirm `OK` with exit 0 |
| The section lists the two paths without describing them, or gives each a single keyword | AC1 exits 1 with `FAIL: pkg/factory line missing wiring description` (or the `pkg/steps` variant) | Add the one-line description naming each package's responsibility (two tokens minimum), then re-run the AC1 block |
| `make format` inside `make precommit` rewrites a Go file under `pkg/`, `cmd/`, or the root `main.go` — import order, line wrapping, or a comment detached by a blank line | AC4 fails and lists the path; the GoDoc check additionally prints that `file:line` and exits 1 when a comment was detached | `git checkout -- <path>` to drop the formatting-only rewrite, then re-run `make precommit` and the check |
| `make precommit` fails after its `generate` target rewrote `mocks/` or `go mod tidy` touched `go.mod`/`go.sum` | AC4 fails on the leftover modifications | `git checkout -- mocks go.mod go.sum`, then re-run `make precommit` |
| The implementer also "improves" the stale `## Related` paths while editing nearby | Out-of-scope content change beyond the inserted section | Revert those lines; keep only the inserted `## Package layout` section |
| `make precommit` fails inside an unrelated target that needs the network (tool download, vulnerability database) | AC3 cannot be evaluated; the failing target and its output are recorded | Re-run `make precommit` once; if it fails again, report the failing target rather than editing Go source to work around it |

## Do-Nothing Option

The GoDoc half needs nothing — it is already done. The only cost of doing nothing is onboarding: a new agent author reads a "Related" list pointing at four paths that do not exist in this repository and must open the Go source to find where the dependency wiring and the step framework live. The fix is one README section, which is cheaper than the confusion it removes; deferring it buys nothing and the confusion recurs with every template copy.
