package builtins

import (
	"fmt"
	"io"
)

type BackgroundJob struct {
	JobID   int
	Pid     int
	Command string
	Status  string
}

var JobList []*BackgroundJob
var nextJobId = 1

func HandleJobs(stdout io.Writer) {
	for i := len(JobList) - 1; i >= 0; i-- {
		job := JobList[i]
		marker := " "
		if i == len(JobList)-1 {
			marker = "+"
		} else if i == len(JobList)-2 {
			marker = "-"
		}

		fmt.Fprintf(stdout,
			"[%d]%s  %-24s%s\n",
			job.JobID,
			marker,
			job.Status,
			job.Command,
		)
		if job.Status == "Done" {
			JobList = append(JobList[:i], JobList[i+1:]...)
		}
	}
}

func AddJob(pid int, command string) *BackgroundJob {
	job := &BackgroundJob{
		JobID:   nextJobId,
		Pid:     pid,
		Command: command,
		Status:  "Running",
	}
	JobList = append(JobList, job)
	nextJobId++
	return job
}
