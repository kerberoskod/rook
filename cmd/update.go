package cmd

import (
	"fmt"

	"github.com/kerberoskod/rook/scanner"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dependencies to latest versions",
	Long:  `Scan dependencies, check for updates, and rewrite dependency files with the latest versions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		sc := scanner.New()
		deps, err := sc.Scan(path)
		if err != nil {
			return err
		}

		updated, err := sc.CheckUpdates(deps)
		if err != nil {
			return err
		}

		var toUpdate []scanner.Dependency
		for _, d := range updated {
			if d.Outdated {
				toUpdate = append(toUpdate, d)
			}
		}

		if len(toUpdate) == 0 {
			fmt.Println("All dependencies are up-to-date.")
			return nil
		}

		fmt.Printf("Found %d outdated dependencies:\n", len(toUpdate))
		for _, d := range toUpdate {
			fmt.Printf("  %s %s → %s\n", d.Name, d.Version, d.Latest)
		}

		if dryRun {
			fmt.Println("\nDry-run: no files were modified.")
			return nil
		}

		if err := sc.ApplyUpdates(path, toUpdate); err != nil {
			return err
		}

		fmt.Printf("\nUpdated %d dependencies successfully.\n", len(toUpdate))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().Bool("dry-run", false, "Show what would be updated without modifying files")
}
