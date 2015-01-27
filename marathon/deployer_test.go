package marathon

import (
	"testing"
	"fmt"
	"time"
	"math/rand"

	"github.com/stretchr/testify/assert"
)

func preCondition(*app) int {
	fmt.Println("preCondition Slow")
	time.Sleep(time.Duration(rand.Intn(100) * 100) * time.Millisecond)
	fmt.Println("preCondition After")
	return OK
}

func postCondition(*app) int {
	fmt.Println("postCondition")
	return OK
}

func deploy(*app) int {
	time.Sleep(time.Duration(rand.Intn(20) * 100) * time.Millisecond)
	fmt.Println("deploy")
	return OK
}

func TestGroupDeployment(t *testing.T) {

	var testApp, slowApp app

	slowApp.name = "slowApp"
	slowApp.currentState = 1
	slowApp.preFn = preCondition
	slowApp.deployFn = deploy
	slowApp.postFn = postCondition

	testApp.name = "testApp"
	testApp.currentState = 1
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
