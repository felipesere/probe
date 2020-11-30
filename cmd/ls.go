package cmd

import (
	"github.com/felipesere/probe/lib"
	"github.com/spf13/cobra"
	"sort"
)

var (
	lsCmd = &cobra.Command{
		Use:   "ls",
		Short: "list c",
		RunE: func(cmd *cobra.Command, args []string) error {
			items := db.LoadData()
			lib.Print(items, byKey)

			return nil
		},
	}
)

func byKey(data map[int]lib.GithubData) []int {
	var keys []int
	for _, item := range data {
		keys = append(keys, item.Key)
	}
	sort.Ints(keys)
	return keys
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
