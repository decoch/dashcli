package app

import (
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/decoch/dashcli/internal/config"
	"github.com/decoch/dashcli/internal/exitcode"
	"github.com/decoch/dashcli/internal/output"
)

var (
	loadConfig = config.LoadDefault
	lookupEnv  = os.LookupEnv
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
	flags    *rootFlags
	resolved config.Resolved
	stdout   io.Writer
	stderr   io.Writer
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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "version" {
				return nil
			}
			cfg, err := loadConfig()
			if err != nil {
				return exitcode.WrapRuntime(err)
			}
			resolved, err := config.Resolve(config.ResolveInput{
				Flags: config.Flags{
					BaseURL: state.flags.BaseURL,
					APIKey:  state.flags.APIKey,
					Profile: state.flags.Profile,
					Timeout: state.flags.Timeout,
					Debug:   state.flags.Debug,
				},
				Config:    cfg,
				LookupEnv: lookupEnv,
			})
			if err != nil {
				return exitcode.WrapUsage(err)
			}
			state.resolved = resolved
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&flags.BaseURL, "base-url", "", "Redash base URL")
	rootCmd.PersistentFlags().StringVar(&flags.APIKey, "api-key", "", "Redash API key")
	rootCmd.PersistentFlags().BoolVar(&flags.JSON, "json", false, "Print JSON output")
	rootCmd.PersistentFlags().DurationVar(&flags.Timeout, "timeout", 10*time.Second, "HTTP timeout")
	rootCmd.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringVar(&flags.Profile, "profile", "", "Profile name")

	rootCmd.AddCommand(newVersionCmd(state))
	rootCmd.AddCommand(newMeCmd(state))
	rootCmd.AddCommand(newQueryCmd(state))
	rootCmd.AddCommand(newJobCmd(state))
	rootCmd.AddCommand(newDashboardCmd(state))
	rootCmd.AddCommand(newDataSourceCmd(state))

	return rootCmd
}

func (state *appState) output() *output.Output {
	return output.New(output.Options{
		JSON:   state.flags.JSON,
		Stdout: state.stdout,
		Stderr: state.stderr,
	})
}
