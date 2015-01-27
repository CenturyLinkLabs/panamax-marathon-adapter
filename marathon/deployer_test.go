package marathon

import (
	"testing"
	"fmt"
	"time"
	"math/rand"

	"github.com/jbdalido/gomarathon"
	"github.com/stretchr/testify/assert"
)

func preCondition(*gomarathon.Application, *context) int {
	fmt.Println("preCondition Slow")
	time.Sleep(time.Duration(rand.Intn(10) * 100) * time.Millisecond)
	fmt.Println("preCondition After")
	return OK
}

func postCondition(*gomarathon.Application, *context) int {
	fmt.Println("postCondition")
	return OK
}

func deploy(*gomarathon.Application, *context) int {
	time.Sleep(time.Duration(rand.Intn(100) * 100) * time.Millisecond)
	fmt.Println("deploy")
	return OK
}

func deployFail(*gomarathon.Application, *context) int {
	fmt.Println("failed deployment")
	return FAIL
}

func TestGroupDeployment(t *testing.T) {
	var testApp, slowApp app

	slowApp.name = "slowApp"
	slowApp.preFn = preCondition
	slowApp.deployFn = deploy
	slowApp.postFn = postCondition

	testApp.name = "testApp"
	testApp.preFn = preCondition
	testApp.deployFn = deploy
	testApp.postFn = postCondition

	myGroup := new(group)
	myGroup.apps = []app{slowApp, testApp}

	done := make(chan bool)
	appchan := make(chan *app, len(myGroup.apps))

	go GroupDeployment(done, appchan, myGroup)

	<-done
	assert.Equal(t, true, myGroup.Done())
}

func TestFailedGroupDeployment(t *testing.T) {
	var failApp app

	failApp.name = "failApp"
	failApp.preFn = preCondition
	failApp.deployFn = deployFail
	failApp.postFn = postCondition

	myGroup := new(group)
	myGroup.apps = []app{failApp}

	done := make(chan bool)
	appchan := make(chan *app, len(myGroup.apps))

	go GroupDeployment(done, appchan, myGroup)

	<-done
	assert.Equal(t, true, myGroup.Failed())


}
