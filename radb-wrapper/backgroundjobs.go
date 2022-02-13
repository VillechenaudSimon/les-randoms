package radbwrapper

import (
	"crypto/rand"
	"fmt"
	"les-randoms/backgroundworker"
	"les-randoms/utils"
	"math/big"
	"time"
)

func SetupJobs() {
	AddDBUsersSummonersJob()
	//AddLadderSummonersJob()
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

/*
func AddLadderSummonersJob() {
	backgroundworker.AddJob(time.Minute, func(memory *interface{}) {
		s, _ := database.Summoner_SelectFirst("")
		utils.LogInfo("JOB1: " + s.Name)
	})
}
*/
