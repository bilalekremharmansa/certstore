package cli

import (
	"github.com/spf13/cobra"
)

var (
	clusterName string
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "manage cluster",
}

func init() {
	rootCmd.AddCommand(clusterCmd)
}
