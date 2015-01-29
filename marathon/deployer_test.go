package marathon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testSuccessState(deployment *deployment, ctx *context) stateFn  {
	deployment.status.code = OK
	return nil
}

func testFailState(deployment *deployment, ctx *context) stateFn  {
	deployment.status.code = FAIL
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

	go deployGroupChannel(done, appchan, myGroup)

	<-done
	assert.Equal(t, true, myGroup.Done())
}

func TestFailedGroupDeployment(t *testing.T) {
	var deployment1, deployment2 deployment

	deployment1.name = "slowApp"
	deployment1.startingState = testSuccessState
	deployment2.name = "testApp"
	deployment2.startingState = testFailState

	myGroup := new(deploymentGroup)
	myGroup.deployments = []deployment{deployment1, deployment2}

	done := make(chan status)
	appchan := make(chan *deployment, len(myGroup.deployments))

	go deployGroupChannel(done, appchan, myGroup)

	<-done
	assert.Equal(t, true, myGroup.Failed())
}
