package backgroundworker

import (
	"les-randoms/utils"
	"time"
)

// Ticker delay the worker waits before try to make every work
// -> Any work with less than this as time spacing will be done with this spacing
const tickerUpdateSpacing time.Duration = 5 * time.Minute

var JobAdder chan Job = make(chan Job)
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
	for range time.Tick(tickerUpdateSpacing) {
		//utils.LogDebug("Updating " + fmt.Sprint(len(jobs)) + " jobs")
		for _, j := range jobs {
			//utils.LogDebug("TICKING JOB ID : " + fmt.Sprint(j.Id))
			if time.Now().UTC().Sub(j.LastDoneTime) > j.Spacing {
				c <- j
			}
		}
	}
}

func AddJob(d time.Duration, memory interface{}, f func(*interface{})) {
	JobAdder <- Job{
		Id:           len(jobs),
		Spacing:      d,
		LastDoneTime: time.Now().UTC().Add(-d),
		Memory:       memory,
		Do:           f,
	}
	utils.LogInfo("Job added to background worker")
}
