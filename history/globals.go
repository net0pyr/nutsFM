// Пакет history предоставляет инструменты для работы с историей команд в приложении nutsFM.
package history

import "os"

// HistoryFile представляет путь к файлу истории команд.
var HistoryFile = os.Getenv("HOME") + "/.nutsfm/commands_history"
