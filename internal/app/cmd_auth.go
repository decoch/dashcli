package app

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/99designs/keyring"
	"github.com/spf13/cobra"

	"github.com/decoch/dashcli/internal/exitcode"
	"github.com/decoch/dashcli/internal/secrets"
)

var (
	authSetSecret              = secrets.Set
	authGetSecret              = secrets.Get
	authDeleteSecret           = secrets.Delete
	authInput        io.Reader = os.Stdin
)

func newAuthCmd(state *appState) *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage API keys in keyring",
	}

	authCmd.AddCommand(&cobra.Command{
		Use:   "set",
		Short: "Store API key for profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			profile := profileNameFromGlobalFlag(state.flags.Profile)
			var apiKey string
			if _, err := fmt.Fscan(authInput, &apiKey); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := authSetSecret(profile, strings.TrimSpace(apiKey)); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := state.output().PrintText("API key stored for profile %s\n", profile); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete API key for profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			profile := profileNameFromGlobalFlag(state.flags.Profile)
			if err := authDeleteSecret(profile); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := state.output().PrintText("API key removed for profile %s\n", profile); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Show API key status for profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			profile := profileNameFromGlobalFlag(state.flags.Profile)
			_, err := authGetSecret(profile)
			if err == nil {
				if printErr := state.output().PrintText("API key is set for profile %s\n", profile); printErr != nil {
					return exitcode.WrapRuntime(printErr)
				}
				return nil
			}
			if err == keyring.ErrKeyNotFound {
				if printErr := state.output().PrintText("No API key stored for profile %s\n", profile); printErr != nil {
					return exitcode.WrapRuntime(printErr)
				}
				return nil
			}
			return exitcode.WrapRuntime(err)
		},
	})

	return authCmd
}

func profileNameFromGlobalFlag(profile string) string {
	trimmed := strings.TrimSpace(profile)
	if trimmed == "" {
		return "default"
	}
	return trimmed
}
