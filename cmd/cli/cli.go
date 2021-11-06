package cli

import (
	"os"

	"github.com/spf13/cobra"

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
	rootCmd.Execute()
}

// ----

func error(template string, args ...interface{}) {
	logging.GetLogger().Errorf(template, args)
	os.Exit(1)
}
