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

const LadderSummonersUpdateSpacing time.Duration = time.Hour * 48
const LadderSummonerUpdateBatchSize int = 5 //30

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
	backgroundworker.AddJob(time.Minute*10, make([]string, 0), func(m *interface{}) {
		memory := (*m).([]string)
		if len(memory) > 0 {
			updateSummonersCount := 0
			var i int
			for i = 0; i < len(memory); i++ {
				if updateSummonersCount >= LadderSummonerUpdateBatchSize {
					break
				}
				_, updated, err := GetSummonerFromName(memory[i])
				if updated {
					updateSummonersCount++
				}
				utils.LogNotNilError(err)
			}
			memory = memory[LadderSummonerUpdateBatchSize:]
			utils.LogInfo("LadderSummonersJobUpdate - " + fmt.Sprint(i) + " summoners affected")
		} else {
			challengerLeague, err := riotinterface.GetSoloDuoChallengerLeague()
			if err != nil {
				utils.LogError(err.Error())
				return
			}
			//queryBody := "WHERE"
			for _, entry := range challengerLeague.Entries {
				//queryBody += " name=" + utils.Esc(entry.SummonerName) + " OR"
				memory = append(memory, entry.SummonerName)
			}
			//memory, err = database.Summoner_SelectAll(queryBody[:len(queryBody)-3])
			//utils.LogNotNilError(err)
			utils.LogInfo("LadderSummonersJobUpdate - List refreshed")
		}
		*m = memory
	})
}
