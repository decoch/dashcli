package app

import (
	"io"
	"time"

	"github.com/spf13/cobra"

	"github.com/decoch/dashcli/internal/output"
)

type rootFlags struct {
	BaseURL string
	APIKey  string
	JSON    bool
	Timeout time.Duration
	Debug   bool
	Profile string
}

type appState struct {
	flags  *rootFlags
	stdout io.Writer
	stderr io.Writer
}

func newRootCmd(stdout, stderr io.Writer) *cobra.Command {
	flags := &rootFlags{}
	state := &appState{
		flags:  flags,
		stdout: stdout,
		stderr: stderr,
	}

	rootCmd := &cobra.Command{
		Use:           "dash",
		Short:         "CLI for Redash API",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.PersistentFlags().StringVar(&flags.BaseURL, "base-url", "", "Redash base URL")
	rootCmd.PersistentFlags().StringVar(&flags.APIKey, "api-key", "", "Redash API key")
	rootCmd.PersistentFlags().BoolVar(&flags.JSON, "json", false, "Print JSON output")
	rootCmd.PersistentFlags().DurationVar(&flags.Timeout, "timeout", 10*time.Second, "HTTP timeout")
	rootCmd.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringVar(&flags.Profile, "profile", "", "Profile name")

	rootCmd.AddCommand(newVersionCmd(state))

	return rootCmd
}

func (state *appState) output() *output.Output {
	return output.New(output.Options{
		JSON:   state.flags.JSON,
		Stdout: state.stdout,
		Stderr: state.stderr,
	})
}

