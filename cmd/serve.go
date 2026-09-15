package cmd

import (
	"devfix/internal/utils"
	"net"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Fast local HTTP server in the current directory on the first free port",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()
		
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			utils.PrintError("Failed to find a free port: %v", err)
			return
		}
		port := listener.Addr().(*net.TCPAddr).Port
		
		utils.PrintSuccess("Serving %s on http://127.0.0.1:%d", dir, port)
		
		fs := http.FileServer(http.Dir(dir))
		err = http.Serve(listener, fs)
		if err != nil {
			utils.PrintError("Server error: %v", err)
		}
	},
}
