package marathon

import (
	"testing"
	"fmt"
)

/*
func TestJobStates(t *testing.T) {
	go handleJobs(job_channel)
	go handleApp(app_channel)

	// new template

	// analyze -- create the job struct list
	var testJob job
	testJob.preFn = preCondition
	testJob.deployFn = deploy
	testJob.postFn = postCondition

	myApp := new(app)
	myApp.jobs = []job{testJob}

	app_channel <- myApp

	myApp.Start()

	time.Sleep(4 * 1e9)
}
*/

func GroupDeployment(done chan bool, appchan chan *app, myGroup *group) {

	fmt.Println("Group Deployment service started")
	for i:=0; i < len(myGroup.apps); i++ {
		appchan <- &(myGroup.apps[i])
		go handleApp(appchan)
	}

    	for {
		if (myGroup.Done()) {
			done <- true
		}
	}
}


func TestGroupServiceApproach(t *testing.T) {

	var testApp, nextApp app
	testApp.name = "testApp"
	testApp.currentState = 1
	testApp.preFn = preConditionSlow
	testApp.deployFn = deploy
	testApp.postFn = postCondition

	nextApp.name = "nextApp"
	nextApp.currentState = 1
	nextApp.preFn = preCondition
	nextApp.deployFn = deploy
	nextApp.postFn = postCondition

	myGroup := new(group)
	myGroup.apps = []app{testApp, nextApp}

	done := make(chan bool)
	appchan := make(chan *app, len(myGroup.apps))

	go GroupDeployment(done, appchan, myGroup)

	<-done
	fmt.Println("Clean Up")
}
