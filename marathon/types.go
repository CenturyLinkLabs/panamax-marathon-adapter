package marathon

import (
	"time"

	"github.com/jbdalido/gomarathon"
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

func (j *app) complete() bool {
	if (j.currentState >= DONE) {
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
