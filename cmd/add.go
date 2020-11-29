package cmd

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/dustin/go-humanize"
	"github.com/felipesere/probe/gh"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
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
				getIssue, err := gh.GetIssue(client, t.owner, t.name, t.nr)
				if err != nil {
					return err
				}

				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"idx", "Repository", "Name", "Title", "Status", "Last action", "Last changed", "Link"})
				table.SetAutoWrapText(false)
				table.SetAutoFormatHeaders(true)
				table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
				table.SetAlignment(tablewriter.ALIGN_LEFT)
				table.SetCenterSeparator("")
				table.SetColumnSeparator("")
				table.SetRowSeparator("")
				table.SetHeaderLine(false)
				table.SetBorder(false)
				table.SetTablePadding("\t") // pad with tabs
				table.SetNoWhiteSpace(true)
				table.Append([]string{
					"0",
					t.owner,
					t.name,
					getIssue.Title,
					getIssue.Status,
					getIssue.LastAction,
					humanize.Time(getIssue.LastUpdated),
					targetUrl,
				})
				table.Render()
			} else {
				t, err := extract(prs, targetUrl)
				if err != nil {
					return err
				}

				pr, err := gh.GetPr(client, t.owner, t.name, t.nr)
				if err != nil {
					return err
				}
				spew.Dump(pr)
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
