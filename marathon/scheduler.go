package marathon

import (
	"time"
	"fmt"

	"github.com/jbdalido/gomarathon"
)

var job_channel = make(chan *job)
var app_channel = make(chan *app)

const (
	OK = 0
	FAIL = 1
)

const (
	PRE = 0
	DEPLOY = 1
	POST = 2
	DONE = 3
)

type stateFn func(*job) int

func preCondition(*job) int {
	fmt.Println("preCondition")
	return OK
}

func postCondition(*job) int {
	fmt.Println("postCondition")
	return OK
}

func deploy(*job) int {
	fmt.Println("deploy")
	return OK
}

func noOp(*job) int {
	fmt.Println("NoOp")
	return OK
}

func handleJobs(jobs chan *job) {
	var state stateFn
	j := <-jobs
	fmt.Println("Before Switch: ", j.currentState)
	switch (j.currentState) {
		case DONE:
			return
		case PRE:
			state = j.preFn
			break
		case DEPLOY:
			state = j.deployFn
			break
		case POST:
			state = j.postFn
			break
		default:
			state = noOp

	}

	if (state(j) == OK) {
		j.currentState +=1
	}

	fmt.Println("put it back")

	job_channel <- j
}

type job struct {
	preFn stateFn
	deployFn stateFn
	postFn stateFn
	currentState int
	application gomarathon.Application
	submitted time.Time
}

func (j *job) Complete() bool {
	fmt.Println("job state: %d", j.currentState)
	if (j.currentState == DONE) {
		return true
	}

	return false
}



type app struct {
	id string
	jobs []job
}

func handleApp(apps chan *app) {
	app := <-apps
	if (!app.Done()) {
		app_channel <- app
	}
}

func (a *app) Start() {
	for _, job := range a.jobs {
		job_channel <- &job
	}
}

func (a *app) Done() bool {
	completed := true
	for _, job := range a.jobs {
		completed = completed && job.Complete()
    	}
	if (completed) {
		//clean up
		fmt.Println("Clean up app")
	} else {
		fmt.Println("Try app again")
	}

	return completed
}


