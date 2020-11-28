package cmd

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/felipesere/probe/gh"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const issuePattern = "https://github.com/([^/]+)/([^/]+)/issues/(.+)"

var (
	issue bool

	addCmd = &cobra.Command{
		Use:   "add",
		Args: cobra.ExactArgs(1),
		Short: "adds a new MR based on the URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			issues := regexp.MustCompile(issuePattern)

			if issues.MatchString(args[0]) {
				submatch := issues.FindStringSubmatch(args[0])
				fmt.Println(strings.Join(submatch[1:], ", "))


				src := oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
				)
				httpClient := oauth2.NewClient(context.Background(), src)

				client := githubv4.NewClient(httpClient)

				nr, _ := strconv.ParseInt(submatch[3], 10, 63)
				getIssue, err := gh.GetIssue(*client, submatch[1], submatch[2], int32(nr))
				if err != nil {
					return err
				}

				spew.Dump(getIssue)
			}

			return nil
		},
	}

)

func init() {
	addCmd.Flags().BoolVarP(&issue, "issue", "", false, "add an issue instead of an PR")
	rootCmd.AddCommand(addCmd)
}

