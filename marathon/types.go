package marathon

import (
	"time"

	"github.com/jbdalido/gomarathon"
)

const (
	FAILED = -1
	PRE = 0
	DEPLOY = 1
	POST = 2
	DONE = 3

	OK = 0
	WAIT = 1
	FAIL = 2
)


type stateFn func(*app) int

type app struct {
	name string
	preFn stateFn
	deployFn stateFn
	postFn stateFn
	currentState int
	application *gomarathon.Application
	submitted time.Time
}

func (a *app) complete() bool {
	if (a.currentState >= DONE) {
		return true
	}

	return false
}

func (a *app) failed() bool {
	if (a.currentState == FAILED) {
		return true
	}
	return false
}

type group struct {
	id string
	apps []app
}

func (g *group) Done() bool {
	completed := true
	for _, app := range g.apps {
		completed = completed && app.complete()
    	}

	return completed
}

func (g *group) Failed() bool {
	failed := true
	for _, app := range g.apps {
		failed = failed && app.failed()
    	}

	return failed
}
