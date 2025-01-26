package main

import (
	"fmt"
	"os/exec"
)

func BackupToFile(path string) error {
	fmt.Println("Backing up to", path)
	cmd := exec.Command("tar", "--directory="+pvDirectory, "-czvf", path, ".")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to backup to file: %w", err)
	}
	return nil
}

func BackupToURL(url string, method string) error {
	tarCmd := exec.Command("tar", "--directory="+pvDirectory, "-czvf", "-", ".")
	uploadCmd := exec.Command("curl", "-X", method, "--data-binary", "@-", url)
	tarStdout, err := tarCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get tar stdout pipe: %w", err)
	}
	uploadCmd.Stdin = tarStdout

	defer tarStdout.Close()

	// Start the commands
	if err := tarCmd.Start(); err != nil {
		return fmt.Errorf("failed to start tar command: %w", err)
	}
	if err := uploadCmd.Start(); err != nil {
		return fmt.Errorf("failed to start upload command: %w", err)
	}
	// Wait for the commands to finish
	if err := tarCmd.Wait(); err != nil {
		return fmt.Errorf("tar command failed: %w", err)
	}
	if err := uploadCmd.Wait(); err != nil {
		return fmt.Errorf("upload command failed: %w", err)
	}
	return nil
}
