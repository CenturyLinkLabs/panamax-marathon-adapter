package marathon

import (
	"testing"
)


func TestJobStates(t *testing.T) {
	go handleJobs(job_channel)
	go handleApp(app_channel)

	// new template

	// analyze -- create the job struct list
	var testJob job
  testJob.stateFns[0] = preCondition
	testJob.stateFns[1] = deploy
	testJob.stateFns[2] = postCondition

	myApp := new(app)
	myApp.jobs = []job{testJob}

	app_channel <- myApp

	myApp.Start()
}
