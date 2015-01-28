package marathon

import (
	"log"

	"github.com/jbdalido/gomarathon"
)


func GroupDeployment(done chan bool, appchan chan *app, myGroup *group) {
	log.Printf("Group Deployment started")
	var ctx = NewContext()

	for i:=0; i < len(myGroup.apps); i++ {
		appchan <- &(myGroup.apps[i])
		go deployApp(appchan, &ctx)
	}

    	for {
		if (myGroup.Done() || myGroup.Failed()) {
			done <- true
		}
	}
}

func deployApp(apps chan *app, ctx *context) {
	var state stateFn
	app := <-apps

	for {
		switch (app.currentState) {
			case DONE:
				return
			case FAILED:
				return
			case PRE:
				state = app.preFn
				break
			case DEPLOY:
				state = app.deployFn
				break
			case POST:
				state = app.postFn
				break
			default:
				state = func(*gomarathon.Application, *context) int { return OK }

		}

		status := state(app.application, ctx)

		switch (status) {
			case OK:
				app.currentState +=1
				break
			case FAIL:
				app.currentState = FAILED
				return
		}
	}
}



