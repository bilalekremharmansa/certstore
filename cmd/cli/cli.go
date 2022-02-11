package cli

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"

	"bilalekrem.com/certstore/cmd/cli/cluster"
	"bilalekrem.com/certstore/cmd/cli/server"
	"bilalekrem.com/certstore/cmd/cli/agent"
	"bilalekrem.com/certstore/internal/logging"
)

var (
	rootCmd = &cobra.Command{
		Use:   "certstore",
		Short: "",
		Long:  ``,
	}
)

func Run() {
	addCommands()

	logging.ChangeLogLevel(zapcore.DebugLevel)
	rootCmd.Execute()
}

// ----

func addCommands() {
	rootCmd.AddCommand(cluster.NewCommand())
	rootCmd.AddCommand(agent.NewCommand())
	rootCmd.AddCommand(server.NewCommand())
}
