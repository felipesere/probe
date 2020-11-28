package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "adds a new MR based on the URL",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running add.")
	},
}
