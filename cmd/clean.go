// discord @domco00 for help
package cmd

import (
	"devfix/internal/db"
	"devfix/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	cleanDryRun    bool
	cleanOlderThan string
)

func init() {
	cleanCmd.Flags().BoolVarP(&cleanDryRun, "dry-run", "d", true, "Dry run, don't delete")
	cleanCmd.Flags().StringVarP(&cleanOlderThan, "older-than", "o", "30d", "Older than (e.g. 30d, 24h)")
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Recursive disk scan for old unused dev folders",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()
		
		var duration time.Duration
		durStr := cleanOlderThan
		if strings.HasSuffix(durStr, "d") {
			days := 0
			fmt.Sscanf(durStr, "%dd", &days)
			duration = time.Duration(days) * 24 * time.Hour
		} else {
			var err error
			duration, err = time.ParseDuration(durStr)
			if err != nil {
				utils.PrintError("Invalid duration format: %v", err)
				return
			}
		}
		
		cutoff := time.Now().Add(-duration)

		utils.PrintInfo("Scanning %s for dev folders older than %v...", dir, duration)

		targets := map[string]bool{
			"node_modules": true, 
			".venv":        true,
			"target":       true, 
			"vendor":       true,
			".next":        true,
		}

		var totalFreed int64 = 0
		var found []string

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() && targets[info.Name()] {
				if info.ModTime().Before(cutoff) {
					size := getDirSize(path, info) // re-using from nuke.go if in same package cmd
					totalFreed += size
					found = append(found, path)
					
					if !cleanDryRun {
						os.RemoveAll(path)
						utils.PrintSuccess("Deleted %s (%.2f MB)", path, float64(size)/1024/1024)
					} else {
						utils.PrintInfo("[DRY RUN] Would delete %s (%.2f MB)", path, float64(size)/1024/1024)
					}
				}
				return filepath.SkipDir
			}
			return nil
		})

		mb := float64(totalFreed) / 1024 / 1024
		if cleanDryRun {
			utils.PrintInfo("Dry run finished. Found %d folders, total %.2f MB to free. Run with --dry-run=false to execute.", len(found), mb)
		} else {
			msg := fmt.Sprintf("Cleaned %d folders. Freed %.2f MB", len(found), mb)
			utils.PrintSuccess(msg)
			db.LogAction("CLEAN", msg, totalFreed)
		}
	},
}

