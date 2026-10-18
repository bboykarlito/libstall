// Package stall defines the behaviour and execution contracts for agents.
// Applications own orchestration and explicitly pass data between executions.
package stall

import "context"

// Task is the input to one execution.
type Task struct {
	Input string
}

// Request is a task prepared by a Behaviour for execution by a Harness.
type Request struct {
	Task         Task
	Instructions string
}

// Result contains execution output and reported usage.
// A result may accompany an error, preserving available output and usage.
type Result struct {
	Output string
	Usage  Usage
}

// Usage contains token counts for one execution.
// This initial representation does not distinguish unavailable counts from zero;
// usage availability and monetary cost remain to be designed before adapters land.
type Usage struct {
	InputTokens  int64
	OutputTokens int64
}

// Behaviour defines task preparation and domain-specific output validation.
type Behaviour interface {
	// Prepare produces instructions or rejects an unsuitable task with an error.
	Prepare(ctx context.Context, task Task) (Request, error)

	// Validate checks the output against the prepared task and instructions.
	// A validation error must not cause execution output or usage to be discarded.
	Validate(ctx context.Context, request Request, result Result) error
}

// Harness adapts an underlying agent runtime.
// Runtime-specific configuration belongs in the adapter's constructor.
type Harness interface {
	// Execute starts an isolated conversational session without implicit history.
	// Session isolation does not imply filesystem or process sandboxing.
	// On failure, Execute returns any available partial output and usage with the error.
	Execute(ctx context.Context, request Request) (Result, error)
}
