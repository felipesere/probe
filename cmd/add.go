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

			var result lib.GithubData
			if isIssue {
				target, err := extract(issues, targetUrl)
				if err != nil {
					return err
				}
				result, err = lib.GetIssue(githubClient, target)
				if err != nil {
					return err
				}
			} else {
				target, err := extract(prs, targetUrl)
				if err != nil {
					return err
				}

				result, err = lib.GetPr(githubClient, target)
			}

			return db.StoreData(result)
		},
	}
)

func extract(pattern *regexp.Regexp, url string) (lib.Target, error) {
	if !pattern.MatchString(url) {
		return lib.Target{}, fmt.Errorf("url did not match expected pattern: %s", url)
	}
	subMatch := pattern.FindStringSubmatch(url)
	nr, err := strconv.ParseInt(subMatch[3], 10, 63)
	if err != nil {
		return lib.Target{}, fmt.Errorf("could not extract number component of URL: %s", err.Error())
	}

	return lib.Target{owner: subMatch[1], name: subMatch[2], nr: int32(nr)}, nil
}

func init() {
	addCmd.Flags().BoolVarP(&isIssue, "issue", "", false, "add an issue instead of an PR")
	rootCmd.AddCommand(addCmd)
}
