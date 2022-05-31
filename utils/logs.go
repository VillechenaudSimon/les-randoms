package utils

import (
	"fmt"
	"os"
	"time"
)

type logger struct{}

var Logger logger

func (l *logger) Write(p []byte) (n int, err error) {
	_log(string(p), "0", false)
	return 0, nil
}

func init() {
	os.Mkdir("logs", 0755)
	Logger = logger{}
}

func _log(message string, color string, lineBreak bool) {
	if !TestMode || DebugMode {
		message = "\033[1;" + color + "m" + time.Now().Format("[01-02-2006 15:04:05 MST] ") + message + "\033[0m"
		if lineBreak {
			fmt.Println(message)
			logFile(message + "\n")
		} else {
			fmt.Print(message)
			logFile(message)
		}
	}
}

func log(message string, color string) {
	_log(message, color, true)
}

func logFile(message string) {
	logFilePath := "logs/" + time.Now().Format("02-01-2006") + ".log"
	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()
	file.WriteString(message)
}

func LogClassic(message string) {
	log("[LOG] "+message, "0")
}

func LogDebug(message string) {
	if DebugMode {
		log("[DEBUG] "+message, "36")
	}
}

func LogError(message string) {
	log("[ERROR] "+message, "31")
}

func LogInfo(message string) {
	log("[INFO] "+message, "34")
}

func LogSuccess(message string) {
	log("[SUCCESS] "+message, "32")
}

func LogWarning(message string) {
	log("[WARNING] "+message, "33")
}
