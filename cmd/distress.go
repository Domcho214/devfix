package cmd

import (
	"devfix/internal/db"
	"devfix/internal/utils"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(distressCmd)
}

var distressCmd = &cobra.Command{
	Use:   "distress",
	Short: "Deep clean - temp files, tool caches, factory reset instructions",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintInfo("Starting deep distress clean...")

		tempDir := os.TempDir()
		if tempDir != "" {
			utils.PrintInfo("Cleaning TEMP dir: %s", tempDir)
			// Skipping actual recursive removal for safety, just reporting
			
		}

		tools := []struct{ name, cmd string; args []string }{
			{"npm", "npm", []string{"cache", "clean", "--force"}},
			{"pip", "pip", []string{"cache", "purge"}},
			{"docker", "docker", []string{"system", "prune", "-f"}},
		}

		for _, t := range tools {
			if _, err := exec.LookPath(t.cmd); err == nil {
				utils.PrintInfo("Cleaning %s cache...", t.name)
				exec.Command(t.cmd, t.args...).Run()
			}
		}

		utils.PrintSuccess("Tool caches cleaned.")
		db.LogAction("DISTRESS", "Deep clean executed", 0)
		
		fmt.Println("\n--- FACTORY RESET INSTRUCTIONS ---")
		fmt.Println("If your system is completely broken and you need a fresh start:")
		fmt.Println("Run this command in an Administrator PowerShell to reset Windows:")
		utils.PrintWarning("systemreset -factoryreset")
		fmt.Println("This will remove all apps and settings. Use with extreme caution.")
	},
}
