package cmd

import (
	"fmt"

	"github.com/kerberoskod/rook/output"
	"github.com/kerberoskod/rook/scanner"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for outdated dependencies",
	Long:  `Scan dependencies and check for newer versions available online.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		useJSON, _ := cmd.Flags().GetBool("json")
		strict, _ := cmd.Flags().GetBool("strict")

		sc := scanner.New()
		deps, err := sc.Scan(path)
		if err != nil {
			return err
		}

		outdated, err := sc.CheckUpdates(deps)
		if err != nil {
			return err
		}

		if useJSON {
			return output.PrintJSON(outdated)
		}

		current := 0
		for _, d := range outdated {
			if !d.Outdated {
				current++
			}
		}
		fmt.Printf("Checked %d dependencies: %d up-to-date, %d outdated\n\n",
			len(outdated), current, len(outdated)-current)

		var outdatedOnly []scanner.Dependency
		for _, d := range outdated {
			if d.Outdated {
				outdatedOnly = append(outdatedOnly, d)
			}
		}
		if len(outdatedOnly) > 0 {
			output.PrintTable(outdatedOnly)
		}

		if strict && len(outdatedOnly) > 0 {
			fmt.Println("\nStrict mode: found outdated dependencies")
			return fmt.Errorf("rook check failed: %d outdated", len(outdatedOnly))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().Bool("strict", false, "Exit with error if any dependency is outdated")
}
