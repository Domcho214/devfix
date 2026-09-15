package cmd

import (
	"devfix/internal/db"
	"devfix/internal/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(flushCmd)
}

var flushCmd = &cobra.Command{
	Use:   "flush",
	Short: "Flush local DNS cache and check hosts file",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintInfo("Flushing DNS cache...")
		out, err := exec.Command("ipconfig", "/flushdns").CombinedOutput()
		if err != nil {
			utils.PrintError("Failed to flush DNS: %v", err)
			return
		}
		utils.PrintSuccess(strings.TrimSpace(string(out)))
		db.LogAction("FLUSH", "Flushed DNS cache", 0)

		sysRoot := os.Getenv("SystemRoot")
		if sysRoot == "" {
			sysRoot = `C:\Windows`
		}
		hostsPath := filepath.Join(sysRoot, `System32\drivers\etc\hosts`)
		content, err := os.ReadFile(hostsPath)
		if err != nil {
			utils.PrintError("Failed to read hosts file: %v", err)
			return
		}

		utils.PrintInfo("Checking hosts file for localhost blocking...")
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "#") && strings.Contains(line, "localhost") {
				if !strings.HasPrefix(line, "127.0.0.1") && !strings.HasPrefix(line, "::1") {
					utils.PrintWarning("Potential blocked localhost routing: %s", line)
				}
			}
		}
	},
}
