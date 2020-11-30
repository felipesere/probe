package cmd

import (
	"fmt"
	"github.com/felipesere/probe/lib"
	"github.com/spf13/cobra"
)

var (
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "update all items",
		RunE: func(cmd *cobra.Command, args []string) error {
			for id, gh := range db.LoadData() {

				var newItem lib.GithubData
				var err error
				switch gh.Kind {
				case lib.IssueKind:
					newItem, err = lib.GetIssue(githubClient, gh.Owner, gh.Repository, gh.Number)
				case lib.PullRequestKind:
					newItem, err = lib.GetPr(githubClient, gh.Owner, gh.Repository, gh.Number)
				}
				if err != nil {
					return fmt.Errorf("unable to update entry: %s", err.Error())
				}
				db.Update(id, newItem)
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
}
