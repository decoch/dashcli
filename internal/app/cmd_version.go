package app

import "github.com/spf13/cobra"

type versionOutput struct {
	Version string `json:"version"`
}

func newVersionCmd(state *appState) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print CLI version",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := state.output()
			if out.JSONEnabled() {
				return out.Print(versionOutput{Version: Version})
			}
			return out.Print(Version)
		},
	}
}
