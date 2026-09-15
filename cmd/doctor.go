package cmd

import (
	"devfix/internal/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Holistic system diagnosis and Health Score",
	Run: func(cmd *cobra.Command, args []string) {
		score := 100
		utils.PrintInfo("Running system diagnosis...")

		pathVar := os.Getenv("PATH")
		paths := strings.Split(pathVar, string(os.PathListSeparator))
		invalidPaths := 0
		for _, p := range paths {
			p = strings.TrimSpace(p)
			if p != "" {
				if _, err := os.Stat(p); os.IsNotExist(err) {
					invalidPaths++
				}
			}
		}
		if invalidPaths > 0 {
			utils.PrintError("Found %d invalid paths in PATH", invalidPaths)
			score -= invalidPaths * 2
		} else {
			utils.PrintSuccess("PATH variable is clean")
		}

		tempDir := os.TempDir()
		var tempSize int64 = 0
		filepath.Walk(tempDir, func(_ string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				tempSize += info.Size()
			}
			return nil
		})
		tempMB := float64(tempSize) / 1024 / 1024
		if tempMB > 1000 {
			utils.PrintWarning("Temp folder is quite large: %.2f MB", tempMB)
			score -= int(tempMB / 500)
		} else {
			utils.PrintSuccess("Temp folder size is OK (%.2f MB)", tempMB)
		}

		out, err := exec.Command("git", "status", "--porcelain").Output()
		if err == nil {
			if len(strings.TrimSpace(string(out))) > 0 {
				utils.PrintWarning("Current repo has uncommitted changes")
				score -= 5
			} else {
				utils.PrintSuccess("Current git repo is clean")
			}
		}

		if score < 0 {
			score = 0
		}
		fmt.Printf("\n--- Health Score: %d/100 ---\n", score)
		if score > 90 {
			utils.PrintSuccess("System is in great shape!")
		} else if score > 70 {
			utils.PrintWarning("System is okay, but could use some maintenance.")
		} else {
			utils.PrintError("System needs attention. Run devfix clean and distress.")
		}
	},
}
