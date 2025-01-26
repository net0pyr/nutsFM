// Пакет git предоставляет функции для работы с репозиториями Git.
package git

import (
	"os"
	"path/filepath"
)

// IsGitRepo проверяет, является ли указанный путь репозиторием Git.
func IsGitRepo(path string) bool {
	gitPath := filepath.Join(path, ".git")
	info, err := os.Stat(gitPath)
	return err == nil && info.IsDir()
}
