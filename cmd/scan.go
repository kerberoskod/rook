package cmd

import (
	"fmt"

	"github.com/kerberoskod/rook/output"
	"github.com/kerberoskod/rook/scanner"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan project dependencies",
	Long:  `Scan the project directory for dependency files and list all dependencies with their current versions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		useJSON, _ := cmd.Flags().GetBool("json")

		sc := scanner.New()
		deps, err := sc.Scan(path)
		if err != nil {
			return err
		}

		if useJSON {
			return output.PrintJSON(deps)
		}

		fmt.Printf("Found %d dependencies:\n\n", len(deps))
		return output.PrintTable(deps)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
