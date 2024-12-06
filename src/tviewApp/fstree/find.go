package fstree

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func Find(startDir, query string) ([]string, error) {
	var results []string
	var wg sync.WaitGroup
	fileChan := make(chan string)

	go func() {
		for file := range fileChan {
			results = append(results, file)
		}
	}()

	var search func(string)
	search = func(dir string) {
		defer wg.Done()
		files, err := os.ReadDir(dir)
		if err != nil {
			// log.Printf("Failed to read directory %s: %v\n", dir, err)
			return
		}

		for _, file := range files {
			childPath := filepath.Join(dir, file.Name())
			if strings.Contains(strings.ToLower(file.Name()), strings.ToLower(query)) {
				fileChan <- childPath
			}
			if file.IsDir() {
				wg.Add(1)
				go search(childPath)
			}
		}
	}

	wg.Add(1)
	go search(startDir)

	wg.Wait()
	close(fileChan)

	return results, nil
}
