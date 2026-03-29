package app

import (
	"github.com/decoch/dashcli/internal/exitcode"
	"github.com/spf13/cobra"
)

func newDashboardCmd(state *appState) *cobra.Command {
	dashboardCmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Manage dashboards",
	}

	dashboardCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List dashboards",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}
			dashboards, err := client.ListDashboards(cmd.Context())
			if err != nil {
				return exitcode.WrapRuntime(err)
			}

			out := state.output()
			if out.JSONEnabled() {
				return out.Print(dashboards)
			}
			for _, dashboard := range dashboards {
				if err := out.PrintText("%s\t%s\t%s\n", asString(dashboard["id"]), asString(dashboard["slug"]), asString(dashboard["name"])); err != nil {
					return exitcode.WrapRuntime(err)
				}
			}
			return nil
		},
	})

	dashboardCmd.AddCommand(&cobra.Command{
		Use:   "get <slug-or-id>",
		Short: "Get dashboard detail",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}
			dashboard, err := client.GetDashboard(cmd.Context(), args[0])
			if err != nil {
				return exitcode.WrapRuntime(err)
			}

			out := state.output()
			if out.JSONEnabled() {
				return out.Print(dashboard)
			}
			if err := out.PrintText("id: %s\n", asString(dashboard["id"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := out.PrintText("slug: %s\n", asString(dashboard["slug"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := out.PrintText("name: %s\n", asString(dashboard["name"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	})

	return dashboardCmd
}

