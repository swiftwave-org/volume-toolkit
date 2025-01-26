package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

/*
If the path is started with http , assume it is a URL and download/upload the file.
*/
func init() {
	importCmd.Flags().String("path", "/app/backup.tar.gz", "Path of tar to import")
	importCmd.Flags().String("http-method", "GET", "HTTP method to use")
	rootCmd.AddCommand(importCmd)
	exportCmd.Flags().String("path", "/app/backup.tar.gz", "Path of tar to export") // Path or URL where the backup will be stored.
	exportCmd.Flags().String("http-method", "PUT", "HTTP method to use")
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(sizeCmd)
	rootCmd.AddCommand(destroyCmd)
}

var importCmd = &cobra.Command{
	Use: "import",
	Run: func(cmd *cobra.Command, args []string) {
		// If the path is started with http , assume it is a URL and download the file.
		path := cmd.Flag("path").Value.String()
		var err error
		if len(path) > 4 && path[:4] == "http" {
			method := cmd.Flag("http-method").Value.String()
			err = RestoreFromURL(path, method)
		} else {
			err = RestoreFromFile(path)
		}
		if err != nil {
			PrintError(err.Error())
		}
		PrintData("Data has been imported")
	},
}

var exportCmd = &cobra.Command{
	Use: "export",
	Run: func(cmd *cobra.Command, args []string) {
		path := cmd.Flag("path").Value.String()
		var err error
		// If the path is started with http , assume it is a URL and upload the file.
		if len(path) > 4 && path[:4] == "http" {
			method := cmd.Flag("http-method").Value.String()
			err = BackupToURL(path, method)
		} else {
			err = BackupToFile(path)
		}
		if err != nil {
			PrintError(err.Error())
		}
		PrintData("Data has been exported")
	},
}

var sizeCmd = &cobra.Command{
	Use: "size",
	Run: func(cmd *cobra.Command, args []string) {
		size, err := FetchFileSize(pvDirectory)
		if err != nil {
			size = 0
		}
		// Write size to /app/size.txt for backward compatibility.
		file := "/app/size.txt"
		f, err := os.Create(file)
		if err != nil {
			PrintError(err.Error())
		}
		defer f.Close()
		_, err = f.WriteString(fmt.Sprintf("%d", size))
		if err != nil {
			PrintError(err.Error())
		}
		PrintData(size)
	},
}

var destroyCmd = &cobra.Command{
	Use: "destroy",
	Run: func(cmd *cobra.Command, args []string) {
		// Destroy the backup.
		_ = os.RemoveAll(pvDirectory)
		PrintData("All data has been destroyed")
	},
}
