package main

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

func PVSize() int64 {
	var totalSize int64
	var wg sync.WaitGroup
	var mu sync.Mutex
	fileSizes := make(chan int64)

	// Start a fixed number of goroutines to process file sizes
	numWorkers := runtime.NumCPU()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for size := range fileSizes {
				mu.Lock()
				totalSize += size
				mu.Unlock()
			}
		}()
	}

	// Walk the directory and send file sizes to the channel
	go func() {
		defer close(fileSizes)
		filepath.Walk(pvDirectory, func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fileSizes <- info.Size()
			}
			return nil
		})
	}()

	// Wait for all workers to finish
	wg.Wait()

	return totalSize
}
