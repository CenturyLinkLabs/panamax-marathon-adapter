package marathon

import (
	"log"
	"time"
)

func timeoutChannel(duration time.Duration) chan bool {
	// make a timeout channel
	timeout := make(chan bool, 1)
	go func() {
		select {
		case <-time.After(duration):
			timeout <- true
		}
	}()
	return timeout
}

func deployGroupChannel(done chan status, deploychan chan *deployment, myGroup *deploymentGroup, timeoutDuration time.Duration) {
	log.Printf("Deploying Group: %s", myGroup.id)
	var ctx = NewContext()

	// use a timeout channel
	timeoutchan := timeoutChannel(timeoutDuration)

	for i:=0; i < len(myGroup.deployments); i++ {
		deploychan <- &(myGroup.deployments[i])
		go deployChannel(deploychan, &ctx)
	}

	for {
		select {
		case <-timeoutchan:
			log.Printf("Deployment timed out")
			done <- status{code: TIMEOUT}
		default:
			if (myGroup.Done()) {
				done <- status{code: OK}
			}

			if (myGroup.Failed()) {
				done <- status{code: FAIL}
			}
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

