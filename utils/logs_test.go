package utils

import (
	"testing"
)

func TestLogs(t *testing.T) {
	TestMode = false //Needed because logs are desactivated during test but we want them only here
	LogClassic("Classic")
	LogDebug("Debug")
	LogError("Error")
	LogInfo("Info")
	LogSuccess("Success")
	LogWarning("Warning")
	TestMode = true
}
