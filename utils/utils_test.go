package utils

import (
	"testing"
)

func TestLogs(t *testing.T) {
	LogClassic("Classic")
	LogDebug("Debug")
	LogError("Error")
	LogInfo("Info")
	LogSuccess("Success")
	LogWarning("Warning")
}
