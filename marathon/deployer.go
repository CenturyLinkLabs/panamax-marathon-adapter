package marathon

import (
	"log"
	"time"
)

func deployGroup(myGroup *deploymentGroup, timeout time.Duration) status {
	log.Printf("Deploying Group: %s", myGroup.id)

	var ctx = NewContext()

	// use a timeout channel
	timeoutchan := timeoutChannel(timeout)

	// set up deployment channel
	deploymentChannel := deployGroupChannel(myGroup, &ctx)

	for {
		select {
		case <-timeoutchan:
			log.Printf("Deployment timed out")
			return status{code: TIMEOUT}
		case <-deploymentChannel:
			if myGroup.Done() {
				return status{code: OK}
			}

			if myGroup.Failed() {
				log.Printf("Deployment Failed")
				return status{code: FAIL}
			}
		}
	}

}

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

func deployGroupChannel(myGroup *deploymentGroup, ctx *context) chan status {

	deploymentChannel := make(chan status, len(myGroup.deployments))
	for i := 0; i < len(myGroup.deployments); i++ {
		go deploy(deploymentChannel, &myGroup.deployments[i], ctx)
	}

	return deploymentChannel
}

func deploy(done chan status, deployment *deployment, ctx *context) {
	log.Printf("Starting Deployment: %s", deployment.name)

	for state := deployment.startingState; state != nil; {
		state = state(deployment, ctx)
	}

	done <- deployment.status
}
