package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	importCmd.Flags().String("path", "/app/backup.tar.gz", "Path of tar to import")
	importCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(importCmd)
	exportCmd.Flags().String("path", "/app/backup.tar.gz", "Path of tar to export")
	exportCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(sizeCmd)
	rootCmd.AddCommand(destroyCmd)
}

var importCmd = &cobra.Command{
	Use: "import",
	Run: func(cmd *cobra.Command, args []string) {
		err := RestoreFromFile(cmd.Flag("path").Value.String())
		if err != nil {
			PrintError(err.Error())
		}
		PrintData("Data has been imported")
	},
}

var exportCmd = &cobra.Command{
	Use: "export",
	Run: func(cmd *cobra.Command, args []string) {
		err := BackupToFile(cmd.Flag("file").Value.String())
		if err != nil {
			PrintError(err.Error())
		}
		PrintData("Data has been exported")
	},
}

var sizeCmd = &cobra.Command{
	Use: "size",
	Run: func(cmd *cobra.Command, args []string) {
		size := PVSize()
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
		PrintData(map[string]int64{"size": size})
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
