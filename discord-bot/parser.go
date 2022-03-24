package discordbot

import (
	"errors"
	"strings"
)

func removeCommandName(input string) string {
	if strings.Index(input, " ") > 0 {
		return input[strings.Index(input, " ")+1:]
	} else {
		return input
	}
}

type Args struct {
	Options map[string]string
	Params  []string
}

func newArgs() Args {
	return Args{
		Options: make(map[string]string),
		Params:  make([]string, 0),
	}
}

func parseArgs(input string) (Args, error) {
	//return strings.Split(removeCommandName(input), " ")
	input = removeCommandName(input)
	result := newArgs()
	buffer := ""
	escape := false
	optionMode := false // By default is false, detecting a '-' make it true
	buffering := false
	optionNameBuffer := ""
	for _, v := range input {
		if escape {
			if !buffering {
				return Args{}, errors.New("escape character ('\\') detected outside of string")
			}
			buffer += string(v)
			escape = false
		} else {
			switch v {
			case '-':
				if buffer == "" {
					optionMode = true
				} else {
					buffer += string(v)
				}
			case '"':
				buffering = true
			case ' ':
				if buffering {
					buffer += string(v)
				} else {
					if optionMode {
						if optionNameBuffer == "" {
							optionNameBuffer = buffer
							buffer = ""
						} else {
							result.Options[optionNameBuffer] = buffer
							buffer = ""
							optionNameBuffer = ""
							optionMode = false
						}
					} else {
						result.Params = append(result.Params, buffer)
						buffer = ""
					}
				}
			case '\\':
				escape = true
			default:
				buffer += string(v)
			}
		}
	}
	if buffer != "" {
		if buffering {
			return Args{}, errors.New("error during arguments parsing: Unfinished string")
		} else {
			if optionMode {
				result.Options[optionNameBuffer] = buffer
			} else {
				result.Params = append(result.Params, buffer)
			}
		}
	}
	return result, nil
}

/*
for k, v := range args.Params {
	utils.LogDebug(fmt.Sprint(k) + " -> " + v)
}
for k, v := range args.Options {
	utils.LogDebug(k + " -> " + v)
}
*/
