package cmd

import (
	"devfix/internal/utils"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(portsCmd)
}

var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: "List all listening TCP ports with PIDs and process names",
	Run: func(cmd *cobra.Command, args []string) {
		tasklistOut, _ := exec.Command("tasklist", "/FO", "CSV", "/NH").Output()
		procMap := make(map[string]string)
		for _, line := range strings.Split(string(tasklistOut), "\n") {
			parts := strings.Split(line, `","`)
			if len(parts) >= 2 {
				pName := strings.TrimLeft(parts[0], `"`)
				pID := parts[1]
				procMap[pID] = pName
			}
		}

		out, err := exec.Command("netstat", "-ano").Output()
		if err != nil {
			utils.PrintError("Failed to run netstat: %v", err)
			return
		}

		lines := strings.Split(string(out), "\n")
		table := utils.NewTable([]string{"Proto", "Local Address", "PID", "Process Name"})
		
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 4 && fields[0] == "TCP" && strings.Contains(fields[3], "LISTENING") {
				pid := fields[len(fields)-1]
				addr := fields[1]
				
				procName := "Unknown"
				if name, ok := procMap[pid]; ok {
					procName = name
				}
				table.Append([]string{"TCP", addr, pid, procName})
			}
		}
		table.Render()
	},
}
