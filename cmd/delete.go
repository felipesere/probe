package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

var (
	id int

	deleteCmd = &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del", "rm"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires an ID")
			}
			maybeID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("id needs to be a number: %s", args[0])
			}
			id = maybeID
			return nil
		},
		Short: "deletes an entry",
		RunE: func(cmd *cobra.Command, args []string) error {

			db.Delete(id)

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}
