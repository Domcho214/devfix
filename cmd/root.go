// discord @domco00 for help
package cmd

import (
	"bufio"
	"devfix/internal/db"
	"devfix/internal/utils"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devfix",
	Short: "A comprehensive developer system toolkit ",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		
		return db.InitDB()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		
		db.CloseDB()
	},
	Run: func(cmd *cobra.Command, args []string) {
		
		interactiveMenu()
	},
}

func clearScreen() {
	c := exec.Command("cmd", "/c", "cls")
	c.Stdout = os.Stdout
	c.Run()
}

func interactiveMenu() {
	for {
		clearScreen()
		color.Green("================================================================")
		color.Green("                    devfix - Developer Toolkit                  ")
		color.Green("================================================================")
		fmt.Println("  System Modules:")
		color.HiWhite("  [1] doctor      - Holistic system diagnosis and Health Score")
		color.HiWhite("  [2] ports       - List all listening TCP ports")
		color.HiWhite("  [3] kill        - Kill a process by PID or port")
		color.HiWhite("  [4] sys         - Audit PATH environment variable")
		color.HiWhite("  [5] flush       - Flush local DNS cache and check hosts file")
		fmt.Println("\n  Cleanup & Fixes:")
		color.HiCyan("  [6] clean       - Recursive scan for old unused dev folders")
		color.HiCyan("  [7] nuke        - Quick fix for crashed projects (rm node_modules)")
		color.HiCyan("  [8] distress    - Deep clean - temp files, tool caches")
		fmt.Println("\n  Project & Git:")
		color.HiYellow("  [9] env         - Audit local .env files against .env.example")
		color.HiYellow("  [A] drift       - Detect tool version mismatch")
		color.HiYellow("  [B] repos       - Recursive scan for Git repos")
		color.HiYellow("  [C] guard       - Security check git diff before commit")
		color.HiYellow("  [D] serve       - Fast local HTTP server")
		fmt.Println("\n  Misc:")
		color.HiMagenta("  [H] history     - Export action history")
		color.HiRed("  [0] Exit")
		color.Green("================================================================")
		color.Green("Choose a menu option [1,2...A,B...H,0] and press Enter : ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToUpper(input))

		fmt.Println()
		switch input {
		case "1":
			doctorCmd.Run(doctorCmd, []string{})
		case "2":
			portsCmd.Run(portsCmd, []string{})
		case "3":
			fmt.Print("Enter PID or PORT to kill: ")
			target, _ := reader.ReadString('\n')
			target = strings.TrimSpace(target)
			if target != "" {
				killCmd.Run(killCmd, []string{target})
			}
		case "4":
			sysCmd.Run(sysCmd, []string{})
		case "5":
			flushCmd.Run(flushCmd, []string{})
		case "6":
			cleanCmd.Run(cleanCmd, []string{})
		case "7":
			nukeCmd.Run(nukeCmd, []string{})
		case "8":
			distressCmd.Run(distressCmd, []string{})
		case "9":
			envCmd.Run(envCmd, []string{})
		case "A":
			driftCmd.Run(driftCmd, []string{})
		case "B":
			reposCmd.Run(reposCmd, []string{})
		case "C":
			guardCmd.Run(guardCmd, []string{})
		case "D":
			serveCmd.Run(serveCmd, []string{})
		case "H":
			historyCmd.Run(historyCmd, []string{})
		case "0":
			return
		default:
			utils.PrintError("Invalid option!")
		}

		fmt.Println("\nPress Enter to return to menu...")
		reader.ReadString('\n')
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

