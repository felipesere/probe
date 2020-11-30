package cmd

import (
	"fmt"
	"github.com/bradfitz/slice"
	"github.com/felipesere/probe/lib"
	"github.com/spf13/cobra"
	"sort"
	"time"
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
	var keys []int
	for _, item := range data {
		keys = append(keys, item.Key)
	}
	sort.Ints(keys)
	return keys
}

func byLastUpdate(data map[int]lib.GithubData) []int {
	type byLastChange struct {
		id          int
		lastChanged time.Time
	}

	var toBeSorted []byLastChange
	for _, item := range data {
		toBeSorted = append(toBeSorted, byLastChange{
			id:          item.Key,
			lastChanged: item.LastUpdated,
		})
	}
	sort.Sort(slice.SortInterface(toBeSorted[:], func(i, j int) bool {
		return toBeSorted[i].lastChanged.After(toBeSorted[j].lastChanged)
	}))

	var keys []int
	for _, x := range toBeSorted {
		keys = append(keys, x.id)
	}
	return keys
}

func byStatus(data map[int]lib.GithubData) []int {
	type sortByStatus struct {
		id     int
		status string
	}

	var toBeSorted []sortByStatus
	for _, item := range data {
		toBeSorted = append(toBeSorted, sortByStatus{
			id:     item.Key,
			status: item.Status,
		})
	}
	sort.Sort(slice.SortInterface(toBeSorted[:], func(i, j int) bool {
		return toBeSorted[i].status < toBeSorted[j].status
	}))

	var keys []int
	for _, x := range toBeSorted {
		keys = append(keys, x.id)
	}
	return keys
}

func init() {
	lsCmd.Flags().StringVar(&sortingBy, "sort-by", "key", "choose from 'key', 'changed', or 'status' to sort the table")
	rootCmd.AddCommand(lsCmd)
}
