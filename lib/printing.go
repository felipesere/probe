package lib

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func Print(data [][]string)  {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"idx", "Owner", "Repository", "Title", "Status", "Last action", "Last changed", "Link"})
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
	table.AppendBulk(data)
	table.Render()
}

