package cmd

import (
	"devfix/internal/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sysCmd)
}

var sysCmd = &cobra.Command{
	Use:   "sys",
	Short: "Audit PATH environment variable for blind paths and duplicates",
	Run: func(cmd *cobra.Command, args []string) {
		pathVar := os.Getenv("PATH")
		paths := strings.Split(pathVar, string(os.PathListSeparator))

		table := utils.NewTable([]string{"Status", "Path", "Issue"})
		
		seen := make(map[string]bool)
		
		for _, p := range paths {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			
			p = filepath.Clean(p)
			
			if seen[strings.ToLower(p)] {
				table.Append([]string{"WARN", p, "Duplicate entry"})
				continue
			}
			seen[strings.ToLower(p)] = true
			
			if _, err := os.Stat(p); os.IsNotExist(err) {
				table.Append([]string{"ERROR", p, "Directory does not exist"})
			}
		}
		
		if table.NumLines() > 0 {
			utils.PrintWarning("Found issues in PATH variable:")
			table.Render()
		} else {
			utils.PrintSuccess("PATH variable is clean!")
		}
	},
}
