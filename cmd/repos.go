package cmd

import (
	"devfix/internal/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(reposCmd)
}

var reposCmd = &cobra.Command{
	Use:   "repos",
	Short: "Recursive scan for Git repos",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()
		utils.PrintInfo("Scanning for git repos in %s...", dir)
		
		table := utils.NewTable([]string{"Repo", "Status", "Details"})

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() && info.Name() == ".git" {
				repoDir := filepath.Dir(path)
				
				statusCmd := exec.Command("git", "status", "--porcelain")
				statusCmd.Dir = repoDir
				out, _ := statusCmd.Output()
				
				status := "Clean"
				details := ""
				if len(strings.TrimSpace(string(out))) > 0 {
					status = "Dirty"
					details = "Uncommitted changes"
				} else {
					logCmd := exec.Command("git", "log", "@{u}..")
					logCmd.Dir = repoDir
					out, err := logCmd.Output()
					if err == nil && len(strings.TrimSpace(string(out))) > 0 {
						status = "Unpushed"
						details = "Commits ahead of origin"
					}
				}
				
				if status != "Clean" {
					table.Append([]string{repoDir, status, details})
				}
				
				return filepath.SkipDir
			}
			return nil
		})
		
		if table.NumLines() > 0 {
			table.Render()
		} else {
			utils.PrintSuccess("All repos are clean and pushed!")
		}
	},
}
