package main

import (
	"fmt"
	"os"
	"strconv"

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

var filePath string

func init() {
	rootCmd.Flags().SortFlags = false
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(sizeCmd)
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(fileOpCmd)

	importCmd.Flags().String("path", "/app/dump.tar.gz", "Path of tar to import")
	importCmd.Flags().String("http-method", "GET", "HTTP method to use")

	exportCmd.Flags().String("path", "/app/dump.tar.gz", "Path of tar to export") // Path or URL where the backup will be stored.
	exportCmd.Flags().String("http-method", "PUT", "HTTP method to use")

	fileOpCmd.PersistentFlags().StringVar(&filePath, "path", "", "Path of the file")
	fileOpCmd.MarkPersistentFlagRequired("path")

	fileOpCmd.Flags().SortFlags = false
	fileOpCmd.AddCommand(lsCmd)
	fileOpCmd.AddCommand(catCmd)
	fileOpCmd.AddCommand(cpCmd)
	fileOpCmd.AddCommand(mvCmd)
	fileOpCmd.AddCommand(rmCmd)
	fileOpCmd.AddCommand(mkdirCmd)
	fileOpCmd.AddCommand(chmodCmd)

	chownCmd.Flags().String("uid", "", "User ID")
	chownCmd.Flags().String("gid", "", "Group ID")
	chownCmd.MarkFlagRequired("uid")
	chownCmd.MarkFlagRequired("gid")
	fileOpCmd.AddCommand(chownCmd)

	fileOpCmd.AddCommand(downloadCmd)
}

func main() {
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

/*
file-op command
 - ls <full_path>
 - cat <full_path>
 - cp <full_path> <full_path_dest>
 - mv <full_path> <full_path_dest>
 - rm <full_path> [Do recursive delete anyway]
 - mkdir <full_path> [Will create all the directories in the path]
 - chmod <mode> <full_path>
 - chown <uid> <gid> <full_path>
 - download <url> <full_path>
*/

var fileOpCmd = &cobra.Command{
	Use: "file-op",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var lsCmd = &cobra.Command{
	Use: "ls",
	Run: func(cmd *cobra.Command, args []string) {
		path := filePath
		files, err := ListFiles(path)
		if err != nil {
			PrintError(err.Error())
		}
		PrintData(files)
	},
}

var catCmd = &cobra.Command{
	Use: "cat",
	Run: func(cmd *cobra.Command, args []string) {
		path := filePath
		data, err := os.ReadFile(path)
		if err != nil {
			PrintError(err.Error())
		}
		// Stream the bytes to the stdout.
		os.Stdout.Write(data)
	},
}

var cpCmd = &cobra.Command{
	Use: "cp",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			PrintError("Invalid number of arguments")
		}
		src := filePath
		dest := args[0]
		// copy file without changing the ownership and permissions.
		err := CopyFile(src, dest)
		if err != nil {
			PrintError(err.Error())
		}
		PrintData("file has been copied")
	},
}

var mvCmd = &cobra.Command{
	Use: "mv",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			PrintError("Invalid number of arguments")
		}
		src := filePath
		dest := args[0]
		// move file without changing the ownership and permissions.
		err := MoveFile(src, dest)
		if err != nil {
			PrintError(err.Error())
		}
		PrintData("file has been moved")
	},
}

var rmCmd = &cobra.Command{
	Use: "rm",
	Run: func(cmd *cobra.Command, args []string) {
		path := filePath
		err := os.RemoveAll(path)
		if err != nil {
			PrintError(err.Error())
		}
		PrintData("file has been removed")
	},
}

var mkdirCmd = &cobra.Command{
	Use: "mkdir",
	Run: func(cmd *cobra.Command, args []string) {
		path := filePath
		err := os.MkdirAll(path, 0777)
		if err != nil {
			PrintError(err.Error())
		}
		PrintData("directory has been created")
	},
}

var chmodCmd = &cobra.Command{
	Use: "chmod",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			PrintError("Invalid number of arguments")
		}
		mode := args[0] // mode should be in octal format like 777, 755 etc.
		perm, err := strconv.ParseUint(mode, 8, 32)
		if err != nil {
			PrintError(err.Error())
		}
		path := filePath
		err = os.Chmod(path, os.FileMode(perm))
		if err != nil {
			PrintError(err.Error())
		}
		PrintData("file permissions have been changed")
	},
}

var chownCmd = &cobra.Command{
	Use: "chown",
	Run: func(cmd *cobra.Command, args []string) {
		uid, err := strconv.Atoi(cmd.Flag("uid").Value.String())
		if err != nil {
			PrintError(err.Error())
		}
		gid, err := strconv.Atoi(cmd.Flag("gid").Value.String())
		if err != nil {
			PrintError(err.Error())
		}
		err = os.Chown(filePath, uid, gid)
		if err != nil {
			PrintError(err.Error())
		}
		PrintData("file ownership has been changed")
	},
}

var downloadCmd = &cobra.Command{
	Use: "download",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			PrintError("Invalid number of arguments")
		}
		url := args[0]
		err := DownloadFile(url, filePath)
		if err != nil {
			PrintError(err.Error())
		}
		PrintData("file has been downloaded")
	},
}
