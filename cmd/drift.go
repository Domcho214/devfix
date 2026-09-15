package cmd

import (
	"devfix/internal/utils"
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(driftCmd)
}

var driftCmd = &cobra.Command{
	Use:   "drift",
	Short: "Detect tool version mismatch",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("package.json"); err == nil {
			content, _ := os.ReadFile("package.json")
			var pkg map[string]interface{}
			json.Unmarshal(content, &pkg)
			
			if engines, ok := pkg["engines"].(map[string]interface{}); ok {
				if nodeReq, ok := engines["node"].(string); ok {
					out, _ := exec.Command("node", "-v").Output()
					nodeVer := strings.TrimSpace(string(out))
					utils.PrintInfo("Node required: %s, Local: %s", nodeReq, nodeVer)
				}
			}
		}
		
		if _, err := os.Stat(".nvmrc"); err == nil {
			content, _ := os.ReadFile(".nvmrc")
			req := strings.TrimSpace(string(content))
			out, _ := exec.Command("node", "-v").Output()
			nodeVer := strings.TrimSpace(string(out))
			utils.PrintInfo(".nvmrc requires: %s, Local: %s", req, nodeVer)
		}

		if _, err := os.Stat("pyproject.toml"); err == nil {
			content, _ := os.ReadFile("pyproject.toml")
			if strings.Contains(string(content), "python =") {
				out, _ := exec.Command("python", "--version").Output()
				pyVer := strings.TrimSpace(string(out))
				utils.PrintInfo("pyproject.toml has python requirements. Local: %s", pyVer)
			}
		}
		
		utils.PrintSuccess("Drift check completed.")
	},
}
