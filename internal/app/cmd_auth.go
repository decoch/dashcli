package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/99designs/keyring"
	"github.com/spf13/cobra"

	"github.com/decoch/dashcli/internal/exitcode"
	"github.com/decoch/dashcli/internal/redash"
	"github.com/decoch/dashcli/internal/secrets"
)

var (
	authSetAPIKey               = secrets.SetAPIKey
	authGetAPIKey               = secrets.GetAPIKey
	authDeleteAPIKey            = secrets.DeleteAPIKey
	authSetBaseURL              = secrets.SetBaseURL
	authGetBaseURL              = secrets.GetBaseURL
	authDeleteBaseURL           = secrets.DeleteBaseURL
	authInput         io.Reader = os.Stdin
)

func newAuthCmd(state *appState) *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage credentials in keyring",
	}

	authCmd.AddCommand(&cobra.Command{
		Use:   "set",
		Short: "Store base URL and API key",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := bufio.NewReader(authInput)

			if err := state.output().PrintText("Base URL: "); err != nil {
				return exitcode.WrapRuntime(err)
			}
			baseURL, err := reader.ReadString('\n')
			if err != nil {
				return exitcode.WrapRuntime(err)
			}
			baseURL = strings.TrimSpace(baseURL)
			if baseURL == "" {
				return exitcode.Usagef("base URL cannot be empty")
			}
			if _, err := redash.NewClient(baseURL, "dummy", time.Second, false); err != nil {
				return exitcode.WrapUsage(err)
			}

			if err := state.output().PrintText("API key: "); err != nil {
				return exitcode.WrapRuntime(err)
			}
			var apiKey string
			if _, err := fmt.Fscan(reader, &apiKey); err != nil {
				return exitcode.WrapRuntime(err)
			}
			apiKey = strings.TrimSpace(apiKey)
			if apiKey == "" {
				return exitcode.Usagef("API key cannot be empty")
			}

			if err := authSetBaseURL(baseURL); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := authSetAPIKey(apiKey); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := state.output().PrintText("Credentials stored in keyring\n"); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete base URL and API key",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := authDeleteBaseURL(); err != nil && err != keyring.ErrKeyNotFound {
				return exitcode.WrapRuntime(err)
			}
			if err := authDeleteAPIKey(); err != nil && err != keyring.ErrKeyNotFound {
				return exitcode.WrapRuntime(err)
			}
			if err := state.output().PrintText("Credentials removed from keyring\n"); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Show credential status",
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURLSet, err := isBaseURLSet()
			if err != nil {
				return exitcode.WrapRuntime(err)
			}
			apiKeySet, err := isAPIKeySet()
			if err != nil {
				return exitcode.WrapRuntime(err)
			}

			if baseURLSet {
				if printErr := state.output().PrintText("Base URL is set\n"); printErr != nil {
					return exitcode.WrapRuntime(printErr)
				}
			} else {
				if printErr := state.output().PrintText("No base URL stored\n"); printErr != nil {
					return exitcode.WrapRuntime(printErr)
				}
			}

			if apiKeySet {
				if printErr := state.output().PrintText("API key is set\n"); printErr != nil {
					return exitcode.WrapRuntime(printErr)
				}
			} else {
				if printErr := state.output().PrintText("No API key stored\n"); printErr != nil {
					return exitcode.WrapRuntime(printErr)
				}
			}

			return nil
		},
	})

	return authCmd
}

func isBaseURLSet() (bool, error) {
	_, err := authGetBaseURL()
	if err == nil {
		return true, nil
	}
	if err == keyring.ErrKeyNotFound {
		return false, nil
	}
	return false, err
}

func isAPIKeySet() (bool, error) {
	_, err := authGetAPIKey()
	if err == nil {
		return true, nil
	}
	if err == keyring.ErrKeyNotFound {
		return false, nil
	}
	return false, err
}
