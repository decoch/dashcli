package app

import (
	"github.com/decoch/dashcli/internal/exitcode"
	"github.com/spf13/cobra"
)

func newDataSourceCmd(state *appState) *cobra.Command {
	dataSourceCmd := &cobra.Command{
		Use:   "datasource",
		Short: "Manage data sources",
	}

	dataSourceCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List data sources",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}
			dataSources, err := client.ListDataSources(cmd.Context())
			if err != nil {
				return exitcode.WrapRuntime(err)
			}

			out := state.output()
			if out.JSONEnabled() {
				return out.Print(dataSources)
			}
			for _, dataSource := range dataSources {
				if err := out.PrintText("%s\t%s\t%s\n", asString(dataSource["id"]), asString(dataSource["name"]), asString(dataSource["type"])); err != nil {
					return exitcode.WrapRuntime(err)
				}
			}
			return nil
		},
	})

	return dataSourceCmd
}

