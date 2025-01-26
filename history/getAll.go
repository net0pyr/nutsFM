package history

import (
	"bufio"
	"os"
)

// GetAll возвращает все команды из файла истории команд.
func GetAll() []string {
	var commands []string
	f, err := os.Open(HistoryFile)
	if err != nil {
		return commands
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		commands = append(commands, scanner.Text())
	}
	return commands
}
