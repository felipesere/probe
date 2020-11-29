package cmd

import (
	"github.com/felipesere/probe/lib"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"os"
)

var (
	githubClient githubv4.Client
	db           *lib.Storage

	rootCmd = &cobra.Command{
		Use:   "init",
		Short: "tracks PRs and issues on GitHub",
		Long: `A little tool to keep track of individual Github issues and PRs
and see when changes occur.`,
		SilenceUsage: true,
	}
)

func Execute(client githubv4.Client, storage *lib.Storage) {
	githubClient = client
	db = storage
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
