package agent

import (
	cliutils "bilalekrem.com/certstore/cmd/cli/utils"
	wrk "bilalekrem.com/certstore/internal/cluster/agent"
	"github.com/spf13/cobra"
)

func newRunPipelineCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "runPipeline",
		Short: "run a pipeline",
		Run: func(cmd *cobra.Command, args []string) {
			configPath, _ := cmd.Flags().GetString("config")
			pipelineToRun, _ := cmd.Flags().GetString("pipeline")

			// -----

			agent, err := wrk.NewFromFileWithSkipJobInitialization(configPath, true)
			cliutils.ValidateNotError(err)

			// ---

			err = agent.RunPipeline(pipelineToRun)
			cliutils.ValidateNotError(err)
		},
	}

	// ----

	cmd.Flags().String("config", "", "agent config file path")
	cmd.MarkFlagRequired("config")

	cmd.Flags().String("pipeline", "", "pipeline name to run")
	cmd.MarkFlagRequired("pipeline")

	return cmd
}
