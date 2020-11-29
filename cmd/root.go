package cmd

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"os"
)

var (
	client = githubv4.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)))

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