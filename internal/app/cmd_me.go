package app

import (
	"github.com/decoch/dashcli/internal/exitcode"
	"github.com/spf13/cobra"
)

func newMeCmd(state *appState) *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "Show current Redash user",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}

			user, err := client.Me(cmd.Context())
			if err != nil {
				return exitcode.WrapRuntime(err)
			}

			out := state.output()
			if out.JSONEnabled() {
				return out.Print(user)
			}

			if err := out.PrintText("id: %s\n", asString(user["id"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := out.PrintText("name: %s\n", asString(user["name"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := out.PrintText("email: %s\n", asString(user["email"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := out.PrintText("is_admin: %s\n", asBoolString(user["is_admin"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	}
}

