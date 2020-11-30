package cmd

import (
	"fmt"
	"github.com/felipesere/probe/lib"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"regexp"
	"strconv"
)

type Retrieval struct {
	ExtractionPattern string
	Getter            func(client githubv4.Client, target lib.Target) (lib.GithubData, error)
}

var (
	config = []Retrieval{
		{
			"https://github.com/([^/]+)/([^/]+)/issues/(.+)",
			lib.GetIssue,
		},
		{
			"https://github.com/([^/]+)/([^/]+)/pull/(.+)",
			lib.GetPr,
		},
	}

	addCmd = &cobra.Command{
		Use:   "add",
		Short: "adds a new MR based on the URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			targetUrl := args[0]

			for _, option := range config {
				target, matched, err := extract(option.ExtractionPattern, targetUrl)
				if err != nil {
					return err
				}
				if !matched {
					continue
				}
				result, err := option.Getter(githubClient, target)
				if err != nil {
					return err
				}
				db.StoreData(result)
				return nil
			}

			return fmt.Errorf("did not figure out what to do with %s", targetUrl)
		},
	}
)

func extract(p string, url string) (lib.Target, bool, error) {
	pattern := regexp.MustCompile(p)
	if !pattern.MatchString(url) {
		return lib.Target{}, false, nil
	}
	subMatch := pattern.FindStringSubmatch(url)
	nr, err := strconv.ParseInt(subMatch[3], 10, 63)
	if err != nil {
		return lib.Target{}, false, fmt.Errorf("could not extract number component of URL: %s", err.Error())
	}

	return lib.Target{Owner: subMatch[1], Name: subMatch[2], Nr: int32(nr)}, true, nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}
