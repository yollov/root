package step

import (
	"context"
	"fmt"
	"strings"

	"github.com/boomfunc/base/tools/executor"
)

// Parallel is group of container running in parallel (asynchronous)
type Parallel []Interface

// NewParallel creates group of steps, running in parallel
// (just utility function)
func NewParallel(steps ...Interface) Interface {
	switch len(steps) {
	// TODO: what? what empty response?
	case 0:
		return nil
	case 1: // there is single step, no need to group
		return steps[0]
	default: // many steps, grouping this
		return Parallel(steps)
	}
}

// Run implements Step interface
// just run one by one
func (group Parallel) Run(ctx context.Context) error {
	stepsUp := make([]executor.OperationFunc, len(group))

	// generate list of functions for executor
	for i, step := range group {
		stepsUp[i] = step.Run
	}

	// run concurrently
	return executor.New(
		ctx,
		executor.Operation(stepsUp, nil, false),
	).Run()
}

// String implements fmt.Stringer interface
func (group Parallel) String() string {
	parts := make([]string, len(group))

	for i, step := range group {
		parts[i] = fmt.Sprintf("%s", step)
	}

	return fmt.Sprintf("PARALLEL(\n%s\n)", strings.Join(parts, ",\n"))
}