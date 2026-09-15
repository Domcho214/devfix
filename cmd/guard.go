package cmd

import (
	"devfix/internal/utils"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(guardCmd)
}

var guardCmd = &cobra.Command{
	Use:   "guard",
	Short: "Security check git diff before commit",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "diff", "--cached").Output()
		if err != nil || len(out) == 0 {
			out, err = exec.Command("git", "diff").Output()
			if err != nil || len(out) == 0 {
				utils.PrintInfo("No git changes to check.")
				return
			}
		}

		diff := string(out)
		
		patterns := map[string]*regexp.Regexp{
			"OpenAI API Key": regexp.MustCompile(`sk-[a-zA-Z0-9]{48}`),
			"Stripe Key":     regexp.MustCompile(`[rs]k_(test|live)_[a-zA-Z0-9]{24}`),
			"GitHub Token":   regexp.MustCompile(`(ghp|gho|ghu|ghs|ghr)_[a-zA-Z0-9]{36}`),
			"Generic Secret": regexp.MustCompile(`(?i)(password|secret|api_key|token)\s*[:=]\s*['"]?[a-zA-Z0-9\-_]{8,}['"]?`),
		}

		found := false
		lines := strings.Split(diff, "\n")
		for i, line := range lines {
			if !strings.HasPrefix(line, "+") || strings.HasPrefix(line, "+++") {
				continue
			}
			
			for name, re := range patterns {
				if re.MatchString(line) {
					utils.PrintError("Possible %s found at line %d:\n%s", name, i+1, line)
					found = true
				}
			}
		}

		if found {
			
			utils.PrintError("Secrets detected in git diff! Commit aborted (conceptually).")
			os.Exit(1)
		} else {
			utils.PrintSuccess("Git diff looks clean. ")
		}
	},
}
