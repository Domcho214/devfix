package cmd

import (
	"devfix/internal/db"
	"devfix/internal/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(nukeCmd)
}

var nukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "Quick fix for crashed projects (deletes node_modules, .venv, locks in current dir)",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()
		targets := []string{
			"node_modules",
			".venv",
			"package-lock.json",
			"pnpm-lock.yaml",
			"yarn.lock",
			"poetry.lock",
		}

		var totalFreed int64 = 0
		var removed []string

		for _, t := range targets {
			path := filepath.Join(dir, t)
			info, err := os.Stat(path)
			if err == nil {
				size := getDirSize(path, info)
				if err := os.RemoveAll(path); err == nil {
					totalFreed += size
					removed = append(removed, t)
					utils.PrintSuccess("Removed %s", t)
				} else {
					utils.PrintError("Failed to remove %s: %v", t, err)
				}
			}
		}

		if len(removed) > 0 {
			mb := float64(totalFreed) / 1024 / 1024
			msg := fmt.Sprintf("Nuked %d targets in %s. Freed %.2f MB", len(removed), dir, mb)
			utils.PrintInfo(msg)
			db.LogAction("NUKE", msg, totalFreed)
		} else {
			utils.PrintInfo("Nothing to nuke in this directory.")
		}
	},
}

func getDirSize(path string, info os.FileInfo) int64 {
	if !info.IsDir() {
		return info.Size()
	}
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}
