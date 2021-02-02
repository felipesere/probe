package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var (
	defaultDBname = "./probe_db.json"

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialises a new tracking file for PRs and MRs",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := os.Stat(defaultDBname)

			if os.IsNotExist(err) {
				db, err := os.Create(defaultDBname)
				if err != nil {
					return fmt.Errorf("unable to create DB: %s", err.Error())
				}
				defer db.Close()
			}

			return nil
		},
	}
)
