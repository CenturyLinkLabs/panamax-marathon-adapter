package marathon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var timeoutDuration = time.Second * 10

func testSuccessState(deployment *deployment, ctx *context) stateFn  {
	deployment.status.code = OK
	return nil
}

func testFailState(deployment *deployment, ctx *context) stateFn {
	deployment.status.code = FAIL
	return nil
}

func testTimeoutState(deployment *deployment, ctx *context) stateFn {
	time.Sleep(time.Second * 15)
	deployment.status.code = OK
	return nil
}

func TestGroupDeployment(t *testing.T) {
	var deployment1, deployment2 deployment

	deployment1.name = "slowApp"
	deployment1.startingState = testSuccessState
	deployment2.name = "testApp"
	deployment2.startingState = testSuccessState

	myGroup := new(deploymentGroup)
	myGroup.deployments = []deployment{deployment1, deployment2}

	done := make(chan status)
	appchan := make(chan *deployment, len(myGroup.deployments))

	go deployGroupChannel(done, appchan, myGroup, timeoutDuration)

	<-done
	assert.Equal(t, true, myGroup.Done())
}

func TestFailedGroupDeployment(t *testing.T) {
	var deployment1, deployment2 deployment

	deployment1.name = "slowApp"
	deployment1.startingState = testSuccessState
	deployment2.name = "failApp"
	deployment2.startingState = testFailState

	myGroup := new(deploymentGroup)
	myGroup.deployments = []deployment{deployment1, deployment2}

	done := make(chan status)
	appchan := make(chan *deployment, len(myGroup.deployments))

	go deployGroupChannel(done, appchan, myGroup, timeoutDuration)

	<-done
	assert.Equal(t, true, myGroup.Failed())
}

func TestTimedoutGroupDeployment(t *testing.T) {
	var deployment1, deployment2 deployment

	deployment1.name = "testApp"
	deployment1.startingState = testSuccessState
	deployment2.name = "timeoutApp"
	deployment2.startingState = testTimeoutState

	myGroup := new(deploymentGroup)
	myGroup.deployments = []deployment{deployment1, deployment2}

	done := make(chan status)
	appchan := make(chan *deployment, len(myGroup.deployments))

	go deployGroupChannel(done, appchan, myGroup, timeoutDuration)

	status := <-done

	assert.Equal(t, OK, myGroup.deployments[0].status.code)
	assert.Equal(t, TIMEOUT, status.code)
}
