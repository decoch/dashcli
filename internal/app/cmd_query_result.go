package app

import (
	"github.com/decoch/dashcli/internal/exitcode"
	"github.com/spf13/cobra"
)

func newQueryResultCmd(state *appState) *cobra.Command {
	queryResultCmd := &cobra.Command{
		Use:   "query-result",
		Short: "Manage query results",
	}

	queryResultCmd.AddCommand(&cobra.Command{
		Use:   "get <id>",
		Short: "Get query result by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}
			result, err := client.GetQueryResult(cmd.Context(), args[0])
			if err != nil {
				return exitcode.WrapRuntime(err)
			}

			if err := state.output().Print(result); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	})

	var dataSourceID int
	var queryText string
	var maxAge int
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Execute SQL and create a query result",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}
			result, err := client.CreateQueryResult(cmd.Context(), dataSourceID, queryText, maxAge)
			if err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := state.output().Print(result); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	}
	createCmd.Flags().IntVar(&dataSourceID, "datasource", 0, "Data source ID")
	createCmd.Flags().StringVar(&queryText, "query", "", "SQL query")
	createCmd.Flags().IntVar(&maxAge, "max-age", 0, "Maximum cached result age in seconds")
	_ = createCmd.MarkFlagRequired("datasource")
	_ = createCmd.MarkFlagRequired("query")
	queryResultCmd.AddCommand(createCmd)

	return queryResultCmd
}
