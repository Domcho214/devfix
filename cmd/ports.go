package cmd

import (
	"devfix/internal/utils"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var format string

func init() {
	portsCmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json)")
	rootCmd.AddCommand(portsCmd)
}

type PortInfo struct {
	Proto   string `json:"proto"`
	Address string `json:"address"`
	Port    string `json:"port"`
	PID     string `json:"pid"`
	Process string `json:"process"`
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
		var results []PortInfo
		
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 4 && fields[0] == "TCP" && strings.Contains(fields[3], "LISTENING") {
				pid := fields[len(fields)-1]
				addr := fields[1]
				
				port := ""
				lastColon := strings.LastIndex(addr, ":")
				if lastColon != -1 {
					port = addr[lastColon+1:]
				}

				procName := "Unknown"
				if name, ok := procMap[pid]; ok {
					procName = name
				}
				results = append(results, PortInfo{
					Proto:   "TCP",
					Address: addr,
					Port:    port,
					PID:     pid,
					Process: procName,
				})
			}
		}

		if format == "json" {
			j, _ := json.MarshalIndent(results, "", "  ")
			fmt.Println(string(j))
		} else {
			table := utils.NewTable([]string{"Proto", "Local Address", "Port", "PID", "Process Name"})
			for _, r := range results {
				table.Append([]string{r.Proto, r.Address, r.Port, r.PID, r.Process})
			}
			table.Render()
		}
	},
}
