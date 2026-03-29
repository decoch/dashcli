package app

import "github.com/spf13/cobra"

func newVersionCmd(state *appState) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print CLI version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return state.output().Print(Version)
		},
	}
}

