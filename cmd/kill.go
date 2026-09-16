package cmd

import (
	"devfix/internal/db"
	"devfix/internal/utils"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var forceKill bool

func init() {
	killCmd.Flags().BoolVarP(&forceKill, "force", "f", false, "Force kill")
	rootCmd.AddCommand(killCmd)
}

var killCmd = &cobra.Command{
	Use:   "kill [PID or PORT]",
	Short: "Kill a process by PID or port",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		
		out, _ := exec.Command("netstat", "-ano").Output()
		lines := strings.Split(string(out), "\n")
		pidToKill := target
		
		isPort := false
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 4 && fields[0] == "TCP" && strings.Contains(fields[3], "LISTENING") {
				addr := fields[1]
				lastColon := strings.LastIndex(addr, ":")
				if lastColon != -1 && addr[lastColon+1:] == target {
					pidToKill = fields[len(fields)-1]
					isPort = true
					break
				}
			}
		}

		argsExec := []string{"/PID", pidToKill}
		if forceKill {
			argsExec = append(argsExec, "/F")
		}

		err := exec.Command("taskkill", argsExec...).Run()
		if err != nil {
			utils.PrintError("Failed to kill process %s. Likely an Access Denied error.", pidToKill)
			utils.PrintInfo("Attempting to elevate privileges (UAC prompt may appear)...")
			
			psArgs := fmt.Sprintf("Start-Process taskkill -ArgumentList '/PID %s", pidToKill)
			if forceKill {
				psArgs += " /F"
			}
			psArgs += "' -Verb RunAs"
			
			errElev := exec.Command("powershell", "-NoProfile", "-Command", psArgs).Run()
			if errElev != nil {
				utils.PrintError("Elevated kill failed: %v", errElev)
			} else {
				msg := fmt.Sprintf("Requested elevated kill for process %s", pidToKill)
				if isPort {
					msg += fmt.Sprintf(" (listening on port %s)", target)
				}
				utils.PrintSuccess(msg)
				db.LogAction("KILL", msg, 0)
			}
		} else {
			msg := fmt.Sprintf("Killed process %s", pidToKill)
			if isPort {
				msg += fmt.Sprintf(" (listening on port %s)", target)
			}
			utils.PrintSuccess(msg)
			db.LogAction("KILL", msg, 0)
		}
	},
}
