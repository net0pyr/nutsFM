package git

import (
	"os/exec"
	"strings"
)

func GetGitInfo(path string) (string, error) {
	cmd := exec.Command("git", "-C", path, "status", "-sb")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return lines[0], nil
	}
	return "No Git Information Available", nil
}
