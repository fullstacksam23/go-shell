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
	for i, job := range JobList {
		marker := " "
		if i == len(JobList)-1 {
			marker = "+"
		} else if i == len(JobList)-2 {
			marker = "-"
		}

		printJob(job, marker, stdout)
	}
	//remove "Done" jobs
	reapJobs()
}

func ReapBeforePrompt(stdout io.Writer) {
	for i, job := range JobList {
		if job.Status == "Done" {
			marker := " "
			if i == len(JobList)-1 {
				marker = "+"
			} else if i == len(JobList)-2 {
				marker = "-"
			}
			printJob(job, marker, stdout)
		}
	}
	reapJobs()
}

func AddJob(pid int, command string) *BackgroundJob {
	if len(JobList) == 0 {
		nextJobId = 1
	}
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

func reapJobs() {
	for i := len(JobList) - 1; i >= 0; i-- {
		if JobList[i].Status == "Done" {
			JobList = append(JobList[:i], JobList[i+1:]...)
		}
	}
}

func printJob(job *BackgroundJob, marker string, stdout io.Writer) {
	fmt.Fprintf(stdout,
		"[%d]%s  %-24s%s\n",
		job.JobID,
		marker,
		job.Status,
		job.Command,
	)
}
