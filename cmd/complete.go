package cmd

import (
	"emobodigo/todo-cli/internal"

	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Complete a todo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		internal.Complete(id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
