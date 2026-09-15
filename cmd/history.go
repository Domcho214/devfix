package cmd

import (
	"devfix/internal/db"
	"devfix/internal/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(historyCmd)
}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Export action history as a PowerShell script",
	Run: func(cmd *cobra.Command, args []string) {
		rows, err := db.GetHistory()
		if err != nil {
			utils.PrintError("Failed to get history: %v", err)
			return
		}
		defer rows.Close()

		f, err := os.Create("history.ps1")
		if err != nil {
			utils.PrintError("Failed to create file: %v", err)
			return
		}
		defer f.Close()

		f.WriteString("# devfix history export\n\n")

		count := 0
		for rows.Next() {
			var action, details string
			var freed int64
			var ts string
			rows.Scan(&action, &details, &freed, &ts)
			f.WriteString(fmt.Sprintf("# %s - %s (Freed: %d bytes)\n", ts, details, freed))
			f.WriteString(fmt.Sprintf("Write-Host 'Executed: %s'\n\n", action))
			count++
		}

		utils.PrintSuccess("Exported %d actions to history.ps1", count)
	},
}
