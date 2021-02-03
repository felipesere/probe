package lib

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/kyokomi/emoji/v2"
	"github.com/olekukonko/tablewriter"
	"os"
)

func Print(data map[int]GithubData, updates []string, ordering func(map[int]GithubData) []int) {
	order := ordering(data)

	var toPrint [][]string
	for _, k := range order {
		row := data[k]
		var title string
		if contains(row.ID, updates) {
			title = emoji.Sprintf(":bulb: %s", row.Title)
		} else {
			title = row.Title
		}
		toPrint = append(toPrint, []string{
			fmt.Sprintf("%d", k),
			title,
			row.Status,
			row.LastAction,
			humanize.Time(row.LastUpdated),
			row.Link,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"idx", "Title", "Status", "Last action", "Last changed", "Link"})
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
	table.AppendBulk(toPrint)
	table.Render()
}

func contains(val string, collection []string) bool {
	for _, x := range collection {
		if x == val {
			return true
		}
	}

	return false
}
