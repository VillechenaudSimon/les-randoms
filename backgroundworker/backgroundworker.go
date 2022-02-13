package backgroundworker

import (
	"time"
)

// tickerUpdateSpacing > tickerTickSpacing Or there will be problems
const tickerTickSpacing time.Duration = time.Minute
const tickerUpdateSpacing time.Duration = 10 * tickerTickSpacing

//const tickerTickSpacing time.Duration = 2 * time.Second
//const tickerUpdateSpacing time.Duration = 10 * tickerTickSpacing

var JobAdder chan Job = make(chan Job)
var lastUpdateTime time.Time
var jobs []*Job

type Job struct {
	Id           int
	Spacing      time.Duration
	LastDoneTime time.Time
	Do           func()
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
			job.LastDoneTime = time.Now().UTC()
			//utils.LogDebug("DOING JOB ID : " + fmt.Sprint(job.Id))
			job.Do()
		}
	}
}

func startTicker(c chan *Job) {
	for {
		//utils.LogDebug("Ticking.. " + lastUpdateTime.Format(utils.DateTimeFormat) + " - " + fmt.Sprint(len(jobs)) + " jobs")
		time.Sleep(tickerTickSpacing)
		if time.Now().UTC().Sub(lastUpdateTime) > tickerUpdateSpacing {
			//utils.LogDebug("Updating.. " + lastUpdateTime.Format(utils.DateTimeFormat) + " - " + fmt.Sprint(len(jobs)) + " jobs" + " - J0:" + fmt.Sprint(jobs[0].Id) + " - J1:" + fmt.Sprint(jobs[1].Id))
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

func AddJob(d time.Duration, f func()) {
	JobAdder <- Job{
		Id:           len(jobs),
		Spacing:      d,
		LastDoneTime: time.Now().UTC(),
		Do:           f,
	}
	//utils.LogDebug("Job Added - " + fmt.Sprint(len(jobs)))
}
