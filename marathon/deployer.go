package marathon

import (
	"time"
	"fmt"

	"github.com/jbdalido/gomarathon"
)

const (
	OK = 0
	FAIL = 1
)

const (
	PRE = 1
	DEPLOY = 2
	POST = 3
	DONE = 4
)

type stateFn func(*app) int

func preCondition(*app) int {
	fmt.Println("preCondition")
	return OK
}

func preConditionSlow(*app) int {
	fmt.Println("preCondition Slow")
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("preCondition After")
	return OK
}

func postCondition(*app) int {
	fmt.Println("postCondition")
	return OK
}

func deploy(*app) int {
	fmt.Println("deploy")
	return OK
}

func noOp(*app) int {
	fmt.Println("NoOp")
	return OK
}

func handleApp(apps chan *app) {
	var state stateFn
	j := <-apps

	for {
		fmt.Println("Job: ", j.name)
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
	}
}

type app struct {
	name string
	preFn stateFn
	deployFn stateFn
	postFn stateFn
	currentState int
	previousState int
	application gomarathon.Application
	submitted time.Time
}

func (j *app) Complete() bool {
	if (j.currentState == DONE) {
		return true
	}

	return false
}


type group struct {
	id string
	apps []app
}

func (a *group) Done() bool {
	completed := true
	for _, app := range a.apps {
		completed = completed && app.Complete()
    	}

	return completed
}


