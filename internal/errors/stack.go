package errors

import (
	"context"

	"github.com/rl404/fairy/errors"
)

var stacker errors.ErrStacker

func init() {
	stacker = errors.New()
}

func Init(ctx context.Context) context.Context {
	return stacker.Init(ctx)
}

func Wrap(ctx context.Context, err error, errs ...error) error {
	return stacker.Wrap(ctx, err, errs...)
}

func Get(ctx context.Context) []string {
	stacks := stacker.Get(ctx).([]string)

	for i, j := 0, len(stacks)-1; i < j; i, j = i+1, j-1 {
		stacks[i], stacks[j] = stacks[j], stacks[i]
	}

	return stacks
}
