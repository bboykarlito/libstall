# ADR-001: Thin execution supervision with application-owned orchestration

## Status

Accepted as the initial design direction. Exact Go signatures and implementation details are not yet specified.

## Date

2026-10-18

## Context

libstall is intended for Go engineers building agents, with the author's own agents as its first use cases. It should initially be a minimal library rather than a framework.

The desired entry point is conceptually:

```go
result, err := runner.Run(ctx, behaviour, harness, task)
```

Agents may involve several stages with different instructions, constraints, and runtimes. Existing agent runtimes already own the model/tool interaction loop; duplicating that loop or introducing a workflow language would undermine the thin-library goal.

## Decision

### Responsibilities

| Component | Responsibility |
| --- | --- |
| Application code | Explicitly sequence calls, choose components for each call, pass outputs, branch, and retry. |
| Runner | Hold reusable policy; supervise admission, budget accounting, cancellation, execution, validation, and result reporting. |
| Behaviour | Prepare instructions and execution requirements, reject unsuitable tasks, and validate output against domain-specific success criteria. |
| Harness | Adapt an underlying agent runtime and execute within its model/tool loop. |

Runner is a concrete struct, not an interface. Multiple runner instances can carry different policies without requiring different implementations.

Behaviour and Harness are interfaces. Harness describes execution, not a universal set of configuration setters. Runtime-specific configuration belongs in adapter constructors.

A stage is an application-level concept represented by one `Run` call, not a graph node or workflow object supplied by the library.

### Lifecycle

1. Behaviour prepares the task or rejects it as unsuitable.
2. Runner checks execution policy and whether the harness can satisfy the requested capabilities and constraints.
3. Harness executes in an isolated session.
4. Behaviour validates the completed output.
5. Runner returns the result, usage, and any error.

Rejection or failed admission prevents harness execution. Validation failure preserves the output and usage so the caller can decide what to do next. Exact error types and partial-result semantics for other failures remain to be designed.

### Execution scope

- Every call starts with fresh budget accounting and execution state.
- A configured cap applies independently to each call, not cumulatively across a runner's lifetime.
- Run-specific mutable state does not live on the reusable runner.
- Every call uses an isolated session with no implicit conversation history.
- Application code passes previous outputs explicitly when needed.
- A retry is an explicit new call with a fresh budget and isolated session; the library does not retry automatically.

Session isolation concerns conversational state. It does not by itself promise filesystem, process, or security sandbox isolation.

Concurrent-use guarantees must be specified and tested before being advertised; avoiding runner-local execution state does not automatically make harnesses or behaviours concurrency-safe.

### Admission and budgets

Behaviour decides task suitability. Runner decides policy admission, including budget-related checks using available estimation and runtime information.

Keep these concepts distinct:

- **Estimated cost:** an admission heuristic, not a guarantee.
- **Recorded spend:** usage/cost reported for the execution.
- **Enforceable limit:** the spending boundary the underlying runtime can actually prevent exceeding.

Runner owns budget policy, but enforcement depends on harness capabilities. It must not present retrospective accounting or best-effort cancellation as a hard spend guarantee. How capabilities, estimates, and enforcement strength are represented remains an API design question.

Required constraints must not be silently weakened. For example, a tool allowlist must be enforced by the runtime, not merely expressed as a prompt. Unsupported requirements fail explicitly.

## Alternatives considered

### Runner orchestrates a plan or stage graph

Inspectable workflows could eventually be useful, but introduce workflow state, branching semantics, and a second orchestration abstraction. Rejected for v0: ordinary Go already provides explicit control flow.

### Runner is an interface

Could support alternative supervision implementations. Not needed initially: different policies require multiple configured instances, not multiple implementations. Keep a concrete struct until a real substitution need appears.

### Shared budget across a runner's calls

Useful for workflow-wide accounting, but makes runner reuse stateful and couples independent executions. Rejected for v0 in favour of a fresh per-call cap. Applications must account for aggregate spending themselves if needed.

### Persistent sessions between stages

Could preserve conversational context implicitly, but couples calls and obscures inputs. Rejected for v0 in favour of isolated sessions and explicit data flow.

### Universal harness configuration

A uniform runtime configuration language risks becoming a lowest-common-denominator abstraction. Prefer adapter-specific constructors with a small execution interface.

## Consequences

- Application code remains responsible for orchestration and aggregate workflow policy.
- Behaviours can be separated from execution machinery, subject to harness capability compatibility.
- The library supervises the runtime's agent loop rather than reimplementing it.
- Honest budget semantics and capability checks are essential to avoid misleading guarantees.
- Multi-stage agents can begin as ordinary Go functions composing calls; a workflow abstraction is deferred until concrete needs justify it.

## MVP scope

Design the smallest Go API supporting this lifecycle, then validate it with one real behaviour and one harness adapter.

Not doing in v0:

- Workflow plans, graphs, or a stage DSL.
- Automatic retries or implicit branching.
- Persistent conversational sessions or execution persistence.
- Shared cross-call budgets.
- A universal tool-definition or runtime-configuration format.
- Multiple runner implementations without an observed need.

## Open questions

- What are the minimal Behaviour and Harness method signatures?
- What shapes should task, prepared request, result, usage, and validation errors take?
- How are capabilities and enforceable constraints expressed without recreating runtime configuration?
- How are estimates obtained, and what happens when cost or usage information is unavailable?
- How are money, currency, and enforcement strength represented?
- What are cancellation, cleanup, and partial-result semantics?
- Which real agent and harness adapter should validate the design first?
- What is the published Go module path?
