package history

import (
	"log"
	"os"
)

func AppendHistory(command string) {
	f, err := os.OpenFile(HistoryFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Невозможно открыть или создать файл с историей команд: %v\n", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(command + "\n"); err != nil {
		log.Printf("Невозможно записать в файл с историей команд: %v\n", err)
	}
}
