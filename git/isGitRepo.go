package git

import (
	"os"
	"path/filepath"
)

func IsGitRepo(path string) bool {
	gitPath := filepath.Join(path, ".git")
	info, err := os.Stat(gitPath)
	return err == nil && info.IsDir()
}
