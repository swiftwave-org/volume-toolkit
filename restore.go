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

func RestoreFromURL(url string, method string) error {
	// download the backup file and stream it to tar
	cmdCurl := exec.Command("curl", "-X", method, url)
	cmdCurlStdout, err := cmdCurl.StdoutPipe()
	if err != nil {
		return errors.New("failed to get curl stdout pipe")
	}

	cmdTar := exec.Command("tar", "-xzvf", "-", "--directory", pvDirectory)
	cmdTar.Stdin = cmdCurlStdout

	// start the curl command
	if err := cmdCurl.Start(); err != nil {
		return errors.New("failed to start curl command")
	}

	// start the tar command
	if err := cmdTar.Start(); err != nil {
		return errors.New("failed to start tar command")
	}

	// wait for the curl command to finish
	if err := cmdCurl.Wait(); err != nil {
		return errors.New("failed to download backup file")
	}

	// wait for the tar command to finish
	if err := cmdTar.Wait(); err != nil {
		return errors.New("failed to extract backup file")
	}

	return nil
}
