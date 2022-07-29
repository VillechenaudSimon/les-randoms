package riotinterface

import (
	"fmt"
	"les-randoms/utils"
	"strings"
)

// This file contains only functions that does NOT use the riot api

func ParseGameModeFromQueueId(id int) string {
	switch id {
	case 400:
		return "Normal Game"
	case 420:
		return "Ranked Solo/Duo"
	case 450:
		return "ARAM"
	default:
		return "Unknown Game Mode (queueId : " + fmt.Sprint(id) + ")"
	}
}

func ParseRiotError(err string) RiotApiError {
	switch err[0:3] {
	case "403":
		return RiotApiErrorForbidden
	case "429":
		return RiotApiErrorTooManyRequests
	default:
		utils.LogError("Unknown RiotAPIError : " + err)
		return RiotApiErrorUnknown
	}
}

func ParseTierRank(tier string, rank string) string {
	switch tier {
	case "MASTER":
		return "Master"
	case "GRANDMASTER":
		return "GrandMaster"
	case "CHALLENGER":
		return "Challenger"
	default:
		return strings.Title(strings.ToLower(tier)) + " " + rank
	}
}
