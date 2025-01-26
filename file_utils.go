package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"
)

type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Mode    string    `json:"mode"`
	ModTime time.Time `json:"mod_time"`
	UID     uint      `json:"uid"`
	GID     uint      `json:"gid"`
	IsDir   bool      `json:"is_dir"`
}

/*
* When some function work for both file and directory, it should be named as Path
* Just refer that with `File` in the function name
* We shouldn't use `Dir` in the function name, unless it's specifically for directory
 */

func ExistsPath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !os.IsNotExist(err)
}

func RemovePath(path string) {
	// Silently try to remove directory and all contents
	_ = os.RemoveAll(path)
}

func CreateDirectoryWithOptions(path string, removeIfExists bool, perm os.FileMode) error {
	// Check if directory exists and remove if requested
	if removeIfExists {
		if ExistsPath(path) {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		}
	}

	return os.MkdirAll(path, perm)
}

func WriteFile(path string, data []byte, perm os.FileMode, uid int, gid int) error {
	// Open file for writing
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write data to file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	// Change the ownership of the file
	if err := os.Chown(path, uid, gid); err != nil {
		return err
	}

	return nil
}

func CopyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy contents of source file to destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Fix permissions of destination file
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Fix ownership of destination file
	srcSysInfo := srcInfo.Sys().(*syscall.Stat_t)
	if err := os.Chown(dst, int(srcSysInfo.Uid), int(srcSysInfo.Gid)); err != nil {
		return err
	}

	return nil
}

func MoveFile(src, dst string) error {
	// Fix ownership of destination file
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	srcSysInfo := srcInfo.Sys().(*syscall.Stat_t)

	// Move file
	if err := os.Rename(src, dst); err != nil {
		return err
	}

	// Fix ownership of destination file
	if err := os.Chown(dst, int(srcSysInfo.Uid), int(srcSysInfo.Gid)); err != nil {
		return err
	}
	return nil
}

func ModifyFile(path string, data []byte) error {
	// Get the current file info
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Open file for writing
	file, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, info.Mode())
	if err != nil {
		return err
	}
	defer file.Close()

	// Write data to file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	// Change the ownership of the file to match the original
	if err := os.Chown(path, int(info.Sys().(*syscall.Stat_t).Uid), int(info.Sys().(*syscall.Stat_t).Gid)); err != nil {
		return err
	}

	return nil
}

func FetchFileSize(path string) (int64, error) {
	// Create a channel to receive file sizes
	sizeChan := make(chan int64)
	errChan := make(chan error)
	done := make(chan bool)

	var totalSize int64

	go func() {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				sizeChan <- info.Size()
			}
			return nil
		})

		if err != nil {
			errChan <- err
		}
		done <- true
	}()

	// Collect results
	for {
		select {
		case size := <-sizeChan:
			atomic.AddInt64(&totalSize, size)
		case err := <-errChan:
			return 0, err
		case <-done:
			return totalSize, nil
		}
	}
}

func ListFiles(fullPath string) ([]FileInfo, error) {
	var fileInfos []FileInfo
	files, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return nil, err
		}
		sysInfo := info.Sys().(*syscall.Stat_t)
		fileInfos = append(fileInfos, FileInfo{
			Name:    info.Name(),
			Size:    info.Size(),
			Mode:    fmt.Sprintf("%o", info.Mode().Perm()),
			ModTime: info.ModTime(),
			UID:     uint(sysInfo.Uid),
			GID:     uint(sysInfo.Gid),
			IsDir:   info.IsDir(),
		})
	}

	return fileInfos, nil
}

// DownloadFile downloads a file from the given URL and saves it to the specified path
func DownloadFile(url string, path string) error {
	// Create the file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the data to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
