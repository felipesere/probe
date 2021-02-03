package cmd

import (
	"fmt"
	"github.com/felipesere/probe/lib"
	"github.com/spf13/cobra"
	"sort"
)

func init() {
	lsCmd.Flags().StringVar(&sortingBy, "sort-by", "key", "choose from 'key', 'changed', or 'status' to sort the table")
	rootCmd.AddCommand(lsCmd)
}

var (
	sortingBy string

	lsCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "show update to date status of the tracked items",
		RunE: func(cmd *cobra.Command, args []string) error {
			updates, err := updateAll()

			if err != nil {
				return err
			}

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

			lib.Print(items, updates, sortingFunc)

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

	sort.Slice(toBeSorted, func(i, j int) bool {
		return compare(toBeSorted[i], toBeSorted[j])
	})

	var keys []int
	for _, x := range toBeSorted {
		keys = append(keys, x.Key)
	}
	return keys
}

func updateAll() ([]string, error) {
	currentItems := db.LoadData()

	N := len(currentItems)
	semaphore := make(chan bool, N)

	var updates []string
	for id, gh := range currentItems {
		go func(i int, g lib.GithubData) {
			newItem, updated, err := getNewItem(g)
			if err != nil {
				panic(err)
			}
			if updated {
				updates = append(updates, g.ID)
				db.Update(i, newItem)
			}
			semaphore <- true
		}(id, gh)
	}
	for i := 0; i < N; i++ {
		<-semaphore
	}
	return updates, nil
}

func getNewItem(item lib.GithubData) (lib.GithubData, bool, error) {
	var newItem lib.GithubData
	var err error
	target := lib.Target{Owner: item.Owner, Name: item.Repository, Nr: item.Number}
	switch item.Kind {
	case lib.IssueKind:
		newItem, err = lib.GetIssue(githubClient, target)

	case lib.PullRequestKind:
		newItem, err = lib.GetPr(githubClient, target)
	}

	if err != nil {
		return newItem, false, fmt.Errorf("unable to update entry: %s", err.Error())
	}

	return newItem, item.Changed(newItem), nil
}
