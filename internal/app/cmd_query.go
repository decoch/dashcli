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

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List queries",
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
			queries, err := client.ListQueries(cmd.Context(), page, pageSize, order, search)
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
	}
	listCmd.Flags().Int("page", 1, "Page number")
	listCmd.Flags().Int("page-size", 20, "Results per page")
	listCmd.Flags().String("order", "-updated_at", "Sort order")
	listCmd.Flags().String("search", "", "Filter by query name")
	queryCmd.AddCommand(listCmd)

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

	var createName string
	var createSQL string
	var createDataSource int
	var createDescription string
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create query",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}
			query, err := client.CreateQuery(cmd.Context(), createName, createSQL, createDataSource, createDescription)
			if err != nil {
				return exitcode.WrapRuntime(err)
			}

			out := state.output()
			if out.JSONEnabled() {
				return out.Print(query)
			}
			if err := out.PrintText("Query created: id=%s name=%s\n", asString(query["id"]), asString(query["name"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	}
	createCmd.Flags().StringVar(&createName, "name", "", "Query name")
	createCmd.Flags().StringVar(&createSQL, "sql", "", "SQL text")
	createCmd.Flags().IntVar(&createDataSource, "datasource", 0, "Data source ID")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Query description")
	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("sql")
	_ = createCmd.MarkFlagRequired("datasource")
	queryCmd.AddCommand(createCmd)

	var updateName string
	var updateSQL string
	var updateDataSource int
	var updateDescription string
	updateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update query",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fields := map[string]any{}
			if cmd.Flags().Changed("name") {
				fields["name"] = updateName
			}
			if cmd.Flags().Changed("sql") {
				fields["query"] = updateSQL
			}
			if cmd.Flags().Changed("datasource") {
				fields["data_source_id"] = updateDataSource
			}
			if cmd.Flags().Changed("description") {
				fields["description"] = updateDescription
			}

			client, err := state.apiClient()
			if err != nil {
				return err
			}
			query, err := client.UpdateQuery(cmd.Context(), args[0], fields)
			if err != nil {
				return exitcode.WrapRuntime(err)
			}

			out := state.output()
			if out.JSONEnabled() {
				return out.Print(query)
			}
			if err := out.PrintText("Query updated: id=%s\n", args[0]); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	}
	updateCmd.Flags().StringVar(&updateName, "name", "", "Query name")
	updateCmd.Flags().StringVar(&updateSQL, "sql", "", "SQL text")
	updateCmd.Flags().IntVar(&updateDataSource, "datasource", 0, "Data source ID")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Query description")
	queryCmd.AddCommand(updateCmd)

	archiveCmd := &cobra.Command{
		Use:   "archive <id>",
		Short: "Archive query",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}
			if err := client.ArchiveQuery(cmd.Context(), args[0]); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := state.output().PrintText("Query archived: id=%s\n", args[0]); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	}
	queryCmd.AddCommand(archiveCmd)

	resultsCmd := &cobra.Command{
		Use:   "results <id>",
		Short: "Get query results",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}
			results, err := client.GetQueryResults(cmd.Context(), args[0])
			if err != nil {
				return exitcode.WrapRuntime(err)
			}

			if err := state.output().Print(results); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	}
	queryCmd.AddCommand(resultsCmd)

	return queryCmd
}
