package cmder

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Cmder interface {
	Command() string
	Description() string
	ProvideArgs() cobra.PositionalArgs
	FlagSet(f *pflag.FlagSet)
	Run(ctx CmderContext)
}

type Cmders = []Cmder
