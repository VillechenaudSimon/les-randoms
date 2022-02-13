package backgroundworker

import (
	"les-randoms/utils"
	"time"
)

// tickerUpdateSpacing > tickerTickSpacing Or there will be problems

const tickerTickSpacing time.Duration = 2 * time.Second

//const tickerTickSpacing time.Duration = time.Minute
const tickerUpdateSpacing time.Duration = 10 * tickerTickSpacing

var currentTickerTickSpacing = tickerTickSpacing * 9

var JobAdder chan Job = make(chan Job)
var lastUpdateTime time.Time
var jobs []*Job

type Job struct {
	Id           int
	Spacing      time.Duration
	LastDoneTime time.Time
	Do           func(*interface{})
	Memory       interface{}
}

func Start() {
	jobs = make([]*Job, 0)

	lastUpdateTime = time.Now().UTC()
	jobQueue := make(chan *Job)

	go startTicker(jobQueue)

	for {
		select {
		case job := <-JobAdder:
			jobs = append(jobs, &job)
		case job := <-jobQueue:
			//utils.LogDebug("DOING JOB ID : " + fmt.Sprint(job.Id))
			job.Do(&job.Memory)
			job.LastDoneTime = time.Now().UTC()
		}
	}
}

func startTicker(c chan *Job) {
	for {
		time.Sleep(currentTickerTickSpacing)
		currentTickerTickSpacing = tickerTickSpacing / 3
		//utils.LogDebug("Ticking.. " + lastUpdateTime.Format(utils.DateTimeFormat) + " - " + fmt.Sprint(len(jobs)) + " jobs")
		if time.Now().UTC().Sub(lastUpdateTime) > tickerUpdateSpacing {
			currentTickerTickSpacing = 9 * tickerTickSpacing
			//utils.LogDebug("Updating.. " + lastUpdateTime.Format(utils.DateTimeFormat) + " - " + fmt.Sprint(len(jobs)) + " jobs")
			lastUpdateTime = lastUpdateTime.Add(tickerUpdateSpacing)
			for _, j := range jobs {
				//utils.LogDebug("TICKING JOB ID : " + fmt.Sprint(j.Id))
				if time.Now().UTC().Sub(j.LastDoneTime) > j.Spacing {
					c <- j
				}
			}
		}
	}
}

func AddJob(d time.Duration, memory interface{}, f func(*interface{})) {
	JobAdder <- Job{
		Id:           len(jobs),
		Spacing:      d,
		LastDoneTime: time.Now().UTC(),
		Memory:       memory,
		Do:           f,
	}
	utils.LogInfo("Job added to background worker")
}
