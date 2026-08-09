---
status: draft
created: "2026-08-09T16:35:00Z"
---

<summary>
- A factory function currently returns an error, which the project's factory rule forbids
- Factories are pure composition; construction that can fail belongs at the entry point
- The call site already has correct error handling, so the error path is preserved exactly
- The wrapper is a one-line pass-through, so removing it loses no behaviour
- No change to logging, metrics, or shutdown ordering
</summary>

<objective>
The `Create*` factory in the factory package no longer returns an error, and the Kafka sync producer is constructed at the entry point where its failure is already handled.
</objective>

<context>
Read `CLAUDE.md` for project conventions if present (this repo has none; conventions come from sibling bborbe repos and the rule cited below).

Files to read before making changes (read ALL first):
- `pkg/factory/factory.go` — `CreateSyncProducer` is declared immediately after the imports
- `main.go` — the only call site, inside the `Run` method; the surrounding `if err != nil` block records failure metrics and wraps the error

The governing rule is `go-factory/no-error-return`: "Factory functions (Create*) must not return error. Factories are pure composition — errors belong in main.go Run or behind a Provider interface." `main.go` is explicitly exempt from that rule, which is why the entry point is the correct home for this construction.
</context>

<requirements>
1. Delete the `CreateSyncProducer` function from `pkg/factory/factory.go` entirely. It is a single-statement pass-through to `libkafka.NewSyncProducerWithName(ctx, brokers, agentName)` with exactly one caller.
2. In `main.go`, replace the `factory.CreateSyncProducer(ctx, a.KafkaBrokers, agentName)` call with a direct call to `libkafka.NewSyncProducerWithName(ctx, a.KafkaBrokers, agentName)`.
3. Leave the surrounding error handling exactly as it is — the metrics recording and `errors.Wrap` call must be unchanged.
4. `main.go` already imports `libkafka`. Keep the `factory` import — other `factory.Create*` calls remain in that file; verify before touching it.
5. In `pkg/factory/factory.go`: keep the `libkafka` import (still used by `CreateKafkaResultDeliverer`'s `SyncProducer` parameter), and remove the `context` import — `CreateSyncProducer` is its only consumer in that file.
6. Add a bullet under `## Unreleased` in `CHANGELOG.md` using a conventional prefix.
</requirements>

<constraints>
- Only change `pkg/factory/factory.go`, `main.go`, and `CHANGELOG.md`
- Do NOT commit — dark-factory handles git
- Existing tests must still pass
- Use `errors.Wrap`/`errors.Errorf` from `github.com/bborbe/errors` — never `fmt.Errorf` or bare `return err`
- Do not introduce a Provider interface; direct construction at the entry point is the chosen approach
</constraints>

<verification>
make precommit
</verification>
