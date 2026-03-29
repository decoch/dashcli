package app

import (
	"github.com/decoch/dashcli/internal/exitcode"
	"github.com/spf13/cobra"
)

func newSQLCmd(state *appState) *cobra.Command {
	sqlCmd := &cobra.Command{
		Use:   "sql",
		Short: "Run ad-hoc SQL",
	}

	var dataSourceID int
	var queryText string
	var maxAge int

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Execute SQL without polling",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}
			response, err := client.ExecuteSQL(cmd.Context(), dataSourceID, queryText, maxAge)
			if err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := state.output().Print(response); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	}
	runCmd.Flags().IntVar(&dataSourceID, "datasource", 0, "Data source ID")
	runCmd.Flags().StringVar(&queryText, "query", "", "SQL query")
	runCmd.Flags().IntVar(&maxAge, "max-age", 0, "Maximum cached result age in seconds")
	_ = runCmd.MarkFlagRequired("datasource")
	_ = runCmd.MarkFlagRequired("query")

	sqlCmd.AddCommand(runCmd)

	return sqlCmd
}
