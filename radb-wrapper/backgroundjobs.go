package radbwrapper

import (
	"crypto/rand"
	"fmt"
	"les-randoms/backgroundworker"
	"les-randoms/riotinterface"
	"les-randoms/utils"
	"math/big"
	"time"
)

// Time after which the program consider that a summoner in the db is 'spoiled'
const LadderSummonersUpdateSpacing time.Duration = time.Hour * 12

// Time elapsed between each summoners batch update
const LadderSummonerBatchUpdateSpacing time.Duration = time.Minute * 30

// Count of summoners to update
const LadderSummonerUpdateBatchSize int = 10

func SetupJobs() {
	//AddDBUsersSummonersJob()
	AddLadderSummonersJob()
}

func AddDBUsersSummonersJob() {
	/*
		type Memory struct {
			array []int
		}
		mem := Memory{
			array: make([]int, 0),
		}
	*/
	backgroundworker.AddJob(time.Second*15, make([]int, 0), func(m *interface{}) {
		memory := (*m).([]int)
		if len(memory) > 0 {
			utils.LogInfo("JOB0 - " + fmt.Sprint(memory[0]))
			memory = memory[1:]
		} else {
			for i := 0; i < 10; i++ {
				j, _ := rand.Int(rand.Reader, big.NewInt(100))
				memory = append(memory, int(j.Int64()))
			}
			utils.LogInfo("JOB0 - Numbers regenerated")
		}
		*m = memory
	})
}

func AddLadderSummonersJob() {
	backgroundworker.AddJob(LadderSummonerBatchUpdateSpacing, make([]string, 0), func(m *interface{}) {
		memory := (*m).([]string)
		if len(memory) > 0 {
			updateSummonersCount := 0
			var i int
			for i = 0; i < len(memory); i++ {
				if updateSummonersCount >= LadderSummonerUpdateBatchSize {
					break
				}
				_, updated, err := GetSummonerFromId(memory[i])
				if updated {
					updateSummonersCount++
				}
				utils.LogNotNilError(err)
			}
			if len(memory) < LadderSummonerUpdateBatchSize+1 {
				memory = memory[len(memory)-1:]
			} else {
				memory = memory[LadderSummonerUpdateBatchSize:]
			}
			utils.LogInfo("LadderSummonersJobUpdate - " + fmt.Sprint(i) + " summoners affected")
		} else {
			challengerLeague, err := riotinterface.GetSoloDuoChallengerLeague()
			if err != nil {
				utils.LogError(err.Error())
				return
			}
			for _, entry := range challengerLeague.Entries {
				memory = append(memory, entry.SummonerId)
			}
			//utils.LogNotNilError(err)
			utils.LogInfo("LadderSummonersJobUpdate - List refreshed")
		}
		*m = memory
	})
}
