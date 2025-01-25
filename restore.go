package main

import (
	"errors"
	"os"
	"os/exec"
)

func RestoreFromFile(backupPath string) error {
	// check if the backup file exists
	_, err := os.Stat(backupPath)
	if err != nil {
		return errors.New("backup file does not exist")
	}
	// delete everything in pvDirectory
	cmd := exec.Command("rm", "-rf", pvDirectory+"/*")
	err = cmd.Run()
	if err != nil {
		return errors.New("failed to delete pvDirectory")
	}
	// restore it to pvDirectory
	cmd = exec.Command("tar", "-xzvf", backupPath, "--directory", pvDirectory)
	return cmd.Run()
}
