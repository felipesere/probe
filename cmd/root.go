package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "init",
		Short: "tracks PRs and issues on GitHub",
		Long: `A little tool to keep track of individual Github issues and PRs
and see when changes occur.`,
		SilenceUsage: true,
	}
)

func Execute()  {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}


func init() {
	cobra.OnInitialize()
}