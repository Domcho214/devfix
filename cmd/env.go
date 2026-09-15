package cmd

import (
	"devfix/internal/utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	envCmd.AddCommand(envFixCmd)
	rootCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Audit local .env files against .env.example",
	Run: func(cmd *cobra.Command, args []string) {
		auditEnv(false)
	},
}

var envFixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Automatically add missing keys to .env",
	Run: func(cmd *cobra.Command, args []string) {
		auditEnv(true)
	},
}

func parseEnv(content string) map[string]string {
	res := make(map[string]string)
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			res[parts[0]] = parts[1]
		}
	}
	return res
}

func auditEnv(fix bool) {
	exampleContent, err := os.ReadFile(".env.example")
	if err != nil {
		utils.PrintError("Could not read .env.example: %v", err)
		return
	}

	envContent, err := os.ReadFile(".env")
	if err != nil {
		if fix {
			os.WriteFile(".env", exampleContent, 0644)
			utils.PrintSuccess("Created .env from .env.example")
			return
		}
		utils.PrintError("Could not read .env: %v", err)
		return
	}

	exampleMap := parseEnv(string(exampleContent))
	envMap := parseEnv(string(envContent))

	missing := []string{}
	defaultVals := []string{}

	for k, v := range exampleMap {
		envVal, exists := envMap[k]
		if !exists {
			missing = append(missing, k)
		} else if envVal == "" || envVal == "YOUR_KEY_HERE" || envVal == v {
			defaultVals = append(defaultVals, k)
		}
	}

	if len(missing) > 0 {
		utils.PrintWarning("Missing keys in .env: %s", strings.Join(missing, ", "))
		if fix {
			f, _ := os.OpenFile(".env", os.O_APPEND|os.O_WRONLY, 0644)
			defer f.Close()
			for _, m := range missing {
				f.WriteString("\n" + m + "=" + exampleMap[m])
			}
			utils.PrintSuccess("Appended missing keys to .env")
		}
	} else {
		utils.PrintSuccess("No missing keys in .env")
	}

	if len(defaultVals) > 0 {
		utils.PrintWarning("Keys with default/empty values in .env: %s", strings.Join(defaultVals, ", "))
	}
}
