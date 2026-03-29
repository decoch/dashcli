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

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List dashboards",
		RunE: func(cmd *cobra.Command, args []string) error {
			page, err := cmd.Flags().GetInt("page")
			if err != nil {
				return exitcode.WrapUsage(err)
			}
			pageSize, err := cmd.Flags().GetInt("page-size")
			if err != nil {
				return exitcode.WrapUsage(err)
			}
			order, err := cmd.Flags().GetString("order")
			if err != nil {
				return exitcode.WrapUsage(err)
			}
			search, err := cmd.Flags().GetString("search")
			if err != nil {
				return exitcode.WrapUsage(err)
			}

			client, err := state.apiClient()
			if err != nil {
				return err
			}
			dashboards, err := client.ListDashboards(cmd.Context(), page, pageSize, order, search)
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
	}
	listCmd.Flags().Int("page", 1, "Page number")
	listCmd.Flags().Int("page-size", 20, "Results per page")
	listCmd.Flags().String("order", "-updated_at", "Sort order")
	listCmd.Flags().String("search", "", "Filter by dashboard name")
	dashboardCmd.AddCommand(listCmd)

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
