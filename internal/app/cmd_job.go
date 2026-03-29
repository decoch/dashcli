package app

import (
	"context"
	"time"

	"github.com/decoch/dashcli/internal/exitcode"
	"github.com/decoch/dashcli/internal/redash"
	"github.com/spf13/cobra"
)

var jobGetJob = func(ctx context.Context, client *redash.Client, id string) (map[string]any, error) {
	return client.GetJob(ctx, id)
}

func newJobCmd(state *appState) *cobra.Command {
	jobCmd := &cobra.Command{
		Use:   "job",
		Short: "Inspect jobs",
	}

	jobCmd.AddCommand(&cobra.Command{
		Use:   "get <job-id>",
		Short: "Get job status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}
			jobResponse, err := jobGetJob(cmd.Context(), client, args[0])
			if err != nil {
				return exitcode.WrapRuntime(err)
			}

			out := state.output()
			if out.JSONEnabled() {
				return out.Print(jobResponse)
			}

			job := extractJobObject(jobResponse)
			if err := out.PrintText("id: %s\n", asString(job["id"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := out.PrintText("status: %s\n", asString(job["status"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			if err := out.PrintText("error: %s\n", asString(job["error"])); err != nil {
				return exitcode.WrapRuntime(err)
			}
			return nil
		},
	})

	waitCmd := &cobra.Command{
		Use:   "wait <job-id>",
		Short: "Wait until job finishes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := state.apiClient()
			if err != nil {
				return err
			}

			interval, err := cmd.Flags().GetDuration("interval")
			if err != nil {
				return exitcode.WrapUsage(err)
			}
			maxWait, err := cmd.Flags().GetDuration("max-wait")
			if err != nil {
				return exitcode.WrapUsage(err)
			}

			deadline := time.Now().Add(maxWait)
			for {
				jobResponse, err := jobGetJob(cmd.Context(), client, args[0])
				if err != nil {
					return exitcode.WrapRuntime(err)
				}

				job := extractJobObject(jobResponse)
				status, _ := asInt(job["status"])
				if status >= 3 {
					if state.output().JSONEnabled() {
						return state.output().Print(jobResponse)
					}
					if err := state.output().PrintText("id: %s\nstatus: %s\nerror: %s\n", asString(job["id"]), asString(job["status"]), asString(job["error"])); err != nil {
						return exitcode.WrapRuntime(err)
					}
					return nil
				}

				if time.Now().After(deadline) {
					return exitcode.Runtimef("timed out waiting for job %s", args[0])
				}
				select {
				case <-time.After(interval):
				case <-cmd.Context().Done():
					return exitcode.WrapRuntime(cmd.Context().Err())
				}
			}
		},
	}
	waitCmd.Flags().Duration("interval", 2*time.Second, "Polling interval")
	waitCmd.Flags().Duration("max-wait", 60*time.Second, "Maximum wait duration")

	jobCmd.AddCommand(waitCmd)

	return jobCmd
}
