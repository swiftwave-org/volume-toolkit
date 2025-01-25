package main

import "os/exec"

func BackupToFile(path string) error {
	cmd := exec.Command("tar", "--directory="+pvDirectory, "-czvf", path, ".")
	return cmd.Run()
}
