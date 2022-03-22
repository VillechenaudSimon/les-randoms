package discordbot

import "strings"

func removeCommandName(input string) string {
	if strings.Index(input, " ") > 0 {
		return input[strings.Index(input, " ")+1:]
	} else {
		return input
	}
}

func parseArgs(input string) []string {
	return strings.Split(removeCommandName(input), " ")
}
