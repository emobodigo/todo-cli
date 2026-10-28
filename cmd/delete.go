package cmd

import (
	"emobodigo/todo-cli/internal"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a. todo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		internal.Delete(id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
