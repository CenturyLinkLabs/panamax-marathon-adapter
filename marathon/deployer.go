package marathon

import (
	"log"
)

const (
	OK = 0
	WAIT = 1
	FAIL = 2

	PRE = 0
	DEPLOY = 1
	POST = 2
	DONE = 3
)

func GroupDeployment(done chan bool, appchan chan *app, myGroup *group) {

	log.Printf("Group Deployment started")
	for i:=0; i < len(myGroup.apps); i++ {
		appchan <- &(myGroup.apps[i])
		go deployApp(appchan)
	}

    	for {
		if (myGroup.Done()) {
			done <- true
		}
	}
}

func deployApp(apps chan *app) {
	var state stateFn
	j := <-apps

	for {
		log.Printf("Job: %s", j.name)
		log.Printf("Before Switch: %d", j.currentState)
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
				state = func(*app) int { return OK }

		}

		if (state(j) == OK) {
			j.currentState +=1
		}
	}
}



