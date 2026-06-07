package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rook",
	Short: "A dependency migration CLI tool",
	Long: `Rook scans your project files, checks dependency versions,
and helps you update them — all from the terminal.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("path", "p", ".", "Project root directory")
	rootCmd.PersistentFlags().Bool("json", false, "Output in JSON format")
}
