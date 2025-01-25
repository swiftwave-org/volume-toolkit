package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const pvDirectory = "/data"
const appDirectory = "/app" // `app` means `volume-toolkit` in this context.

var rootCmd = &cobra.Command{
	Use: "volume-toolkit",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	// Ensure we have tar/gzip/curl installed.
	CheckForTool("tar")
	CheckForTool("gzip")
	CheckForTool("curl")

	// Ensure /data directory exists.
	if _, err := os.Stat(pvDirectory); os.IsNotExist(err) {
		fmt.Printf("Directory %s does not exist", pvDirectory)
		os.Exit(1)
	}
	// Create /app directory if it does not exist. Will act as a temporary directory.
	if _, err := os.Stat(appDirectory); os.IsNotExist(err) {
		os.Mkdir(appDirectory, 0777)
	}

	rootCmd.Execute()
}

func CheckForTool(tool string) {
	_, err := exec.LookPath(tool)
	if err != nil {
		panic(fmt.Sprintf("Could not find %s in PATH", tool))
	}
}
