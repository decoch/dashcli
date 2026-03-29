package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/99designs/keyring"
	"github.com/spf13/cobra"
	"golang.org/x/term"

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
	authReadPassword            = term.ReadPassword
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
			if _, err := redash.NewClient(baseURL, "dummy", "", time.Second); err != nil {
				return exitcode.WrapUsage(err)
			}

			if err := state.output().PrintText("API key: "); err != nil {
				return exitcode.WrapRuntime(err)
			}
			var apiKey string
			if authInput == os.Stdin {
				raw, err := authReadPassword(0)
				if err != nil {
					return exitcode.WrapRuntime(err)
				}
				if _, err := fmt.Fprintln(state.stdout); err != nil {
					return exitcode.WrapRuntime(err)
				}
				apiKey = string(raw)
			} else {
				raw, err := reader.ReadString('\n')
				if err != nil {
					return exitcode.WrapRuntime(err)
				}
				apiKey = raw
			}
			apiKey = strings.TrimSpace(apiKey)
			if apiKey == "" {
				return exitcode.Usagef("API key cannot be empty")
			}

			if err := authSetBaseURL(baseURL); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := authSetAPIKey(apiKey); err != nil {
				if rollbackErr := authDeleteBaseURL(); rollbackErr != nil && !errors.Is(rollbackErr, keyring.ErrKeyNotFound) {
					return exitcode.WrapRuntime(rollbackErr)
				}
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
			if err := authDeleteBaseURL(); err != nil && !errors.Is(err, keyring.ErrKeyNotFound) {
				return exitcode.WrapRuntime(err)
			}
			if err := authDeleteAPIKey(); err != nil && !errors.Is(err, keyring.ErrKeyNotFound) {
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
	if errors.Is(err, keyring.ErrKeyNotFound) {
		return false, nil
	}
	return false, err
}

func isAPIKeySet() (bool, error) {
	_, err := authGetAPIKey()
	if err == nil {
		return true, nil
	}
	if errors.Is(err, keyring.ErrKeyNotFound) {
		return false, nil
	}
	return false, err
}
