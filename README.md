# libstall

A thin Go library for building agents on top of existing agent runtimes.

The name comes from **harness**: a **stall** sits one level above it. libstall separates an agent's behaviour from the runtime executing it, while leaving orchestration in ordinary Go code.

## Status

Early API skeleton. The module is `github.com/bboykarlito/libstall`, with package name `stall` (Go 1.24.4 or later). Public types are experimental.

Implemented: `Task`, `Request`, `Result`, `Usage`, and the `Behaviour` and `Harness` interfaces. Runner execution, constraints, capability checks, budget enforcement, and adapters are not implemented yet. Token usage currently cannot distinguish unavailable counts from zero; availability and monetary cost semantics remain open.

## Development

```sh
git clone https://github.com/bboykarlito/libstall.git
cd libstall
go test ./...
go vet ./...
```

There are no runtime tests yet; these commands currently check package compilation and static diagnostics.

## Intended usage

Conceptual API, not executable code yet:

```go
result, err := runner.Run(ctx, behaviour, harness, task)
```

- **Runner:** a reusable concrete struct holding execution policy, including a per-run spend cap.
- **Behaviour:** an interface for preparing or rejecting a task and validating its output.
- **Harness:** an interface for executing through an underlying agent runtime; runtime configuration belongs in adapter constructors.

One call represents one stage, with an isolated session and fresh budget accounting. Application code explicitly chooses the runner, behaviour, and harness for each call, passes outputs to subsequent calls, and handles branching and retries.

The library governs execution, not orchestration.

## Execution lifecycle

```text
Prepare → check policy/capabilities → Execute → Validate → Result
```

Unsupported requirements fail explicitly. Validation failure preserves output and usage. Cost estimates are not guarantees, and budget enforcement cannot exceed the underlying harness's capabilities.

## Design documentation

- [ADR-001: Thin execution supervision with application-owned orchestration](docs/decisions/001-execution-model.md) — agreed responsibilities, guarantees, trade-offs, and open questions.

## Initial scope

Validate the design with one real behaviour and one harness adapter before expanding the API. No workflow engine, automatic retries, persistence, or universal tool-definition format in v0.

## Contributing

Start with the decision record before proposing API or implementation changes. Keep changes small and ground abstractions in real agent use cases. Record significant changes to the agreed design in a new ADR.
