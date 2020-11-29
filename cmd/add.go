package cmd

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/felipesere/probe/gh"
	"github.com/spf13/cobra"
	"regexp"
	"strconv"
)

const issuePattern = "https://github.com/([^/]+)/([^/]+)/issues/(.+)"
const prPattern = "https://github.com/([^/]+)/([^/]+)/pull/(.+)"

var (
	issue bool

	addCmd = &cobra.Command{
		Use:   "add",
		Args: cobra.ExactArgs(1),
		Short: "adds a new MR based on the URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			targetUrl := args[0]

			if issue {
				issues := regexp.MustCompile(issuePattern)

				if !issues.MatchString(targetUrl) {
					return fmt.Errorf("url did not match expected pattern: %s", targetUrl)
				}
				subMatch := issues.FindStringSubmatch(targetUrl)
				nr, _ := strconv.ParseInt(subMatch[3], 10, 63)
				getIssue, err := gh.GetIssue(*client, subMatch[1], subMatch[2], int32(nr))
				if err != nil {
					return err
				}
				spew.Dump(getIssue)
			} else {
				prs := regexp.MustCompile(prPattern)

				if !prs.MatchString(targetUrl) {
					return fmt.Errorf("url did not match expected pattern: %s", targetUrl)
				}
				subMatch := prs.FindStringSubmatch(targetUrl)
				nr, _ := strconv.ParseInt(subMatch[3], 10, 63)
				pr, err := gh.GetPr(*client, subMatch[1], subMatch[2], int32(nr))
				if err != nil {
					return err
				}
				spew.Dump(pr)
			}

			return nil
		},
	}

)

func init() {
	addCmd.Flags().BoolVarP(&issue, "issue", "", false, "add an issue instead of an PR")
	rootCmd.AddCommand(addCmd)
}

