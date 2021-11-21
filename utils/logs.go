package utils

import (
	"fmt"
	"time"
)

func log(message string, color string) {
	if !TestMode || DebugMode {
		fmt.Println("\033[1;" + color + "m" + time.Now().Format("[01-02-2006 15:04:05 MST] ") + message + "\033[0m")
	}
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
