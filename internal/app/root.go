package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/99designs/keyring"
	"github.com/spf13/cobra"

	"github.com/decoch/dashcli/internal/exitcode"
	"github.com/decoch/dashcli/internal/output"
	"github.com/decoch/dashcli/internal/secrets"
)

var (
	lookupEnv  = os.LookupEnv
	getAPIKey  = secrets.GetAPIKey
	getBaseURL = secrets.GetBaseURL
)

type rootFlags struct {
	BaseURL string
	APIKey  string
	JSON    bool
	Timeout time.Duration
	Debug   bool
}

type appState struct {
	flags    *rootFlags
	resolved resolvedConfig
	stdout   io.Writer
	stderr   io.Writer
}

type resolvedConfig struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
	Debug   bool
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
			if cmd.Name() == "version" || isAuthCommand(cmd) {
				return nil
			}

			resolvedBaseURL, err := resolveBaseURL(state.flags.BaseURL)
			if err != nil {
				return err
			}
			resolvedAPIKey, err := resolveAPIKey(state.flags.APIKey)
			if err != nil {
				return err
			}

			state.resolved = resolvedConfig{
				BaseURL: resolvedBaseURL,
				APIKey:  resolvedAPIKey,
				Timeout: state.flags.Timeout,
				Debug:   state.flags.Debug,
			}

			flagAPIKey := strings.TrimSpace(state.flags.APIKey)
			if flagAPIKey != "" && resolvedAPIKey == flagAPIKey {
				_, _ = fmt.Fprintln(state.stderr, "Warning: passing API key via --api-key is insecure; prefer keyring or environment variable")
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&flags.BaseURL, "base-url", "", "Redash base URL")
	rootCmd.PersistentFlags().StringVar(&flags.APIKey, "api-key", "", "Redash API key")
	rootCmd.PersistentFlags().BoolVar(&flags.JSON, "json", false, "Print JSON output")
	rootCmd.PersistentFlags().DurationVar(&flags.Timeout, "timeout", 10*time.Second, "HTTP timeout")
	rootCmd.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "Enable debug logging")

	rootCmd.AddCommand(newVersionCmd(state))
	rootCmd.AddCommand(newAuthCmd(state))
	rootCmd.AddCommand(newMeCmd(state))
	rootCmd.AddCommand(newQueryCmd(state))
	rootCmd.AddCommand(newSQLCmd(state))
	rootCmd.AddCommand(newJobCmd(state))
	rootCmd.AddCommand(newDashboardCmd(state))
	rootCmd.AddCommand(newDataSourceCmd(state))

	return rootCmd
}

func resolveBaseURL(flagValue string) (string, error) {
	baseURL := strings.TrimSpace(flagValue)
	if baseURL == "" {
		secretValue, err := getBaseURL()
		if err == nil {
			baseURL = strings.TrimSpace(secretValue)
		} else if !errors.Is(err, keyring.ErrKeyNotFound) {
			return "", exitcode.WrapRuntime(err)
		}
	}
	if baseURL == "" {
		if value, ok := lookupEnv("REDASH_BASE_URL"); ok {
			baseURL = strings.TrimSpace(value)
		}
	}
	return baseURL, nil
}

func resolveAPIKey(flagValue string) (string, error) {
	apiKey := strings.TrimSpace(flagValue)
	if apiKey == "" {
		secretValue, err := getAPIKey()
		if err == nil {
			apiKey = strings.TrimSpace(secretValue)
		} else if !errors.Is(err, keyring.ErrKeyNotFound) {
			return "", exitcode.WrapRuntime(err)
		}
	}
	if apiKey == "" {
		if value, ok := lookupEnv("REDASH_API_KEY"); ok {
			apiKey = strings.TrimSpace(value)
		}
	}
	return apiKey, nil
}

func isAuthCommand(cmd *cobra.Command) bool {
	for current := cmd; current != nil; current = current.Parent() {
		if current.Name() == "auth" {
			return true
		}
	}
	return false
}

func (state *appState) output() *output.Output {
	return output.New(output.Options{
		JSON:   state.flags.JSON,
		Stdout: state.stdout,
		Stderr: state.stderr,
	})
}
