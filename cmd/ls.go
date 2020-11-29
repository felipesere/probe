package cmd

import (
	"github.com/felipesere/probe/lib"
	"github.com/spf13/cobra"
)

var (
	lsCmd = &cobra.Command{
		Use:   "ls",
		Short: "lists current items",
		RunE: func(cmd *cobra.Command, args []string) error {
			lib.Print(db.Load())
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(lsCmd)
}