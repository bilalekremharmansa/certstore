package cli

import (
	"github.com/spf13/cobra"

	"bilalekrem.com/certstore/internal/certstore"
)

var (
	rootCmd = &cobra.Command{
		Use:   "certstore",
		Short: "",
		Long:  ``,
	}

	certStore, _ = certstore.New()
)

func Run() {
	rootCmd.Execute()
}
