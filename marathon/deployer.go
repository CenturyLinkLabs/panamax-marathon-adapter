package marathon

import (
	"log"
)

func deployGroupChannel(done chan status, deploychan chan *deployment, myGroup *deploymentGroup) {
	log.Printf("Deploying Group: %s", myGroup.id)
	var ctx = NewContext()

	for i:=0; i < len(myGroup.deployments); i++ {
		deploychan <- &(myGroup.deployments[i])
		go deployChannel(deploychan, &ctx)
	}

    	for {
		if (myGroup.Done()) {
			done <- status{code: OK}
		}

		if (myGroup.Failed()) {
			done <- status{code: FAIL}
		}
	}
}

func deployChannel(deployments chan *deployment, ctx *context) {
	deployment := <-deployments
	log.Printf("Starting Deployment: %s", deployment.name)

	for state := deployment.startingState; state != nil; {
        	state = state(deployment, ctx)
    	}
}

