package app

import (
	"context"
	"io"

	"github.com/decoch/dashcli/internal/exitcode"
)

var Version = "dev"

func Run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	root := newRootCmd(stdout, stderr)
	root.SetArgs(args)
	root.SetContext(ctx)

	if err := root.Execute(); err != nil {
		_, _ = io.WriteString(stderr, err.Error()+"\n")
		return exitcode.Code(err)
	}
	return exitcode.CodeSuccess
}
