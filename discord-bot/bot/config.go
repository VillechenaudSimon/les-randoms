package bot

import "les-randoms/utils"

var Prefix string

func init() {
	if utils.DebugMode {
		Prefix = "k!"
	} else {
		Prefix = "!"
	}
}
