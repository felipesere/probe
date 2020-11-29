package cmd

import (
	"fmt"
	"github.com/felipesere/probe/lib"
	"github.com/spf13/cobra"
	"regexp"
	"strconv"
)

const issuePattern = "https://github.com/([^/]+)/([^/]+)/issues/(.+)"
const prPattern = "https://github.com/([^/]+)/([^/]+)/pull/(.+)"

var (
	isIssue bool

	issues = regexp.MustCompile(issuePattern)
	prs    = regexp.MustCompile(prPattern)

	addCmd = &cobra.Command{
		Use:   "add",
		Args:  cobra.ExactArgs(1),
		Short: "adds a new MR based on the URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			targetUrl := args[0]

			if isIssue {
				t, err := extract(issues, targetUrl)
				if err != nil {
					return err
				}
				issue, err := lib.GetIssue(githubClient, t.owner, t.name, t.nr)
				if err != nil {
					return err
				}
				issue.Link = targetUrl

				db.AddIssue(issue)
			} else {
				t, err := extract(prs, targetUrl)
				if err != nil {
					return err
				}

				pr, err := lib.GetPr(githubClient, t.owner, t.name, t.nr)
				if err != nil {
					return err
				}
				pr.Link = targetUrl
				db.AddPullRequest(pr)
			}

			return nil
		},
	}
)

type Target struct {
	owner string
	name  string
	nr    int32
}

func extract(pattern *regexp.Regexp, url string) (Target, error) {
	if !pattern.MatchString(url) {
		return Target{}, fmt.Errorf("url did not match expected pattern: %s", url)
	}
	subMatch := pattern.FindStringSubmatch(url)
	nr, err := strconv.ParseInt(subMatch[3], 10, 63)
	if err != nil {
		return Target{}, fmt.Errorf("could not extract number component of URL: %s", err.Error())
	}

	return Target{owner: subMatch[1], name: subMatch[2], nr: int32(nr)}, nil
}

func init() {
	addCmd.Flags().BoolVarP(&isIssue, "issue", "", false, "add an issue instead of an PR")
	rootCmd.AddCommand(addCmd)
}
