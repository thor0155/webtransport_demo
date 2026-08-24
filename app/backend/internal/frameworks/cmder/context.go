package cmder

import (
	"context"

	"github.com/spf13/cobra"
)

type CmderContext interface {
	context.Context
	Args() []string
	Error(err error)
	GetErrors() []error
}

type cmderContext struct {
	context.Context
	args   []string
	errors []error
	cmd    *cobra.Command
}

func (ctx *cmderContext) Args() []string {
	return ctx.args
}

func (ctx *cmderContext) Error(err error) {
	ctx.errors = append(ctx.errors, err)
}

func (ctx *cmderContext) GetErrors() []error {
	return ctx.errors
}

func newCmderContext(ctx context.Context, cmd *cobra.Command, args []string) CmderContext {
	return &cmderContext{
		Context: ctx,
		cmd:     cmd,
		args:    args,
		errors:  make([]error, 0, 1),
	}
}
