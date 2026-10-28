package cmd

import (
	"emobodigo/todo-cli/internal"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [task]",
	Short: "Menambahkan Todo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		task := args[0]

		internal.Add(task)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
