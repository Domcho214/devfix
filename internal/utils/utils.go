// discord @domco00 for help
package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

func GetAppDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, ".config")
	}
	appDir := filepath.Join(configDir, "devfix")
	os.MkdirAll(appDir, 0755)
	return appDir
}

func PrintSuccess(msg string, args ...interface{}) {
	color.Green("✓ "+msg, args...)
}

func PrintError(msg string, args ...interface{}) {
	color.Red("✗ "+msg, args...)
}

func PrintWarning(msg string, args ...interface{}) {
	color.Yellow("! "+msg, args...)
}

func PrintInfo(msg string, args ...interface{}) {
	color.Blue("ℹ "+msg, args...)
}

type Table struct {
	writer  *tabwriter.Writer
	headers []string
	lines   [][]string
}

func NewTable(headers []string) *Table {
	return &Table{
		writer:  tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0),
		headers: headers,
	}
}

func (t *Table) Append(row []string) {
	t.lines = append(t.lines, row)
}

func (t *Table) Render() {
	fmt.Fprintln(t.writer, strings.Join(t.headers, "\t"))
	for _, row := range t.lines {
		fmt.Fprintln(t.writer, strings.Join(row, "\t"))
	}
	t.writer.Flush()
}

func (t *Table) NumLines() int {
	return len(t.lines)
}

func GetDirSize(path string, info os.FileInfo) int64 {
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


