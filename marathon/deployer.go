package marathon

import (
	"log"
)


func GroupDeployment(done chan bool, appchan chan *app, myGroup *group) {

	log.Printf("Group Deployment started")
	for i:=0; i < len(myGroup.apps); i++ {
		appchan <- &(myGroup.apps[i])
		go deployApp(appchan)
	}

    	for {
		if (myGroup.Done() || myGroup.Failed()) {
			done <- true
		}
	}
}

func deployApp(apps chan *app) {
	var state stateFn
	j := <-apps

	for {
		log.Printf("Job: %s", j.name)
		switch (j.currentState) {
			case DONE:
				return
			case FAILED:
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

		status := state(j)

		switch (status) {
			case OK:
				j.currentState +=1
				break
			case FAIL:
				j.currentState = FAILED
				return
		}
	}
}



