package cmd

import (
	"fmt"
	"github.com/bradfitz/slice"
	"github.com/felipesere/probe/lib"
	"github.com/spf13/cobra"
	"sort"
)

var (
	sortingBy string

	lsCmd = &cobra.Command{
		Use:   "ls",
		Short: "list c",
		RunE: func(cmd *cobra.Command, args []string) error {
			items := db.LoadData()
			var sortingFunc func(map[int]lib.GithubData) []int

			switch sortingBy {
			case "changed":
				sortingFunc = byLastUpdate
			case "key":
				sortingFunc = byKey
			case "status":
				sortingFunc = byStatus
			default:
				return fmt.Errorf("unknown sorting function: %s", sortingBy)
			}

			lib.Print(items, sortingFunc)

			return nil
		},
	}
)

func byKey(data map[int]lib.GithubData) []int {
	return sortBy(data, func(left, right lib.GithubData) bool {
		return left.Key < right.Key
	})
}

func byLastUpdate(data map[int]lib.GithubData) []int {
	return sortBy(data, func(left, right lib.GithubData) bool {
		return left.LastUpdated.After(right.LastUpdated)
	})
}

func byStatus(data map[int]lib.GithubData) []int {
	return sortBy(data, func(left, right lib.GithubData) bool {
		return left.Status < right.Status
	})
}

func sortBy(data map[int]lib.GithubData, compare func(left, right lib.GithubData) bool) []int {
	var toBeSorted []lib.GithubData
	for _, item := range data {
		toBeSorted = append(toBeSorted, item)
	}

	sort.Sort(slice.SortInterface(toBeSorted[:], func(i, j int) bool {
		return compare(toBeSorted[i], toBeSorted[j])
	}))

	var keys []int
	for _, x := range toBeSorted {
		keys = append(keys, x.Key)
	}
	return keys
}

func init() {
	lsCmd.Flags().StringVar(&sortingBy, "sort-by", "key", "choose from 'key', 'changed', or 'status' to sort the table")
	rootCmd.AddCommand(lsCmd)
}
