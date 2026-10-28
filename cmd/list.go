package cmd

import (
	"emobodigo/todo-cli/internal"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Menampilkan list todo",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}
		internal.List(all)
		return nil
	},
}

func init() {
	listCmd.Flags().BoolP("all", "a", false, "Tampilkan semua todo")
	rootCmd.AddCommand(listCmd)
}
