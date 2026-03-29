package app

import (
	"github.com/decoch/dashcli/internal/exitcode"
	"github.com/spf13/cobra"
)

func newQueryCmd(state *appState) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "Manage queries",
	}

	queryCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List queries",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}
			queries, err := client.ListQueries(cmd.Context())
			if err != nil {
				return exitcode.WrapRuntime(err)
			}

			out := state.output()
			if out.JSONEnabled() {
				return out.Print(queries)
			}
			for _, query := range queries {
				if err := out.PrintText("%s\t%s\n", asString(query["id"]), asString(query["name"])); err != nil {
					return exitcode.WrapRuntime(err)
				}
			}
			return nil
		},
	})

	queryCmd.AddCommand(&cobra.Command{
		Use:   "get <id>",
		Short: "Get query detail",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}
			query, err := client.GetQuery(cmd.Context(), args[0])
			if err != nil {
				return exitcode.WrapRuntime(err)
			}

			out := state.output()
			if out.JSONEnabled() {
				return out.Print(query)
			}
			if err := out.PrintText("id: %s\n", asString(query["id"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := out.PrintText("name: %s\n", asString(query["name"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := out.PrintText("query: %s\n", asString(query["query"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	})

	queryCmd.AddCommand(&cobra.Command{
		Use:   "run <id>",
		Short: "Run query",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}
			result, err := client.RunQuery(cmd.Context(), args[0])
			if err != nil {
				return exitcode.WrapRuntime(err)
			}

			out := state.output()
			if out.JSONEnabled() {
				return out.Print(result)
			}

			job := extractJobObject(result)
			if err := out.PrintText("job_id: %s\n", asString(job["id"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := out.PrintText("status: %s\n", asString(job["status"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	})

	return queryCmd
}

