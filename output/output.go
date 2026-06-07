package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kerberoskod/rook/scanner"
)

func PrintJSON(deps []scanner.Dependency) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(deps)
}

func PrintTable(deps []scanner.Dependency) error {
	if len(deps) == 0 {
		fmt.Println("No dependencies found.")
		return nil
	}

	nameW := nameWidth(deps)

	header := fmt.Sprintf("  %-*s  %-16s  %-16s  %s", nameW, "Name", "Installed", "Latest", "Manager")
	sep := fmt.Sprintf("  %s", strings.Repeat("-", len(header)))

	fmt.Println(sep)
	fmt.Println(header)
	fmt.Println(sep)

	for _, d := range deps {
		latest := d.Latest
		if latest == "" {
			if d.Outdated {
				latest = "?"
			} else {
				latest = "✓"
			}
		}

		name := d.Name
		if len(name) > nameW {
			name = name[:nameW-3] + "..."
		}

		installed := d.Version
		if installed == "*" {
			installed = "any"
		}

		fmt.Printf("  %-*s  %-16s  %-16s  %s\n",
			nameW, name,
			truncate(installed, 16),
			truncate(latest, 16),
			truncate(d.Manager, 10))
	}
	fmt.Println(sep)

	return nil
}

func nameWidth(deps []scanner.Dependency) int {
	max := 40
	for _, d := range deps {
		if len(d.Name) > max {
			max = len(d.Name)
		}
	}
	if max < 20 {
		return 20
	}
	if max > 60 {
		return 60
	}
	return max
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
